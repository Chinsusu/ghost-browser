package browser

import (
	"fmt"

	"github.com/user/ghost-browser/internal/fingerprint"
)

// generateSpoofScript creates JavaScript code for fingerprint spoofing
// This script is injected via Edge's --user-script parameter for pre-load execution

func generateSpoofScript(fp *fingerprint.Fingerprint) string {
	return fmt.Sprintf(`
// Ghost Browser ChromeDP Fingerprint Spoofing Script
// This script runs BEFORE page load via page.AddScriptToEvaluateOnNewDocument
(function() {
	'use strict';
	
	console.log('[Ghost Browser] ChromeDP fingerprint spoofing initializing...');
	
	// ========== Navigator Spoofing ==========
	// CRITICAL: Must happen before page reads these values
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
	
	Object.defineProperty(Navigator.prototype, 'appVersion', {
		get: function() { return '%s'; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'vendor', {
		get: function() { return '%s'; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'language', {
		get: function() { return '%s'; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'maxTouchPoints', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'webdriver', {
		get: function() { return false; },
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'cookieEnabled', {
		get: function() { return %t; },
		configurable: true
	});
	
	// ========== Screen Spoofing ==========
	Object.defineProperty(Screen.prototype, 'width', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'height', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availWidth', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'availHeight', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'colorDepth', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(Screen.prototype, 'pixelDepth', {
		get: function() { return %d; },
		configurable: true
	});
	
	Object.defineProperty(window, 'devicePixelRatio', {
		get: function() { return %f; },
		configurable: true
	});
	
	// ========== WebGL Spoofing ==========
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
	
	// ========== Canvas Noise ==========
	const noise = %f;
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
	
	// ========== Audio Noise ==========
	if (typeof AudioContext !== 'undefined') {
		const audioNoise = %f;
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
	
	// ========== Timezone Spoofing ==========
	const targetTZ = '%s';
	const targetOffset = %d;
	
	Date.prototype.getTimezoneOffset = function() {
		return targetOffset;
	};
	
	const origDateTimeFormat = Intl.DateTimeFormat;
	Intl.DateTimeFormat = function(locales, options) {
		options = options || {};
		if (!options.timeZone) options.timeZone = targetTZ;
		return new origDateTimeFormat(locales, options);
	};
	Intl.DateTimeFormat.prototype = origDateTimeFormat.prototype;
	
	// ========== WebRTC Protection ==========
	if ('%s' === 'disable') {
		['RTCPeerConnection', 'webkitRTCPeerConnection', 'mozRTCPeerConnection'].forEach(name => {
			try {
				Object.defineProperty(window, name, {
					value: undefined,
					writable: false,
					configurable: false,
				});
			} catch (e) {}
		});
		
		if (navigator.mediaDevices) {
			navigator.mediaDevices.getUserMedia = () => Promise.reject(new Error('Permission denied'));
			navigator.mediaDevices.enumerateDevices = () => Promise.resolve([]);
		}
	}
	
	// ========== Remove Automation Flags ==========
	try { delete Object.getPrototypeOf(navigator).webdriver; } catch (e) {}
	
	const autoProps = [
		'__webdriver_evaluate', '__selenium_evaluate', '__webdriver_script_function',
		'__driver_evaluate', '_selenium', '_Selenium_IDE_Recorder', 'callSelenium',
		'$cdc_asdjflasutopfhvcZLmcfl_', '$chrome_asyncScriptInfo',
	];
	autoProps.forEach(prop => {
		try { if (window[prop]) delete window[prop]; } catch (e) {}
	});
	
	console.log('[Ghost Browser] ✅ Navigator spoofed (hardwareConcurrency: %d, deviceMemory: %d)');
	console.log('[Ghost Browser] ✅ Screen spoofed (%dx%d)');
	console.log('[Ghost Browser] ✅ WebGL spoofed (%s)');
	console.log('[Ghost Browser] Fingerprint spoofing active - ChromeDP pre-load injection');
	
})();
`, 
		fp.Navigator.HardwareConcurrency,
		fp.Navigator.DeviceMemory,
		fp.Navigator.Platform,
		fp.Navigator.UserAgent,
		fp.Navigator.AppVersion,
		fp.Navigator.Vendor,
		fp.Navigator.Language,
		fp.Navigator.MaxTouchPoints,
		fp.Navigator.CookieEnabled,
		fp.Screen.Width,
		fp.Screen.Height,
		fp.Screen.AvailWidth,
		fp.Screen.AvailHeight,
		fp.Screen.ColorDepth,
		fp.Screen.PixelDepth,
		fp.Screen.PixelRatio,
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
		fp.WebGL.Vendor,
		fp.WebGL.Renderer,
		fp.Canvas.Noise,
		fp.Audio.Noise,
		fp.Timezone.Timezone,
		fp.Timezone.TimezoneOffset,
		fp.Network.WebRTCPolicy,
		fp.Navigator.HardwareConcurrency,
		fp.Navigator.DeviceMemory,
		fp.Screen.Width,
		fp.Screen.Height,
		fp.WebGL.Vendor,
	)
}
