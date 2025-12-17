package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	fmt.Println("=== Ghost Browser - Edge Launcher Test ===")
	fmt.Println()

	// Step 1: Find Edge
	fmt.Println("[1] Finding Microsoft Edge...")
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("❌ Failed to find Edge: ", err)
	}
	fmt.Println("✅ Found Edge at:", edgePath)
	fmt.Println()

	// Step 2: Create temp profile directory
	fmt.Println("[2] Creating temp profile directory...")
	tempDir, err := os.MkdirTemp("", "ghost-browser-test-*")
	if err != nil {
		log.Fatal("❌ Failed to create temp dir: ", err)
	}
	defer os.RemoveAll(tempDir)
	fmt.Println("✅ Temp dir:", tempDir)
	fmt.Println()

	// Step 3: Launch Edge with CDP
	fmt.Println("[3] Launching Edge with CDP...")
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
		log.Fatal("❌ Failed to launch Edge: ", err)
	}
	fmt.Println("✅ Edge launched! Control URL:", controlURL)
	fmt.Println()

	// Step 4: Connect with Rod
	fmt.Println("[4] Connecting to browser...")
	browser := rod.New().ControlURL(controlURL).MustConnect()
	fmt.Println("✅ Connected!")
	fmt.Println()

	// Step 5: Inject fingerprint spoofing script
	fmt.Println("[5] Injecting fingerprint spoofing script...")
	
	spoofScript := generateTestSpoofScript()
	
	// Get the page and inject
	page := browser.MustPage("about:blank")
	page.MustEvaluate(spoofScript)
	fmt.Println("✅ Script injected!")
	fmt.Println()

	// Step 6: Navigate to test page
	fmt.Println("[6] Navigating to fingerprint test page...")
	page.MustNavigate("https://browserleaks.com/javascript")
	fmt.Println("✅ Opened browserleaks.com - Check the values!")
	fmt.Println()

	// Step 7: Also open CreepJS for comprehensive test
	fmt.Println("[7] Opening CreepJS in new tab...")
	page2 := browser.MustPage("https://abrahamjuliot.github.io/creepjs/")
	_ = page2
	fmt.Println("✅ Opened CreepJS - This shows detailed fingerprint analysis")
	fmt.Println()

	fmt.Println("===========================================")
	fmt.Println("Browser is running. Check the fingerprint values!")
	fmt.Println("Press Ctrl+C to close...")
	fmt.Println("===========================================")

	// Keep running until user stops
	select {}
}

func findEdgePath() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("this test is for Windows only")
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

	// Try via PATH
	if path, err := exec.LookPath("msedge"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Edge not found in common locations")
}

