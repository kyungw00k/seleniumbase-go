package sb

import (
	"fmt"
	"testing"

	"github.com/playwright-community/playwright-go"
)

// SB manages playwright instance, browser, and context lifecycle.
type SB struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	config  *Config
	process *stealthProcess
}

func newSB(opts ...Option) (*SB, error) {
	cfg := newDefaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	pw, err := playwright.Run()
	if err != nil {
		// Playwright driver not found — attempt auto-install
		channel, installErr := EnsureBrowser()
		if installErr != nil {
			return nil, fmt.Errorf("sb: could not start playwright: %w (auto-install also failed: %v)", err, installErr)
		}
		// Set channel if system Chrome was found and no explicit browser was set
		if channel != "" && cfg.Browser == "chromium" {
			cfg.Channel = channel
		}
		// Retry
		pw, err = playwright.Run()
		if err != nil {
			return nil, fmt.Errorf("sb: could not start playwright after install: %w", err)
		}
	}

	if cfg.Stealth {
		return newStealthSB(pw, cfg)
	}

	if cfg.RemoteCDPURL != "" {
		return newRemoteCDPSB(pw, cfg)
	}
	if cfg.RemoteWSURL != "" {
		return newRemoteWSSB(pw, cfg)
	}

	var browserType playwright.BrowserType
	switch cfg.Browser {
	case "firefox":
		browserType = pw.Firefox
	case "webkit":
		browserType = pw.WebKit
	default:
		browserType = pw.Chromium
	}

	launchOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(cfg.Headless),
	}
	if cfg.SlowMo > 0 {
		launchOpts.SlowMo = playwright.Float(cfg.SlowMo)
	}
	if cfg.Proxy != "" {
		launchOpts.Proxy = &playwright.Proxy{Server: cfg.Proxy}
	}
	if cfg.Channel != "" {
		launchOpts.Channel = playwright.String(cfg.Channel)
	}

	browser, err := browserType.Launch(launchOpts)
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("sb: could not launch browser: %w", err)
	}

	contextOpts := buildContextOpts(pw, cfg)

	context, err := browser.NewContext(contextOpts)
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, fmt.Errorf("sb: could not create context: %w", err)
	}

	return &SB{
		pw:      pw,
		browser: browser,
		context: context,
		config:  cfg,
	}, nil
}

func (s *SB) newPage() (*Page, error) {
	page, err := s.context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("sb: could not create page: %w", err)
	}
	if s.config.Timeout > 0 {
		page.SetDefaultTimeout(float64(s.config.Timeout.Milliseconds()))
	}
	if s.config.Stealth {
		if err := applyStealthCDP(s.context, page, s.config); err != nil {
			page.Close()
			return nil, err
		}
	}
	return newPage(page, s.context, s.config), nil
}

func (s *SB) close() {
	if s.context != nil {
		s.context.Close()
	}
	if s.browser != nil {
		s.browser.Close()
	}
	if s.process != nil {
		s.process.stop()
	}
	if s.pw != nil {
		s.pw.Stop()
	}
}

// Run creates an SB instance, a page, calls fn, and cleans up.
func Run(fn func(p *Page) error, opts ...Option) error {
	s, err := newSB(opts...)
	if err != nil {
		return err
	}
	defer s.close()

	page, err := s.newPage()
	if err != nil {
		return err
	}

	return fn(page)
}

// RunTest creates an SB instance integrated with testing.T.
func RunTest(t *testing.T, fn func(p *Page), opts ...Option) {
	t.Helper()

	// Check if ScreenshotOnFailure is requested
	cfg := newDefaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	err := Run(func(p *Page) error {
		fn(p)
		if cfg.ScreenshotOnFailure && t.Failed() {
			_ = p.Screenshot(fmt.Sprintf("failure_%s.png", t.Name()))
		}
		return nil
	}, opts...)
	if err != nil {
		t.Fatal(err)
	}
}

// NewPage creates a standalone Page with its own browser instance.
// Returns the page and a cleanup function that must be called (typically via defer).
func NewPage(opts ...Option) (*Page, func(), error) {
	s, err := newSB(opts...)
	if err != nil {
		return nil, nil, err
	}

	page, err := s.newPage()
	if err != nil {
		s.close()
		return nil, nil, err
	}

	cleanup := func() {
		s.close()
	}

	return page, cleanup, nil
}

// buildContextOpts creates BrowserNewContextOptions from the config.
func buildContextOpts(pw *playwright.Playwright, cfg *Config) playwright.BrowserNewContextOptions {
	contextOpts := playwright.BrowserNewContextOptions{}
	if cfg.ViewportWidth > 0 && cfg.ViewportHeight > 0 {
		contextOpts.Viewport = &playwright.Size{
			Width:  cfg.ViewportWidth,
			Height: cfg.ViewportHeight,
		}
	}
	if cfg.UserAgent != "" {
		contextOpts.UserAgent = playwright.String(cfg.UserAgent)
	}
	if cfg.Locale != "" {
		contextOpts.Locale = playwright.String(cfg.Locale)
	}
	if cfg.IgnoreHTTPSErrors {
		contextOpts.IgnoreHttpsErrors = playwright.Bool(true)
	}
	if cfg.ColorScheme != "" {
		switch cfg.ColorScheme {
		case "dark":
			contextOpts.ColorScheme = playwright.ColorSchemeDark
		case "light":
			contextOpts.ColorScheme = playwright.ColorSchemeLight
		case "no-preference":
			contextOpts.ColorScheme = playwright.ColorSchemeNoPreference
		}
	}
	if cfg.Mobile {
		contextOpts.IsMobile = playwright.Bool(true)
		contextOpts.HasTouch = playwright.Bool(true)
	}
	if cfg.DeviceName != "" {
		device := pw.Devices[cfg.DeviceName]
		if device != nil {
			contextOpts.UserAgent = playwright.String(device.UserAgent)
			contextOpts.Viewport = device.Viewport
			contextOpts.DeviceScaleFactor = playwright.Float(device.DeviceScaleFactor)
			contextOpts.IsMobile = playwright.Bool(device.IsMobile)
			contextOpts.HasTouch = playwright.Bool(device.HasTouch)
		}
	}
	if cfg.RecordVideo != "" {
		contextOpts.RecordVideo = &playwright.RecordVideo{Dir: cfg.RecordVideo}
	}
	if cfg.RecordHAR != "" {
		contextOpts.RecordHarPath = playwright.String(cfg.RecordHAR)
	}
	if cfg.DisableCSP {
		contextOpts.BypassCSP = playwright.Bool(true)
	}
	return contextOpts
}
