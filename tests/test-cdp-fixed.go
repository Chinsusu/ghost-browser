package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/user/ghost-browser/internal/database"
	"github.com/user/ghost-browser/internal/profile"
)

type CDPResponse struct {
	ID     int         `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

type TabInfo struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	URL            string `json:"url"`
	WebSocketURL   string `json:"webSocketDebuggerUrl"`
}

func main() {
	fmt.Println("🧪 === FIXED CDP Browser Test (Pre-load Injection) ===")
	fmt.Println()

	// Initialize database
	fmt.Println("🔧 Initializing database...")
	db, err := database.New("test-cdp-fixed.db")
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
	fmt.Println("\n🎭 === SPOOFED VALUES ===")
	fp := testProfile.Fingerprint
	fmt.Printf("hardwareConcurrency: %d (should be %d)\n", fp.Navigator.HardwareConcurrency, fp.Navigator.HardwareConcurrency)
	fmt.Printf("deviceMemory: %d (should be %d)\n", fp.Navigator.DeviceMemory, fp.Navigator.DeviceMemory)
	fmt.Printf("screen: %dx%d (should be %dx%d)\n", fp.Screen.Width, fp.Screen.Height, fp.Screen.Width, fp.Screen.Height)
	fmt.Printf("userAgent: %s\n", fp.Navigator.UserAgent)
	fmt.Printf("platform: %s\n", fp.Navigator.Platform)
	fmt.Printf("webGL vendor: %s\n", fp.WebGL.Vendor)
	fmt.Printf("timezone: %s (offset: %d)\n", fp.Timezone.Timezone, fp.Timezone.TimezoneOffset)
	fmt.Println("======================")

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

	// Launch Edge with CDP
	fmt.Println("\n🚀 Launching Edge with CDP...")
	cdpPort := "9222"
	args := []string{
		"--user-data-dir=" + userDataDir,
		"--remote-debugging-port=" + cdpPort,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"--disable-blink-features=AutomationControlled",
		"--disable-web-security",
		"--allow-running-insecure-content",
		"about:blank",
	}

	cmd := exec.Command(edgePath, args...)
	err = cmd.Start()
	if err != nil {
		log.Fatal("Failed to launch Edge:", err)
	}

	fmt.Println("✅ Edge launched successfully!")
	fmt.Println("⏳ Waiting for CDP to be available...")
	time.Sleep(5 * time.Second)

	// Test CDP connection
	fmt.Println("\n🔗 Testing CDP connection...")
	cdpURL := fmt.Sprintf("http://localhost:%s", cdpPort)
	
	resp, err := http.Get(cdpURL + "/json/version")
	if err != nil {
		fmt.Printf("❌ CDP not available: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("✅ CDP connected: %s\n", string(body))

	// Get tabs
	fmt.Println("\n📑 Getting browser tabs...")
	tabs, err := getTabs(cdpURL)
	if err != nil {
		fmt.Printf("❌ Failed to get tabs: %v\n", err)
		return
	}

	if len(tabs) == 0 {
		fmt.Println("❌ No tabs found")
		return
	}

	fmt.Printf("✅ Found %d tab(s)\n", len(tabs))

	// CRITICAL FIX: Add script to evaluate on NEW DOCUMENT (before page load)
	fmt.Println("\n💉 Adding spoofing script to evaluate on new document...")
	err = addScriptToEvaluateOnNewDocument(cdpURL, tabs[0].ID, testProfile)
	if err != nil {
		fmt.Printf("❌ Failed to add script: %v\n", err)
		return
	} else {
		fmt.Println("✅ Spoofing script added to new document evaluation!")
	}

	// Navigate to test pages
	fmt.Println("\n🌐 Opening test pages...")
	
	// Navigate first tab to browserleaks
	err = navigateTab(cdpURL, tabs[0].ID, "https://browserleaks.com/javascript")
	if err != nil {
		fmt.Printf("❌ Failed to navigate to browserleaks: %v\n", err)
	} else {
		fmt.Println("✅ Opened browserleaks.com")
	}

	// Wait for page to load
	time.Sleep(3 * time.Second)

	// Create new tab for CreepJS with pre-injected script
	fmt.Println("\n🆕 Creating new tab for CreepJS...")
	newTabID, err := createNewTab(cdpURL, "about:blank")
	if err != nil {
		fmt.Printf("❌ Failed to create new tab: %v\n", err)
	} else {
		fmt.Println("✅ Created new tab")
		
		// Add script to new tab BEFORE navigating
		err = addScriptToEvaluateOnNewDocument(cdpURL, newTabID, testProfile)
		if err != nil {
			fmt.Printf("❌ Failed to add script to new tab: %v\n", err)
		} else {
			fmt.Println("✅ Added spoofing script to new tab")
			
			// Now navigate to CreepJS
			time.Sleep(1 * time.Second)
			err = navigateTab(cdpURL, newTabID, "https://creepjs.com")
			if err != nil {
				fmt.Printf("❌ Failed to navigate to CreepJS: %v\n", err)
			} else {
				fmt.Println("✅ Opened creepjs.com")
			}
		}
	}

	fmt.Println("\n📋 Test Instructions:")
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
	fmt.Println("⚠️  IMPORTANT: Script now runs BEFORE page load!")
	fmt.Println("   This should fix the fingerprint spoofing issue.")
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

	fmt.Println("✅ Test completed!")
	fmt.Println()
	fmt.Println("📊 Results to report:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Was CDP connection successful? (Yes/No)")
	fmt.Println("3. Were fingerprint values spoofed correctly? (Yes/No)")
	fmt.Printf("   - hardwareConcurrency = %d? (should NOT be your real CPU cores)\n", fp.Navigator.HardwareConcurrency)
	fmt.Printf("   - deviceMemory = %d? (should NOT be your real RAM)\n", fp.Navigator.DeviceMemory)
	fmt.Println("4. Any errors in terminal? (Yes/No)")
	fmt.Println("5. Console messages visible in DevTools? (Yes/No)")
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

func getTabs(cdpURL string) ([]TabInfo, error) {
	resp, err := http.Get(cdpURL + "/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tabs []TabInfo
	err = json.Unmarshal(body, &tabs)
	return tabs, err
}

func sendCDPCommand(cdpURL, tabID string, method string, params interface{}) (*CDPResponse, error) {
	command := map[string]interface{}{
		"id":     1,
		"method": method,
		"params": params,
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	// Use WebSocket or direct tab endpoint for CDP commands
	// For now, we'll use a simpler approach with direct HTTP
	url := fmt.Sprintf("%s/json", cdpURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cdpResp CDPResponse
	err = json.Unmarshal(body, &cdpResp)
	return &cdpResp, err
}

// CRITICAL FIX: Use Page.addScriptToEvaluateOnNewDocument instead of Runtime.evaluate
func addScriptToEvaluateOnNewDocument(cdpURL, tabID string, testProfile *profile.Profile) error {
	fp := testProfile.Fingerprint
	
	script := fmt.Sprintf(`
(function() {
	'use strict';
	
	console.log('[Ghost Browser] PRE-LOAD fingerprint spoofing initializing...');
	
	// Navigator spoofing - MUST happen before page reads these values
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
	
	Object.defineProperty(Screen.prototype, 'availWidth', {
		get: () => %d,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availHeight', {
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
	console.log('[Ghost Browser] Fingerprint spoofing active - PRE-LOAD injection');
	
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

	params := map[string]interface{}{
		"source": script,
	}

	// Use Page.addScriptToEvaluateOnNewDocument instead of Runtime.evaluate
	_, err := sendCDPCommand(cdpURL, tabID, "Page.addScriptToEvaluateOnNewDocument", params)
	return err
}

func navigateTab(cdpURL, tabID, url string) error {
	params := map[string]interface{}{
		"url": url,
	}

	_, err := sendCDPCommand(cdpURL, tabID, "Page.navigate", params)
	return err
}

func createNewTab(cdpURL, url string) (string, error) {
	resp, err := http.Get(cdpURL + "/json/new?" + url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tab TabInfo
	err = json.Unmarshal(body, &tab)
	return tab.ID, err
}