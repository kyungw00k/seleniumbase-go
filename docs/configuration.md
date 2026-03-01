# Configuration

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
| `WithScreenshotOnFailure` | `WithScreenshotOnFailure(enabled bool)` | `false` | Auto-screenshot when test fails |
| `WithMobile` | `WithMobile(mobile bool)` | `false` | Emulate mobile device (touch events) |
| `WithDevice` | `WithDevice(name string)` | — | Emulate named device (e.g. `"iPhone 13"`) |
| `WithRecordVideo` | `WithRecordVideo(dir string)` | — | Record test video to directory |
| `WithRecordHAR` | `WithRecordHAR(path string)` | — | Record HAR network log |
| `WithDisableCSP` | `WithDisableCSP(disable bool)` | `false` | Disable Content Security Policy |
| `WithIncognito` | `WithIncognito(incognito bool)` | `false` | Use incognito context |
| `WithChannel` | `WithChannel(channel string)` | — | Browser channel: `"chrome"`, `"msedge"` (use system browser) |
| `WithRemoteCDPURL` | `WithRemoteCDPURL(url string)` | — | Connect to remote browser via CDP |
| `WithRemoteWSURL` | `WithRemoteWSURL(url string)` | — | Connect to remote Playwright server via WebSocket |

## Timeout constants

```go
sb.MiniTimeout     // 2s
sb.SmallTimeout    // 7s
sb.LargeTimeout    // 10s  (default)
sb.ExtremeTimeout  // 30s
sb.PageLoadTimeout // 120s
```
