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
	fmt.Println("🧪 === SIMPLE FIXED BROWSER TEST ===")
	fmt.Println()

	// Initialize database
	fmt.Println("🔧 Initializing database...")
	db, err := database.New("test-simple-fixed.db")
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
	fmt.Println("\n🎭 === EXPECTED SPOOFED VALUES ===")
	fp := testProfile.Fingerprint
	fmt.Printf("hardwareConcurrency: %d (should NOT be your real CPU cores)\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("deviceMemory: %d (should NOT be your real RAM GB)\n", fp.Navigator.DeviceMemory)
	fmt.Printf("screen: %dx%d\n", fp.Screen.Width, fp.Screen.Height)
	fmt.Printf("userAgent: %s\n", fp.Navigator.UserAgent)
	fmt.Printf("platform: %s\n", fp.Navigator.Platform)
	fmt.Printf("webGL vendor: %s\n", fp.WebGL.Vendor)
	fmt.Printf("timezone: %s (offset: %d)\n", fp.Timezone.Timezone, fp.Timezone.TimezoneOffset)
	fmt.Println("=====================================")

	// Launch browser
	fmt.Println("\n🚀 Launching Edge browser with FIXED spoofing...")
	fmt.Println("⚠️  Note: This uses the updated injection method")
	fmt.Println("Script should now run BEFORE page load")
	fmt.Println()

	err = browserManager.Launch(testProfile.ID)
	if err != nil {
		log.Printf("❌ Failed to launch browser: %v", err)
		fmt.Println("\n🔄 POSSIBLE ISSUES:")
		fmt.Println("1. Rod library blocked by antivirus")
		fmt.Println("2. Edge not found or CDP not available")
		fmt.Println("3. Port conflicts")
		fmt.Println()
		fmt.Println("You can still test manually by:")
		fmt.Println("1. Opening Edge browser")
		fmt.Println("2. Going to browserleaks.com/javascript")
		fmt.Println("3. Checking if values match the fingerprint above")
		
		// Cleanup and exit
		profileManager.Delete(testProfile.ID)
		return
	}

	fmt.Println("✅ Browser launched successfully!")
	fmt.Println("✅ Test pages should have opened automatically:")
	fmt.Println("   - browserleaks.com/javascript")
	fmt.Println("   - creepjs.com")
	fmt.Println()
	fmt.Println("📋 CRITICAL TEST INSTRUCTIONS:")
	fmt.Println("1. Check browserleaks.com tab:")
	fmt.Printf("   - hardwareConcurrency should be: %d (NOT your real CPU cores)\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("   - deviceMemory should be: %d (NOT your real RAM)\n", fp.Navigator.DeviceMemory)
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
	fmt.Println("⚠️  IMPORTANT: Fixed version uses PRE-LOAD injection!")
	fmt.Println("   If values are still not spoofed, the fix needs more work.")
	fmt.Println()
	fmt.Println("⏳ Browser will stay open for 60 seconds for testing...")
	
	// Wait longer for testing
	for i := 60; i > 0; i-- {
		fmt.Printf("\rTime remaining: %d seconds ", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()

	// Cleanup
	fmt.Println("\n🧹 Cleaning up...")
	err = profileManager.Delete(testProfile.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete test profile: %v", err)
	}

	fmt.Println("✅ Test completed!")
	fmt.Println()
	fmt.Println("📊 RESULTS TO REPORT:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Did test pages open automatically? (Yes/No)")
	fmt.Println("3. Were fingerprint values spoofed correctly? (Yes/No)")
	fmt.Printf("   - hardwareConcurrency = %d? (should NOT be your real CPU)\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("   - deviceMemory = %d? (should NOT be your real RAM)\n", fp.Navigator.DeviceMemory)
	fmt.Println("4. Any errors in terminal? (Yes/No)")
	fmt.Println("5. Console messages visible in DevTools? (Yes/No)")
	fmt.Println("6. Was Rod library blocked by antivirus? (Yes/No)")
	fmt.Println()
	fmt.Println("🎯 KEY QUESTION: Are the fingerprint values now spoofed correctly?")
}