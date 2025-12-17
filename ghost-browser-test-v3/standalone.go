package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Fingerprint struct {
	HardwareConcurrency int
	DeviceMemory        int
	Platform            string
	UserAgent           string
	ScreenWidth         int
	ScreenHeight        int
	WebGLVendor         string
	WebGLRenderer       string
}

func main() {
	fmt.Println("🧪 === Ghost Browser Test v3 (ChromeDP Standalone) ===")
	fmt.Println()

	// Create test fingerprint
	testFingerprint := &Fingerprint{
		HardwareConcurrency: 8,  // Should NOT be 16 (real CPU)
		DeviceMemory:        16, // Should NOT be 8 (real RAM)
		Platform:            "Win32",
		UserAgent:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.86 Safari/537.36 Edg/131.0.2903.86",
		ScreenWidth:         1920,
		ScreenHeight:        1080,
		WebGLVendor:         "Google Inc. (AMD)",
		WebGLRenderer:       "ANGLE (AMD, AMD Radeon RX 6600 XT (0x000073FF) Direct3D11 vs_5_0 ps_5_0, D3D11)",
	}

	// Display expected spoofed values
	fmt.Println("🎭 === EXPECTED SPOOFED VALUES ===")
	fmt.Printf("hardwareConcurrency: %d (should NOT be 16)\n", testFingerprint.HardwareConcurrency)
	fmt.Printf("deviceMemory: %d (should NOT be 8)\n", testFingerprint.DeviceMemory)
	fmt.Printf("platform: %s\n", testFingerprint.Platform)
	fmt.Printf("userAgent: %s\n", testFingerprint.UserAgent)
	fmt.Printf("screen: %dx%d\n", testFingerprint.ScreenWidth, testFingerprint.ScreenHeight)
	fmt.Printf("webGL vendor: %s\n", testFingerprint.WebGLVendor)
	fmt.Println("=====================================")

	// Launch browser with ChromeDP
	fmt.Println("\n🚀 Launching Edge with ChromeDP...")
	fmt.Println("⚠️  This uses ChromeDP instead of Rod")
	fmt.Println("Should avoid antivirus blocking issues")
	fmt.Println()

	err := launchBrowserWithSpoofing(testFingerprint)
	if err != nil {
		log.Printf("❌ Failed to launch browser: %v", err)
		fmt.Println("\n🔄 POSSIBLE ISSUES:")
		fmt.Println("1. ChromeDP blocked by antivirus")
		fmt.Println("2. Edge not found or CDP not available")
		fmt.Println("3. Port conflicts")
		fmt.Println()
		fmt.Println("SOLUTIONS:")
		fmt.Println("1. Add Go folder to antivirus exclusions")
		fmt.Println("2. Temporarily disable real-time protection")
		fmt.Println("3. Run as administrator")
	} else {
		fmt.Println("✅ Browser launched successfully!")
		fmt.Println()
		fmt.Println("📋 VERIFICATION INSTRUCTIONS:")
		fmt.Println("1. Check the opened browserleaks.com tab")
		fmt.Println("2. Look for the verification results in terminal")
		fmt.Println("3. Open DevTools (F12) and check Console for:")
		fmt.Println("   [Ghost Browser] ✅ Navigator spoofed")
		fmt.Println("   [Ghost Browser] ✅ Screen spoofed")
		fmt.Println("   [Ghost Browser] Fingerprint spoofing active")
	}

	fmt.Println("\n✅ Test completed!")
	fmt.Println()
	fmt.Println("📊 RESULTS TO REPORT:")
	fmt.Println("1. Did Edge browser open? (Yes/No)")
	fmt.Println("2. Were fingerprint values spoofed correctly? (Yes/No)")
	fmt.Printf("   - hardwareConcurrency = %d? (should NOT be 16)\n", testFingerprint.HardwareConcurrency)
	fmt.Printf("   - deviceMemory = %d? (should NOT be 8)\n", testFingerprint.DeviceMemory)
	fmt.Println("3. Any errors in terminal? (Yes/No)")
	fmt.Println("4. Console messages visible in DevTools? (Yes/No)")
	fmt.Println("5. Was ChromeDP blocked by antivirus? (Yes/No)")
}

