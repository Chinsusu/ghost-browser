package browser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/user/ghost-browser/internal/profile"
	"github.com/user/ghost-browser/internal/proxy"
)

type Instance struct {
	ProfileID string
	Browser   *rod.Browser
	Launcher  *launcher.Launcher
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

	l := launcher.New().
		Bin(edgePath).
		UserDataDir(p.DataDir).
		Headless(false).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars").
		Set("no-first-run").
		Set("no-default-browser-check")

	if p.ProxyID != nil {
		if proxyConfig, err := m.proxyManager.GetByID(*p.ProxyID); err == nil {
			l.Set("proxy-server", proxyConfig.ToURL())
		}
	}

	controlURL := l.MustLaunch()
	browser := rod.New().ControlURL(controlURL).MustConnect()

	// Inject spoofing scripts
	injectSpoofingScripts(browser, p)

	m.instances[profileID] = &Instance{
		ProfileID: profileID,
		Browser:   browser,
		Launcher:  l,
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

	instance.Browser.MustClose()
	delete(m.instances, profileID)
	return nil
}

func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, instance := range m.instances {
		instance.Browser.MustClose()
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
