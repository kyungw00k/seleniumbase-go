**English** | [한국어](README.ko.md) | [中文](README.zh.md) | [日本語](README.ja.md)

# seleniumbase-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kyungw00k/seleniumbase-go.svg)](https://pkg.go.dev/github.com/kyungw00k/seleniumbase-go)
[![CI](https://github.com/kyungw00k/seleniumbase-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kyungw00k/seleniumbase-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go port of [SeleniumBase](https://github.com/seleniumbase/SeleniumBase) — simple, powerful browser automation built on [playwright-go](https://github.com/playwright-community/playwright-go).

---

## Installation

```bash
go get github.com/kyungw00k/seleniumbase-go
```

Playwright browser binaries are **auto-installed** on first run. Or install manually:

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5700.1 install --with-deps
```

---

## Quick Start

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

See also: [`sb.RunTest`](docs/features.md#sbruntest--testingt-integration) for `testing.T` integration, [`sb.NewPage`](docs/features.md#sbnewpage--manual-lifecycle) for manual lifecycle control.

---

## Usage Patterns

seleniumbase-go provides three ways to use the browser:

| Pattern | Use Case | Lifecycle |
|---------|----------|-----------|
| `sb.Run(func(p *sb.Page) error { ... })` | Scripts & automation | Automatic |
| `sb.RunTest(t, func(p *sb.Page) { ... })` | Go tests with `testing.T` | Automatic, failures go to `t.Fatal` |
| `sb.NewPage(opts...)` → `page, cleanup, err` | Full control | Manual via `defer cleanup()` |

---

## Features

- **SeleniumBase-compatible API** — familiar method names from the Python library
- **Selector translation** — use `id=`, `name=`, `link=`, `css=`, `xpath=` prefixes or bare CSS ([details](docs/selectors.md))
- **Auto-waiting assertions** — Playwright's retry engine for reliable tests
- **90+ Page methods** — navigation, interaction, assertions, queries, cookies, storage, JS eval, scroll, network, and more ([full API](docs/api.md))
- **25+ configuration options** — browser, headless, proxy, viewport, timeout, locale, and more ([all options](docs/configuration.md))
- **CDP stealth mode** — launch Chrome with anti-detection flags ([guide](docs/features.md#stealth-mode))
- **Recorder** — capture browser interactions and generate Go test code ([guide](docs/features.md#recorder))
- **Visual testing** — pixel-level screenshot comparison ([guide](docs/features.md#visual-testing))
- **Parallel test runner** — concurrent execution with isolated browsers ([guide](docs/features.md#parallel-test-runner))
- **Report generation** — JUnit XML and styled HTML reports ([guide](docs/features.md#report-generation))
- **Remote browser** — connect to BrowserStack, LambdaTest via CDP or WebSocket ([guide](docs/features.md#remote-browser))
- **MasterQA** — hybrid automated/manual verification ([guide](docs/features.md#masterqa))
- **Auto browser detection** — finds system Chrome or installs Playwright Chromium ([guide](docs/features.md#auto-browser-detection))
- **Network interception** — route, abort, and mock API responses
- **Highlight / demo mode** — visual element highlighting for debugging
- **Deferred assertions** — queue failures and process together
- **MFA / TOTP support** — generate and enter time-based one-time passwords
- **Cookie persistence** — save and load browser state
- **Tab management, frames, dialogs** — full browser control
- **Escape hatches** — drop down to `playwright.Page`, `playwright.Locator`, `playwright.BrowserContext` anytime

---

## Documentation

| Document | Description |
|----------|-------------|
| [API Reference](docs/api.md) | All 90+ Page methods with signatures |
| [Configuration](docs/configuration.md) | 25+ options and timeout constants |
| [Selector Syntax](docs/selectors.md) | SeleniumBase selector translation rules |
| [Feature Guide](docs/features.md) | Detailed guides for stealth, recorder, visual testing, and more |

---

## Examples

| Example | Description |
|---------|-------------|
| [`examples/simple/main.go`](examples/simple/main.go) | `sb.Run` — title and assertion checks |
| [`examples/basic_test.go`](examples/basic_test.go) | `sb.RunTest` — login/cart/logout flow |
| [`examples/demo_site_test.go`](examples/demo_site_test.go) | Demo mode with highlighting |
| [`examples/visual_test.go`](examples/visual_test.go) | Visual regression testing |
| [`examples/stealth_coupang_test.go`](examples/stealth_coupang_test.go) | CDP stealth mode |

Run integration tests:

```bash
go test -tags integration ./examples/...
```

---

## Internationalization (i18n)

SeleniumBase Python supports [method name translations](https://github.com/seleniumbase/SeleniumBase/tree/master/seleniumbase/translate) for 10 languages via static subclassing (e.g., `self.클릭()` in Korean). This is **not supported** in the Go port — Go method names are resolved at compile time, making runtime or alias-based translation impractical.

---

## Inspired By

This project is inspired by [SeleniumBase](https://github.com/seleniumbase/SeleniumBase), the comprehensive Python browser automation framework. seleniumbase-go brings the same simple, powerful API design to the Go ecosystem using [playwright-go](https://github.com/playwright-community/playwright-go) as the underlying engine.

---

## License

MIT. See [LICENSE](LICENSE).
