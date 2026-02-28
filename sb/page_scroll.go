package sb

import (
	"fmt"
	"math"
	"time"
)

// ScrollTo scrolls the page until the element matching sel is visible.
func (p *Page) ScrollTo(sel string) error {
	return p.locator(sel).ScrollIntoViewIfNeeded()
}

// ScrollToTop scrolls the page to the very top.
func (p *Page) ScrollToTop() error {
	_, err := p.pw.Evaluate("window.scrollTo(0, 0)")
	return err
}

// ScrollToBottom scrolls the page to the very bottom.
func (p *Page) ScrollToBottom() error {
	_, err := p.pw.Evaluate("window.scrollTo(0, document.body.scrollHeight)")
	return err
}

// ScrollToY scrolls the page to the given vertical pixel offset.
func (p *Page) ScrollToY(y int) error {
	_, err := p.pw.Evaluate(fmt.Sprintf("window.scrollTo(0, %d)", y))
	return err
}

// ScrollUp scrolls the page up by px pixels. Defaults to 100 if px is 0.
func (p *Page) ScrollUp(px ...float64) error {
	delta := 100.0
	if len(px) > 0 && px[0] != 0 {
		delta = px[0]
	}
	return p.pw.Mouse().Wheel(0, -delta)
}

// ScrollDown scrolls the page down by px pixels. Defaults to 100 if px is 0.
func (p *Page) ScrollDown(px ...float64) error {
	delta := 100.0
	if len(px) > 0 && px[0] != 0 {
		delta = px[0]
	}
	return p.pw.Mouse().Wheel(0, delta)
}

// SlowScrollTo scrolls smoothly to the element matching sel using animated steps.
func (p *Page) SlowScrollTo(sel string) error {
	loc := p.locator(sel)
	box, err := loc.BoundingBox()
	if err != nil || box == nil {
		// Element not in viewport yet, try direct scroll
		return loc.ScrollIntoViewIfNeeded()
	}

	// Get current scroll position
	result, err := p.pw.Evaluate("() => ({ x: window.scrollX, y: window.scrollY })")
	if err != nil {
		return loc.ScrollIntoViewIfNeeded()
	}
	pos, ok := result.(map[string]interface{})
	if !ok {
		return loc.ScrollIntoViewIfNeeded()
	}
	currentY, _ := pos["y"].(float64)
	targetY := currentY + box.Y

	// Animate scroll in steps
	distance := targetY - currentY
	steps := int(math.Max(10, math.Abs(distance)/50))
	for i := 1; i <= steps; i++ {
		y := currentY + distance*float64(i)/float64(steps)
		_, err := p.pw.Evaluate(fmt.Sprintf("window.scrollTo(0, %f)", y))
		if err != nil {
			break
		}
		time.Sleep(30 * time.Millisecond)
	}

	return loc.ScrollIntoViewIfNeeded()
}