func launchBrowserWithSpoofing(fp *Fingerprint) error {
	// Create spoofing script
	spoofScript := fmt.Sprintf(`
(function() {
	'use strict';
	
	console.log('[Ghost Browser] ChromeDP fingerprint spoofing initializing...');
	
	// Navigator spoofing - CRITICAL: Must happen before page reads these values
	Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'deviceMemory', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'platform', {
		get: function() { return '%s'; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'userAgent', {
		get: function() { return '%s'; },
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
	
	// WebGL spoofing
	const origGetParam = WebGLRenderingContext.prototype.getParameter;
	WebGLRenderingContext.prototype.getParameter = function(param) {
		if (param === 37445) return '%s'; // UNMASKED_VENDOR_WEBGL
		if (param === 37446) return '%s'; // UNMASKED_RENDERER_WEBGL
		return origGetParam.call(this, param);
	};
	
	console.log('[Ghost Browser] ✅ Navigator spoofed (hardwareConcurrency: %d, deviceMemory: %d)');
	console.log('[Ghost Browser] ✅ Screen spoofed (%dx%d)');
	console.log('[Ghost Browser] ✅ WebGL spoofed (%s)');
	console.log('[Ghost Browser] Fingerprint spoofing active - ChromeDP injection');
	
})();
`, 
		fp.HardwareConcurrency,
		fp.DeviceMemory,
		fp.Platform,
		fp.UserAgent,
		fp.ScreenWidth,
		fp.ScreenHeight,
		fp.WebGLVendor,
		fp.WebGLRenderer,
		fp.HardwareConcurrency,
		fp.DeviceMemory,
		fp.ScreenWidth,
		fp.ScreenHeight,
		fp.WebGLVendor,
	)

	// ChromeDP options to use existing Edge installation
	opts := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("allow-running-insecure-content", true),
		chromedp.Flag("remote-debugging-port", "9222"),
	}

	// Create context
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Navigate and inject spoofing script
	var hardwareConcurrency, deviceMemory int64
	var platform string
	var webdriver bool

	err := chromedp.Run(ctx,
		// Navigate to about:blank first
		chromedp.Navigate("about:blank"),
		
		// Inject spoofing script BEFORE navigating to test page
		chromedp.Evaluate(spoofScript, nil),
		
		// Navigate to browserleaks.com
		chromedp.Navigate("https://browserleaks.com/javascript"),
		
		// Wait for page to load
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		
		// Extract values to verify spoofing
		chromedp.Evaluate(`navigator.hardwareConcurrency`, &hardwareConcurrency),
		chromedp.Evaluate(`navigator.deviceMemory`, &deviceMemory),
		chromedp.Evaluate(`navigator.platform`, &platform),
		chromedp.Evaluate(`navigator.webdriver`, &webdriver),
	)

	if err != nil {
		return fmt.Errorf("ChromeDP execution failed: %w", err)
	}

	// Verify spoofing results
	fmt.Println("\n=== VERIFICATION ===")
	fmt.Printf("hardwareConcurrency: %d (expected: %d)\n", hardwareConcurrency, fp.HardwareConcurrency)
	fmt.Printf("deviceMemory: %d (expected: %d)\n", deviceMemory, fp.DeviceMemory)
	fmt.Printf("platform: %s\n", platform)
	fmt.Printf("webdriver: %t (expected: false)\n", webdriver)
	fmt.Println("====================")

	// Check if spoofing worked
	if hardwareConcurrency == int64(fp.HardwareConcurrency) && 
	   deviceMemory == int64(fp.DeviceMemory) {
		fmt.Println("✅ SPOOFING SUCCESS!")
	} else {
		fmt.Println("❌ SPOOFING FAILED!")
		fmt.Printf("Expected hardwareConcurrency: %d, got: %d\n", fp.HardwareConcurrency, hardwareConcurrency)
		fmt.Printf("Expected deviceMemory: %d, got: %d\n", fp.DeviceMemory, deviceMemory)
	}

	// Keep browser open for manual verification
	fmt.Println("\n⏳ Browser will stay open for 30 seconds for manual verification...")
	time.Sleep(30 * time.Second)

	return nil
}