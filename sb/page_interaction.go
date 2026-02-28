package sb

import "github.com/playwright-community/playwright-go"

func (p *Page) Click(sel string) error {
	return p.locator(sel).Click()
}

func (p *Page) DoubleClick(sel string) error {
	return p.locator(sel).Dblclick()
}

func (p *Page) RightClick(sel string) error {
	return p.locator(sel).Click(playwright.LocatorClickOptions{
		Button: playwright.MouseButtonRight,
	})
}

func (p *Page) Type(sel, text string) error {
	return p.locator(sel).Fill(text)
}

func (p *Page) SendKeys(sel, text string) error {
	return p.locator(sel).PressSequentially(text)
}

func (p *Page) Press(sel, key string) error {
	return p.locator(sel).Press(key)
}

func (p *Page) Clear(sel string) error {
	return p.locator(sel).Clear()
}

func (p *Page) Focus(sel string) error {
	return p.locator(sel).Focus()
}

func (p *Page) Hover(sel string) error {
	return p.locator(sel).Hover()
}

func (p *Page) Check(sel string) error {
	return p.locator(sel).Check()
}

func (p *Page) Uncheck(sel string) error {
	return p.locator(sel).Uncheck()
}

func (p *Page) SelectOption(sel string, values playwright.SelectOptionValues) error {
	_, err := p.locator(sel).SelectOption(values)
	return err
}

func (p *Page) SetInputFiles(sel string, files any) error {
	return p.locator(sel).SetInputFiles(files)
}

func (p *Page) DragAndDrop(srcSel, dstSel string) error {
	return p.locator(srcSel).DragTo(p.locator(dstSel))
}
