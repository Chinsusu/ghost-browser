package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("=== Ghost Browser - ChromeDP Launcher ===")
	fmt.Println()

	// Find Edge
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("❌ Edge not found:", err)
	}
	fmt.Println("✅ Found Edge:", edgePath)

	// Create temp profile
	tempDir, err := os.MkdirTemp("", "ghost-profile-*")
	if err != nil {
		log.Fatal("❌ Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)
	fmt.Println("✅ Profile dir:", tempDir)

	// Setup Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(edgePath),
		chromedp.UserDataDir(tempDir),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-web-security", false),
		chromedp.Flag("disable-site-isolation-trials", true),
		// Window size
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Fingerprint spoofing script
	spoofScript := generateSpoofScript()

	fmt.Println("✅ Launching browser with fingerprint spoofing...")
	fmt.Println()

	// Add script to run before page load
	err = chromedp.Run(ctx,
		// Add script that runs on every new document
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(spoofScript).Do(ctx)
			return err
		}),
		
		// Navigate to test page
		chromedp.Navigate("https://browserleaks.com/javascript"),
		chromedp.WaitReady("body"),
	)
	if err != nil {
		log.Fatal("❌ Failed:", err)
	}

	// Wait for page to fully load
	fmt.Println("⏳ Waiting for page to load...")
	
	// Verify spoofing
	var hwConcurrency, devMemory int
	var platform string
	var webdriver bool

	err = chromedp.Run(ctx,
		chromedp.Sleep(3000), // Wait 3 seconds
		chromedp.Evaluate(`navigator.hardwareConcurrency || 0`, &hwConcurrency),
		chromedp.Evaluate(`navigator.deviceMemory || 0`, &devMemory),
		chromedp.Evaluate(`navigator.platform || 'unknown'`, &platform),
		chromedp.Evaluate(`navigator.webdriver || false`, &webdriver),
	)
	if err != nil {
		log.Println("⚠️ Could not verify values:", err)
	} else {
		fmt.Println("=== VERIFICATION ===")
		fmt.Printf("hardwareConcurrency: %d (expected: 8)\n", hwConcurrency)
		fmt.Printf("deviceMemory: %d (expected: 16)\n", devMemory)
		fmt.Printf("platform: %s\n", platform)
		fmt.Printf("webdriver: %v (expected: false)\n", webdriver)
		fmt.Println("====================")
		
		if hwConcurrency == 8 && devMemory == 16 {
			fmt.Println("✅ SPOOFING SUCCESS!")
		} else {
			fmt.Println("❌ SPOOFING FAILED - values not changed")
		}
	}

	fmt.Println()
	fmt.Println("Browser is running. Press Ctrl+C to close...")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nClosing browser...")
}

func findEdgePath() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Windows only")
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

	return "", fmt.Errorf("not found")
}

func generateSpoofScript() string {
	return `
(function() {
	'use strict';

	const SPOOF = {
		hardwareConcurrency: 8,
		deviceMemory: 16,
		platform: 'Win32',
		vendor: 'Google Inc.',
		language: 'en-US',
		languages: Object.freeze(['en-US', 'en']),
		maxTouchPoints: 0,
		screen: {
			width: 1920,
			height: 1080,
			availWidth: 1920,
			availHeight: 1040,
			colorDepth: 24,
			pixelDepth: 24,
		},
		webgl: {
			vendor: 'Google Inc. (NVIDIA)',
			renderer: 'ANGLE (NVIDIA, NVIDIA GeForce RTX 3080 Direct3D11 vs_5_0 ps_5_0)',
		}
	};

	// ===== NAVIGATOR SPOOFING =====
	const navOverrides = {
		hardwareConcurrency: { get: () => SPOOF.hardwareConcurrency },
		deviceMemory: { get: () => SPOOF.deviceMemory },
		platform: { get: () => SPOOF.platform },
		vendor: { get: () => SPOOF.vendor },
		language: { get: () => SPOOF.language },
		languages: { get: () => SPOOF.languages },
		maxTouchPoints: { get: () => SPOOF.maxTouchPoints },
		webdriver: { get: () => false },
	};

	for (const [prop, descriptor] of Object.entries(navOverrides)) {
		try {
			Object.defineProperty(Navigator.prototype, prop, {
				...descriptor,
				configurable: true,
				enumerable: true,
			});
		} catch (e) {}
	}

	// ===== SCREEN SPOOFING =====
	for (const [prop, value] of Object.entries(SPOOF.screen)) {
		try {
			Object.defineProperty(Screen.prototype, prop, {
				get: () => value,
				configurable: true,
			});
		} catch (e) {}
	}

	try {
		Object.defineProperty(window, 'devicePixelRatio', {
			get: () => 1.0,
			configurable: true,
		});
	} catch (e) {}

	// ===== WEBGL SPOOFING =====
	const webglHandler = {
		apply(target, thisArg, args) {
			const param = args[0];
			if (param === 37445) return SPOOF.webgl.vendor;
			if (param === 37446) return SPOOF.webgl.renderer;
			return Reflect.apply(target, thisArg, args);
		}
	};

	try {
		WebGLRenderingContext.prototype.getParameter = new Proxy(
			WebGLRenderingContext.prototype.getParameter,
			webglHandler
		);
		if (typeof WebGL2RenderingContext !== 'undefined') {
			WebGL2RenderingContext.prototype.getParameter = new Proxy(
				WebGL2RenderingContext.prototype.getParameter,
				webglHandler
			);
		}
	} catch (e) {}

	// ===== CANVAS NOISE =====
	const origToDataURL = HTMLCanvasElement.prototype.toDataURL;
	HTMLCanvasElement.prototype.toDataURL = function() {
		try {
			const ctx = this.getContext('2d');
			if (ctx && this.width && this.height) {
				const imgData = ctx.getImageData(0, 0, Math.min(this.width, 50), Math.min(this.height, 50));
				for (let i = 0; i < imgData.data.length; i += 100) {
					imgData.data[i] ^= 1;
				}
				ctx.putImageData(imgData, 0, 0);
			}
		} catch (e) {}
		return origToDataURL.apply(this, arguments);
	};

	// ===== TIMEZONE =====
	Date.prototype.getTimezoneOffset = () => 300;

	// ===== WEBRTC DISABLE =====
	window.RTCPeerConnection = undefined;
	window.webkitRTCPeerConnection = undefined;

	// ===== REMOVE AUTOMATION FLAGS =====
	try { delete Object.getPrototypeOf(navigator).webdriver; } catch (e) {}
	
	['__webdriver_evaluate', '__selenium_evaluate', '__driver_evaluate', 
	 '_selenium', '$cdc_asdjflasutopfhvcZLmcfl_', '$chrome_asyncScriptInfo'
	].forEach(p => { try { delete window[p]; } catch(e) {} });

	console.log('[Ghost] ✅ Spoofing active | CPU:', navigator.hardwareConcurrency, '| RAM:', navigator.deviceMemory);
})();
`
}
