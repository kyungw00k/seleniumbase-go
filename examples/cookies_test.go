//go:build integration

package examples

import (
	"os"
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestCookiePersistence(t *testing.T) {
	cookiePath := "test_cookies.json"
	defer os.Remove(cookiePath)

	// First session: login and save cookies
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://www.saucedemo.com")
		p.WaitForElement("div.login_logo")
		p.Type("#user-name", "standard_user")
		p.Type("#password", "secret_sauce")
		p.Click(`input[type="submit"]`)
		p.AssertElement("div.inventory_list")
		p.Highlight("div.inventory_list")
		if err := p.SaveCookies(cookiePath); err != nil {
			t.Fatal(err)
		}
	}, sb.WithHeadless(true))

	// Second session: load cookies to bypass login
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://www.saucedemo.com")
		if err := p.LoadCookies(cookiePath); err != nil {
			t.Fatal(err)
		}
		p.Open("https://www.saucedemo.com/inventory.html")
		p.AssertElement("div.inventory_list")
		p.Highlight("div.inventory_list")
	}, sb.WithHeadless(true))
}
