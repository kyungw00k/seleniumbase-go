package sb

import (
	"sync"

	"github.com/kyungw00k/seleniumbase-go/selector"
	"github.com/playwright-community/playwright-go"
)

// ConsoleMessage represents a captured browser console message.
type ConsoleMessage struct {
	Type string // "log", "error", "warning", "info"
	Text string
}

// Page wraps a playwright.Page with SeleniumBase-style convenience methods.
type Page struct {
	pw               playwright.Page
	context          playwright.BrowserContext
	config           *Config
	expect           playwright.PlaywrightAssertions
	deferredFailures []error
	recording        bool
	consoleMu        sync.Mutex
	consoleMessages  []ConsoleMessage
}

func newPage(pw playwright.Page, ctx playwright.BrowserContext, cfg *Config) *Page {
	return &Page{
		pw:      pw,
		context: ctx,
		config:  cfg,
		expect:  playwright.NewPlaywrightAssertions(float64(cfg.Timeout.Milliseconds())),
	}
}

// Playwright returns the underlying playwright.Page for direct access.
func (p *Page) Playwright() playwright.Page {
	return p.pw
}

// Locator returns a playwright.Locator for the given SeleniumBase-style selector.
func (p *Page) Locator(sel string) playwright.Locator {
	return p.pw.Locator(selector.Parse(sel))
}

// Context returns the underlying playwright.BrowserContext.
func (p *Page) Context() playwright.BrowserContext {
	return p.context
}

// locator is an internal helper that creates a locator from a SeleniumBase selector.
func (p *Page) locator(sel string) playwright.Locator {
	return p.pw.Locator(selector.Parse(sel))
}
