//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestURLAsserts(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/help_docs/how_it_works/")
		p.AssertURL("https://seleniumbase.io/help_docs/how_it_works/")
		p.AssertTitleContains("How it Works")
		p.Click(`a:contains("Coffee Cart")`)
		p.AssertURLContains("/coffee")
		p.AssertTitle("Coffee Cart")
	}, sb.WithHeadless(true))
}
