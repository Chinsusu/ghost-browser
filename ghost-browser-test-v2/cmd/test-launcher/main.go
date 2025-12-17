package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func main() {
	fmt.Println("=== Ghost Browser - Fixed Launcher Test ===")
	fmt.Println()

	// Step 1: Find Edge
	fmt.Println("[1] Finding Microsoft Edge...")
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("❌ Failed to find Edge: ", err)
	}
	fmt.Println("✅ Found Edge at:", edgePath)

	// Step 2: Create temp profile directory
	fmt.Println("[2] Creating temp profile directory...")
	tempDir, err := os.MkdirTemp("", "ghost-browser-test-*")
	if err != nil {
		log.Fatal("❌ Failed to create temp dir: ", err)
	}
	defer os.RemoveAll(tempDir)
	fmt.Println("✅ Temp dir:", tempDir)

	// Step 3: Launch Edge
	fmt.Println("[3] Launching Edge...")
	l := launcher.New().
		Bin(edgePath).
		UserDataDir(tempDir).
		Headless(false).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars").
		Set("no-first-run").
		Set("no-default-browser-check").
		Set("disable-extensions")

	controlURL, err := l.Launch()
	if err != nil {
		log.Fatal("❌ Failed to launch: ", err)
	}
	fmt.Println("✅ Edge launched!")

	// Step 4: Connect
	fmt.Println("[4] Connecting...")
	browser := rod.New().ControlURL(controlURL).MustConnect()
	fmt.Println("✅ Connected!")

	// Step 5: Add script to evaluate on new document (CRITICAL!)
	fmt.Println("[5] Setting up fingerprint spoofing (before page load)...")
	
	spoofScript := generateSpoofScript()
	
	// This is the KEY - inject script BEFORE any page loads
	browser.MustSetCookies() // Initialize
	
	// Get default page
	pages, _ := browser.Pages()
	var page *rod.Page
	if len(pages) > 0 {
		page = pages[0]
	} else {
		page = browser.MustPage("")
	}

	// Add script to run on every new document BEFORE page scripts run
	_, err = page.AddScriptToEvaluateOnNewDocument(spoofScript)
	if err != nil {
		log.Fatal("❌ Failed to add script: ", err)
	}
	fmt.Println("✅ Fingerprint spoofing script installed!")

	// Step 6: Navigate to test page
	fmt.Println("[6] Navigating to browserleaks.com...")
	page.MustNavigate("https://browserleaks.com/javascript")
	page.MustWaitLoad()
	fmt.Println("✅ Page loaded!")

	// Step 7: Verify spoofing worked
	fmt.Println("[7] Verifying spoofing...")
	
	// Check values via JavaScript
	result := page.MustEval(`() => {
		return {
			hardwareConcurrency: navigator.hardwareConcurrency,
			deviceMemory: navigator.deviceMemory,
			platform: navigator.platform,
			webdriver: navigator.webdriver,
			userAgent: navigator.userAgent.substring(0, 50) + '...',
			screenWidth: screen.width,
			screenHeight: screen.height,
		}
	}`)
	
	fmt.Println()
	fmt.Println("=== SPOOFED VALUES ===")
	fmt.Printf("hardwareConcurrency: %v (should be 8)\n", result.Get("hardwareConcurrency").Int())
	fmt.Printf("deviceMemory: %v (should be 16)\n", result.Get("deviceMemory").Int())
	fmt.Printf("platform: %v\n", result.Get("platform").String())
	fmt.Printf("webdriver: %v (should be false)\n", result.Get("webdriver").Bool())
	fmt.Printf("screen: %vx%v (should be 1920x1080)\n", result.Get("screenWidth").Int(), result.Get("screenHeight").Int())
	fmt.Println("======================")
	fmt.Println()

	// Step 8: Open CreepJS in new tab
	fmt.Println("[8] Opening CreepJS for full analysis...")
	page2 := browser.MustPage("")
	page2.AddScriptToEvaluateOnNewDocument(spoofScript)
	page2.MustNavigate("https://abrahamjuliot.github.io/creepjs/")
	
	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("Browser is running with fingerprint spoofing!")
	fmt.Println("Check the values in browserleaks.com")
	fmt.Println("Press Ctrl+C to close...")
	fmt.Println("===========================================")

	select {}
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

	return "", fmt.Errorf("Edge not found")
}

