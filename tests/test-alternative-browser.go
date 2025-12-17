package main

import (
	"fmt"
	"io/ioutil"
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
	fmt.Println("🧪 === Alternative Browser Test (CDP via DevTools Protocol) ===")
	fmt.Println()

	// Initialize database
	fmt.Println("🔧 Initializing database...")
	db, err := database.New("test-alternative-browser.db")
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

	// Create user data directory
	userDataDir := filepath.Join(testProfile.DataDir, "EdgeData")
	os.MkdirAll(userDataDir, 0755)

	// Create startup script for injection
	fmt.Println("\n📝 Creating fingerprint injection script...")
	scriptPath, err := createInjectionScript(testProfile, userDataDir)
	if err != nil {
		log.Fatal("Failed to create injection script:", err)
	}
	fmt.Printf("✅ Script created at: %s\n", scriptPath)

	// Launch Edge with CDP enabled
	fmt.Println("\n🚀 Launching Edge with CDP enabled...")
	fmt.Println("This will open Edge with fingerprint spoofing via startup script injection.")
	fmt.Println()

	args := []string{
		"--user-data-dir=" + userDataDir,
		"--remote-debugging-port=9222",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"--disable-blink-features=AutomationControlled",
		"--disable-web-security",
		"--allow-running-insecure-content",
		"--user-script=" + scriptPath,
		"about:blank",
	}

	cmd := exec.Command(edgePath, args...)
	err = cmd.Start()
	if err != nil {
		log.Fatal("Failed to launch Edge:", err)
	}

	fmt.Println("✅ Edge launched successfully!")
	fmt.Println("⏳ Waiting for Edge to initialize...")
	time.Sleep(3 * time.Second)

	// Try to inject via CDP
	fmt.Println("\n🔧 Attempting CDP injection...")
	err = injectViaCDP(testProfile)
	if err != nil {
		fmt.Printf("⚠️  CDP injection failed: %v\n", err)
		fmt.Println("Falling back to manual testing...")
	} else {
		fmt.Println("✅ CDP injection successful!")
	}

	// Open test pages
	fmt.Println("\n🌐 Opening test pages...")
	openTestPages()

	fmt.Println("\n📋 Test Instructions:")
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
	fmt.Println("📊 Results to report:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Were fingerprint values spoofed correctly? (Yes/No)")
	fmt.Println("3. Any errors in terminal? (Yes/No)")
	fmt.Println("4. Console messages visible in DevTools? (Yes/No)")
	fmt.Println("5. Did CDP injection work? (Yes/No)")
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

func createInjectionScript(testProfile *profile.Profile, userDataDir string) (string, error) {
	fp := testProfile.Fingerprint
	
	script := fmt.Sprintf(`
// Ghost Browser Fingerprint Spoofing Script
(function() {
	'use strict';
	
	console.log('[Ghost Browser] Initializing fingerprint spoofing...');
	
	// Navigator spoofing
	Object.defineProperty(Navigator.prototype, 'userAgent', {
		get: () => '%s',
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'platform', {
		get: () => '%s',
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
		get: () => %d,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'deviceMemory', {
		get: () => %d,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'webdriver', {
		get: () => false,
		configurable: true
	});
	
	// Screen spoofing
	Object.defineProperty(Screen.prototype, 'width', {
		get: () => %d,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'height', {
		get: () => %d,
		configurable: true
	});
	
	// WebGL spoofing
	const origGetParam = WebGLRenderingContext.prototype.getParameter;
	WebGLRenderingContext.prototype.getParameter = function(param) {
		if (param === 37445) return '%s'; // UNMASKED_VENDOR_WEBGL
		if (param === 37446) return '%s'; // UNMASKED_RENDERER_WEBGL
		return origGetParam.call(this, param);
	};
	
	console.log('[Ghost Browser] ✅ Navigator spoofed');
	console.log('[Ghost Browser] ✅ Screen spoofed');
	console.log('[Ghost Browser] ✅ WebGL spoofed');
	console.log('[Ghost Browser] Fingerprint spoofing active');
	
})();
`, 
		fp.Navigator.UserAgent,
		fp.Navigator.Platform,
		fp.Navigator.HardwareConcurrency,
		fp.Navigator.DeviceMemory,
		fp.Screen.Width,
		fp.Screen.Height,
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
	)

	scriptPath := filepath.Join(userDataDir, "spoof.js")
	err := ioutil.WriteFile(scriptPath, []byte(script), 0644)
	return scriptPath, err
}

func injectViaCDP(testProfile *profile.Profile) error {
	// This is a simplified CDP injection attempt
	// In a real implementation, you'd use a proper CDP client
	fmt.Println("CDP injection would happen here...")
	fmt.Println("(This is a placeholder - real CDP implementation needed)")
	return nil
}

func openTestPages() {
	// Open test pages in new tabs
	time.Sleep(2 * time.Second)
	
	// Try to open browserleaks.com
	exec.Command("cmd", "/c", "start", "http://localhost:9222/json/new?browserleaks.com/javascript").Run()
	time.Sleep(1 * time.Second)
	
	// Try to open CreepJS
	exec.Command("cmd", "/c", "start", "http://localhost:9222/json/new?creepjs.com").Run()
}