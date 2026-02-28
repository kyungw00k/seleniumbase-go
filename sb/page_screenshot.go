package sb

import "github.com/playwright-community/playwright-go"

func (p *Page) Screenshot(path string) error {
	_, err := p.pw.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(path),
	})
	return err
}

func (p *Page) FullPageScreenshot(path string) error {
	_, err := p.pw.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
	return err
}

func (p *Page) PDF(path string) error {
	_, err := p.pw.PDF(playwright.PagePdfOptions{
		Path: playwright.String(path),
	})
	return err
}

func (p *Page) ElementScreenshot(sel, path string) error {
	_, err := p.locator(sel).Screenshot(playwright.LocatorScreenshotOptions{
		Path: playwright.String(path),
	})
	return err
}
