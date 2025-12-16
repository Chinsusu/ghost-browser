package browser

import (
	"encoding/json"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/user/ghost-browser/internal/fingerprint"
	"github.com/user/ghost-browser/internal/profile"
)

func injectSpoofingScripts(browser *rod.Browser, p *profile.Profile) {
	script := generateSpoofScript(p.Fingerprint)

	// Inject on default page
	pages, _ := browser.Pages()
	for _, page := range pages {
		page.MustEvaluate(rod.Eval(script))
	}
}

func generateSpoofScript(fp *fingerprint.Fingerprint) string {
	fpJSON, _ := json.Marshal(fp)

	return fmt.Sprintf(`
(function() {
	'use strict';
	const CONFIG = %s;

	// ========== Navigator Spoofing ==========
	const navProps = {
		userAgent: CONFIG.navigator.userAgent,
		appVersion: CONFIG.navigator.appVersion,
		platform: CONFIG.navigator.platform,
		vendor: CONFIG.navigator.vendor,
		language: CONFIG.navigator.language,
		languages: Object.freeze(CONFIG.navigator.languages),
		hardwareConcurrency: CONFIG.navigator.hardwareConcurrency,
		deviceMemory: CONFIG.navigator.deviceMemory,
		maxTouchPoints: CONFIG.navigator.maxTouchPoints,
		doNotTrack: CONFIG.navigator.doNotTrack,
		cookieEnabled: CONFIG.navigator.cookieEnabled,
		webdriver: false,
	};

	for (const [prop, value] of Object.entries(navProps)) {
		try {
			Object.defineProperty(Navigator.prototype, prop, {
				get: () => value,
				configurable: true,
			});
		} catch (e) {}
	}

	// ========== Screen Spoofing ==========
	const screenProps = {
		width: CONFIG.screen.width,
		height: CONFIG.screen.height,
		availWidth: CONFIG.screen.availWidth,
		availHeight: CONFIG.screen.availHeight,
		colorDepth: CONFIG.screen.colorDepth,
		pixelDepth: CONFIG.screen.pixelDepth,
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

	// ========== WebGL Spoofing ==========
	const origGetParam = WebGLRenderingContext.prototype.getParameter;
	WebGLRenderingContext.prototype.getParameter = function(param) {
		if (param === 37445) return CONFIG.webgl.unmaskedVendor;
		if (param === 37446) return CONFIG.webgl.unmaskedRenderer;
		return origGetParam.call(this, param);
	};

	if (typeof WebGL2RenderingContext !== 'undefined') {
		const origGetParam2 = WebGL2RenderingContext.prototype.getParameter;
		WebGL2RenderingContext.prototype.getParameter = function(param) {
			if (param === 37445) return CONFIG.webgl.unmaskedVendor;
			if (param === 37446) return CONFIG.webgl.unmaskedRenderer;
			return origGetParam2.call(this, param);
		};
	}

	// ========== Canvas Noise ==========
	const noise = CONFIG.canvas.noise || 0.0001;
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
		const audioNoise = CONFIG.audio.noise || 0.0001;
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
	const targetTZ = CONFIG.timezone.timezone;
	const targetOffset = CONFIG.timezone.timezoneOffset;

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
	if (CONFIG.network.webRTCPolicy === 'disable') {
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

	// ========== Plugins Spoofing ==========
	const fakePlugins = CONFIG.misc.plugins.map(p => ({
		name: p.name,
		filename: p.filename,
		description: p.description,
		length: 1,
	}));

	Object.defineProperty(Navigator.prototype, 'plugins', {
		get: () => {
			const arr = Object.create(PluginArray.prototype);
			fakePlugins.forEach((p, i) => arr[i] = p);
			arr.length = fakePlugins.length;
			arr.item = i => fakePlugins[i] || null;
			arr.namedItem = name => fakePlugins.find(p => p.name === name) || null;
			arr.refresh = () => {};
			return arr;
		},
		configurable: true,
	});

	console.log('[Ghost Browser] Fingerprint spoofing active');
})();
`, string(fpJSON))
}
