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
	"github.com/user/ghost-browser/internal/profile"
	"github.com/user/ghost-browser/internal/proxy"
)

type Instance struct {
	ProfileID string
	Context   context.Context
	Cancel    context.CancelFunc
}

type Manager struct {
	profileManager *profile.Manager
	proxyManager   *proxy.Manager
	instances      map[string]*Instance
	mu             sync.RWMutex
}

func NewManager(pm *profile.Manager, proxym *proxy.Manager) *Manager {
	return &Manager{
		profileManager: pm,
		proxyManager:   proxym,
		instances:      make(map[string]*Instance),
	}
}

func (m *Manager) Launch(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.instances[profileID]; exists {
		return fmt.Errorf("browser already running for profile %s", profileID)
	}

	p, err := m.profileManager.GetByID(profileID)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	edgePath, err := findEdgePath()
	if err != nil {
		return fmt.Errorf("failed to find Edge: %w", err)
	}

	// Create user data directory
	userDataDir := filepath.Join(p.DataDir, "EdgeData")
	os.MkdirAll(userDataDir, 0755)

	// Setup ChromeDP options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(edgePath),
		chromedp.UserDataDir(userDataDir),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.WindowSize(1920, 1080),
	)

	// Add proxy if configured
	if p.ProxyID != nil {
		if proxyConfig, err := m.proxyManager.GetByID(*p.ProxyID); err == nil {
			opts = append(opts, chromedp.Flag("proxy-server", proxyConfig.ToURL()))
		}
	}

	// Create ChromeDP context
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Generate spoofing script
	spoofScript := generateSpoofScript(p.Fingerprint)

	// Launch browser with pre-load script injection
	err = chromedp.Run(ctx,
		// CRITICAL: Add script BEFORE any navigation
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(spoofScript).Do(ctx)
			return err
		}),
		// Navigate to about:blank to initialize
		chromedp.Navigate("about:blank"),
	)

	if err != nil {
		cancel()
		allocCancel()
		return fmt.Errorf("failed to launch browser with ChromeDP: %w", err)
	}

	// Store instance with combined cancel function
	combinedCancel := func() {
		cancel()
		allocCancel()
	}

	m.instances[profileID] = &Instance{
		ProfileID: profileID,
		Context:   ctx,
		Cancel:    combinedCancel,
	}

	m.profileManager.UpdateLastUsed(profileID)
	return nil
}

func (m *Manager) Close(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance, exists := m.instances[profileID]
	if !exists {
		return fmt.Errorf("no browser running for profile %s", profileID)
	}

	// Cancel ChromeDP context
	if instance.Cancel != nil {
		instance.Cancel()
	}

	delete(m.instances, profileID)
	return nil
}

func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, instance := range m.instances {
		// Cancel ChromeDP context
		if instance.Cancel != nil {
			instance.Cancel()
		}
		delete(m.instances, id)
	}
}

func (m *Manager) GetRunning() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.instances))
	for id := range m.instances {
		ids = append(ids, id)
	}
	return ids
}

func (m *Manager) IsRunning(profileID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.instances[profileID]
	return exists
}

func findEdgePath() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Edge only supported on Windows")
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

	return "", fmt.Errorf("Edge browser not found")
}

// NavigateToURL navigates the browser to a specific URL
func (m *Manager) NavigateToURL(profileID, url string) error {
	m.mu.RLock()
	instance, exists := m.instances[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no browser running for profile %s", profileID)
	}

	return chromedp.Run(instance.Context, chromedp.Navigate(url))
}
