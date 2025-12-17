package fingerprint

// Fingerprint represents a complete browser fingerprint
type Fingerprint struct {
	Navigator NavigatorFP `json:"navigator"`
	Screen    ScreenFP    `json:"screen"`
	WebGL     WebGLFP     `json:"webgl"`
	Canvas    CanvasFP    `json:"canvas"`
	Audio     AudioFP     `json:"audio"`
	Fonts     FontFP      `json:"fonts"`
	Hardware  HardwareFP  `json:"hardware"`
	Network   NetworkFP   `json:"network"`
	Timezone  TimezoneFP  `json:"timezone"`
	Misc      MiscFP      `json:"misc"`
}

type NavigatorFP struct {
	UserAgent           string   `json:"userAgent"`
	AppVersion          string   `json:"appVersion"`
	Platform            string   `json:"platform"`
	Vendor              string   `json:"vendor"`
	Language            string   `json:"language"`
	Languages           []string `json:"languages"`
	HardwareConcurrency int      `json:"hardwareConcurrency"`
	DeviceMemory        int      `json:"deviceMemory"`
	MaxTouchPoints      int      `json:"maxTouchPoints"`
	ProductSub          string   `json:"productSub"`
	DoNotTrack          string   `json:"doNotTrack"`
	CookieEnabled       bool     `json:"cookieEnabled"`
	Webdriver           bool     `json:"webdriver"`
}

type ScreenFP struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	AvailWidth  int     `json:"availWidth"`
	AvailHeight int     `json:"availHeight"`
	ColorDepth  int     `json:"colorDepth"`
	PixelDepth  int     `json:"pixelDepth"`
	PixelRatio  float64 `json:"pixelRatio"`
}

type WebGLFP struct {
	Vendor           string  `json:"vendor"`
	Renderer         string  `json:"renderer"`
	UnmaskedVendor   string  `json:"unmaskedVendor"`
	UnmaskedRenderer string  `json:"unmaskedRenderer"`
	Noise            float64 `json:"noise"`
}

type CanvasFP struct {
	Noise float64 `json:"noise"`
}

type AudioFP struct {
	Noise      float64 `json:"noise"`
	SampleRate int     `json:"sampleRate"`
}

type FontFP struct {
	InstalledFonts []string `json:"installedFonts"`
}

type HardwareFP struct {
	BatteryCharging *bool    `json:"batteryCharging,omitempty"`
	BatteryLevel    *float64 `json:"batteryLevel,omitempty"`
}

type NetworkFP struct {
	WebRTCPolicy   string   `json:"webRTCPolicy"`
	PublicIP       string   `json:"publicIP,omitempty"`
	LocalIPs       []string `json:"localIPs,omitempty"`
	ConnectionType string   `json:"connectionType"`
	EffectiveType  string   `json:"effectiveType"`
	Downlink       float64  `json:"downlink"`
	RTT            int      `json:"rtt"`
}

type TimezoneFP struct {
	Timezone       string `json:"timezone"`
	TimezoneOffset int    `json:"timezoneOffset"`
	Locale         string `json:"locale"`
}

type MiscFP struct {
	Plugins     []Plugin          `json:"plugins"`
	MimeTypes   []MimeType        `json:"mimeTypes"`
	Permissions map[string]string `json:"permissions"`
}

type Plugin struct {
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
}

type MimeType struct {
	Type     string `json:"type"`
	Suffixes string `json:"suffixes"`
}
