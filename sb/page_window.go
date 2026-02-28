package sb

import "fmt"

func (p *Page) OpenNewTab() (*Page, error) {
	newPwPage, err := p.context.NewPage()
	if err != nil {
		return nil, err
	}
	return newPage(newPwPage, p.context, p.config), nil
}

func (p *Page) SwitchToTab(index int) error {
	pages := p.context.Pages()
	if index < 0 || index >= len(pages) {
		return fmt.Errorf("sb: tab index %d out of range (have %d tabs)", index, len(pages))
	}
	p.pw = pages[index]
	return p.pw.BringToFront()
}

func (p *Page) SetViewportSize(width, height int) error {
	return p.pw.SetViewportSize(width, height)
}

func (p *Page) BringToFront() error {
	return p.pw.BringToFront()
}
