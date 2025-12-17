package fingerprint

import (
	"crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"
)

type Generator struct {
	rng *mathrand.Rand
}

type GenerateOptions struct {
	OS      string
	Browser string
}

func NewGenerator() *Generator {
	// Use crypto/rand for truly unique seed
	var seedBytes [8]byte
	rand.Read(seedBytes[:])
	
	// Convert to int64 seed
	seed := int64(0)
	for i, b := range seedBytes {
		seed |= int64(b) << (i * 8)
	}
	
	// Combine with nanosecond timestamp for extra entropy
	seed ^= time.Now().UnixNano()
	
	return &Generator{rng: mathrand.New(mathrand.NewSource(seed))}
}

// Generate creates a random realistic fingerprint
func (g *Generator) Generate(opts *GenerateOptions) *Fingerprint {
	if opts == nil {
		opts = &GenerateOptions{OS: "windows", Browser: "edge"}
	}

	// Screen resolutions (common)
	screens := []struct{ w, h int }{
		{1920, 1080}, {2560, 1440}, {1366, 768}, {1536, 864},
		{1440, 900}, {1680, 1050}, {1280, 720}, {1600, 900},
	}
	screen := screens[g.rng.Intn(len(screens))]

	// GPUs (realistic mix)
	gpus := []struct{ vendor, renderer string }{
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce RTX 4090 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce RTX 4080 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce RTX 3080 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce RTX 3070 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce RTX 3060 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 Ti Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (NVIDIA)", "ANGLE (NVIDIA, NVIDIA GeForce GTX 1650 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (Intel)", "ANGLE (Intel, Intel(R) UHD Graphics 630 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (Intel)", "ANGLE (Intel, Intel(R) UHD Graphics 770 Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (Intel)", "ANGLE (Intel, Intel(R) Iris(R) Xe Graphics Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (AMD)", "ANGLE (AMD, AMD Radeon RX 6800 XT Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (AMD)", "ANGLE (AMD, AMD Radeon RX 6700 XT Direct3D11 vs_5_0 ps_5_0, D3D11)"},
		{"Google Inc. (AMD)", "ANGLE (AMD, AMD Radeon RX 5700 XT Direct3D11 vs_5_0 ps_5_0, D3D11)"},
	}
	gpu := gpus[g.rng.Intn(len(gpus))]

	// Timezones with correct offsets
	timezones := []struct {
		tz     string
		offset int
		locale string
	}{
		{"America/New_York", 300, "en-US"},
		{"America/Chicago", 360, "en-US"},
		{"America/Denver", 420, "en-US"},
		{"America/Los_Angeles", 480, "en-US"},
		{"Europe/London", 0, "en-GB"},
		{"Europe/Paris", -60, "fr-FR"},
		{"Europe/Berlin", -60, "de-DE"},
		{"Asia/Tokyo", -540, "ja-JP"},
		{"Asia/Shanghai", -480, "zh-CN"},
		{"Asia/Singapore", -480, "en-SG"},
		{"Australia/Sydney", -660, "en-AU"},
	}
	tz := timezones[g.rng.Intn(len(timezones))]

	// Common values
	cores := []int{4, 6, 8, 12, 16}
	rams := []int{4, 8, 16, 32}
	ratios := []float64{1.0, 1.25, 1.5, 2.0}

	// Languages based on timezone locale
	languages := g.getLanguagesForLocale(tz.locale)

	// User agent
	userAgent := g.generateUserAgent(opts)

	return &Fingerprint{
		Navigator: NavigatorFP{
			UserAgent:           userAgent,
			AppVersion:          userAgent[8:],
			Platform:            "Win32",
			Vendor:              "Google Inc.",
			Language:            languages[0],
			Languages:           languages,
			HardwareConcurrency: cores[g.rng.Intn(len(cores))],
			DeviceMemory:        rams[g.rng.Intn(len(rams))],
			MaxTouchPoints:      0,
			ProductSub:          "20030107",
			DoNotTrack:          g.pick([]string{"1", "null", "unspecified"}),
			CookieEnabled:       true,
			Webdriver:           false,
		},
		Screen: ScreenFP{
			Width:       screen.w,
			Height:      screen.h,
			AvailWidth:  screen.w,
			AvailHeight: screen.h - 40,
			ColorDepth:  24,
			PixelDepth:  24,
			PixelRatio:  ratios[g.rng.Intn(len(ratios))],
		},
		WebGL: WebGLFP{
			Vendor:           gpu.vendor,
			Renderer:         gpu.renderer,
			UnmaskedVendor:   gpu.vendor,
			UnmaskedRenderer: gpu.renderer,
			Noise:            g.rng.Float64() * 0.0001,
		},
		Canvas: CanvasFP{
			Noise: g.rng.Float64() * 0.0001,
		},
		Audio: AudioFP{
			Noise:      g.rng.Float64() * 0.0001,
			SampleRate: 48000,
		},
		Fonts: FontFP{
			InstalledFonts: g.generateFonts(),
		},
		Hardware: HardwareFP{},
		Network: NetworkFP{
			WebRTCPolicy:   "disable",
			ConnectionType: g.pick([]string{"wifi", "ethernet"}),
			EffectiveType:  "4g",
			Downlink:       float64(g.rng.Intn(100) + 10),
			RTT:            g.rng.Intn(100) + 20,
		},
		Timezone: TimezoneFP{
			Timezone:       tz.tz,
			TimezoneOffset: tz.offset,
			Locale:         tz.locale,
		},
		Misc: MiscFP{
			Plugins: []Plugin{
				{Name: "PDF Viewer", Filename: "internal-pdf-viewer", Description: "Portable Document Format"},
				{Name: "Chrome PDF Viewer", Filename: "internal-pdf-viewer", Description: "Portable Document Format"},
				{Name: "Chromium PDF Viewer", Filename: "internal-pdf-viewer", Description: "Portable Document Format"},
			},
			MimeTypes: []MimeType{
				{Type: "application/pdf", Suffixes: "pdf"},
			},
			Permissions: map[string]string{
				"geolocation": "prompt", "notifications": "prompt",
			},
		},
	}
}

func (g *Generator) generateUserAgent(opts *GenerateOptions) string {
	edgeVersions := []string{"131.0.2903.86", "131.0.2903.70", "130.0.2849.80", "130.0.2849.68"}
	chromeVersions := []string{"131.0.6778.86", "131.0.6778.70", "130.0.6723.91", "130.0.6723.70"}
	windowsVersions := []string{"10.0; Win64; x64"}

	chromeVer := chromeVersions[g.rng.Intn(len(chromeVersions))]
	edgeVer := edgeVersions[g.rng.Intn(len(edgeVersions))]
	winVer := windowsVersions[g.rng.Intn(len(windowsVersions))]

	if opts.Browser == "chrome" {
		return fmt.Sprintf(
			"Mozilla/5.0 (Windows NT %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
			winVer, chromeVer,
		)
	}

	return fmt.Sprintf(
		"Mozilla/5.0 (Windows NT %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s",
		winVer, chromeVer, edgeVer,
	)
}

func (g *Generator) getLanguagesForLocale(locale string) []string {
	switch locale {
	case "en-US":
		return []string{"en-US", "en"}
	case "en-GB":
		return []string{"en-GB", "en"}
	case "en-AU":
		return []string{"en-AU", "en"}
	case "en-SG":
		return []string{"en-SG", "en"}
	case "fr-FR":
		return []string{"fr-FR", "fr", "en"}
	case "de-DE":
		return []string{"de-DE", "de", "en"}
	case "ja-JP":
		return []string{"ja-JP", "ja", "en"}
	case "zh-CN":
		return []string{"zh-CN", "zh", "en"}
	default:
		return []string{"en-US", "en"}
	}
}

func (g *Generator) generateFonts() []string {
	allFonts := []string{
		"Arial", "Arial Black", "Calibri", "Cambria", "Cambria Math",
		"Comic Sans MS", "Consolas", "Constantia", "Corbel", "Courier New",
		"Georgia", "Impact", "Lucida Console", "Lucida Sans Unicode",
		"Microsoft Sans Serif", "Palatino Linotype", "Segoe UI",
		"Segoe UI Symbol", "Tahoma", "Times New Roman", "Trebuchet MS",
		"Verdana", "Wingdings",
	}

	fonts := make([]string, 0)
	for _, f := range allFonts {
		if g.rng.Float64() > 0.1 {
			fonts = append(fonts, f)
		}
	}
	return fonts
}

func (g *Generator) pick(items []string) string {
	return items[g.rng.Intn(len(items))]
}

func (g *Generator) GenerateRandomName() string {
	adjectives := []string{"Swift", "Shadow", "Ghost", "Stealth", "Cyber", "Neo", "Phantom", "Silent"}
	nouns := []string{"Fox", "Wolf", "Eagle", "Dragon", "Phoenix", "Ninja", "Hawk", "Tiger"}
	return fmt.Sprintf("%s%s%03d",
		adjectives[g.rng.Intn(len(adjectives))],
		nouns[g.rng.Intn(len(nouns))],
		g.rng.Intn(1000))
}
