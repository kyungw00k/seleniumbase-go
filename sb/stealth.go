package sb

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// getChromeVersion runs "chrome --version" and returns the version string (e.g. "145.0.7449.84").
// Returns an empty string on failure.
func getChromeVersion(chromePath string) string {
	out, err := exec.Command(chromePath, "--version").Output()
	if err != nil {
		return ""
	}
	// Output: "Google Chrome 145.0.7449.84\n" or "Chromium 145.0.7449.84\n"
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) >= 1 {
		return fields[len(fields)-1]
	}
	return ""
}

// buildHeadlessUA returns a realistic Chrome user agent for the current OS.
// Chrome 112+ uses a "frozen" UA format (major.0.0.0) to prevent fingerprinting.
// This replaces the "HeadlessChrome/major.0.0.0" marker that Chrome headless=new
// still reports, making navigator.userAgent indistinguishable from headed Chrome.
func buildHeadlessUA(chromePath string) string {
	full := getChromeVersion(chromePath)
	// Parse major version from "145.0.7632.117" → "145"
	major := full
	if idx := strings.Index(full, "."); idx >= 0 {
		major = full[:idx]
	}
	if major == "" {
		major = "131"
	}
	// Chrome frozen UA format: major.0.0.0
	frozenVersion := major + ".0.0.0"

	var osPart string
	switch runtime.GOOS {
	case "darwin":
		osPart = "Macintosh; Intel Mac OS X 10_15_7"
	case "linux":
		osPart = "X11; Linux x86_64"
	case "windows":
		osPart = "Windows NT 10.0; Win64; x64"
	default:
		osPart = "Macintosh; Intel Mac OS X 10_15_7"
	}
	return fmt.Sprintf(
		"Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
		osPart, frozenVersion,
	)
}

var defaultStealthArgs = []string{
	"--disable-blink-features=AutomationControlled",
	"--no-first-run",
	"--no-service-autorun",
	"--disable-auto-reload",
	"--no-default-browser-check",
	"--homepage=about:blank",
	"--no-pings",
	"--password-store=basic",
	"--deny-permission-prompts",
	"--disable-breakpad",
	"--disable-setuid-sandbox",
	"--disable-prompt-on-repost",
	"--disable-application-cache",
	"--disable-password-generation",
	"--disable-save-password-bubble",
	"--disable-single-click-autofill",
	"--disable-ipc-flooding-protection",
	"--disable-background-timer-throttling",
	"--disable-search-engine-choice-screen",
	"--disable-backgrounding-occluded-windows",
	"--disable-client-side-phishing-detection",
	"--disable-device-discovery-notifications",
	"--disable-top-sites",
	"--disable-translate",
	"--dns-prefetch-disable",
	"--disable-renderer-backgrounding",
	"--disable-dev-shm-usage",
	"--test-type",
	"--ash-no-nudges",
	"--wm-window-animations-disabled",
	"--animation-duration-scale=0",
	"--safebrowsing-disable-download-protection",
	`--simulate-outdated-no-au=Tue, 31 Dec 2099 23:59:59 GMT`,
	"--enable-privacy-sandbox-ads-apis",
	"--enable-unsafe-extension-debugging",
	"--use-mock-keychain",
	"--disable-features=IsolateOrigins,site-per-process,Translate," +
		"InsecureDownloadWarnings,DownloadBubble,DownloadBubbleV2," +
		"OptimizationTargetPrediction,OptimizationGuideModelDownloading," +
		"SidePanelPinning,UserAgentClientHint,PrivacySandboxSettings4," +
		"OptimizationHintsFetching,InterestFeedContentSuggestions," +
		"Bluetooth,WebBluetooth,UnifiedWebBluetooth,ComponentUpdater," +
		"DisableLoadExtensionCommandLineSwitch," +
		"WebAuthentication,PasskeyAuth",
}

