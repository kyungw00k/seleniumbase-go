//go:build integration

package examples

import (
	"testing"
	"time"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestCoupangAccess(t *testing.T) {
	// Stealth mode (headless) is the default — no options needed.
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://www.coupang.com/")
		p.WaitForLoadState("domcontentloaded")
		p.Sleep(3 * time.Second)

		title, err := p.GetTitle()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Page title: %s", title)

		source, err := p.GetPageSource()
		if err != nil {
			t.Fatal(err)
		}

		if len(source) < 1000 {
			t.Error("page source too short — likely blocked by bot detection")
		}

		p.Screenshot("coupang_stealth.png")
		t.Logf("Current URL: %s", p.GetCurrentURL())
	})
}
