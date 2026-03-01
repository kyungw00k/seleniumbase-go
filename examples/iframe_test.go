//go:build integration

package examples

import (
	"strings"
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestIframeBasics(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/w3schools/iframes.html")
		p.AssertTitle("iframe Testing")
		p.Click("button#runbtn")

		// Use FrameLocator to interact with iframe content
		frame := p.FrameLocator("#iframeResult")

		// Assert text inside iframe
		textContent, err := frame.Locator("h2").TextContent()
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(textContent, "HTML Iframes") {
			t.Errorf("expected text containing 'HTML Iframes', got %q", textContent)
		}

		// Assert nested iframe content
		innerFrame := frame.FrameLocator(`[title*="Iframe"]`)
		innerText, err := innerFrame.Locator("h1").TextContent()
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(innerText, "This page is displayed in an iframe") {
			t.Errorf("expected text containing 'This page is displayed in an iframe', got %q", innerText)
		}
	}, sb.WithHeadless(true))
}
