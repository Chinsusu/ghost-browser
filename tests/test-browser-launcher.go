package main

import (
	"fmt"
	"log"
	"time"

	"github.com/user/ghost-browser/internal/browser"
	"github.com/user/ghost-browser/internal/database"
	"github.com/user/ghost-browser/internal/profile"
	"github.com/user/ghost-browser/internal/proxy"
)

func main() {
	fmt.Println("🧪 === Ghost Browser Launcher Test ===")
	fmt.Println()

	// Initialize database
	fmt.Println("🔧 Initializing database...")
	db, err := database.New("test-ghost-browser.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("✅ Database initialized")

	// Initialize managers
	profileManager := profile.NewManager(db)
	proxyManager := proxy.NewManager(db)
	browserManager := browser.NewManager(profileManager, proxyManager)

	// Create test profile
	fmt.Println("\n👤 Creating test profile...")
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

	// Launch browser
	fmt.Println("\n🚀 Launching Edge browser...")
	fmt.Println("This will open Edge with spoofed fingerprint.")
	fmt.Println("Two tabs will open:")
	fmt.Println("  1. browserleaks.com/javascript - Check Navigator values")
	fmt.Println("  2. creepjs.com - Full fingerprint analysis")
	fmt.Println()

	err = browserManager.Launch(testProfile.ID)
	if err != nil {
		log.Printf("❌ Failed to launch browser: %v", err)
		fmt.Println("\nThis might be normal if Edge is not installed or CDP is not available.")
		fmt.Println("You can still test manually by:")
		fmt.Println("1. Opening Edge browser")
		fmt.Println("2. Going to browserleaks.com/javascript")
		fmt.Println("3. Checking if values match the fingerprint above")
	} else {
		fmt.Println("✅ Browser launched successfully!")
		fmt.Println()
		fmt.Println("📋 Test Instructions:")
		fmt.Println("1. Check browserleaks.com tab:")
		fmt.Printf("   - hardwareConcurrency should be: %d\n", fp.Navigator.HardwareConcurrency)
		fmt.Printf("   - deviceMemory should be: %d\n", fp.Navigator.DeviceMemory)
		fmt.Printf("   - platform should be: %s\n", fp.Navigator.Platform)
		fmt.Println("   - webdriver should be: false or undefined")
		fmt.Println()
		fmt.Println("2. Check CreepJS tab for full fingerprint analysis")
		fmt.Println()
		fmt.Println("3. Open DevTools (F12) and check Console for:")
		fmt.Println("   [Ghost Browser] ✅ Navigator spoofed")
		fmt.Println("   [Ghost Browser] ✅ Screen spoofed")
		fmt.Println("   [Ghost Browser] ✅ WebGL spoofed")
		fmt.Println("   [Ghost Browser] Fingerprint spoofing active")
		fmt.Println()
		fmt.Println("⏳ Browser will stay open for 30 seconds for testing...")
		time.Sleep(30 * time.Second)
	}

	// Cleanup
	fmt.Println("\n🧹 Cleaning up...")
	err = profileManager.Delete(testProfile.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete test profile: %v", err)
	}

	// Remove test database
	dataDir, _ := database.GetDataDir()
	fmt.Printf("Test database location: %s/test-ghost-browser.db\n", dataDir)

	fmt.Println("✅ Test completed!")
	fmt.Println()
	fmt.Println("📊 Results to report:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Were fingerprint values spoofed correctly? (Yes/No)")
	fmt.Println("3. Any errors in terminal? (Yes/No)")
	fmt.Println("4. Console messages visible in DevTools? (Yes/No)")
}