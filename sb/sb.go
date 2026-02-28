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
}

func newSB(opts ...Option) (*SB, error) {
	cfg := newDefaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("sb: could not start playwright: %w", err)
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

	browser, err := browserType.Launch(launchOpts)
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("sb: could not launch browser: %w", err)
	}

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
	return newPage(page, s.context, s.config), nil
}

func (s *SB) close() {
	if s.context != nil {
		s.context.Close()
	}
	if s.browser != nil {
		s.browser.Close()
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
	err := Run(func(p *Page) error {
		fn(p)
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
