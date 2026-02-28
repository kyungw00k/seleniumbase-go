package sb

import "github.com/playwright-community/playwright-go"

func (p *Page) Open(url string) error {
	_, err := p.pw.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	return err
}

func (p *Page) GoBack() error {
	_, err := p.pw.GoBack()
	return err
}

func (p *Page) GoForward() error {
	_, err := p.pw.GoForward()
	return err
}

func (p *Page) Refresh() error {
	_, err := p.pw.Reload()
	return err
}

func (p *Page) GetCurrentURL() string {
	return p.pw.URL()
}

func (p *Page) GetTitle() (string, error) {
	return p.pw.Title()
}

func (p *Page) GetPageSource() (string, error) {
	return p.pw.Content()
}

func (p *Page) SetContent(html string) error {
	return p.pw.SetContent(html)
}
