# Feature Guide

## Usage Patterns

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

## Recorder

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

## Stealth Mode

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

## Visual Testing

Compare screenshots against stored baselines for pixel-level regression testing.

```go
sb.Run(func(p *sb.Page) error {
    p.Open("https://example.com")

    // First run saves a baseline; subsequent runs compare against it
    result, err := p.CheckWindow("homepage", 0.5) // 0.5% threshold
    if err != nil {
        return err
    }
    fmt.Printf("Match: %v, Diff: %.2f%%\n", result.Match, result.DiffPercent)

    // Or use the assertion form
    p.AssertVisualMatch("homepage", 0.5)

    // Force-update the baseline
    p.UpdateBaseline("homepage")
    return nil
}, sb.WithHeadless(true))
```

Baselines are stored in `visual_baseline/` with `_baseline.png`, `_latest.png`, and `_diff.png` suffixes.

## Parallel Test Runner

Run multiple test functions concurrently, each with its own browser instance.

```go
// Standalone usage
tests := []sb.TestFunc{
    {Name: "Login", Fn: func(p *sb.Page) error {
        p.Open("https://example.com/login")
        return nil
    }},
    {Name: "Search", Fn: func(p *sb.Page) error {
        p.Open("https://example.com/search")
        return nil
    }},
}

results := sb.RunParallel(tests, sb.WithHeadless(true))
fmt.Println(sb.ParallelSummary(results))

// testing.T integration
func TestAll(t *testing.T) {
    sb.RunParallelTest(t, tests, sb.WithHeadless(true))
}
```

## Report Generation

Generate JUnit XML or styled HTML reports from parallel test results.

```go
results := sb.RunParallel(tests, sb.WithHeadless(true))

// JUnit XML report (CI-compatible)
sb.GenerateJUnitReport("report.xml", "MySuite", results)

// HTML report with visual summary
sb.GenerateHTMLReport("report.html", "Test Results", results)
```

## Remote Browser

Connect to remote browsers for cloud testing (BrowserStack, LambdaTest, Sauce Labs).

```go
// Connect via Chrome DevTools Protocol
sb.Run(func(p *sb.Page) error {
    p.Open("https://example.com")
    return nil
}, sb.WithRemoteCDPURL("wss://cdp.browserstack.com/playwright?caps=..."))

// Connect via WebSocket to a Playwright server
sb.Run(func(p *sb.Page) error {
    p.Open("https://example.com")
    return nil
}, sb.WithRemoteWSURL("ws://localhost:3000"))

// Use system Chrome instead of bundled Chromium
sb.Run(func(p *sb.Page) error {
    p.Open("https://example.com")
    return nil
}, sb.WithChannel("chrome"))
```

## MasterQA

Hybrid manual/automated testing — run automated steps, then pause for human verification.

```go
func TestManualCheck(t *testing.T) {
    sb.RunMasterQA(t, func(m *sb.MasterQAPage) {
        m.Open("https://example.com")
        m.AssertElement("h1")

        // Pauses and asks tester via terminal
        m.Verify("Does the page layout look correct?")
        m.Verify("Are all images loading properly?")

        // Save verification report
        m.SaveReport("masterqa_report.html")
    }, sb.WithHeadless(false))
}
```

In headless mode, all `Verify` calls automatically pass (useful for CI).

## Auto Browser Detection

seleniumbase-go automatically detects system Chrome and falls back to installing
Playwright's bundled Chromium if no browser is found.

```go
// Manually trigger browser detection and installation
channel, err := sb.EnsureBrowser()

// Find system Chrome path
path := sb.FindSystemChrome()
```

The auto-install runs automatically when `playwright.Run()` fails on first use.
No manual `playwright install` step is needed.