func generateTestSpoofScript() string {
	return `
(function() {
	'use strict';
	
	console.log('[Ghost Browser] Initializing fingerprint spoofing...');

	// ========== TEST CONFIG ==========
	const CONFIG = {
		navigator: {
			userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0',
			platform: 'Win32',
			vendor: 'Google Inc.',
			language: 'en-US',
			languages: ['en-US', 'en'],
			hardwareConcurrency: 8,
			deviceMemory: 16,
			maxTouchPoints: 0,
		},
		screen: {
			width: 1920,
			height: 1080,
			availWidth: 1920,
			availHeight: 1040,
			colorDepth: 24,
			pixelRatio: 1.0,
		},
		webgl: {
			vendor: 'Google Inc. (NVIDIA)',
			renderer: 'ANGLE (NVIDIA, NVIDIA GeForce RTX 3080 Direct3D11 vs_5_0 ps_5_0, D3D11)',
		},
		timezone: {
			name: 'America/New_York',
			offset: 300,
		},
		canvas: {
			noise: 0.0001,
		}
	};

	// ========== NAVIGATOR SPOOFING ==========
	const navProps = {
		userAgent: CONFIG.navigator.userAgent,
		appVersion: CONFIG.navigator.userAgent.replace('Mozilla/', ''),
		platform: CONFIG.navigator.platform,
		vendor: CONFIG.navigator.vendor,
		language: CONFIG.navigator.language,
		languages: Object.freeze(CONFIG.navigator.languages),
		hardwareConcurrency: CONFIG.navigator.hardwareConcurrency,
		deviceMemory: CONFIG.navigator.deviceMemory,
		maxTouchPoints: CONFIG.navigator.maxTouchPoints,
		webdriver: false,
	};

	for (const [prop, value] of Object.entries(navProps)) {
		try {
			Object.defineProperty(Navigator.prototype, prop, {
				get: () => value,
				configurable: true,
			});
		} catch (e) {
			console.warn('[Ghost] Failed to spoof navigator.' + prop, e);
		}
	}
	console.log('[Ghost Browser] ✅ Navigator spoofed');

	// ========== SCREEN SPOOFING ==========
	const screenProps = {
		width: CONFIG.screen.width,
		height: CONFIG.screen.height,
		availWidth: CONFIG.screen.availWidth,
		availHeight: CONFIG.screen.availHeight,
		colorDepth: CONFIG.screen.colorDepth,
		pixelDepth: CONFIG.screen.colorDepth,
	};

	for (const [prop, value] of Object.entries(screenProps)) {
		try {
			Object.defineProperty(Screen.prototype, prop, {
				get: () => value,
				configurable: true,
			});
		} catch (e) {}
	}

	Object.defineProperty(window, 'devicePixelRatio', {
		get: () => CONFIG.screen.pixelRatio,
		configurable: true,
	});
	console.log('[Ghost Browser] ✅ Screen spoofed');

	// ========== WEBGL SPOOFING ==========
	const getParameterProxyHandler = {
		apply: function(target, thisArg, args) {
			const param = args[0];
			// UNMASKED_VENDOR_WEBGL
			if (param === 37445) {
				return CONFIG.webgl.vendor;
			}
			// UNMASKED_RENDERER_WEBGL
			if (param === 37446) {
				return CONFIG.webgl.renderer;
			}
			return Reflect.apply(target, thisArg, args);
		}
	};

	try {
		WebGLRenderingContext.prototype.getParameter = new Proxy(
			WebGLRenderingContext.prototype.getParameter,
			getParameterProxyHandler
		);
		if (typeof WebGL2RenderingContext !== 'undefined') {
			WebGL2RenderingContext.prototype.getParameter = new Proxy(
				WebGL2RenderingContext.prototype.getParameter,
				getParameterProxyHandler
			);
		}
		console.log('[Ghost Browser] ✅ WebGL spoofed');
	} catch (e) {
		console.warn('[Ghost] WebGL spoof failed', e);
	}

	// ========== CANVAS NOISE ==========
	const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
	HTMLCanvasElement.prototype.toDataURL = function(type, quality) {
		const ctx = this.getContext('2d');
		if (ctx && this.width > 0 && this.height > 0) {
			try {
				const imageData = ctx.getImageData(0, 0, Math.min(this.width, 100), Math.min(this.height, 100));
				for (let i = 0; i < imageData.data.length; i += 4) {
					if (Math.random() < 0.01) {
						imageData.data[i] = Math.max(0, Math.min(255, imageData.data[i] + (Math.random() > 0.5 ? 1 : -1)));
					}
				}
				ctx.putImageData(imageData, 0, 0);
			} catch (e) {}
		}
		return originalToDataURL.call(this, type, quality);
	};
	console.log('[Ghost Browser] ✅ Canvas noise added');

	// ========== TIMEZONE SPOOFING ==========
	const targetOffset = CONFIG.timezone.offset;
	Date.prototype.getTimezoneOffset = function() {
		return targetOffset;
	};
	console.log('[Ghost Browser] ✅ Timezone spoofed');

	// ========== WEBRTC PROTECTION ==========
	try {
		Object.defineProperty(window, 'RTCPeerConnection', {
			value: undefined,
			writable: false,
			configurable: false,
		});
		Object.defineProperty(window, 'webkitRTCPeerConnection', {
			value: undefined,
			writable: false,
			configurable: false,
		});
		console.log('[Ghost Browser] ✅ WebRTC disabled');
	} catch (e) {}

	// ========== REMOVE AUTOMATION FLAGS ==========
	try {
		delete Object.getPrototypeOf(navigator).webdriver;
	} catch (e) {}

	// Remove Chrome automation properties
	const propsToDelete = [
		'__webdriver_evaluate',
		'__selenium_evaluate', 
		'__webdriver_script_function',
		'__driver_evaluate',
		'__selenium_unwrapped',
		'__fxdriver_unwrapped',
		'_Selenium_IDE_Recorder',
		'_selenium',
		'calledSelenium',
		'$cdc_asdjflasutopfhvcZLmcfl_',
		'$chrome_asyncScriptInfo',
	];
	propsToDelete.forEach(prop => {
		try { delete window[prop]; } catch (e) {}
	});
	console.log('[Ghost Browser] ✅ Automation flags removed');

	console.log('[Ghost Browser] ========================================');
	console.log('[Ghost Browser] Fingerprint spoofing complete!');
	console.log('[Ghost Browser] Check values:');
	console.log('[Ghost Browser] - navigator.userAgent:', navigator.userAgent);
	console.log('[Ghost Browser] - navigator.hardwareConcurrency:', navigator.hardwareConcurrency);
	console.log('[Ghost Browser] - navigator.deviceMemory:', navigator.deviceMemory);
	console.log('[Ghost Browser] - screen.width x height:', screen.width, 'x', screen.height);
	console.log('[Ghost Browser] - navigator.webdriver:', navigator.webdriver);
	console.log('[Ghost Browser] ========================================');
})();
`
}
