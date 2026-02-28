package sb

import "github.com/playwright-community/playwright-go"

func (p *Page) OnDialog(fn func(playwright.Dialog)) {
	p.pw.OnDialog(fn)
}

func (p *Page) AcceptDialogs() {
	p.pw.OnDialog(func(d playwright.Dialog) {
		d.Accept()
	})
}

func (p *Page) DismissDialogs() {
	p.pw.OnDialog(func(d playwright.Dialog) {
		d.Dismiss()
	})
}
