# Selector Syntax

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
