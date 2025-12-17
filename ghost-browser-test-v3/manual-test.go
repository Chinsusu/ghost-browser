package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	fmt.Println("🧪 === Manual Script Injection Test ===")
	fmt.Println()

	// Create spoofing script file
	fmt.Println("📝 Creating spoofing script...")
	scriptContent := `
// Ghost Browser Manual Injection Script
(function() {
	'use strict';
	
	console.log('[Ghost Browser] MANUAL fingerprint spoofing initializing...');
	
	// Navigator spoofing - MUST happen before page reads values
	Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
		get: function() { return 8; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'deviceMemory', {
		get: function() { return 16; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'webdriver', {
		get: function() { return false; },
		configurable: true
	});
	
	console.log('[Ghost Browser] ✅ Navigator spoofed (hardwareConcurrency: 8, deviceMemory: 16)');
	console.log('[Ghost Browser] Fingerprint spoofing active - MANUAL injection');
	
	// Test values
	console.log('Current hardwareConcurrency:', navigator.hardwareConcurrency);
	console.log('Current deviceMemory:', navigator.deviceMemory);
	console.log('Current webdriver:', navigator.webdriver);
	
})();
`

	scriptPath := "ghost-spoof-manual.js"
	err := ioutil.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		fmt.Printf("❌ Failed to create script: %v\n", err)
		return
	}
	fmt.Printf("✅ Script created: %s\n", scriptPath)

	// Find Edge
	fmt.Println("\n🔍 Looking for Edge...")
	edgePath, err := findEdgePath()
	if err != nil {
		fmt.Printf("❌ Edge not found: %v\n", err)
		return
	}
	fmt.Printf("✅ Found Edge: %s\n", edgePath)

	// Get absolute script path
	absScriptPath, _ := filepath.Abs(scriptPath)

	// Launch Edge with user script
	fmt.Println("\n🚀 Launching Edge with manual script injection...")
	fmt.Println("Method: --user-script parameter")
	fmt.Println()

	args := []string{
		"--user-script=" + absScriptPath,
		"--disable-blink-features=AutomationControlled",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"https://browserleaks.com/javascript",
	}

	cmd := exec.Command(edgePath, args...)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("❌ Failed to launch Edge: %v\n", err)
		return
	}

	fmt.Println("✅ Edge launched with manual script injection!")
	fmt.Println()
	fmt.Println("📋 MANUAL VERIFICATION STEPS:")
	fmt.Println("1. Wait for browserleaks.com to load")
	fmt.Println("2. Open DevTools (F12)")
	fmt.Println("3. Go to Console tab")
	fmt.Println("4. Look for these messages:")
	fmt.Println("   [Ghost Browser] MANUAL fingerprint spoofing initializing...")
	fmt.Println("   [Ghost Browser] ✅ Navigator spoofed")
	fmt.Println("   Current hardwareConcurrency: 8")
	fmt.Println("   Current deviceMemory: 16")
	fmt.Println()
	fmt.Println("5. Check Navigator values on the page:")
	fmt.Println("   - hardwareConcurrency should be 8 (NOT 16)")
	fmt.Println("   - deviceMemory should be 16 (NOT 8)")
	fmt.Println()
	fmt.Println("6. If values are still wrong, try:")
	fmt.Println("   - Refresh the page (F5)")
	fmt.Println("   - Check if script loaded in Sources tab")
	fmt.Println("   - Try different Edge flags")

	fmt.Println("\nPress Enter to cleanup...")
	fmt.Scanln()

	// Cleanup
	os.Remove(scriptPath)
	fmt.Println("✅ Cleanup completed!")
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