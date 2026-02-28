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
	DemoMode          bool          // enable demo mode (highlight elements before interaction)
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
