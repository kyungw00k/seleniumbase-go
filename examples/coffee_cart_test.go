//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestCoffeeCart(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/coffee/")
		p.AssertTitle("Coffee Cart")
		p.AssertElement(`button:contains("Total: $0.00")`)
		p.Click(`div[data-sb="Cappuccino"]`)
		p.AssertExactText("cart (1)", `a[aria-label="Cart page"]`)
		p.Click(`div[data-sb="Flat-White"]`)
		p.AssertExactText("cart (2)", `a[aria-label="Cart page"]`)
		p.Click(`div[data-sb="Cafe-Latte"]`)
		p.AssertExactText("cart (3)", `a[aria-label="Cart page"]`)
		p.Click(`a[aria-label="Cart page"]`)
		p.AssertExactText("Total: $53.00", "button.pay")
		p.Click("button.pay")
		p.Type("input#name", "Selenium Coffee")
		p.Type("input#email", "test@test.test")
		p.Click("button#submit-payment")
		p.AssertText("Thanks for your purchase.", "#app .success")
	}, sb.WithHeadless(true))
}
