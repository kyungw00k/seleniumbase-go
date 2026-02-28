package sb

import (
	"time"

	"github.com/playwright-community/playwright-go"
)

func (p *Page) WaitForElement(sel string) error {
	return p.locator(sel).WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
}

func (p *Page) WaitForElementPresent(sel string) error {
	return p.locator(sel).WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateAttached,
	})
}

func (p *Page) WaitForElementAbsent(sel string) error {
	return p.locator(sel).WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateDetached,
	})
}

func (p *Page) WaitForElementNotVisible(sel string) error {
	return p.locator(sel).WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateHidden,
	})
}

func (p *Page) WaitForText(text, sel string) error {
	return p.expect.Locator(p.locator(sel)).ToContainText(text)
}

func (p *Page) WaitForLoadState(state string) error {
	var ls *playwright.LoadState
	switch state {
	case "load":
		ls = playwright.LoadStateLoad
	case "domcontentloaded":
		ls = playwright.LoadStateDomcontentloaded
	case "networkidle":
		ls = playwright.LoadStateNetworkidle
	default:
		ls = playwright.LoadStateLoad
	}
	return p.pw.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: ls,
	})
}

func (p *Page) WaitForURL(url string) error {
	return p.pw.WaitForURL(url)
}

func (p *Page) Sleep(d time.Duration) {
	time.Sleep(d)
}
