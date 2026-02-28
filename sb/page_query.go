package sb

import "github.com/playwright-community/playwright-go"

func (p *Page) GetText(sel string) (string, error) {
	return p.locator(sel).TextContent()
}

func (p *Page) GetAttribute(sel, attr string) (string, error) {
	return p.locator(sel).GetAttribute(attr)
}

func (p *Page) GetValue(sel string) (string, error) {
	return p.locator(sel).InputValue()
}

func (p *Page) IsVisible(sel string) (bool, error) {
	return p.locator(sel).IsVisible()
}

func (p *Page) IsHidden(sel string) (bool, error) {
	return p.locator(sel).IsHidden()
}

func (p *Page) IsEnabled(sel string) (bool, error) {
	return p.locator(sel).IsEnabled()
}

func (p *Page) IsChecked(sel string) (bool, error) {
	return p.locator(sel).IsChecked()
}

func (p *Page) IsTextVisible(text, sel string) (bool, error) {
	loc := p.locator(sel).Filter(playwright.LocatorFilterOptions{
		HasText: text,
	})
	return loc.IsVisible()
}

func (p *Page) FindElements(sel string) ([]playwright.Locator, error) {
	return p.locator(sel).All()
}

func (p *Page) Count(sel string) (int, error) {
	return p.locator(sel).Count()
}
