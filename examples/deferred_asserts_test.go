//go:build integration

package examples

import (
	"strings"
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestDeferredAsserts(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://xkcd.com/993/")
		p.WaitForElement("#comic")

		// These should pass
		p.DeferredAssertElement(`img[alt="Brand Identity"]`)
		p.DeferredAssertText("Random", "ul.comicNav")
		p.DeferredAssertExactText("Brand Identity", "#ctitle")

		// These should fail
		p.DeferredAssertElement(`img[alt="Rocket Ship"]`)       // Will fail
		p.DeferredAssertText("Fake Item", "ul.comicNav")        // Will fail
		p.DeferredAssertElement(`a[name="Super Fake !!!"]`)     // Will fail
		p.DeferredAssertExactText("Fake Food", "#comic")        // Will fail

		// Process should return error with multiple failures
		err := p.ProcessDeferredAsserts()
		if err == nil {
			t.Error("expected deferred asserts to produce errors")
		}
		// Should contain multiple failure messages
		errMsg := err.Error()
		if !strings.Contains(errMsg, "deferred") || !strings.Contains(errMsg, "4") {
			t.Logf("Deferred assert error (expected ~4 failures): %v", err)
		}
	}, sb.WithHeadless(true))
}
