//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestTabSwitching(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		// Set content in first tab
		p.SetContent("<html><body><h1>Page A</h1></body></html>")
		p.AssertText("Page A", "h1")

		// Open new tab and set different content
		tab2, err := p.OpenNewTab()
		if err != nil {
			t.Fatal(err)
		}
		tab2.SetContent("<html><body><h1>Page B</h1></body></html>")
		tab2.AssertText("Page B", "h1")

		// Switch back to first tab and verify
		p.BringToFront()
		p.AssertText("Page A", "h1")

		// Switch to second tab and verify
		tab2.BringToFront()
		tab2.AssertText("Page B", "h1")
	}, sb.WithHeadless(true))
}
