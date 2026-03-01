//go:build integration

package examples

import (
	"testing"
	"time"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestAlertHandling(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("about:blank")

		// Test accepting an alert
		p.AcceptDialogs()
		_, err := p.Evaluate(`window.alert("ALERT!!!")`)
		if err != nil {
			t.Fatal(err)
		}
		p.Sleep(500 * time.Millisecond)

		// Test dismissing a prompt
		p.DismissDialogs()
		_, err = p.Evaluate(`window.prompt("My Prompt", "defaultText")`)
		if err != nil {
			t.Fatal(err)
		}
		p.Sleep(500 * time.Millisecond)
	}, sb.WithHeadless(true))
}