func findChrome() (string, error) {
	var candidates []string

	switch runtime.GOOS {
	case "darwin":
		candidates = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"google-chrome",
			"chromium",
		}
	case "linux":
		candidates = []string{
			"google-chrome",
			"google-chrome-stable",
			"google-chrome-beta",
			"chromium",
			"chromium-browser",
			"/usr/bin/google-chrome",
			"/snap/bin/chromium",
		}
	case "windows":
		candidates = []string{"chrome.exe"}
		absoluteWindows := []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			absoluteWindows = append(absoluteWindows, filepath.Join(localAppData, `Google\Chrome\Application\chrome.exe`))
		}
		for _, p := range absoluteWindows {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	for _, candidate := range candidates {
		if filepath.IsAbs(candidate) {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		} else {
			if path, err := exec.LookPath(candidate); err == nil {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("sb: could not find Chrome; install Chrome or set WithChromePath")
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port, nil
}

func buildStealthArgs(cfg *Config, port int, userDataDir string) []string {
	args := append([]string{}, defaultStealthArgs...)
	args = append(args, fmt.Sprintf("--remote-debugging-host=127.0.0.1"))
	args = append(args, fmt.Sprintf("--remote-debugging-port=%d", port))
	args = append(args, fmt.Sprintf("--user-data-dir=%s", userDataDir))
	args = append(args, fmt.Sprintf("--window-size=%d,%d", cfg.ViewportWidth, cfg.ViewportHeight))

	if cfg.Headless {
		args = append(args, "--headless=new")
	}
	if cfg.Proxy != "" {
		args = append(args, fmt.Sprintf("--proxy-server=%s", cfg.Proxy))
		args = append(args, "--ignore-certificate-errors")
		args = append(args, "--ignore-ssl-errors=yes")
	}
	if cfg.UserAgent != "" {
		args = append(args, fmt.Sprintf("--user-agent=%s", cfg.UserAgent))
	}

	for _, extra := range cfg.ExtraArgs {
		prefix := extra
		if idx := strings.Index(extra, "="); idx >= 0 {
			prefix = extra[:idx]
		}
		duplicate := false
		for _, existing := range args {
			existingPrefix := existing
			if idx := strings.Index(existing, "="); idx >= 0 {
				existingPrefix = existing[:idx]
			}
			if existingPrefix == prefix {
				duplicate = true
				break
			}
		}
		if !duplicate {
			args = append(args, extra)
		}
	}

	return args
}

func waitForCDP(host string, port int, timeout time.Duration) error {
	url := fmt.Sprintf("http://%s:%d/json/version", host, port)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.After(timeout)

	for {
		select {
		case <-deadline:
			return fmt.Errorf("sb: Chrome did not start within %s", timeout)
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

type stealthProcess struct {
	cmd         *exec.Cmd
	userDataDir string
	tempDir     bool
}

func (sp *stealthProcess) stop() {
	if sp.cmd != nil && sp.cmd.Process != nil {
		sp.cmd.Process.Signal(os.Interrupt)
		done := make(chan error, 1)
		go func() { done <- sp.cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			sp.cmd.Process.Kill()
		}
	}
	if sp.tempDir && sp.userDataDir != "" {
		os.RemoveAll(sp.userDataDir)
	}
}

// applyStealthCDP sends CDP commands that apply per-page stealth settings.
// User-Agent is handled via the --user-agent= Chrome launch flag (not CDP), because
// Network.setUserAgentOverride without userAgentMetadata suppresses sec-ch-ua headers,
// which WAFs like Akamai detect as an automation signal. The Chrome flag changes
// navigator.userAgent while leaving sec-ch-ua intact.
func applyStealthCDP(ctx playwright.BrowserContext, page playwright.Page, cfg *Config) error {
	session, err := ctx.NewCDPSession(page)
	if err != nil {
		return fmt.Errorf("sb: could not create CDP session: %w", err)
	}
	// Do not detach — overrides must persist for the lifetime of the page.

	// Enable network domain — mirrors Python SeleniumBase cdp_driver behavior.
	if _, err = session.Send("Network.enable", nil); err != nil {
		return fmt.Errorf("sb: Network.enable failed: %w", err)
	}

	// Bypass Content-Security-Policy to allow scripts to run on strict pages.
	if _, err = session.Send("Page.setBypassCSP", map[string]any{"enabled": true}); err != nil {
		return fmt.Errorf("sb: Page.setBypassCSP failed: %w", err)
	}

	return nil
}

func newStealthSB(pw *playwright.Playwright, cfg *Config) (*SB, error) {
	chromePath := cfg.ChromePath
	if chromePath == "" {
		var err error
		chromePath, err = findChrome()
		if err != nil {
			return nil, err
		}
	}

	// Chrome headless=new still reports "HeadlessChrome/major.0.0.0" in its User-Agent,
	// which WAFs like Akamai detect via navigator.userAgent. Override with the real
	// Chrome UA using the --user-agent= flag. This changes navigator.userAgent without
	// suppressing sec-ch-ua (unlike Network.setUserAgentOverride without metadata).
	if cfg.Headless && cfg.UserAgent == "" {
		cfg.UserAgent = buildHeadlessUA(chromePath)
	}

	port, err := freePort()
	if err != nil {
		return nil, err
	}

	var userDataDir string
	tempDir := false
	if cfg.UserDataDir != "" {
		userDataDir = cfg.UserDataDir
	} else {
		userDataDir, err = os.MkdirTemp("", "sb-stealth-*")
		if err != nil {
			return nil, err
		}
		tempDir = true
	}

	args := buildStealthArgs(cfg, port, userDataDir)
	cmd := exec.Command(chromePath, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		if tempDir {
			os.RemoveAll(userDataDir)
		}
		return nil, err
	}

	proc := &stealthProcess{cmd: cmd, userDataDir: userDataDir, tempDir: tempDir}

	if err := waitForCDP("127.0.0.1", port, 15*time.Second); err != nil {
		proc.stop()
		pw.Stop()
		return nil, err
	}

	browser, err := pw.Chromium.ConnectOverCDP(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		proc.stop()
		pw.Stop()
		return nil, err
	}

	var ctx playwright.BrowserContext
	if contexts := browser.Contexts(); len(contexts) > 0 {
		ctx = contexts[0]
	} else {
		ctx, err = browser.NewContext(buildContextOpts(pw, cfg))
		if err != nil {
			proc.stop()
			pw.Stop()
			return nil, err
		}
	}

	return &SB{pw: pw, browser: browser, context: ctx, config: cfg, process: proc}, nil
}
