# seleniumbase-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kyungw00k/seleniumbase-go.svg)](https://pkg.go.dev/github.com/kyungw00k/seleniumbase-go)
[![CI](https://github.com/kyungw00k/seleniumbase-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kyungw00k/seleniumbase-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

SeleniumBase-style browser automation for Go, built on [playwright-go](https://github.com/playwright-community/playwright-go).

---

## Features

- SeleniumBase-compatible API — familiar method names for anyone coming from the Python library
- Three usage patterns: script runner, `testing.T` integration, and manual lifecycle
- Automatic Playwright lifecycle management (launch, context, page, teardown)
- Selector translation layer — use `id=`, `name=`, `link=`, `css=`, `xpath=` prefixes or bare CSS
- Playwright assertion engine wired in for reliable auto-waiting assertions
- Cookie persistence — save and load browser state to disk
- Local/session storage helpers
- Tab management, frame access, dialog handling, PDF and screenshot capture
- Full escape hatches to the underlying `playwright.Page`, `playwright.Locator`, and `playwright.BrowserContext`
- Scroll methods — smooth, animated, and directional scrolling
- Network interception — route, abort, and mock API responses
- Highlight / demo mode — visual element highlighting for debugging
- Deferred assertions — queue assertion failures and process them together
- MFA / TOTP support — generate and enter time-based one-time passwords
- CDP stealth mode — launch Chrome externally with anti-detection flags, connect via CDP
- Recorder — capture browser interactions and generate Go test code

---

## Installation

```bash
go get github.com/kyungw00k/seleniumbase-go
```

Install Playwright browser binaries (one-time):

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5700.1 install --with-deps
```

---

## Quick Start

### `sb.Run` — script usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/kyungw00k/seleniumbase-go/sb"
)

func main() {
    err := sb.Run(func(p *sb.Page) error {
        p.Open("https://example.com")

        title, err := p.GetTitle()
        if err != nil {
            return err
        }
        fmt.Println("Page title:", title)

        p.AssertTitle("Example Domain")
        p.AssertElement("h1")
        p.AssertText("Example Domain", "h1")
        return nil
    }, sb.WithBrowser("chromium"), sb.WithHeadless(true))

    if err != nil {
        log.Fatal(err)
    }
}
```

### `sb.RunTest` — `testing.T` integration

```go
//go:build integration

package mypackage_test

import (
    "testing"

    "github.com/kyungw00k/seleniumbase-go/sb"
)

func TestLogin(t *testing.T) {
    sb.RunTest(t, func(p *sb.Page) {
        p.Open("https://www.saucedemo.com")
        p.Type("#user-name", "standard_user")
        p.Type("#password", "secret_sauce")
        p.Click("#login-button")
        p.AssertElement("div.inventory_list")
        p.AssertExactText("Products", "span.title")
    }, sb.WithHeadless(true))
}
```

Any error or panic inside the function is forwarded to `t.Fatal`.

### `sb.NewPage` — manual lifecycle

```go
page, cleanup, err := sb.NewPage(sb.WithHeadless(false), sb.WithBrowser("firefox"))
if err != nil {
    log.Fatal(err)
}
defer cleanup()

page.Open("https://example.com")
page.AssertTitle("Example Domain")
```

---

## Configuration Options

Options are passed as variadic arguments to `sb.Run`, `sb.RunTest`, or `sb.NewPage`.

| Option | Signature | Default | Description |
|--------|-----------|---------|-------------|
| `WithBrowser` | `WithBrowser(browser string)` | `"chromium"` | Browser engine: `"chromium"`, `"firefox"`, `"webkit"` |
| `WithHeadless` | `WithHeadless(headless bool)` | `true` | Run browser without a visible window |
| `WithProxy` | `WithProxy(server string)` | — | HTTP/SOCKS proxy server URL |
| `WithUserAgent` | `WithUserAgent(agent string)` | — | Custom User-Agent header |
| `WithViewportSize` | `WithViewportSize(width, height int)` | `1280x720` | Browser viewport dimensions |
| `WithSlowMo` | `WithSlowMo(ms float64)` | `0` | Slow down all actions by the given milliseconds |
| `WithTimeout` | `WithTimeout(d time.Duration)` | `10s` | Default timeout for all page operations |
| `WithLocale` | `WithLocale(code string)` | — | Browser locale, e.g. `"en-US"`, `"ko-KR"` |
| `WithIgnoreHTTPSErrors` | `WithIgnoreHTTPSErrors(ignore bool)` | `false` | Accept invalid TLS certificates |
| `WithColorScheme` | `WithColorScheme(scheme string)` | — | Emulate color scheme: `"dark"`, `"light"`, `"no-preference"` |
| `WithDemoMode` | `WithDemoMode(enabled bool)` | `false` | Enable demo mode (highlight elements before interaction) |
| `WithStealth` | `WithStealth(enabled bool)` | `false` | Enable CDP stealth mode (Chromium only) |
| `WithChromePath` | `WithChromePath(path string)` | auto-detected | Custom Chrome executable path |
| `WithUserDataDir` | `WithUserDataDir(dir string)` | temp dir | Custom user data directory for stealth mode |
| `WithExtraArgs` | `WithExtraArgs(args ...string)` | — | Additional Chrome launch arguments |

### Timeout constants

```go
sb.MiniTimeout     // 2s
sb.SmallTimeout    // 7s
sb.LargeTimeout    // 10s  (default)
sb.ExtremeTimeout  // 30s
sb.PageLoadTimeout // 120s
```

---

## Selector Syntax

All methods that accept a `sel string` parameter run the selector through the translation layer before passing it to Playwright.

| SeleniumBase selector | Playwright equivalent | Example |
|-----------------------|-----------------------|---------|
| `id=myId` | `#myId` | `p.Click("id=login-button")` |
| `name=fieldName` | `[name="fieldName"]` | `p.Type("name=q", "go")` |
| `link=Exact Link Text` | `a:has-text("Exact Link Text")` | `p.Click("link=Sign in")` |
| `partial_link=Partial` | `a:has-text("Partial")` | `p.Click("partial_link=Sign")` |
| `text=Hello` | `text=Hello` (Playwright native) | `p.AssertText("Hello", "text=Hello")` |
| `role=button` | `role=button` (Playwright native) | `p.Click("role=button")` |
| `label=Email` | `label=Email` (Playwright native) | `p.Type("label=Email", "test@example.com")` |
| `css=div.header` | `div.header` | `p.AssertElement("css=div.header")` |
| `xpath=//h1` | `//h1` | `p.AssertElement("xpath=//h1")` |
| bare CSS (default) | passed through as-is | `p.Click("button.submit")` |

XPath expressions starting with `//`, `./`, or `(//` are also passed through directly.

---

## API Reference

All methods are on `*sb.Page`. Methods that can fail return an `error` as the last return value.

### Navigation

```go
func (p *Page) Open(url string) error
func (p *Page) GoBack() error
func (p *Page) GoForward() error
func (p *Page) Refresh() error
func (p *Page) GetCurrentURL() string
func (p *Page) GetTitle() (string, error)
func (p *Page) GetPageSource() (string, error)
func (p *Page) SetContent(html string) error
```

### Interaction

```go
func (p *Page) Click(sel string) error
func (p *Page) DoubleClick(sel string) error
func (p *Page) RightClick(sel string) error
func (p *Page) Type(sel, text string) error          // fills the field (replaces existing value)
func (p *Page) SendKeys(sel, text string) error      // types character by character
func (p *Page) Press(sel, key string) error          // presses a named key, e.g. "Enter", "Tab"
func (p *Page) Clear(sel string) error
func (p *Page) Focus(sel string) error
func (p *Page) Hover(sel string) error
func (p *Page) Check(sel string) error
func (p *Page) Uncheck(sel string) error
func (p *Page) SelectOption(sel string, values playwright.SelectOptionValues) error
func (p *Page) SetInputFiles(sel string, files any) error
func (p *Page) DragAndDrop(srcSel, dstSel string) error
```

### Assertions

Assertions use Playwright's built-in auto-retry engine and respect the configured timeout.

```go
func (p *Page) AssertElement(sel string) error                        // element is visible
func (p *Page) AssertElementPresent(sel string) error                 // element is in the DOM
func (p *Page) AssertElementAbsent(sel string) error                  // element is not in the DOM
func (p *Page) AssertElementNotVisible(sel string) error              // element is hidden
func (p *Page) AssertText(text, sel string) error                     // element contains text
func (p *Page) AssertExactText(text, sel string) error                // element text matches exactly
func (p *Page) AssertTextNotVisible(text, sel string) error           // element does not contain text
func (p *Page) AssertTitle(title string) error                        // page title matches exactly
func (p *Page) AssertTitleContains(sub string) error                  // page title contains substring
func (p *Page) AssertURL(url string) error                            // URL matches exactly
func (p *Page) AssertURLContains(sub string) error                    // URL contains substring
func (p *Page) AssertAttribute(sel, attr, val string) error           // attribute equals value
```

### Wait

```go
func (p *Page) WaitForElement(sel string) error                       // waits until visible
func (p *Page) WaitForElementPresent(sel string) error                // waits until attached
func (p *Page) WaitForElementAbsent(sel string) error                 // waits until detached
func (p *Page) WaitForElementNotVisible(sel string) error             // waits until hidden
func (p *Page) WaitForText(text, sel string) error                    // waits until element contains text
func (p *Page) WaitForLoadState(state string) error                   // "load", "domcontentloaded", "networkidle"
func (p *Page) WaitForURL(url string) error
func (p *Page) Sleep(d time.Duration)
```

### Query

```go
func (p *Page) GetText(sel string) (string, error)
func (p *Page) GetAttribute(sel, attr string) (string, error)
func (p *Page) GetValue(sel string) (string, error)
func (p *Page) IsVisible(sel string) (bool, error)
func (p *Page) IsHidden(sel string) (bool, error)
func (p *Page) IsEnabled(sel string) (bool, error)
func (p *Page) IsChecked(sel string) (bool, error)
func (p *Page) IsTextVisible(text, sel string) (bool, error)
func (p *Page) FindElements(sel string) ([]playwright.Locator, error)
func (p *Page) Count(sel string) (int, error)
```

### Cookie

```go
func (p *Page) GetCookies() ([]playwright.Cookie, error)
func (p *Page) AddCookie(cookie playwright.OptionalCookie) error
func (p *Page) ClearCookies() error
func (p *Page) SaveCookies(path string) error   // saves full browser storage state to a JSON file
func (p *Page) LoadCookies(path string) error   // loads cookies from a saved storage state file
```

### JavaScript

```go
func (p *Page) Evaluate(expr string, arg ...any) (any, error)
func (p *Page) EvalOnSelector(sel, expr string) (any, error)
```

### Window / Tab

```go
func (p *Page) OpenNewTab() (*Page, error)
func (p *Page) SwitchToTab(index int) error     // 0-based index within the current context
func (p *Page) SetViewportSize(width, height int) error
func (p *Page) BringToFront() error
```

### Frame

```go
func (p *Page) FrameLocator(sel string) playwright.FrameLocator
func (p *Page) MainFrame() playwright.Frame
```

### Screenshot / PDF

```go
func (p *Page) Screenshot(path string) error
func (p *Page) FullPageScreenshot(path string) error
func (p *Page) PDF(path string) error                          // Chromium only
func (p *Page) ElementScreenshot(sel, path string) error
```

### Alert / Dialog

```go
func (p *Page) OnDialog(fn func(playwright.Dialog))   // register a custom handler
func (p *Page) AcceptDialogs()                        // auto-accept all future dialogs
func (p *Page) DismissDialogs()                       // auto-dismiss all future dialogs
```

### Storage

```go
func (p *Page) SetLocalStorage(key, val string) error
func (p *Page) GetLocalStorage(key string) (string, error)
func (p *Page) ClearLocalStorage() error
func (p *Page) SetSessionStorage(key, val string) error
func (p *Page) GetSessionStorage(key string) (string, error)
func (p *Page) ClearSessionStorage() error
```

### Scroll

```go
func (p *Page) ScrollTo(sel string) error
func (p *Page) ScrollToTop() error
func (p *Page) ScrollToBottom() error
func (p *Page) ScrollToY(y int) error
func (p *Page) ScrollUp(px ...float64) error
func (p *Page) ScrollDown(px ...float64) error
func (p *Page) SlowScrollTo(sel string) error
```

### Network

```go
func (p *Page) Route(url string, handler func(playwright.Route)) error
func (p *Page) Unroute(url string) error
func (p *Page) RouteAbort(urlPattern string) error
func (p *Page) MockAPI(urlPattern string, body string, opts ...MockAPIOptions) error
```

### Highlight

```go
func (p *Page) Highlight(sel string, loops ...int) error
func (p *Page) HighlightClick(sel string) error
func (p *Page) HighlightType(sel, text string) error
func (p *Page) RemoveHighlights() error
```

### Deferred Assertions

```go
func (p *Page) DeferredAssertElement(sel string) bool
func (p *Page) DeferredAssertElementPresent(sel string) bool
func (p *Page) DeferredAssertText(text, sel string) bool
func (p *Page) DeferredAssertExactText(text, sel string) bool
func (p *Page) ProcessDeferredAsserts() error
func (p *Page) ClearDeferredAsserts()
```

### MFA / TOTP

```go
func (p *Page) GetMFACode(totpKey string) (string, error)
func (p *Page) EnterMFACode(sel, totpKey string) error
```

### Recorder

Record browser interactions and generate Go test code automatically.

```go
// Manual control
sb.Run(func(p *sb.Page) error {
    p.StartRecording()
    p.Open("https://example.com")
    // ... interact with the page ...
    actions, _ := p.StopRecording()
    code := sb.GenerateGoCode(actions)
    os.WriteFile("recorded_test.go", []byte(code), 0644)
    return nil
}, sb.WithHeadless(false))

// Convenience: opens browser, records until closed, writes Go file
sb.RunRecorder("recorded_test.go", sb.WithHeadless(false))
```

Keyboard shortcuts during recording:
- `Shift+S` — Assert text on hovered element
- `Shift+E` — Assert element exists
- `Shift+P` — Assert element present in DOM
- `Shift+V` — Assert element visible
- `Shift+N` — Assert element not visible
- `Shift+H` — Highlight element
- `Shift+G` — Save screenshot
- `Escape` — Pause recording
- `` ` `` — Resume recording

### Stealth Mode

Stealth mode launches Chrome externally with anti-detection flags and connects
Playwright via CDP. This avoids the automation fingerprints that Playwright's
built-in launch adds.

```go
sb.Run(func(p *sb.Page) error {
    p.Open("https://bot.sannysoft.com")
    p.Screenshot("stealth-check.png")
    return nil
}, sb.WithStealth(true), sb.WithHeadless(false))
```

When stealth mode is enabled:
- Only Chromium is supported (not Firefox or WebKit)
- Chrome must be installed on the system (or specify path with `WithChromePath`)
- A temporary user data directory is created and cleaned up on close
- All standard `sb.Page` methods work normally

---

## Escape Hatches

When you need Playwright features not yet wrapped by seleniumbase-go, drop down to the underlying objects:

```go
// Raw playwright.Page — full Playwright API
pwPage := page.Playwright()

// playwright.Locator with selector translation applied
loc := page.Locator("id=submit-btn")
loc.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible})

// playwright.BrowserContext — cookies, storage state, new pages
ctx := page.Context()
state, _ := ctx.StorageState("state.json")
```

---

## Examples

| Example | Location | Description |
|---------|----------|-------------|
| Simple script | `examples/simple/main.go` | `sb.Run` usage with title and assertion checks |
| Integration test | `examples/basic_test.go` | `sb.RunTest` with a full login/cart/logout flow on saucedemo.com |

Run the integration test:

```bash
go test -tags integration ./examples/...
```

---

## Roadmap

**Phase 1 (complete)** — Core API: navigation, interaction, assertions, wait, query, cookies, JS, window/tab, frames, screenshots, dialogs, storage

**Phase 2 (complete)** — Extension features: scroll methods, network interception, highlight/demo mode, deferred assertions, extended selectors (`text=`, `role=`, `label=`), MFA/TOTP support

**Phase 3 (in progress)**

- ~~CDP stealth mode for bot detection bypass~~ (complete)
- ~~Recorder — capture and replay browser sessions~~ (complete)
- Visual testing — screenshot comparison and diffing
- Parallel test runner utilities
- Report generation (HTML, JUnit)

---

## License

MIT. See [LICENSE](LICENSE).
