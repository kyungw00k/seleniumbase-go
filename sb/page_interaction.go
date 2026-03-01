package sb

import (
	"time"

	"github.com/kyungw00k/seleniumbase-go/selector"
	"github.com/playwright-community/playwright-go"
)

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

// JsClick clicks an element via JavaScript, bypassing overlay issues.
func (p *Page) JsClick(sel string) error {
	_, err := p.pw.Evaluate(`(selector) => document.querySelector(selector).click()`, selector.Parse(sel))
	return err
}

// ClickIfVisible clicks the element only if it is currently visible.
func (p *Page) ClickIfVisible(sel string) error {
	visible, err := p.IsVisible(sel)
	if err != nil || !visible {
		return err
	}
	return p.Click(sel)
}

// HoverAndClick hovers over one element and then clicks another.
func (p *Page) HoverAndClick(hoverSel, clickSel string) error {
	if err := p.Hover(hoverSel); err != nil {
		return err
	}
	return p.Click(clickSel)
}

// SlowClick clicks an element after waiting for the given duration.
func (p *Page) SlowClick(sel string, d time.Duration) error {
	if err := p.locator(sel).WaitFor(); err != nil {
		return err
	}
	time.Sleep(d)
	return p.Click(sel)
}

// SelectOptionByText selects a <select> option by its visible text.
func (p *Page) SelectOptionByText(sel, text string) error {
	_, err := p.locator(sel).SelectOption(playwright.SelectOptionValues{Labels: &[]string{text}})
	return err
}

// SelectOptionByValue selects a <select> option by its value attribute.
func (p *Page) SelectOptionByValue(sel, value string) error {
	_, err := p.locator(sel).SelectOption(playwright.SelectOptionValues{Values: &[]string{value}})
	return err
}

// SelectOptionByIndex selects a <select> option by its zero-based index.
func (p *Page) SelectOptionByIndex(sel string, index int) error {
	_, err := p.locator(sel).SelectOption(playwright.SelectOptionValues{Indexes: &[]int{index}})
	return err
}
