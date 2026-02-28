package sb

import "time"

// Option configures an SB instance.
type Option func(*Config)

func WithBrowser(browser string) Option {
	return func(c *Config) { c.Browser = browser }
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
