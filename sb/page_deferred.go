package sb

import (
	"fmt"
	"strings"
)

// DeferredAssertElement queues a visibility check for the given selector.
// Returns true if the assertion passed, false if it was queued as a failure.
func (p *Page) DeferredAssertElement(sel string) bool {
	if err := p.AssertElement(sel); err != nil {
		p.deferredFailures = append(p.deferredFailures, fmt.Errorf("DeferredAssertElement(%q): %w", sel, err))
		return false
	}
	return true
}

// DeferredAssertElementPresent queues a DOM presence check for the given selector.
// Returns true if the assertion passed, false if it was queued as a failure.
func (p *Page) DeferredAssertElementPresent(sel string) bool {
	if err := p.AssertElementPresent(sel); err != nil {
		p.deferredFailures = append(p.deferredFailures, fmt.Errorf("DeferredAssertElementPresent(%q): %w", sel, err))
		return false
	}
	return true
}

// DeferredAssertText queues a text-contains check for the given text within the element matching sel.
// Returns true if the assertion passed, false if it was queued as a failure.
func (p *Page) DeferredAssertText(text, sel string) bool {
	if err := p.AssertText(text, sel); err != nil {
		p.deferredFailures = append(p.deferredFailures, fmt.Errorf("DeferredAssertText(%q, %q): %w", text, sel, err))
		return false
	}
	return true
}

// DeferredAssertExactText queues an exact text match check for the given text within the element matching sel.
// Returns true if the assertion passed, false if it was queued as a failure.
func (p *Page) DeferredAssertExactText(text, sel string) bool {
	if err := p.AssertExactText(text, sel); err != nil {
		p.deferredFailures = append(p.deferredFailures, fmt.Errorf("DeferredAssertExactText(%q, %q): %w", text, sel, err))
		return false
	}
	return true
}

// ProcessDeferredAsserts returns a combined error of all queued assertion failures
// and clears the queue. Returns nil if no failures were queued.
func (p *Page) ProcessDeferredAsserts() error {
	if len(p.deferredFailures) == 0 {
		return nil
	}
	var msgs []string
	for i, err := range p.deferredFailures {
		msgs = append(msgs, fmt.Sprintf("  %d) %s", i+1, err.Error()))
	}
	combined := fmt.Errorf("deferred assertions failed (%d):\n%s", len(p.deferredFailures), strings.Join(msgs, "\n"))
	p.deferredFailures = nil
	return combined
}

// ClearDeferredAsserts discards all queued assertion failures without reporting them.
func (p *Page) ClearDeferredAsserts() {
	p.deferredFailures = nil
}
