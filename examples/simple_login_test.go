//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestSimpleLogin(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/simple/login")
		p.Type("#username", "demo_user")
		p.Type("#password", "secret_pass")
		p.Click(`a:contains("Sign in")`)
		p.AssertExactText("Welcome!", "h1")
		p.AssertElement("img#image1")
		p.Highlight("img#image1")
		p.Click("link=Sign out")
		p.AssertText("signed out", "#top_message")
	}, sb.WithHeadless(true))
}
