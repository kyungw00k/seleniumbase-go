package sb

import (
	"encoding/json"
	"os"

	"github.com/playwright-community/playwright-go"
)

func (p *Page) GetCookies() ([]playwright.Cookie, error) {
	return p.context.Cookies()
}

func (p *Page) AddCookie(cookie playwright.OptionalCookie) error {
	return p.context.AddCookies([]playwright.OptionalCookie{cookie})
}

func (p *Page) ClearCookies() error {
	return p.context.ClearCookies()
}

func (p *Page) SaveCookies(path string) error {
	_, err := p.context.StorageState(path)
	return err
}

func (p *Page) LoadCookies(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var state playwright.StorageState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	for _, c := range state.Cookies {
		if err := p.context.AddCookies([]playwright.OptionalCookie{{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   playwright.String(c.Domain),
			Path:     playwright.String(c.Path),
			Secure:   playwright.Bool(c.Secure),
			HttpOnly: playwright.Bool(c.HttpOnly),
			SameSite: c.SameSite,
		}}); err != nil {
			return err
		}
	}
	return nil
}
