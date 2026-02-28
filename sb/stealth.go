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

var defaultStealthArgs = []string{
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

func newStealthSB(pw *playwright.Playwright, cfg *Config) (*SB, error) {
	chromePath := cfg.ChromePath
	if chromePath == "" {
		var err error
		chromePath, err = findChrome()
		if err != nil {
			return nil, err
		}
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
		contextOpts := playwright.BrowserNewContextOptions{}
		if cfg.ViewportWidth > 0 && cfg.ViewportHeight > 0 {
			contextOpts.Viewport = &playwright.Size{
				Width:  cfg.ViewportWidth,
				Height: cfg.ViewportHeight,
			}
		}
		if cfg.UserAgent != "" {
			contextOpts.UserAgent = playwright.String(cfg.UserAgent)
		}
		if cfg.Locale != "" {
			contextOpts.Locale = playwright.String(cfg.Locale)
		}
		if cfg.IgnoreHTTPSErrors {
			contextOpts.IgnoreHttpsErrors = playwright.Bool(true)
		}
		if cfg.ColorScheme != "" {
			switch cfg.ColorScheme {
			case "dark":
				contextOpts.ColorScheme = playwright.ColorSchemeDark
			case "light":
				contextOpts.ColorScheme = playwright.ColorSchemeLight
			case "no-preference":
				contextOpts.ColorScheme = playwright.ColorSchemeNoPreference
			}
		}
		ctx, err = browser.NewContext(contextOpts)
		if err != nil {
			proc.stop()
			pw.Stop()
			return nil, err
		}
	}

	return &SB{pw: pw, browser: browser, context: ctx, config: cfg, process: proc}, nil
}
