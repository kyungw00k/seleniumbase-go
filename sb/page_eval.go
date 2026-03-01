package sb

import "github.com/kyungw00k/seleniumbase-go/selector"

func (p *Page) Evaluate(expr string, arg ...any) (any, error) {
	return p.pw.Evaluate(expr, arg...)
}

func (p *Page) EvalOnSelector(sel, expr string) (any, error) {
	return p.pw.EvalOnSelector(selector.Parse(sel), expr, nil)
}

// SetAttribute sets an attribute on the first element matching sel.
func (p *Page) SetAttribute(sel, attr, val string) error {
	_, err := p.pw.Evaluate(`([s, a, v]) => document.querySelector(s).setAttribute(a, v)`, []string{selector.Parse(sel), attr, val})
	return err
}

// RemoveAttribute removes an attribute from the first element matching sel.
func (p *Page) RemoveAttribute(sel, attr string) error {
	_, err := p.pw.Evaluate(`([s, a]) => document.querySelector(s).removeAttribute(a)`, []string{selector.Parse(sel), attr})
	return err
}

// HideElement hides an element by setting display:none.
func (p *Page) HideElement(sel string) error {
	_, err := p.pw.Evaluate(`(s) => document.querySelector(s).style.display = 'none'`, selector.Parse(sel))
	return err
}

// ShowElement shows a hidden element by clearing the inline display style.
func (p *Page) ShowElement(sel string) error {
	_, err := p.pw.Evaluate(`(s) => document.querySelector(s).style.display = ''`, selector.Parse(sel))
	return err
}

// RemoveElement removes an element from the DOM.
func (p *Page) RemoveElement(sel string) error {
	_, err := p.pw.Evaluate(`(s) => { const el = document.querySelector(s); if (el) el.remove(); }`, selector.Parse(sel))
	return err
}

// DisableBeforeunload removes onbeforeunload handlers to prevent "leave page?" dialogs.
func (p *Page) DisableBeforeunload() error {
	_, err := p.pw.Evaluate(`() => window.onbeforeunload = null`)
	return err
}

// SetValue sets the value property on an input element via JavaScript.
func (p *Page) SetValue(sel, value string) error {
	_, err := p.pw.Evaluate(`([s, v]) => { const el = document.querySelector(s); el.value = v; el.dispatchEvent(new Event('input', {bubbles: true})); el.dispatchEvent(new Event('change', {bubbles: true})); }`, []string{selector.Parse(sel), value})
	return err
}
