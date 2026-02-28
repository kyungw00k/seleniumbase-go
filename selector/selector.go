package selector

import "strings"

// Parse converts a SeleniumBase-style selector to a Playwright-compatible selector.
func Parse(sel string) string {
	switch {
	case strings.HasPrefix(sel, "link="):
		text := strings.TrimPrefix(sel, "link=")
		return `a:has-text("` + text + `")`
	case strings.HasPrefix(sel, "partial_link="):
		text := strings.TrimPrefix(sel, "partial_link=")
		return `a:has-text("` + text + `")`
	case strings.HasPrefix(sel, "name="):
		name := strings.TrimPrefix(sel, "name=")
		return `[name="` + name + `"]`
	case strings.HasPrefix(sel, "css="):
		return strings.TrimPrefix(sel, "css=")
	case strings.HasPrefix(sel, "xpath="):
		return strings.TrimPrefix(sel, "xpath=")
	case strings.HasPrefix(sel, "id="):
		id := strings.TrimPrefix(sel, "id=")
		return "#" + id
	default:
		return sel
	}
}

// IsXPath returns true if the selector is an XPath expression.
func IsXPath(sel string) bool {
	return strings.HasPrefix(sel, "//") ||
		strings.HasPrefix(sel, "./") ||
		strings.HasPrefix(sel, "(//") ||
		strings.HasPrefix(sel, "xpath=")
}

// IsLinkText returns true if the selector uses link= prefix.
func IsLinkText(sel string) bool {
	return strings.HasPrefix(sel, "link=") || strings.HasPrefix(sel, "partial_link=")
}
