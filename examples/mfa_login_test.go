//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestMFALogin(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/realworld/login")
		p.Type("#username", "demo_user")
		p.Type("#password", "secret_pass")
		p.EnterMFACode("#totpcode", "GAXG2MTEOR3DMMDG")
		p.AssertText("Welcome!", "h1")
		p.Highlight("img#image1")
		p.Click(`a:contains("This Page")`)
		p.Click("link=Sign out")
		p.AssertElement(`a:contains("Sign in")`)
		p.AssertExactText("You have been signed out!", "#top_message")
	}, sb.WithHeadless(true))
}
