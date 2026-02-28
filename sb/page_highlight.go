package sb

import (
	"github.com/kyungw00k/seleniumbase-go/selector"
)

// highlightJS is the JavaScript snippet that animates an element's outline.
const highlightJS = `
async ({selector, loops}) => {
    const el = document.querySelector(selector);
    if (!el) return;
    const orig = el.style.outline;
    for (let i = 0; i < loops; i++) {
        el.style.outline = '3px solid red';
        await new Promise(r => setTimeout(r, 100));
        el.style.outline = orig;
        await new Promise(r => setTimeout(r, 100));
    }
}
`

// Highlight flashes a red outline on the element matching sel.
// The optional loops parameter controls the number of flash cycles (default: 4).
func (p *Page) Highlight(sel string, loops ...int) error {
	n := 4
	if len(loops) > 0 && loops[0] > 0 {
		n = loops[0]
	}
	pwSel := selector.Parse(sel)
	_, err := p.pw.Evaluate(highlightJS, map[string]any{"selector": pwSel, "loops": n})
	return err
}

// HighlightClick highlights the element, then clicks it.
func (p *Page) HighlightClick(sel string) error {
	_ = p.Highlight(sel, 3)
	return p.Click(sel)
}

// HighlightType highlights the element, then types text into it.
func (p *Page) HighlightType(sel, text string) error {
	_ = p.Highlight(sel, 3)
	return p.Type(sel, text)
}

// RemoveHighlights removes all highlight overlays added by Highlight.
func (p *Page) RemoveHighlights() error {
	_, err := p.pw.Evaluate(`() => document.querySelectorAll('[data-sb-highlight]').forEach(e => e.remove())`)
	return err
}
