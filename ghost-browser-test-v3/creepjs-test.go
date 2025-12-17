package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("🕵️ === CreepJS Fingerprint Detection Test ===")
	fmt.Println("Testing Ghost Browser spoofing against advanced detection")
	fmt.Println()

	// Find Edge
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("❌ Edge not found:", err)
	}
	fmt.Println("✅ Found Edge:", edgePath)

	// Create temp profile
	tempDir, err := os.MkdirTemp("", "ghost-creepjs-*")
	if err != nil {
		log.Fatal("❌ Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Setup Chrome options with enhanced stealth
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(edgePath),
		chromedp.UserDataDir(tempDir),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Enhanced spoofing script for CreepJS
	spoofScript := `
(function() {
	console.log('[Ghost CreepJS Test] Advanced fingerprint spoofing initializing...');
	
	// ========== Navigator Spoofing ==========
	Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
		get: () => 8,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'deviceMemory', {
		get: () => 16,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'platform', {
		get: () => 'Win32',
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'webdriver', {
		get: () => false,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'maxTouchPoints', {
		get: () => 0,
		configurable: true
	});
	
	// ========== Screen Spoofing ==========
	Object.defineProperty(Screen.prototype, 'width', {
		get: () => 1920,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'height', {
		get: () => 1080,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availWidth', {
		get: () => 1920,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availHeight', {
		get: () => 1040,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'colorDepth', {
		get: () => 24,
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'pixelDepth', {
		get: () => 24,
		configurable: true
	});
	
	Object.defineProperty(window, 'devicePixelRatio', {
		get: () => 1,
		configurable: true
	});
	
	// ========== WebGL Spoofing ==========
	const origGetParam = WebGLRenderingContext.prototype.getParameter;
	WebGLRenderingContext.prototype.getParameter = function(param) {
		if (param === 37445) return 'Intel Inc.'; // UNMASKED_VENDOR_WEBGL
		if (param === 37446) return 'Intel(R) HD Graphics 630'; // UNMASKED_RENDERER_WEBGL
		return origGetParam.call(this, param);
	};
	
	// WebGL2 spoofing
	if (typeof WebGL2RenderingContext !== 'undefined') {
		const origGetParam2 = WebGL2RenderingContext.prototype.getParameter;
		WebGL2RenderingContext.prototype.getParameter = function(param) {
			if (param === 37445) return 'Intel Inc.';
			if (param === 37446) return 'Intel(R) HD Graphics 630';
			return origGetParam2.call(this, param);
		};
	}
	
	// ========== Canvas Noise ==========
	const noise = 0.1;
	const origToDataURL = HTMLCanvasElement.prototype.toDataURL;
	HTMLCanvasElement.prototype.toDataURL = function(type, quality) {
		if (this.width > 0 && this.height > 0) {
			try {
				const ctx = this.getContext('2d');
				if (ctx) {
					const imageData = ctx.getImageData(0, 0, this.width, this.height);
					for (let i = 0; i < imageData.data.length; i += 4) {
						imageData.data[i] = Math.max(0, Math.min(255,
							imageData.data[i] + Math.floor((Math.random() - 0.5) * noise * 255)));
					}
					ctx.putImageData(imageData, 0, 0);
				}
			} catch (e) {}
		}
		return origToDataURL.call(this, type, quality);
	};
	
	// ========== Audio Spoofing ==========
	if (typeof AudioContext !== 'undefined') {
		const audioNoise = 0.001;
		const origCreateAnalyser = AudioContext.prototype.createAnalyser;
		AudioContext.prototype.createAnalyser = function() {
			const analyser = origCreateAnalyser.call(this);
			const origGetFloat = analyser.getFloatFrequencyData.bind(analyser);
			analyser.getFloatFrequencyData = function(array) {
				origGetFloat(array);
				for (let i = 0; i < array.length; i++) {
					array[i] += (Math.random() - 0.5) * audioNoise;
				}
			};
			return analyser;
		};
	}
	
	// ========== Remove Automation Flags ==========
	try { delete Object.getPrototypeOf(navigator).webdriver; } catch (e) {}
	
	const autoProps = [
		'__webdriver_evaluate', '__selenium_evaluate', '__webdriver_script_function',
		'__driver_evaluate', '_selenium', '_Selenium_IDE_Recorder', 'callSelenium',
		'$cdc_asdjflasutopfhvcZLmcfl_', '$chrome_asyncScriptInfo', '__nightmare',
		'_phantom', 'callPhantom', '__fxdriver_evaluate', '__driver_unwrapped',
		'webdriver', '__webdriver_script_fn', '__selenium_unwrapped'
	];
	autoProps.forEach(prop => {
		try { if (window[prop]) delete window[prop]; } catch (e) {}
	});
	
	// ========== Permissions API Spoofing ==========
	if (navigator.permissions && navigator.permissions.query) {
		const origQuery = navigator.permissions.query;
		navigator.permissions.query = function(params) {
			return origQuery(params).then(result => {
				if (params.name === 'notifications') {
					Object.defineProperty(result, 'state', { value: 'granted', writable: false });
				}
				return result;
			});
		};
	}
	
	console.log('[Ghost CreepJS] ✅ Advanced spoofing active');
	console.log('[Ghost CreepJS] ✅ Navigator: hardwareConcurrency=' + navigator.hardwareConcurrency + ', deviceMemory=' + navigator.deviceMemory);
	console.log('[Ghost CreepJS] ✅ Screen: ' + screen.width + 'x' + screen.height);
	console.log('[Ghost CreepJS] ✅ WebDriver: ' + navigator.webdriver);
	
})();
`

	fmt.Println("🚀 Launching browser with enhanced spoofing...")
	fmt.Println()

	// Launch with pre-load script
	err = chromedp.Run(ctx,
		// Add script BEFORE navigation
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(spoofScript).Do(ctx)
			if err != nil {
				fmt.Printf("⚠️  Script injection failed: %v\n", err)
			} else {
				fmt.Println("✅ Enhanced spoofing script injected!")
			}
			return err
		}),
		
		// Navigate to CreepJS
		chromedp.Navigate("https://abrahamjuliot.github.io/creepjs/"),
	)

	if err != nil {
		log.Printf("❌ Failed: %v", err)
		return
	}

	fmt.Println("✅ Browser launched and navigated to CreepJS!")
	fmt.Println()
	fmt.Println("🕵️ === CREEPJS FINGERPRINT ANALYSIS ===")
	fmt.Println()
	fmt.Println("📋 MANUAL VERIFICATION STEPS:")
	fmt.Println("1. ⏳ Wait for CreepJS to complete scanning (10-30 seconds)")
	fmt.Println("2. 📊 Check the Trust Score at the top (should be GREEN/HIGH)")
	fmt.Println("3. 🔍 Look for these sections:")
	fmt.Println("   • Trust Score: Should be >80% (green)")
	fmt.Println("   • Lies Detected: Should be 0 or minimal")
	fmt.Println("   • Bot Detection: Should show 'Human-like'")
	fmt.Println("   • Canvas: Should show consistent hash")
	fmt.Println("   • WebGL: Should show spoofed GPU info")
	fmt.Println("   • Audio: Should show consistent fingerprint")
	fmt.Println()
	fmt.Println("4. 🖼️ Take screenshot of results")
	fmt.Println("5. 🔧 Open DevTools (F12) → Console to verify spoofing messages")
	fmt.Println()
	fmt.Println("Expected Console Messages:")
	fmt.Println("   [Ghost CreepJS] ✅ Advanced spoofing active")
	fmt.Println("   [Ghost CreepJS] ✅ Navigator: hardwareConcurrency=8, deviceMemory=16")
	fmt.Println("   [Ghost CreepJS] ✅ Screen: 1920x1080")
	fmt.Println("   [Ghost CreepJS] ✅ WebDriver: false")
	fmt.Println()
	fmt.Println("⏳ Browser will stay open for 120 seconds for analysis...")
	
	// Keep browser open longer for CreepJS analysis
	time.Sleep(120 * time.Second)

	fmt.Println("\n🔍 ANALYSIS COMPLETE!")
	fmt.Println("Please report the CreepJS results:")
	fmt.Println("- Trust Score: ___% (Green/Yellow/Red)")
	fmt.Println("- Lies Detected: ___ items")
	fmt.Println("- Bot Detection: Human-like/Suspicious/Bot")
	fmt.Println("- Canvas Hash: Consistent/Inconsistent")
	fmt.Println("- WebGL Vendor: Intel Inc./Other")
	fmt.Println("- Audio Fingerprint: Consistent/Inconsistent")
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