package browser

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/user/ghost-browser/internal/fingerprint"
)

// GenerateSpoofScript creates the JavaScript injection script
func GenerateSpoofScript(fp *fingerprint.Fingerprint) string {
	noiseSeed := generateNoiseSeed()
	languagesJS := formatLanguagesJS(fp.Navigator.Languages)

	return fmt.Sprintf(`
(function() {
	'use strict';
	
	if (window.__ghostInjected) return;
	window.__ghostInjected = true;

	const SPOOF = {
		hardwareConcurrency: %d,
		deviceMemory: %d,
		platform: '%s',
		vendor: '%s',
		language: '%s',
		languages: [%s],
		maxTouchPoints: %d,
		screen: { width: %d, height: %d, availWidth: %d, availHeight: %d, colorDepth: %d, pixelDepth: %d },
		pixelRatio: %.2f,
		webgl: { vendor: '%s', renderer: '%s' },
		timezone: { name: '%s', offset: %d },
		noiseSeed: '%s',
	};

	// ============================================
	// CRITICAL: hardwareConcurrency fix
	// Must use Object.defineProperty on navigator directly
	// AND on Navigator.prototype
	// ============================================
	
	// Method 1: Direct on navigator object
	try {
		Object.defineProperty(navigator, 'hardwareConcurrency', {
			get: function() { return SPOOF.hardwareConcurrency; },
			configurable: true,
			enumerable: true
		});
	} catch(e) { console.warn('[Ghost] navigator.hardwareConcurrency failed:', e); }

	// Method 2: On Navigator.prototype (backup)
	try {
		Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
			get: function() { return SPOOF.hardwareConcurrency; },
			configurable: true,
			enumerable: true
		});
	} catch(e) {}

	// Method 3: Using Proxy (most robust)
	try {
		const navigatorProxy = new Proxy(navigator, {
			get: function(target, prop) {
				if (prop === 'hardwareConcurrency') return SPOOF.hardwareConcurrency;
				if (prop === 'deviceMemory') return SPOOF.deviceMemory;
				if (prop === 'platform') return SPOOF.platform;
				if (prop === 'vendor') return SPOOF.vendor;
				if (prop === 'language') return SPOOF.language;
				if (prop === 'languages') return Object.freeze(SPOOF.languages);
				if (prop === 'maxTouchPoints') return SPOOF.maxTouchPoints;
				if (prop === 'webdriver') return false;
				
				const value = target[prop];
				if (typeof value === 'function') {
					return value.bind(target);
				}
				return value;
			}
		});
		
		// Try to replace navigator (works in some contexts)
		try {
			Object.defineProperty(window, 'navigator', {
				get: function() { return navigatorProxy; },
				configurable: true
			});
		} catch(e) {}
	} catch(e) {}

	// ============================================
	// deviceMemory - Multiple methods
	// ============================================
	try {
		Object.defineProperty(navigator, 'deviceMemory', {
			get: function() { return SPOOF.deviceMemory; },
			configurable: true,
			enumerable: true
		});
	} catch(e) {}

	try {
		Object.defineProperty(Navigator.prototype, 'deviceMemory', {
			get: function() { return SPOOF.deviceMemory; },
			configurable: true,
			enumerable: true
		});
	} catch(e) {}

	// ============================================
	// Other Navigator properties
	// ============================================
	const navProps = {
		platform: SPOOF.platform,
		vendor: SPOOF.vendor,
		language: SPOOF.language,
		languages: Object.freeze(SPOOF.languages),
		maxTouchPoints: SPOOF.maxTouchPoints,
		webdriver: false,
	};

	for (const [prop, value] of Object.entries(navProps)) {
		try {
			Object.defineProperty(navigator, prop, {
				get: () => value,
				configurable: true,
				enumerable: true
			});
		} catch(e) {}
		
		try {
			Object.defineProperty(Navigator.prototype, prop, {
				get: () => value,
				configurable: true,
				enumerable: true
			});
		} catch(e) {}
	}

	// ============================================
	// Screen
	// ============================================
	const screenProps = ['width', 'height', 'availWidth', 'availHeight', 'colorDepth', 'pixelDepth'];
	for (const prop of screenProps) {
		try {
			Object.defineProperty(screen, prop, {
				get: () => SPOOF.screen[prop],
				configurable: true
			});
		} catch(e) {}
		try {
			Object.defineProperty(Screen.prototype, prop, {
				get: () => SPOOF.screen[prop],
				configurable: true
			});
		} catch(e) {}
	}

	try {
		Object.defineProperty(window, 'devicePixelRatio', {
			get: () => SPOOF.pixelRatio,
			configurable: true
		});
	} catch(e) {}

	// ============================================
	// WebRTC Block
	// ============================================
	const rtcProps = ['RTCPeerConnection', 'webkitRTCPeerConnection', 'mozRTCPeerConnection', 
					  'RTCSessionDescription', 'RTCIceCandidate', 'RTCDataChannel'];
	for (const prop of rtcProps) {
		try {
			Object.defineProperty(window, prop, { value: undefined, writable: false, configurable: false });
		} catch(e) {}
	}

	if (navigator.mediaDevices) {
		try {
			navigator.mediaDevices.getUserMedia = () => Promise.reject(new DOMException('Permission denied', 'NotAllowedError'));
			navigator.mediaDevices.enumerateDevices = () => Promise.resolve([]);
		} catch(e) {}
	}

	// ============================================
	// WebGL
	// ============================================
	const webglSpoof = function(origFn) {
		return function(param) {
			if (param === 37445) return SPOOF.webgl.vendor;
			if (param === 37446) return SPOOF.webgl.renderer;
			return origFn.call(this, param);
		};
	};

	try {
		WebGLRenderingContext.prototype.getParameter = webglSpoof(WebGLRenderingContext.prototype.getParameter);
	} catch(e) {}

	try {
		if (typeof WebGL2RenderingContext !== 'undefined') {
			WebGL2RenderingContext.prototype.getParameter = webglSpoof(WebGL2RenderingContext.prototype.getParameter);
		}
	} catch(e) {}

	// ============================================
	// Canvas Noise
	// ============================================
	const rng = (function(seed) {
		let h = 0;
		for (let i = 0; i < seed.length; i++) h = Math.imul(31, h) + seed.charCodeAt(i) | 0;
		return function() {
			h = Math.imul(h ^ (h >>> 15), h | 1);
			h ^= h + Math.imul(h ^ (h >>> 7), h | 61);
			return ((h ^ (h >>> 14)) >>> 0) / 4294967296;
		};
	})(SPOOF.noiseSeed);

	const origToDataURL = HTMLCanvasElement.prototype.toDataURL;
	HTMLCanvasElement.prototype.toDataURL = function() {
		try {
			const ctx = this.getContext('2d');
			if (ctx && this.width > 0 && this.height > 0) {
				const img = ctx.getImageData(0, 0, Math.min(this.width, 100), Math.min(this.height, 100));
				for (let i = 0; i < img.data.length; i += 4) {
					if (rng() < 0.02) img.data[i] ^= 1;
				}
				ctx.putImageData(img, 0, 0);
			}
		} catch(e) {}
		return origToDataURL.apply(this, arguments);
	};

	// ============================================
	// Audio Noise
	// ============================================
	try {
		const origGetChannelData = AudioBuffer.prototype.getChannelData;
		AudioBuffer.prototype.getChannelData = function(ch) {
			const data = origGetChannelData.call(this, ch);
			for (let i = 0; i < data.length; i += 100) data[i] += (rng() - 0.5) * 0.0001;
			return data;
		};
	} catch(e) {}

	// ============================================
	// Timezone
	// ============================================
	Date.prototype.getTimezoneOffset = function() { return SPOOF.timezone.offset; };

	const origDTF = Intl.DateTimeFormat;
	Intl.DateTimeFormat = function(loc, opt) {
		opt = Object.assign({}, opt);
		if (!opt.timeZone) opt.timeZone = SPOOF.timezone.name;
		return new origDTF(loc, opt);
	};
	Intl.DateTimeFormat.prototype = origDTF.prototype;
	Intl.DateTimeFormat.supportedLocalesOf = origDTF.supportedLocalesOf;

	// ============================================
	// Automation flags
	// ============================================
	try { delete Object.getPrototypeOf(navigator).webdriver; } catch(e) {}
	
	['__webdriver_evaluate', '__selenium_evaluate', '__driver_evaluate', '_selenium',
	 '$cdc_asdjflasutopfhvcZLmcfl_', '$chrome_asyncScriptInfo', 'domAutomation'
	].forEach(p => { try { delete window[p]; } catch(e) {} });

	// ============================================
	// Plugins
	// ============================================
	const fakePlugins = [
		{name:'PDF Viewer', filename:'internal-pdf-viewer', description:'PDF', length:1},
		{name:'Chrome PDF Viewer', filename:'internal-pdf-viewer', description:'PDF', length:1},
	];

	try {
		Object.defineProperty(navigator, 'plugins', {
			get: () => {
				const arr = Object.create(PluginArray.prototype);
				fakePlugins.forEach((p,i) => arr[i] = p);
				arr.length = fakePlugins.length;
				arr.item = i => fakePlugins[i] || null;
				arr.namedItem = n => fakePlugins.find(p => p.name === n) || null;
				arr.refresh = () => {};
				return arr;
			},
			configurable: true
		});
	} catch(e) {}

	// ============================================
	// Battery
	// ============================================
	try {
		navigator.getBattery = () => Promise.resolve({
			charging: true, chargingTime: 0, dischargingTime: Infinity, level: 1,
			addEventListener: () => {}, removeEventListener: () => {}
		});
	} catch(e) {}

	// ============================================
	// Verify & Log
	// ============================================
	console.log('[Ghost] ✅ Spoofing active');
	console.log('[Ghost] hardwareConcurrency:', navigator.hardwareConcurrency, '(expected:', SPOOF.hardwareConcurrency, ')');
	console.log('[Ghost] deviceMemory:', navigator.deviceMemory, '(expected:', SPOOF.deviceMemory, ')');
	console.log('[Ghost] screen:', screen.width + 'x' + screen.height);

})();
`,
		// Navigator
		fp.Navigator.HardwareConcurrency,
		fp.Navigator.DeviceMemory,
		fp.Navigator.Platform,
		fp.Navigator.Vendor,
		fp.Navigator.Language,
		languagesJS,
		fp.Navigator.MaxTouchPoints,
		// Screen
		fp.Screen.Width,
		fp.Screen.Height,
		fp.Screen.AvailWidth,
		fp.Screen.AvailHeight,
		fp.Screen.ColorDepth,
		fp.Screen.PixelDepth,
		fp.Screen.PixelRatio,
		// WebGL
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
		// Timezone
		fp.Timezone.Timezone,
		fp.Timezone.TimezoneOffset,
		// Noise
		noiseSeed,
	)
}

func generateNoiseSeed() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func formatLanguagesJS(languages []string) string {
	if len(languages) == 0 {
		return "'en-US', 'en'"
	}
	quoted := make([]string, len(languages))
	for i, lang := range languages {
		quoted[i] = fmt.Sprintf("'%s'", lang)
	}
	return strings.Join(quoted, ", ")
}
