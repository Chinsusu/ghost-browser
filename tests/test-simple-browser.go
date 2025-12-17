package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/user/ghost-browser/internal/database"
	"github.com/user/ghost-browser/internal/profile"
)

func main() {
	fmt.Println("🧪 === Simple Browser Test (No Rod) ===")
	fmt.Println()

	// Initialize database
	fmt.Println("🔧 Initializing database...")
	db, err := database.New("test-simple-browser.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("✅ Database initialized")

	// Create test profile
	fmt.Println("\n👤 Creating test profile...")
	profileManager := profile.NewManager(db)
	testProfile, err := profileManager.GenerateRandom()
	if err != nil {
		log.Fatal("Failed to create profile:", err)
	}
	fmt.Printf("✅ Created profile: %s (ID: %s)\n", testProfile.Name, testProfile.ID)

	// Display fingerprint info
	fmt.Println("\n🎭 Fingerprint details:")
	fp := testProfile.Fingerprint
	fmt.Printf("  • User Agent: %s\n", fp.Navigator.UserAgent)
	fmt.Printf("  • Screen: %dx%d\n", fp.Screen.Width, fp.Screen.Height)
	fmt.Printf("  • Hardware Cores: %d\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("  • Device Memory: %d GB\n", fp.Navigator.DeviceMemory)
	fmt.Printf("  • WebGL Vendor: %s\n", fp.WebGL.Vendor)
	fmt.Printf("  • Timezone: %s (offset: %d)\n", fp.Timezone.Timezone, fp.Timezone.TimezoneOffset)

	// Find Edge
	fmt.Println("\n🔍 Looking for Edge browser...")
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("Edge browser not found:", err)
	}
	fmt.Printf("✅ Found Edge at: %s\n", edgePath)

	// Launch Edge manually
	fmt.Println("\n🚀 Launching Edge browser manually...")
	fmt.Println("This will open Edge with basic parameters.")
	fmt.Println("Note: Fingerprint spoofing requires CDP injection (not available in this test)")
	fmt.Println()

	// Create user data directory
	userDataDir := filepath.Join(testProfile.DataDir, "EdgeData")
	os.MkdirAll(userDataDir, 0755)

	// Launch Edge with basic parameters
	args := []string{
		"--user-data-dir=" + userDataDir,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"--disable-blink-features=AutomationControlled",
		"browserleaks.com/javascript",
	}

	cmd := exec.Command(edgePath, args...)
	err = cmd.Start()
	if err != nil {
		log.Fatal("Failed to launch Edge:", err)
	}

	fmt.Println("✅ Edge launched successfully!")
	fmt.Println()
	fmt.Println("📋 Manual Test Instructions:")
	fmt.Println("1. Edge should have opened with browserleaks.com/javascript")
	fmt.Println("2. Check the Navigator values on the page")
	fmt.Println("3. Note: Values will NOT be spoofed (this is expected)")
	fmt.Println("4. This test only verifies Edge can be launched")
	fmt.Println()
	fmt.Printf("Expected values (if spoofing worked):\n")
	fmt.Printf("  - hardwareConcurrency: %d\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("  - deviceMemory: %d\n", fp.Navigator.DeviceMemory)
	fmt.Printf("  - platform: %s\n", fp.Navigator.Platform)
	fmt.Println()

	// Cleanup
	fmt.Println("🧹 Cleaning up...")
	err = profileManager.Delete(testProfile.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete test profile: %v", err)
	}

	fmt.Println("✅ Test completed!")
	fmt.Println()
	fmt.Println("📊 Results to report:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Did browserleaks.com load? (Yes/No)")
	fmt.Println("3. Any errors in terminal? (Yes/No)")
	fmt.Println("4. Edge process started successfully? (Yes/No)")
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