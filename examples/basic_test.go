//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestBasicDemo(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://www.saucedemo.com")
		p.Type("#user-name", "standard_user")
		p.Type("#password", "secret_sauce")
		p.Click("#login-button")
		p.AssertElement("div.inventory_list")
		p.AssertExactText("Products", "span.title")

		// Add backpack to cart
		p.Click(`button[name*="backpack"]`)

		// Navigate to cart
		p.Click("a.shopping_cart_link")
		p.AssertElement("div.cart_item")

		// Remove item from cart
		p.Click(`button:has-text("Remove")`)
		p.AssertElementAbsent("div.cart_item")

		// Logout
		p.Click("#react-burger-menu-btn")
		p.Click("#logout_sidebar_link")
		p.AssertElement("#login-button")
	}, sb.WithHeadless(true))
}
