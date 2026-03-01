package sb

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// StartConsoleCapture begins capturing browser console messages.
func (p *Page) StartConsoleCapture() {
	p.consoleMu.Lock()
	p.consoleMessages = nil
	p.consoleMu.Unlock()
	p.pw.OnConsole(func(msg playwright.ConsoleMessage) {
		p.consoleMu.Lock()
		p.consoleMessages = append(p.consoleMessages, ConsoleMessage{
			Type: msg.Type(),
			Text: msg.Text(),
		})
		p.consoleMu.Unlock()
	})
}

// GetConsoleMessages returns all captured console messages.
func (p *Page) GetConsoleMessages() []ConsoleMessage {
	p.consoleMu.Lock()
	defer p.consoleMu.Unlock()
	cp := make([]ConsoleMessage, len(p.consoleMessages))
	copy(cp, p.consoleMessages)
	return cp
}

// GetConsoleErrors returns only error-type console messages.
func (p *Page) GetConsoleErrors() []ConsoleMessage {
	p.consoleMu.Lock()
	defer p.consoleMu.Unlock()
	var errors []ConsoleMessage
	for _, m := range p.consoleMessages {
		if m.Type == "error" {
			errors = append(errors, m)
		}
	}
	return errors
}

// AssertNoJsErrors asserts that no JavaScript errors have been logged to the console.
func (p *Page) AssertNoJsErrors() error {
	errs := p.GetConsoleErrors()
	if len(errs) > 0 {
		return fmt.Errorf("sb: %d JavaScript error(s) found: %s", len(errs), errs[0].Text)
	}
	return nil
}
