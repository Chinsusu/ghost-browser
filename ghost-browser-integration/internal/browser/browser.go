package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"ghost-browser/internal/fingerprint"
	"ghost-browser/internal/profile"
	"ghost-browser/internal/proxy"
)

// Instance represents a running browser instance
type Instance struct {
	ProfileID  string
	Cancel     context.CancelFunc
	AllocCancel context.CancelFunc
}

// Manager manages browser instances
type Manager struct {
	profileManager *profile.Manager
	proxyManager   *proxy.Manager
	instances      map[string]*Instance
	mu             sync.RWMutex
}

// NewManager creates a new browser manager
func NewManager(pm *profile.Manager, proxym *proxy.Manager) *Manager {
	return &Manager{
		profileManager: pm,
		proxyManager:   proxym,
		instances:      make(map[string]*Instance),
	}
}

// Launch launches a browser with the given profile
func (m *Manager) Launch(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already running
	if _, exists := m.instances[profileID]; exists {
		return fmt.Errorf("browser already running for profile %s", profileID)
	}

	// Get profile
	p, err := m.profileManager.GetByID(profileID)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// Find Edge
	edgePath, err := findEdgePath()
	if err != nil {
		return fmt.Errorf("failed to find Edge: %w", err)
	}

	// Get fingerprint from profile
	fp := p.Fingerprint
	if fp == nil {
		// Generate new fingerprint if not exists
		gen := fingerprint.NewGenerator()
		fp = gen.Generate(nil)
	}

	// Build Chrome options
	opts := buildChromeOptions(edgePath, p.DataDir, fp)

	// Add proxy if configured
	if p.ProxyID != nil {
		proxyConfig, err := m.proxyManager.GetByID(*p.ProxyID)
		if err == nil && proxyConfig != nil {
			opts = append(opts, chromedp.ProxyServer(proxyConfig.ToURL()))
		}
	}

	// Create contexts
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Generate spoof script
	spoofScript := GenerateSpoofScript(fp)

	// Launch browser
	err = chromedp.Run(ctx,
		page.AddScriptToEvaluateOnNewDocument(spoofScript),
		chromedp.Navigate("about:blank"),
	)
	if err != nil {
		cancel()
		allocCancel()
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	// Store instance
	m.instances[profileID] = &Instance{
		ProfileID:   profileID,
		Cancel:      cancel,
		AllocCancel: allocCancel,
	}

	// Update last used
	m.profileManager.UpdateLastUsed(profileID)

	return nil
}

// Close closes a browser instance
func (m *Manager) Close(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance, exists := m.instances[profileID]
	if !exists {
		return fmt.Errorf("no browser running for profile %s", profileID)
	}

	instance.Cancel()
	instance.AllocCancel()
	delete(m.instances, profileID)

	return nil
}

// CloseAll closes all browser instances
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, instance := range m.instances {
		instance.Cancel()
		instance.AllocCancel()
		delete(m.instances, id)
	}
}

// GetRunning returns list of running profile IDs
func (m *Manager) GetRunning() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.instances))
	for id := range m.instances {
		ids = append(ids, id)
	}
	return ids
}

// IsRunning checks if a profile's browser is running
func (m *Manager) IsRunning(profileID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.instances[profileID]
	return exists
}

// buildChromeOptions builds chromedp options
func buildChromeOptions(edgePath, userDataDir string, fp *fingerprint.Fingerprint) []chromedp.ExecAllocatorOption {
	return []chromedp.ExecAllocatorOption{
		chromedp.ExecPath(edgePath),
		chromedp.UserDataDir(userDataDir),
		chromedp.Flag("headless", false),

		// Anti-detection
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-extensions", true),

		// WebRTC protection
		chromedp.Flag("disable-webrtc", true),
		chromedp.Flag("enforce-webrtc-ip-permission-check", true),
		chromedp.Flag("webrtc-ip-handling-policy", "disable_non_proxied_udp"),
		chromedp.Flag("disable-features", "WebRtcHideLocalIpsWithMdns,WebRTC"),

		// Window size
		chromedp.WindowSize(fp.Screen.Width, fp.Screen.Height),

		// Disable GPU for consistency
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("disable-software-rasterizer", true),
	}
}

// findEdgePath finds the Edge executable
func findEdgePath() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Windows only")
	}

	paths := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("LocalAppData"), "Microsoft", "Edge", "Application", "msedge.exe"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	if path, err := exec.LookPath("msedge"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Edge not found")
}
