package sb

import "regexp"

func (p *Page) AssertElement(sel string) error {
	return p.expect.Locator(p.locator(sel)).ToBeVisible()
}

func (p *Page) AssertElementPresent(sel string) error {
	return p.expect.Locator(p.locator(sel)).ToBeAttached()
}

func (p *Page) AssertElementAbsent(sel string) error {
	return p.expect.Locator(p.locator(sel)).Not().ToBeAttached()
}

func (p *Page) AssertElementNotVisible(sel string) error {
	return p.expect.Locator(p.locator(sel)).ToBeHidden()
}

func (p *Page) AssertText(text, sel string) error {
	return p.expect.Locator(p.locator(sel)).ToContainText(text)
}

func (p *Page) AssertExactText(text, sel string) error {
	return p.expect.Locator(p.locator(sel)).ToHaveText(text)
}

func (p *Page) AssertTextNotVisible(text, sel string) error {
	return p.expect.Locator(p.locator(sel)).Not().ToContainText(text)
}

func (p *Page) AssertTitle(title string) error {
	return p.expect.Page(p.pw).ToHaveTitle(title)
}

func (p *Page) AssertTitleContains(sub string) error {
	return p.expect.Page(p.pw).ToHaveTitle(regexp.MustCompile(regexp.QuoteMeta(sub)))
}

func (p *Page) AssertURL(url string) error {
	return p.expect.Page(p.pw).ToHaveURL(url)
}

func (p *Page) AssertURLContains(sub string) error {
	return p.expect.Page(p.pw).ToHaveURL(regexp.MustCompile(regexp.QuoteMeta(sub)))
}

func (p *Page) AssertAttribute(sel, attr, val string) error {
	return p.expect.Locator(p.locator(sel)).ToHaveAttribute(attr, val)
}
