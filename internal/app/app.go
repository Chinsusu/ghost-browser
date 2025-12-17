package app

import (
	"context"
	"log"

	"github.com/user/ghost-browser/internal/ai"
	"github.com/user/ghost-browser/internal/browser"
	"github.com/user/ghost-browser/internal/database"
	"github.com/user/ghost-browser/internal/profile"
	"github.com/user/ghost-browser/internal/proxy"
)

type App struct {
	ctx            context.Context
	db             *database.Database
	profileManager *profile.Manager
	proxyManager   *proxy.Manager
	browserManager *browser.Manager
	aiEngine       *ai.Engine
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	db, err := database.New("ghost-browser.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	a.db = db

	if err := a.db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	a.profileManager = profile.NewManager(db)
	a.proxyManager = proxy.NewManager(db)
	a.browserManager = browser.NewManager(a.profileManager, a.proxyManager)
	a.aiEngine = ai.NewEngine(db)

	// Ensure at least one default profile exists
	if err := a.ensureDefaultProfile(); err != nil {
		log.Printf("Warning: Failed to create default profile: %v", err)
	}

	log.Println("Ghost Browser started successfully")
}

func (a *App) Shutdown(ctx context.Context) {
	if a.browserManager != nil {
		a.browserManager.CloseAll()
	}
	if a.db != nil {
		a.db.Close()
	}
	log.Println("Ghost Browser shutdown complete")
}

// ==================== Profile API ====================

func (a *App) GetProfiles() ([]*profile.Profile, error) {
	return a.profileManager.GetAll()
}

func (a *App) GetProfile(id string) (*profile.Profile, error) {
	return a.profileManager.GetByID(id)
}

func (a *App) CreateProfile(name string, opts *profile.CreateOptions) (*profile.Profile, error) {
	return a.profileManager.Create(name, opts)
}

func (a *App) UpdateProfile(p *profile.Profile) error {
	return a.profileManager.Update(p)
}

func (a *App) DeleteProfile(id string) error {
	return a.profileManager.Delete(id)
}

func (a *App) GenerateRandomProfile() (*profile.Profile, error) {
	return a.profileManager.GenerateRandom()
}

func (a *App) DuplicateProfile(id string) (*profile.Profile, error) {
	return a.profileManager.Duplicate(id)
}

func (a *App) ExportProfile(id, path string) error {
	return a.profileManager.Export(id, path)
}

func (a *App) ImportProfile(path string) (*profile.Profile, error) {
	return a.profileManager.Import(path)
}

// ==================== Proxy API ====================

func (a *App) GetProxies() ([]*proxy.Proxy, error) {
	return a.proxyManager.GetAll()
}

func (a *App) AddProxy(p *proxy.Proxy) error {
	return a.proxyManager.Add(p)
}

func (a *App) DeleteProxy(id string) error {
	return a.proxyManager.Delete(id)
}

func (a *App) CheckProxy(id string) (*proxy.CheckResult, error) {
	return a.proxyManager.Check(id)
}

func (a *App) CheckAllProxies() ([]*proxy.CheckResult, error) {
	return a.proxyManager.CheckAll()
}

func (a *App) ImportProxies(text, format string) (int, error) {
	return a.proxyManager.ImportFromText(text, format)
}

// ==================== Browser API ====================

func (a *App) LaunchBrowser(profileID string) error {
	return a.browserManager.Launch(profileID)
}

func (a *App) CloseBrowser(profileID string) error {
	return a.browserManager.Close(profileID)
}

func (a *App) GetRunningBrowsers() []string {
	return a.browserManager.GetRunning()
}

// ==================== AI API ====================

func (a *App) GetPersonality(profileID string) (*ai.Personality, error) {
	return a.aiEngine.GetPersonality(profileID)
}

func (a *App) UpdatePersonality(profileID string, p *ai.Personality) error {
	return a.aiEngine.UpdatePersonality(profileID, p)
}

func (a *App) GeneratePersonality() (*ai.Personality, error) {
	return a.aiEngine.GenerateRandom()
}

func (a *App) Chat(profileID, message string) (string, error) {
	return a.aiEngine.Chat(profileID, message)
}

func (a *App) GetSchedule(profileID string) (*ai.Schedule, error) {
	return a.aiEngine.GetSchedule(profileID)
}

func (a *App) UpdateSchedule(profileID string, s *ai.Schedule) error {
	return a.aiEngine.UpdateSchedule(profileID, s)
}

// ==================== Internal Methods ====================

func (a *App) ensureDefaultProfile() error {
	// Check if any profiles exist
	profiles, err := a.profileManager.GetAll()
	if err != nil {
		return err
	}

	// If no profiles exist, create a default one
	if len(profiles) == 0 {
		log.Println("No profiles found, creating default profile...")
		_, err := a.profileManager.Create("Default Profile", &profile.CreateOptions{
			OS:      "windows",
			Browser: "edge",
			Notes:   "Auto-generated default profile with fingerprint spoofing",
		})
		if err != nil {
			return err
		}
		log.Println("✅ Default profile created successfully")
	}

	return nil
}
