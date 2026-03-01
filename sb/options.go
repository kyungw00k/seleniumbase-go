package sb

import "time"

// Option configures an SB instance.
type Option func(*Config)

func WithBrowser(browser string) Option {
	return func(c *Config) { c.Browser = browser }
}

// WithChannel sets the browser channel. Use "chrome" to use system Chrome,
// "msedge" for Edge, etc. When set, Playwright uses the system-installed
// browser instead of its bundled one.
func WithChannel(channel string) Option {
	return func(c *Config) { c.Channel = channel }
}

func WithHeadless(headless bool) Option {
	return func(c *Config) { c.Headless = headless }
}

func WithProxy(server string) Option {
	return func(c *Config) { c.Proxy = server }
}

func WithUserAgent(agent string) Option {
	return func(c *Config) { c.UserAgent = agent }
}

func WithViewportSize(width, height int) Option {
	return func(c *Config) {
		c.ViewportWidth = width
		c.ViewportHeight = height
	}
}

func WithSlowMo(ms float64) Option {
	return func(c *Config) { c.SlowMo = ms }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

func WithLocale(code string) Option {
	return func(c *Config) { c.Locale = code }
}

func WithIgnoreHTTPSErrors(ignore bool) Option {
	return func(c *Config) { c.IgnoreHTTPSErrors = ignore }
}

func WithColorScheme(scheme string) Option {
	return func(c *Config) { c.ColorScheme = scheme }
}

func WithDemoMode(enabled bool) Option {
	return func(c *Config) { c.DemoMode = enabled }
}

func WithStealth(enabled bool) Option {
	return func(c *Config) { c.Stealth = enabled }
}

func WithChromePath(path string) Option {
	return func(c *Config) { c.ChromePath = path }
}

func WithUserDataDir(dir string) Option {
	return func(c *Config) { c.UserDataDir = dir }
}

func WithExtraArgs(args ...string) Option {
	return func(c *Config) { c.ExtraArgs = args }
}

func WithScreenshotOnFailure(enabled bool) Option {
	return func(c *Config) { c.ScreenshotOnFailure = enabled }
}

func WithMobile(mobile bool) Option {
	return func(c *Config) { c.Mobile = mobile }
}

func WithDevice(name string) Option {
	return func(c *Config) { c.DeviceName = name; c.Mobile = true }
}

func WithRecordVideo(dir string) Option {
	return func(c *Config) { c.RecordVideo = dir }
}

func WithRecordHAR(path string) Option {
	return func(c *Config) { c.RecordHAR = path }
}

func WithDisableCSP(disable bool) Option {
	return func(c *Config) { c.DisableCSP = disable }
}

// WithIncognito enables incognito mode. Note: Playwright browser contexts are
// already isolated by default (no shared cookies/storage between contexts).
// This option is provided for API compatibility with SeleniumBase Python.
func WithIncognito(incognito bool) Option {
	return func(c *Config) { c.Incognito = incognito }
}

// WithRemoteCDPURL connects to a remote browser via Chrome DevTools Protocol.
// Use this for BrowserStack, LambdaTest, or self-hosted Chrome instances.
func WithRemoteCDPURL(url string) Option {
	return func(c *Config) { c.RemoteCDPURL = url }
}

// WithRemoteWSURL connects to a remote Playwright server via WebSocket.
func WithRemoteWSURL(url string) Option {
	return func(c *Config) { c.RemoteWSURL = url }
}
