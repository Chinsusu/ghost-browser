package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/user/ghost-browser/internal/database"
	"github.com/user/ghost-browser/internal/profile"
)

func main() {
	fmt.Println("🧪 === PURE CDP BROWSER TEST (No Rod Library) ===")
	fmt.Println()

	// Initialize database
	fmt.Println("🔧 Initializing database...")
	db, err := database.New("test-pure-cdp.db")
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

	// Find Edge
	fmt.Println("\n🔍 Looking for Edge browser...")
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("Edge browser not found:", err)
	}
	fmt.Printf("✅ Found Edge at: %s\n", edgePath)

	// Create user data directory
	userDataDir := filepath.Join(testProfile.DataDir, "EdgeData")
	os.MkdirAll(userDataDir, 0755)

	// Create startup script with fingerprint spoofing
	fmt.Println("\n📝 Creating fingerprint spoofing startup script...")
	scriptPath, err := createStartupScript(testProfile, userDataDir)
	if err != nil {
		log.Fatal("Failed to create startup script:", err)
	}
	fmt.Printf("✅ Script created at: %s\n", scriptPath)

	// Launch Edge with startup script injection
	fmt.Println("\n🚀 Launching Edge with startup script injection...")
	fmt.Println("⚠️  This bypasses Rod library completely")
	fmt.Println("Script will be injected via --user-script parameter")
	fmt.Println()

	args := []string{
		"--user-data-dir=" + userDataDir,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"--disable-blink-features=AutomationControlled",
		"--disable-web-security",
		"--allow-running-insecure-content",
		"--user-script=" + scriptPath,
		"https://browserleaks.com/javascript",
	}

	cmd := exec.Command(edgePath, args...)
	err = cmd.Start()
	if err != nil {
		log.Fatal("Failed to launch Edge:", err)
	}

	fmt.Println("✅ Edge launched successfully!")
	fmt.Println("✅ Startup script injected via --user-script parameter")
	fmt.Println()

	// Wait a bit then open CreepJS in new tab
	fmt.Println("⏳ Waiting 5 seconds then opening CreepJS...")
	time.Sleep(5 * time.Second)

	// Open CreepJS in new tab
	creepArgs := []string{
		"--user-data-dir=" + userDataDir,
		"--new-tab",
		"https://creepjs.com",
	}
	exec.Command(edgePath, creepArgs...).Start()

	fmt.Println("✅ CreepJS opened in new tab")
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
	fmt.Println("⚠️  IMPORTANT: This uses startup script injection!")
	fmt.Println("   Script runs BEFORE any page loads.")
	fmt.Println("   No Rod library = No antivirus blocking!")
	fmt.Println()
	fmt.Println("⏳ Browser will stay open for testing...")
	fmt.Println("Press Enter to continue and cleanup...")
	fmt.Scanln()

	// Cleanup
	fmt.Println("\n🧹 Cleaning up...")
	err = profileManager.Delete(testProfile.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete test profile: %v", err)
	}

	// Remove script file
	os.Remove(scriptPath)

	fmt.Println("✅ Test completed!")
	fmt.Println()
	fmt.Println("📊 RESULTS TO REPORT:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Did browserleaks.com load? (Yes/No)")
	fmt.Println("3. Were fingerprint values spoofed correctly? (Yes/No)")
	fmt.Printf("   - hardwareConcurrency = %d? (should NOT be your real CPU)\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("   - deviceMemory = %d? (should NOT be your real RAM)\n", fp.Navigator.DeviceMemory)
	fmt.Println("4. Any errors in terminal? (Yes/No)")
	fmt.Println("5. Console messages visible in DevTools? (Yes/No)")
	fmt.Println("6. Was Rod library bypassed successfully? (Yes/No)")
	fmt.Println()
	fmt.Println("🎯 KEY QUESTION: Are the fingerprint values now spoofed correctly?")
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

func createStartupScript(testProfile *profile.Profile, userDataDir string) (string, error) {
	fp := testProfile.Fingerprint
	
	script := fmt.Sprintf(`
// Ghost Browser Fingerprint Spoofing Script - Startup Injection
(function() {
	'use strict';
	
	console.log('[Ghost Browser] STARTUP fingerprint spoofing initializing...');
	
	// Navigator spoofing - CRITICAL: Must happen before page reads these values
	Object.defineProperty(Navigator.prototype, 'userAgent', {
		get: function() { return '%s'; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'platform', {
		get: function() { return '%s'; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'deviceMemory', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'webdriver', {
		get: function() { return false; },
		configurable: true
	});
	
	// Screen spoofing
	Object.defineProperty(Screen.prototype, 'width', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'height', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availWidth', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availHeight', {
		get: function() { return %d; },
		configurable: true
	});
	
	// WebGL spoofing
	const origGetParam = WebGLRenderingContext.prototype.getParameter;
	WebGLRenderingContext.prototype.getParameter = function(param) {
		if (param === 37445) return '%s'; // UNMASKED_VENDOR_WEBGL
		if (param === 37446) return '%s'; // UNMASKED_RENDERER_WEBGL
		return origGetParam.call(this, param);
	};
	
	// WebGL2 spoofing
	if (typeof WebGL2RenderingContext !== 'undefined') {
		const origGetParam2 = WebGL2RenderingContext.prototype.getParameter;
		WebGL2RenderingContext.prototype.getParameter = function(param) {
			if (param === 37445) return '%s';
			if (param === 37446) return '%s';
			return origGetParam2.call(this, param);
		};
	}
	
	console.log('[Ghost Browser] ✅ Navigator spoofed (hardwareConcurrency: %d, deviceMemory: %d)');
	console.log('[Ghost Browser] ✅ Screen spoofed (%dx%d)');
	console.log('[Ghost Browser] ✅ WebGL spoofed (%s)');
	console.log('[Ghost Browser] Fingerprint spoofing active - STARTUP injection');
	
})();
`, 
		fp.Navigator.UserAgent,
		fp.Navigator.Platform,
		fp.Navigator.HardwareConcurrency,
		fp.Navigator.DeviceMemory,
		fp.Screen.Width,
		fp.Screen.Height,
		fp.Screen.AvailWidth,
		fp.Screen.AvailHeight,
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
		fp.Navigator.HardwareConcurrency,
		fp.Navigator.DeviceMemory,
		fp.Screen.Width,
		fp.Screen.Height,
		fp.WebGL.Vendor,
	)

	scriptPath := filepath.Join(userDataDir, "ghost-spoof.js")
	err := os.WriteFile(scriptPath, []byte(script), 0644)
	return scriptPath, err
}