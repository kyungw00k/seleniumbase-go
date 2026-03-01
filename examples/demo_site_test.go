//go:build integration

package examples

import (
	"testing"

	"github.com/kyungw00k/seleniumbase-go/sb"
)

func TestDemoSite(t *testing.T) {
	sb.RunTest(t, func(p *sb.Page) {
		// Open the demo page
		p.Open("https://seleniumbase.io/demo_page")

		// Assert the title
		p.AssertTitle("Web Testing Page")

		// Assert element visible
		p.AssertElement("tbody#tbodyId")

		// Assert text in element
		p.AssertText("Demo Page", "h1")

		// Type into text fields and assert
		p.Type("#myTextInput", "This is Automated")
		p.Type("textarea.area1", "Testing Time!\n")
		p.Type(`[name="preText2"]`, "Typing Text!")
		p.AssertText("This is Automated", "#myTextInput")

		// Click a button and verify results
		p.AssertText("This Text is Green", "#pText")
		p.Click(`button:contains("Click Me")`)
		p.AssertText("This Text is Purple", "#pText")

		// Assert SVG element
		p.AssertElement(`svg[name="svgName"]`)

		// Verify checkbox
		p.AssertElementNotVisible("img#logo")
		checked, _ := p.IsChecked("#checkBox1")
		if checked {
			t.Error("checkbox1 should not be checked initially")
		}
		p.Click("#checkBox1")
		checked, _ = p.IsChecked("#checkBox1")
		if !checked {
			t.Error("checkbox1 should be checked after click")
		}
		p.AssertElement("img#logo")

		// Verify drag and drop
		p.AssertElementNotVisible("div#drop2 img#logo")
		p.DragAndDrop("img#logo", "div#drop2")
		p.AssertElement("div#drop2 img#logo")

		// Assert exact text
		p.AssertExactText("Demo Page", "h1")

		// Highlight
		p.Highlight("h2")
	}, sb.WithHeadless(true))
}
