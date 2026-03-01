# API Reference

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
func (p *Page) LoadHTMLFile(path string) error              // open a local HTML file
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
func (p *Page) JsClick(sel string) error               // click via JavaScript (bypasses overlays)
func (p *Page) ClickIfVisible(sel string) error         // click only if element is visible
func (p *Page) HoverAndClick(hoverSel, clickSel string) error  // hover one element, click another
func (p *Page) SlowClick(sel string, d time.Duration) error    // wait then click
func (p *Page) SelectOptionByText(sel, text string) error      // select by visible text
func (p *Page) SelectOptionByValue(sel, value string) error    // select by value attribute
func (p *Page) SelectOptionByIndex(sel string, index int) error // select by index
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
func (p *Page) AssertNoJsErrors() error                     // assert no JS console errors
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
func (p *Page) WaitForAttribute(sel, attr, val string) error  // wait for attribute value
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
func (p *Page) GetUserAgent() (string, error)               // get browser user agent
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
func (p *Page) SetAttribute(sel, attr, val string) error    // set DOM attribute
func (p *Page) RemoveAttribute(sel, attr string) error      // remove DOM attribute
func (p *Page) HideElement(sel string) error                // hide element (display:none)
func (p *Page) ShowElement(sel string) error                // show hidden element
func (p *Page) RemoveElement(sel string) error              // remove element from DOM
func (p *Page) DisableBeforeunload() error                  // suppress "leave page?" dialogs
func (p *Page) SetValue(sel, value string) error            // set input value via JS
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

### Console

```go
func (p *Page) StartConsoleCapture()                        // begin capturing console messages
func (p *Page) GetConsoleMessages() []ConsoleMessage        // get all captured messages
func (p *Page) GetConsoleErrors() []ConsoleMessage          // get only error messages
```

### Download

```go
func (p *Page) WaitForDownload(action func() error) (string, error)  // wait for download, return path
func (p *Page) SaveDownload(path string, action func() error) error  // download and save to path
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
