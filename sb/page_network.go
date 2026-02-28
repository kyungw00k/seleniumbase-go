package sb

import "github.com/playwright-community/playwright-go"

// MockAPIOptions configures mock API responses.
type MockAPIOptions struct {
	Status      int               // HTTP status code (default: 200)
	ContentType string            // Content-Type header (default: "application/json")
	Headers     map[string]string // Additional response headers
}

// Route registers a route handler for URLs matching the given pattern.
func (p *Page) Route(url string, handler func(playwright.Route)) error {
	return p.pw.Route(url, handler)
}

// Unroute removes a previously registered route handler for the given URL pattern.
func (p *Page) Unroute(url string) error {
	return p.pw.Unroute(url)
}

// RouteAbort blocks all requests matching the given URL pattern.
func (p *Page) RouteAbort(urlPattern string) error {
	return p.pw.Route(urlPattern, func(route playwright.Route) {
		_ = route.Abort()
	})
}

// MockAPI mocks API responses for requests matching the given URL pattern.
func (p *Page) MockAPI(urlPattern string, body string, opts ...MockAPIOptions) error {
	opt := MockAPIOptions{
		Status:      200,
		ContentType: "application/json",
	}
	if len(opts) > 0 {
		if opts[0].Status != 0 {
			opt.Status = opts[0].Status
		}
		if opts[0].ContentType != "" {
			opt.ContentType = opts[0].ContentType
		}
		if opts[0].Headers != nil {
			opt.Headers = opts[0].Headers
		}
	}

	return p.pw.Route(urlPattern, func(route playwright.Route) {
		headers := map[string]string{
			"Content-Type": opt.ContentType,
		}
		for k, v := range opt.Headers {
			headers[k] = v
		}
		_ = route.Fulfill(playwright.RouteFulfillOptions{
			Status:      playwright.Int(opt.Status),
			Body:        body,
			ContentType: playwright.String(opt.ContentType),
			Headers:     headers,
		})
	})
}
