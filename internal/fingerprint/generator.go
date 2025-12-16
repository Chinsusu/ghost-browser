package fingerprint

import (
	"fmt"
	"math/rand"
	"time"
)

type Generator struct {
	rng *rand.Rand
}

type GenerateOptions struct {
	OS      string
	Browser string
}

func NewGenerator() *Generator {
	return &Generator{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (g *Generator) Generate(opts *GenerateOptions) *Fingerprint {
	if opts == nil {
		opts = &GenerateOptions{OS: "windows", Browser: "edge"}
	}

	return &Fingerprint{
		Navigator: g.genNavigator(opts),
		Screen:    g.genScreen(),
		WebGL:     g.genWebGL(),
		Canvas:    CanvasFP{Noise: g.rng.Float64() * 0.0001},
		Audio:     AudioFP{Noise: g.rng.Float64() * 0.0001, SampleRate: 48000},
		Fonts:     g.genFonts(),
		Hardware:  HardwareFP{},
		Network:   g.genNetwork(),
		Timezone:  g.genTimezone(),
		Misc:      g.genMisc(),
	}
}

func (g *Generator) genNavigator(opts *GenerateOptions) NavigatorFP {
	ua := g.genUserAgent(opts)
	langs := [][]string{
		{"en-US", "en"}, {"vi-VN", "vi", "en"}, {"ja-JP", "ja"}, {"ko-KR", "ko"},
	}
	lang := langs[g.rng.Intn(len(langs))]
	cores := []int{4, 6, 8, 12, 16}
	mem := []int{4, 8, 16, 32}

	return NavigatorFP{
		UserAgent:           ua,
		AppVersion:          ua[8:],
		Platform:            "Win32",
		Vendor:              "Google Inc.",
		Language:            lang[0],
		Languages:           lang,
		HardwareConcurrency: cores[g.rng.Intn(len(cores))],
		DeviceMemory:        mem[g.rng.Intn(len(mem))],
		MaxTouchPoints:      0,
		ProductSub:          "20030107",
		DoNotTrack:          "1",
		CookieEnabled:       true,
		Webdriver:           false,
	}
}

func (g *Generator) genUserAgent(opts *GenerateOptions) string {
	edgeVer := []string{"131.0.2903.86", "130.0.2849.80", "129.0.2792.89"}
	chromeVer := []string{"131.0.6778.86", "130.0.6723.91", "129.0.6668.90"}
	winVer := []string{"10.0; Win64; x64"}

	if opts.Browser == "edge" || opts.Browser == "" {
		return fmt.Sprintf(
			"Mozilla/5.0 (Windows NT %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s",
			winVer[g.rng.Intn(len(winVer))],
			chromeVer[g.rng.Intn(len(chromeVer))],
			edgeVer[g.rng.Intn(len(edgeVer))],
		)
	}
	return fmt.Sprintf(
		"Mozilla/5.0 (Windows NT %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
		winVer[g.rng.Intn(len(winVer))],
		chromeVer[g.rng.Intn(len(chromeVer))],
	)
}

func (g *Generator) genScreen() ScreenFP {
	res := []struct{ w, h int }{
		{1920, 1080}, {2560, 1440}, {1366, 768}, {3840, 2160},
	}
	r := res[g.rng.Intn(len(res))]
	ratios := []float64{1.0, 1.25, 1.5, 2.0}

	return ScreenFP{
		Width: r.w, Height: r.h,
		AvailWidth: r.w, AvailHeight: r.h - 48,
		ColorDepth: 24, PixelDepth: 24,
		PixelRatio: ratios[g.rng.Intn(len(ratios))],
	}
}

func (g *Generator) genWebGL() WebGLFP {
	gpus := []struct{ vendor, renderer string }{
		{"NVIDIA Corporation", "NVIDIA GeForce RTX 4090"},
		{"NVIDIA Corporation", "NVIDIA GeForce RTX 3080"},
		{"Intel Inc.", "Intel(R) UHD Graphics 630"},
		{"AMD", "AMD Radeon RX 6800 XT"},
	}
	gpu := gpus[g.rng.Intn(len(gpus))]

	return WebGLFP{
		Vendor:           "Google Inc. (" + gpu.vendor + ")",
		Renderer:         "ANGLE (" + gpu.renderer + " Direct3D11)",
		UnmaskedVendor:   gpu.vendor,
		UnmaskedRenderer: gpu.renderer,
		Noise:            g.rng.Float64() * 0.0001,
	}
}

func (g *Generator) genFonts() FontFP {
	fonts := []string{
		"Arial", "Calibri", "Cambria", "Comic Sans MS", "Consolas",
		"Courier New", "Georgia", "Impact", "Segoe UI", "Tahoma",
		"Times New Roman", "Trebuchet MS", "Verdana",
	}
	selected := make([]string, 0)
	for _, f := range fonts {
		if g.rng.Float64() > 0.1 {
			selected = append(selected, f)
		}
	}
	return FontFP{InstalledFonts: selected}
}

func (g *Generator) genNetwork() NetworkFP {
	return NetworkFP{
		WebRTCPolicy:   "disable",
		ConnectionType: "wifi",
		EffectiveType:  "4g",
		Downlink:       float64(g.rng.Intn(100) + 10),
		RTT:            g.rng.Intn(100) + 20,
	}
}

func (g *Generator) genTimezone() TimezoneFP {
	tzs := []struct {
		tz     string
		offset int
		locale string
	}{
		{"America/New_York", -300, "en-US"},
		{"Asia/Ho_Chi_Minh", 420, "vi-VN"},
		{"Asia/Tokyo", 540, "ja-JP"},
		{"Europe/London", 0, "en-GB"},
	}
	tz := tzs[g.rng.Intn(len(tzs))]
	return TimezoneFP{Timezone: tz.tz, TimezoneOffset: tz.offset, Locale: tz.locale}
}

func (g *Generator) genMisc() MiscFP {
	return MiscFP{
		Plugins: []Plugin{
			{Name: "PDF Viewer", Filename: "internal-pdf-viewer", Description: "PDF"},
			{Name: "Chrome PDF Viewer", Filename: "internal-pdf-viewer", Description: "PDF"},
		},
		MimeTypes: []MimeType{
			{Type: "application/pdf", Suffixes: "pdf"},
		},
		Permissions: map[string]string{
			"geolocation": "prompt", "notifications": "prompt",
		},
	}
}

func (g *Generator) GenerateRandomName() string {
	adj := []string{"Swift", "Shadow", "Ghost", "Stealth", "Cyber", "Neo"}
	noun := []string{"Fox", "Wolf", "Eagle", "Dragon", "Phoenix", "Ninja"}
	return fmt.Sprintf("%s%s%03d",
		adj[g.rng.Intn(len(adj))],
		noun[g.rng.Intn(len(noun))],
		g.rng.Intn(1000))
}
