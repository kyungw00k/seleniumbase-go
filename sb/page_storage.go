package sb

import "fmt"

func (p *Page) SetLocalStorage(key, val string) error {
	_, err := p.pw.Evaluate("([k, v]) => localStorage.setItem(k, v)", []string{key, val})
	return err
}

func (p *Page) GetLocalStorage(key string) (string, error) {
	result, err := p.pw.Evaluate("k => localStorage.getItem(k)", key)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return fmt.Sprintf("%v", result), nil
}

func (p *Page) ClearLocalStorage() error {
	_, err := p.pw.Evaluate("() => localStorage.clear()")
	return err
}

func (p *Page) SetSessionStorage(key, val string) error {
	_, err := p.pw.Evaluate("([k, v]) => sessionStorage.setItem(k, v)", []string{key, val})
	return err
}

func (p *Page) GetSessionStorage(key string) (string, error) {
	result, err := p.pw.Evaluate("k => sessionStorage.getItem(k)", key)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return fmt.Sprintf("%v", result), nil
}

func (p *Page) ClearSessionStorage() error {
	_, err := p.pw.Evaluate("() => sessionStorage.clear()")
	return err
}
