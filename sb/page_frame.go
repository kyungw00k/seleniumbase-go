package sb

import (
	"github.com/playwright-community/playwright-go"
	"github.com/kyungw00k/seleniumbase-go/selector"
)

func (p *Page) FrameLocator(sel string) playwright.FrameLocator {
	return p.pw.FrameLocator(selector.Parse(sel))
}

func (p *Page) MainFrame() playwright.Frame {
	return p.pw.MainFrame()
}
