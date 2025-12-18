package browser

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/user/ghost-browser/internal/fingerprint"
)

// GenerateSpoofScript creates the JavaScript injection script
// This script runs BEFORE any page JavaScript
func GenerateSpoofScript(fp *fingerprint.Fingerprint) string {
	noiseSeed := generateNoiseSeed()
	languagesJS := formatLanguagesJS(fp.Navigator.Languages)

	// The script is wrapped in an IIFE that executes immediately
	// and cannot be overridden by page scripts
	return fmt.Sprintf(`
// Ghost Browser Fingerprint Spoofing - Injected before page load
(function(window, document, navigator, screen) {
	'use strict';
	
	// Prevent re-injection
	if (window.__ghostBrowserInjected) return;
	window.__ghostBrowserInjected = true;

	const SPOOF = {
		navigator: {
			hardwareConcurrency: %d,
			deviceMemory: %d,
			platform: '%s',
			vendor: '%s',
			language: '%s',
			languages: Object.freeze([%s]),
			maxTouchPoints: %d,
		},
		screen: {
			width: %d,
			height: %d,
			availWidth: %d,
			availHeight: %d,
			colorDepth: %d,
			pixelDepth: %d,
		},
		window: {
			devicePixelRatio: %.2f,
		},
		webgl: {
			vendor: '%s',
			renderer: '%s',
		},
		timezone: {
			name: '%s',
			offset: %d,
		},
		noise: {
			seed: '%s',
		}
	};

	// ===== UTILS =====
	const defineGetter = (obj, prop, getter) => {
		try {
			Object.defineProperty(obj, prop, {
				get: getter,
				configurable: true,
				enumerable: true,
			});
		} catch (e) {}
	};

	const defineValue = (obj, prop, value) => {
		try {
			Object.defineProperty(obj, prop, {
				value: value,
				writable: false,
				configurable: false,
			});
		} catch (e) {}
	};

	// ===== 1. NAVIGATOR =====
	defineGetter(Navigator.prototype, 'hardwareConcurrency', () => SPOOF.navigator.hardwareConcurrency);
	defineGetter(Navigator.prototype, 'deviceMemory', () => SPOOF.navigator.deviceMemory);
	defineGetter(Navigator.prototype, 'platform', () => SPOOF.navigator.platform);
	defineGetter(Navigator.prototype, 'vendor', () => SPOOF.navigator.vendor);
	defineGetter(Navigator.prototype, 'language', () => SPOOF.navigator.language);
	defineGetter(Navigator.prototype, 'languages', () => SPOOF.navigator.languages);
	defineGetter(Navigator.prototype, 'maxTouchPoints', () => SPOOF.navigator.maxTouchPoints);
	defineGetter(Navigator.prototype, 'webdriver', () => false);

	// ===== 2. SCREEN =====
	defineGetter(Screen.prototype, 'width', () => SPOOF.screen.width);
	defineGetter(Screen.prototype, 'height', () => SPOOF.screen.height);
	defineGetter(Screen.prototype, 'availWidth', () => SPOOF.screen.availWidth);
	defineGetter(Screen.prototype, 'availHeight', () => SPOOF.screen.availHeight);
	defineGetter(Screen.prototype, 'colorDepth', () => SPOOF.screen.colorDepth);
	defineGetter(Screen.prototype, 'pixelDepth', () => SPOOF.screen.pixelDepth);

	// ===== 3. WINDOW =====
	defineGetter(window, 'devicePixelRatio', () => SPOOF.window.devicePixelRatio);
	defineGetter(window, 'innerWidth', () => SPOOF.screen.width);
	defineGetter(window, 'innerHeight', () => SPOOF.screen.height - 100);
	defineGetter(window, 'outerWidth', () => SPOOF.screen.width);
	defineGetter(window, 'outerHeight', () => SPOOF.screen.height);

	// ===== 4. WEBRTC BLOCK =====
	defineValue(window, 'RTCPeerConnection', undefined);
	defineValue(window, 'webkitRTCPeerConnection', undefined);
	defineValue(window, 'mozRTCPeerConnection', undefined);
	defineValue(window, 'RTCSessionDescription', undefined);
	defineValue(window, 'RTCIceCandidate', undefined);

	if (navigator.mediaDevices) {
		navigator.mediaDevices.getUserMedia = () => Promise.reject(new DOMException('Permission denied', 'NotAllowedError'));
		navigator.mediaDevices.enumerateDevices = () => Promise.resolve([]);
	}

	// ===== 5. WEBGL =====
	const webglGetParameter = function(originalFn) {
		return function(param) {
			if (param === 37445) return SPOOF.webgl.vendor;
			if (param === 37446) return SPOOF.webgl.renderer;
			return originalFn.call(this, param);
		};
	};

	try {
		const origWGL = WebGLRenderingContext.prototype.getParameter;
		WebGLRenderingContext.prototype.getParameter = webglGetParameter(origWGL);
	} catch (e) {}

	try {
		if (typeof WebGL2RenderingContext !== 'undefined') {
			const origWGL2 = WebGL2RenderingContext.prototype.getParameter;
			WebGL2RenderingContext.prototype.getParameter = webglGetParameter(origWGL2);
		}
	} catch (e) {}

	// ===== 6. CANVAS NOISE =====
	const seedRandom = (function() {
		let seed = '%s'.split('').reduce((a, c) => (a * 31 + c.charCodeAt(0)) | 0, 0);
		return function() {
			seed = Math.imul(seed ^ (seed >>> 15), seed | 1);
			seed ^= seed + Math.imul(seed ^ (seed >>> 7), seed | 61);
			return ((seed ^ (seed >>> 14)) >>> 0) / 4294967296;
		};
	})();

	const origToDataURL = HTMLCanvasElement.prototype.toDataURL;
	HTMLCanvasElement.prototype.toDataURL = function() {
		try {
			const ctx = this.getContext('2d');
			if (ctx && this.width > 0 && this.height > 0) {
				const imgData = ctx.getImageData(0, 0, Math.min(this.width, 100), Math.min(this.height, 100));
				for (let i = 0; i < imgData.data.length; i += 4) {
					if (seedRandom() < 0.02) imgData.data[i] ^= 1;
				}
				ctx.putImageData(imgData, 0, 0);
			}
		} catch (e) {}
		return origToDataURL.apply(this, arguments);
	};

	// ===== 7. AUDIO NOISE =====
	try {
		const origGetChannelData = AudioBuffer.prototype.getChannelData;
		AudioBuffer.prototype.getChannelData = function(ch) {
			const data = origGetChannelData.call(this, ch);
			for (let i = 0; i < data.length; i += 100) {
				data[i] += (seedRandom() - 0.5) * 0.0001;
			}
			return data;
		};
	} catch (e) {}

	// ===== 8. TIMEZONE =====
	Date.prototype.getTimezoneOffset = function() { return SPOOF.timezone.offset; };

	const origDTF = Intl.DateTimeFormat;
	Intl.DateTimeFormat = function(loc, opt) {
		opt = opt || {};
		if (!opt.timeZone) opt.timeZone = SPOOF.timezone.name;
		return new origDTF(loc, opt);
	};
	Intl.DateTimeFormat.prototype = origDTF.prototype;
	Intl.DateTimeFormat.supportedLocalesOf = origDTF.supportedLocalesOf;

	// ===== 9. AUTOMATION FLAGS =====
	try { delete Object.getPrototypeOf(navigator).webdriver; } catch (e) {}
	
	['__webdriver_evaluate', '__selenium_evaluate', '__driver_evaluate',
	 '_selenium', '$cdc_asdjflasutopfhvcZLmcfl_', '$chrome_asyncScriptInfo',
	 'domAutomation', 'domAutomationController'
	].forEach(p => { try { delete window[p]; } catch(e) {} });

	// ===== 10. PLUGINS =====
	const fakePlugins = [
		{name: 'PDF Viewer', filename: 'internal-pdf-viewer', description: 'Portable Document Format', length: 1},
		{name: 'Chrome PDF Viewer', filename: 'internal-pdf-viewer', description: 'Portable Document Format', length: 1},
		{name: 'Chromium PDF Viewer', filename: 'internal-pdf-viewer', description: 'Portable Document Format', length: 1},
	];

	defineGetter(Navigator.prototype, 'plugins', () => {
		const arr = Object.create(PluginArray.prototype);
		fakePlugins.forEach((p, i) => arr[i] = p);
		arr.length = fakePlugins.length;
		arr.item = i => fakePlugins[i] || null;
		arr.namedItem = n => fakePlugins.find(p => p.name === n) || null;
		arr.refresh = () => {};
		return arr;
	});

	// ===== 11. BATTERY =====
	navigator.getBattery = () => Promise.resolve({
		charging: true, chargingTime: 0, dischargingTime: Infinity, level: 1,
		addEventListener: () => {}, removeEventListener: () => {},
	});

	// ===== 12. MATCHMEDIA FIX =====
	const origMatchMedia = window.matchMedia;
	window.matchMedia = function(q) {
		const r = origMatchMedia.call(this, q);
		if (q.includes('device-width') || q.includes('device-height')) {
			return {...r, matches: true};
		}
		return r;
	};

	console.log('[Ghost] ✅ Spoofing active | CPU:' + navigator.hardwareConcurrency + ' RAM:' + navigator.deviceMemory + 'GB');

})(window, document, navigator, screen);
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
		// Window
		fp.Screen.PixelRatio,
		// WebGL
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
		// Timezone
		fp.Timezone.Timezone,
		fp.Timezone.TimezoneOffset,
		// Noise seed (used twice)
		noiseSeed,
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
