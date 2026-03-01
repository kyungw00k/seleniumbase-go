package sb

import (
	"fmt"
	"strings"
)

// GenerateGoCode transforms recorded browser actions into a Go function body.
func GenerateGoCode(actions []RecordedAction) string {
	var b strings.Builder
	b.WriteString("func recordedTest(p *sb.Page) error {\n")

	for i := 0; i < len(actions); i++ {
		a := actions[i]
		switch a.Action {
		case "begin":
			b.WriteString(fmt.Sprintf("\tp.Open(%q)\n", a.Data))
		case "_url_":
			b.WriteString(fmt.Sprintf("\tp.Open(%q)\n", a.Data))
		case "click":
			// Lookahead: if next action is dbclk on the same selector (or empty selector), skip this click
			if i+1 < len(actions) && actions[i+1].Action == "dbclk" && (actions[i+1].Selector == a.Selector || actions[i+1].Selector == "") {
				continue
			}
			b.WriteString(fmt.Sprintf("\tp.Click(%q)\n", a.Selector))
		case "dbclk":
			sel := a.Selector
			// dbclk from Python recorder may have empty selector — use previous click's selector
			if sel == "" && i > 0 && actions[i-1].Action == "click" {
				sel = actions[i-1].Selector
			}
			if sel != "" {
				b.WriteString(fmt.Sprintf("\tp.DoubleClick(%q)\n", sel))
			}
		case "input":
			text := a.Data
			text = strings.TrimSuffix(text, "\n")
			b.WriteString(fmt.Sprintf("\tp.Type(%q, %q)\n", a.Selector, text))
		case "s_opt":
			b.WriteString(fmt.Sprintf("\tp.SelectOption(%q, %q)\n", a.Selector, a.Data))
		case "set_v":
			b.WriteString(fmt.Sprintf("\tp.Click(%q)\n", a.Selector))
		case "submi":
			b.WriteString(fmt.Sprintf("\tp.Click(%q)\n", a.Selector))
		case "as_te":
			b.WriteString(fmt.Sprintf("\tp.AssertText(%q, %q)\n", a.Data, a.Selector))
		case "as_el":
			b.WriteString(fmt.Sprintf("\tp.AssertElement(%q)\n", a.Selector))
		case "as_ep":
			b.WriteString(fmt.Sprintf("\tp.AssertElementPresent(%q)\n", a.Selector))
		case "as_ev":
			b.WriteString(fmt.Sprintf("\tp.AssertElement(%q)\n", a.Selector))
		case "as_en":
			b.WriteString(fmt.Sprintf("\tp.AssertElementNotVisible(%q)\n", a.Selector))
		case "hi_lt":
			b.WriteString(fmt.Sprintf("\tp.Highlight(%q)\n", a.Selector))
		case "savsc":
			filename := a.Data
			if filename == "" {
				filename = a.Selector
			}
			if filename == "" {
				filename = "screenshot.png"
			}
			b.WriteString(fmt.Sprintf("\tp.Screenshot(%q)\n", filename))
		default:
			// Skip internal/unknown actions silently, add comment
			b.WriteString(fmt.Sprintf("\t// unknown action: %s\n", a.Action))
		}
	}

	b.WriteString("\treturn nil\n")
	b.WriteString("}\n")
	return b.String()
}

// escapeGoString escapes special characters for inclusion in Go string literals.
func escapeGoString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
