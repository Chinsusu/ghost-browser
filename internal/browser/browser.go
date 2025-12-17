package browser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/user/ghost-browser/internal/profile"
	"github.com/user/ghost-browser/internal/proxy"
)

type Instance struct {
	ProfileID string
	Process   *exec.Cmd
	ScriptPath string
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

	// Create fingerprint spoofing script
	scriptPath, err := createSpoofingScript(p, userDataDir)
	if err != nil {
		return fmt.Errorf("failed to create spoofing script: %w", err)
	}

	// Build Edge launch arguments
	args := []string{
		"--user-data-dir=" + userDataDir,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"--disable-blink-features=AutomationControlled",
		"--disable-web-security",
		"--allow-running-insecure-content",
		"--user-script=" + scriptPath,
		"about:blank",
	}

	// Add proxy if configured
	if p.ProxyID != nil {
		if proxyConfig, err := m.proxyManager.GetByID(*p.ProxyID); err == nil {
			args = append(args, "--proxy-server="+proxyConfig.ToURL())
		}
	}

	// Launch Edge process
	cmd := exec.Command(edgePath, args...)
	err = cmd.Start()
	if err != nil {
		os.Remove(scriptPath) // Cleanup script on failure
		return fmt.Errorf("failed to launch Edge: %w", err)
	}

	m.instances[profileID] = &Instance{
		ProfileID:  profileID,
		Process:    cmd,
		ScriptPath: scriptPath,
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

	// Kill Edge process
	if instance.Process != nil {
		instance.Process.Process.Kill()
	}

	// Cleanup script file
	if instance.ScriptPath != "" {
		os.Remove(instance.ScriptPath)
	}

	delete(m.instances, profileID)
	return nil
}

func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, instance := range m.instances {
		// Kill Edge process
		if instance.Process != nil {
			instance.Process.Process.Kill()
		}

		// Cleanup script file
		if instance.ScriptPath != "" {
			os.Remove(instance.ScriptPath)
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

// createSpoofingScript creates the fingerprint spoofing JavaScript file
func createSpoofingScript(p *profile.Profile, userDataDir string) (string, error) {
	script := generateSpoofScript(p.Fingerprint)
	scriptPath := filepath.Join(userDataDir, "ghost-spoof.js")
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	return scriptPath, err
}
