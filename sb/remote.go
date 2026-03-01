package sb

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

func newRemoteCDPSB(pw *playwright.Playwright, cfg *Config) (*SB, error) {
	browser, err := pw.Chromium.ConnectOverCDP(cfg.RemoteCDPURL)
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("sb: could not connect via CDP: %w", err)
	}

	var ctx playwright.BrowserContext
	if contexts := browser.Contexts(); len(contexts) > 0 {
		ctx = contexts[0]
	} else {
		ctx, err = browser.NewContext(buildContextOpts(pw, cfg))
		if err != nil {
			browser.Close()
			pw.Stop()
			return nil, fmt.Errorf("sb: could not create context: %w", err)
		}
	}

	return &SB{pw: pw, browser: browser, context: ctx, config: cfg}, nil
}

func newRemoteWSSB(pw *playwright.Playwright, cfg *Config) (*SB, error) {
	var browserType playwright.BrowserType
	switch cfg.Browser {
	case "firefox":
		browserType = pw.Firefox
	case "webkit":
		browserType = pw.WebKit
	default:
		browserType = pw.Chromium
	}

	browser, err := browserType.Connect(cfg.RemoteWSURL)
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("sb: could not connect via WebSocket: %w", err)
	}

	ctx, err := browser.NewContext(buildContextOpts(pw, cfg))
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, fmt.Errorf("sb: could not create context: %w", err)
	}

	return &SB{pw: pw, browser: browser, context: ctx, config: cfg}, nil
}
