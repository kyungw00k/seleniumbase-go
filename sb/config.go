package sb

import "time"

// Timeout constants matching SeleniumBase defaults
const (
	MiniTimeout     = 2 * time.Second
	SmallTimeout    = 7 * time.Second
	LargeTimeout    = 10 * time.Second
	ExtremeTimeout  = 30 * time.Second
	PageLoadTimeout = 120 * time.Second
)

// Config holds browser and page configuration
type Config struct {
	Browser           string        // "chromium", "firefox", "webkit"
	Channel           string        // browser channel: "chrome", "msedge", etc.
	Headless          bool
	Proxy             string        // proxy server URL
	UserAgent         string
	ViewportWidth     int
	ViewportHeight    int
	SlowMo            float64       // milliseconds
	Timeout           time.Duration // default timeout for operations
	Locale            string        // e.g. "en-US"
	IgnoreHTTPSErrors bool
	ColorScheme       string        // "dark", "light", "no-preference"
	DemoMode            bool          // enable demo mode (highlight elements before interaction)
	Stealth             bool          // enable CDP stealth mode (external Chrome launch)
	ChromePath          string        // custom Chrome executable path (auto-detected if empty)
	UserDataDir         string        // custom user data directory (temp dir if empty)
	ExtraArgs           []string      // additional Chrome launch arguments
	ScreenshotOnFailure bool          // auto-screenshot when RunTest fails
	Mobile              bool          // emulate mobile device
	DeviceName          string        // Playwright device name, e.g. "iPhone 13"
	RecordVideo         string        // directory to save video recordings
	RecordHAR           string        // path to save HAR file
	DisableCSP          bool          // disable Content Security Policy
	Incognito           bool          // use incognito context (no persistent state)
	RemoteCDPURL    string // connect to remote browser via CDP URL
	RemoteWSURL     string // connect to remote browser via WebSocket URL
}

func newDefaultConfig() *Config {
	return &Config{
		Browser:        "chromium",
		Headless:       true,
		ViewportWidth:  1280,
		ViewportHeight: 720,
		Timeout:        LargeTimeout,
	}
}
