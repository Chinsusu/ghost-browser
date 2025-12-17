
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
