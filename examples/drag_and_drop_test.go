//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestDragAndDrop(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		p.Open("https://seleniumbase.io/other/drag_and_drop")
		p.AssertElementNotVisible("#div1 img#drag1")
		p.DragAndDrop("#drag1", "#div1")
		p.AssertElement("#div1 img#drag1")
	}, sb.WithHeadless(true))
}
