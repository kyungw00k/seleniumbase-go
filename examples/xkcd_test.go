//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestXkcd(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://xkcd.com/353/")
		p.AssertTitle("xkcd: Python")
		p.AssertElement(`img[alt="Python"]`)
		p.Click(`a[rel="license"]`)
		p.AssertText("free to copy and reuse", "body")
		p.GoBack()
		p.Click("link=About")
		p.AssertExactText("xkcd.com", "h2")
	}, sb.WithHeadless(true))
}
