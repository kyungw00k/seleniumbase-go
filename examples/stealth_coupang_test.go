//go:build integration

package examples

import (
	"testing"
	"time"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestCoupangStealthAccess(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		// Access Coupang product page
		p.Open("https://www.coupang.com/")

		// Wait for page to load
		p.WaitForLoadState("domcontentloaded")
		p.Sleep(3 * time.Second)

		// Check that we're not blocked — Coupang typically shows homepage elements
		// If blocked, the page would show CAPTCHA or empty content
		title, err := p.GetTitle()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Page title: %s", title)

		// Verify we got actual page content (not a bot detection page)
		source, err := p.GetPageSource()
		if err != nil {
			t.Fatal(err)
		}

		// Check for page content indicators
		if len(source) < 1000 {
			t.Error("page source too short — likely blocked by bot detection")
		}

		// Take a screenshot for visual verification
		p.Screenshot("coupang_stealth.png")
		t.Log("Screenshot saved to coupang_stealth.png")

		// Check page URL — if redirected to a CAPTCHA page, the URL would change
		currentURL := p.GetCurrentURL()
		t.Logf("Current URL: %s", currentURL)

	}, sb.WithStealth(true), sb.WithHeadless(true))
}

func TestCoupangNormalBlocked(t *testing.T) {
	// This test demonstrates that normal Playwright access gets blocked by Coupang.
	// It's expected to encounter bot detection.
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://www.coupang.com/")
		p.WaitForLoadState("domcontentloaded")
		p.Sleep(3 * time.Second)

		title, err := p.GetTitle()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Normal mode title: %s", title)

		// Take screenshot to compare with stealth mode
		p.Screenshot("coupang_normal.png")
		t.Log("Screenshot saved to coupang_normal.png — compare with stealth version")

		currentURL := p.GetCurrentURL()
		t.Logf("Normal mode URL: %s", currentURL)
	}, sb.WithHeadless(true))
}