// AddScriptToEvaluateOnNewDocument adds script that runs before page load
func (p *rod.Page) AddScriptToEvaluateOnNewDocument(script string) (string, error) {
	result, err := proto.PageAddScriptToEvaluateOnNewDocument{
		Source: script,
	}.Call(p)
	if err != nil {
		return "", err
	}
	return string(result.Identifier), nil
}

func generateSpoofScript() string {
	return `
(function() {
	'use strict';

	// ========== CONFIG - Values to spoof ==========
	const SPOOF = {
		navigator: {
			hardwareConcurrency: 8,
			deviceMemory: 16,
			platform: 'Win32',
			vendor: 'Google Inc.',
			language: 'en-US',
			languages: ['en-US', 'en'],
			maxTouchPoints: 0,
			webdriver: false,
		},
		screen: {
			width: 1920,
			height: 1080,
			availWidth: 1920,
			availHeight: 1040,
			colorDepth: 24,
			pixelDepth: 24,
			pixelRatio: 1.0,
		},
		webgl: {
			vendor: 'Google Inc. (NVIDIA)',
			renderer: 'ANGLE (NVIDIA, NVIDIA GeForce RTX 3080 Direct3D11 vs_5_0 ps_5_0, D3D11)',
		},
		timezone: {
			offset: 300, // EST
		}
	};

	// ========== NAVIGATOR ==========
	const navigatorProps = {
		hardwareConcurrency: { value: SPOOF.navigator.hardwareConcurrency },
		deviceMemory: { value: SPOOF.navigator.deviceMemory },
		platform: { value: SPOOF.navigator.platform },
		vendor: { value: SPOOF.navigator.vendor },
		language: { value: SPOOF.navigator.language },
		languages: { value: Object.freeze([...SPOOF.navigator.languages]) },
		maxTouchPoints: { value: SPOOF.navigator.maxTouchPoints },
		webdriver: { value: SPOOF.navigator.webdriver },
	};

	for (const [prop, desc] of Object.entries(navigatorProps)) {
		try {
			Object.defineProperty(Navigator.prototype, prop, {
				get: function() { return desc.value; },
				configurable: true,
				enumerable: true,
			});
		} catch (e) {}
	}

	// Also override on navigator instance directly
	try {
		Object.defineProperty(navigator, 'hardwareConcurrency', {
			get: () => SPOOF.navigator.hardwareConcurrency,
			configurable: true,
		});
		Object.defineProperty(navigator, 'deviceMemory', {
			get: () => SPOOF.navigator.deviceMemory,
			configurable: true,
		});
		Object.defineProperty(navigator, 'webdriver', {
			get: () => false,
			configurable: true,
		});
	} catch (e) {}

	// ========== SCREEN ==========
	const screenProps = ['width', 'height', 'availWidth', 'availHeight', 'colorDepth', 'pixelDepth'];
	for (const prop of screenProps) {
		try {
			Object.defineProperty(Screen.prototype, prop, {
				get: function() { return SPOOF.screen[prop]; },
				configurable: true,
			});
		} catch (e) {}
	}

	try {
		Object.defineProperty(window, 'devicePixelRatio', {
			get: () => SPOOF.screen.pixelRatio,
			configurable: true,
		});
		Object.defineProperty(window, 'innerWidth', {
			get: () => SPOOF.screen.width,
			configurable: true,
		});
		Object.defineProperty(window, 'innerHeight', {
			get: () => SPOOF.screen.height - 100,
			configurable: true,
		});
		Object.defineProperty(window, 'outerWidth', {
			get: () => SPOOF.screen.width,
			configurable: true,
		});
		Object.defineProperty(window, 'outerHeight', {
			get: () => SPOOF.screen.height,
			configurable: true,
		});
	} catch (e) {}

	// ========== WEBGL ==========
	const getParamHandler = {
		apply(target, thisArg, args) {
			const param = args[0];
			// UNMASKED_VENDOR_WEBGL
			if (param === 37445) return SPOOF.webgl.vendor;
			// UNMASKED_RENDERER_WEBGL  
			if (param === 37446) return SPOOF.webgl.renderer;
			return Reflect.apply(target, thisArg, args);
		}
	};

	try {
		WebGLRenderingContext.prototype.getParameter = new Proxy(
			WebGLRenderingContext.prototype.getParameter,
			getParamHandler
		);
	} catch (e) {}

	try {
		if (typeof WebGL2RenderingContext !== 'undefined') {
			WebGL2RenderingContext.prototype.getParameter = new Proxy(
				WebGL2RenderingContext.prototype.getParameter,
				getParamHandler
			);
		}
	} catch (e) {}

	// ========== CANVAS NOISE ==========
	const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
	const originalToBlob = HTMLCanvasElement.prototype.toBlob;
	const originalGetImageData = CanvasRenderingContext2D.prototype.getImageData;

	function addNoise(canvas) {
		try {
			const ctx = canvas.getContext('2d');
			if (!ctx || canvas.width === 0 || canvas.height === 0) return;
			
			const imageData = originalGetImageData.call(ctx, 0, 0, canvas.width, canvas.height);
			const data = imageData.data;
			
			// Add subtle noise
			for (let i = 0; i < data.length; i += 4) {
				if (Math.random() < 0.01) { // 1% of pixels
					const noise = Math.random() > 0.5 ? 1 : -1;
					data[i] = Math.max(0, Math.min(255, data[i] + noise));
				}
			}
			
			ctx.putImageData(imageData, 0, 0);
		} catch (e) {}
	}

	HTMLCanvasElement.prototype.toDataURL = function() {
		addNoise(this);
		return originalToDataURL.apply(this, arguments);
	};

	HTMLCanvasElement.prototype.toBlob = function() {
		addNoise(this);
		return originalToBlob.apply(this, arguments);
	};

	// ========== AUDIO NOISE ==========
	try {
		const originalGetChannelData = AudioBuffer.prototype.getChannelData;
		AudioBuffer.prototype.getChannelData = function() {
			const data = originalGetChannelData.apply(this, arguments);
			for (let i = 0; i < data.length; i += 100) {
				data[i] = data[i] + (Math.random() * 0.0001 - 0.00005);
			}
			return data;
		};
	} catch (e) {}

	// ========== TIMEZONE ==========
	Date.prototype.getTimezoneOffset = function() {
		return SPOOF.timezone.offset;
	};

	// ========== WEBRTC DISABLE ==========
	try {
		Object.defineProperty(window, 'RTCPeerConnection', { value: undefined, writable: false });
		Object.defineProperty(window, 'webkitRTCPeerConnection', { value: undefined, writable: false });
		Object.defineProperty(window, 'mozRTCPeerConnection', { value: undefined, writable: false });
	} catch (e) {}

	// ========== REMOVE AUTOMATION FLAGS ==========
	try {
		delete Object.getPrototypeOf(navigator).webdriver;
	} catch (e) {}

	const automationProps = [
		'__webdriver_evaluate', '__selenium_evaluate', '__webdriver_script_function',
		'__driver_evaluate', '__selenium_unwrapped', '__fxdriver_unwrapped',
		'_Selenium_IDE_Recorder', '_selenium', 'calledSelenium',
		'$cdc_asdjflasutopfhvcZLmcfl_', '$chrome_asyncScriptInfo'
	];
	automationProps.forEach(prop => {
		try { delete window[prop]; } catch (e) {}
	});

	// ========== PLUGINS ==========
	try {
		Object.defineProperty(Navigator.prototype, 'plugins', {
			get: function() {
				const fakePlugins = [
					{ name: 'PDF Viewer', filename: 'internal-pdf-viewer', description: 'Portable Document Format' },
					{ name: 'Chrome PDF Viewer', filename: 'internal-pdf-viewer', description: 'Portable Document Format' },
					{ name: 'Chromium PDF Viewer', filename: 'internal-pdf-viewer', description: 'Portable Document Format' },
				];
				
				const arr = Object.create(PluginArray.prototype);
				fakePlugins.forEach((p, i) => { arr[i] = p; });
				arr.length = fakePlugins.length;
				arr.item = i => fakePlugins[i] || null;
				arr.namedItem = name => fakePlugins.find(p => p.name === name) || null;
				arr.refresh = () => {};
				return arr;
			},
			configurable: true,
		});
	} catch (e) {}

	console.log('[Ghost Browser] ✅ Fingerprint spoofing active');
	console.log('[Ghost Browser] hardwareConcurrency:', navigator.hardwareConcurrency);
	console.log('[Ghost Browser] deviceMemory:', navigator.deviceMemory);
	console.log('[Ghost Browser] webdriver:', navigator.webdriver);
})();
`
}
