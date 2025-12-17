package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"

	"github.com/user/ghost-browser/internal/fingerprint"
	"github.com/user/ghost-browser/internal/profile"
	"github.com/user/ghost-browser/internal/proxy"
)

// Instance represents a running browser instance
type Instance struct {
	ProfileID   string
	Ctx         context.Context
	Cancel      context.CancelFunc
	AllocCancel context.CancelFunc
	Fingerprint *fingerprint.Fingerprint
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
		gen := fingerprint.NewGenerator()
		fp = gen.Generate(nil)
		// Save fingerprint to profile
		p.Fingerprint = fp
		m.profileManager.Update(p)
	}

	// Generate spoof script BEFORE creating browser
	spoofScript := GenerateSpoofScript(fp)

	// Build Chrome options
	opts := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath(edgePath),
		chromedp.UserDataDir(p.DataDir),
		chromedp.Flag("headless", false),

		// Anti-detection
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-popup-blocking", true),

		// WebRTC protection
		chromedp.Flag("disable-webrtc", true),
		chromedp.Flag("enforce-webrtc-ip-permission-check", true),
		chromedp.Flag("webrtc-ip-handling-policy", "disable_non_proxied_udp"),
		chromedp.Flag("disable-features", "WebRtcHideLocalIpsWithMdns,WebRTC,TranslateUI"),

		// Window size matching fingerprint
		chromedp.WindowSize(fp.Screen.Width, fp.Screen.Height),

		// Disable automation flags
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-ipc-flooding-protection", true),
	}

	// Add proxy if configured
	if p.ProxyID != nil {
		proxyConfig, err := m.proxyManager.GetByID(*p.ProxyID)
		if err == nil && proxyConfig != nil {
			opts = append(opts, chromedp.ProxyServer(proxyConfig.ToURL()))
		}
	}

	// Create allocator context
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// Create browser context with logging disabled for production
	ctx, cancel := chromedp.NewContext(allocCtx)

	// CRITICAL: Set up script injection for ALL new targets/pages
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *target.EventTargetCreated:
			// New tab/window created - inject script
			go func() {
				tCtx, tCancel := chromedp.NewContext(ctx, chromedp.WithTargetID(ev.TargetInfo.TargetID))
				defer tCancel()
				
				// Small delay to ensure target is ready
				time.Sleep(50 * time.Millisecond)
				
				chromedp.Run(tCtx,
					page.AddScriptToEvaluateOnNewDocument(spoofScript),
				)
			}()
		}
	})

	// Launch browser and inject script
	err = chromedp.Run(ctx,
		// Add script to evaluate on new document - THIS IS KEY
		page.AddScriptToEvaluateOnNewDocument(spoofScript),
		
		// Enable page events
		page.Enable(),
		
		// Small delay to ensure script is registered
		chromedp.Sleep(100*time.Millisecond),
		
		// Navigate to blank page first
		chromedp.Navigate("about:blank"),
		
		// Wait for page to be ready
		chromedp.WaitReady("body", chromedp.ByQuery),
	)

	if err != nil {
		cancel()
		allocCancel()
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	// Store instance
	m.instances[profileID] = &Instance{
		ProfileID:   profileID,
		Ctx:         ctx,
		Cancel:      cancel,
		AllocCancel: allocCancel,
		Fingerprint: fp,
	}

	// Update last used
	m.profileManager.UpdateLastUsed(profileID)

	return nil
}

// LaunchWithURL launches browser and navigates to URL
func (m *Manager) LaunchWithURL(profileID, url string) error {
	err := m.Launch(profileID)
	if err != nil {
		return err
	}

	m.mu.RLock()
	instance, exists := m.instances[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("instance not found after launch")
	}

	// Navigate to URL
	return chromedp.Run(instance.Ctx,
		chromedp.Navigate(url),
	)
}

// Close closes a browser instance
func (m *Manager) Close(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance, exists := m.instances[profileID]
	if !exists {
		return fmt.Errorf("no browser running for profile %s", profileID)
	}

	// Cancel contexts
	instance.Cancel()
	
	// Small delay before canceling allocator
	time.Sleep(100 * time.Millisecond)
	
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
		time.Sleep(50 * time.Millisecond)
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

// NavigateTo navigates a running browser to a URL
func (m *Manager) NavigateTo(profileID, url string) error {
	m.mu.RLock()
	instance, exists := m.instances[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no browser running for profile %s", profileID)
	}

	return chromedp.Run(instance.Ctx,
		chromedp.Navigate(url),
	)
}

// GetFingerprint returns the fingerprint for a running browser
func (m *Manager) GetFingerprint(profileID string) (*fingerprint.Fingerprint, error) {
	m.mu.RLock()
	instance, exists := m.instances[profileID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no browser running for profile %s", profileID)
	}

	return instance.Fingerprint, nil
}

// findEdgePath finds the Edge executable
func findEdgePath() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Windows only - Edge not available")
	}

	// Check common paths
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

	// Try PATH
	if path, err := exec.LookPath("msedge"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("msedge.exe"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Microsoft Edge not found. Please ensure Edge is installed.")
}
