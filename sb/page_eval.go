package sb

import "github.com/kyungw00k/seleniumbase-go/selector"

func (p *Page) Evaluate(expr string, arg ...any) (any, error) {
	return p.pw.Evaluate(expr, arg...)
}

func (p *Page) EvalOnSelector(sel, expr string) (any, error) {
	return p.pw.EvalOnSelector(selector.Parse(sel), expr, nil)
}
