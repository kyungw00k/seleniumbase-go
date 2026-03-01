//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestCalculator(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/apps/calculator")
		p.Click(`button[id="6"]`)
		p.Click("button#multiply")
		p.Click(`button[id="7"]`)
		p.Click("button#add")
		p.Click(`button[id="1"]`)
		p.Click(`button[id="2"]`)
		p.Click("button#equal")
		p.AssertExactText("54", "input#output")
	}, sb.WithHeadless(true))
}
