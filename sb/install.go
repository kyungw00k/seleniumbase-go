package sb

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/playwright-community/playwright-go"
)

// FindSystemChrome returns the path to a system-installed Chrome/Chromium browser.
// Returns empty string if no system Chrome is found.
func FindSystemChrome() string {
	// Try PATH-based lookup first
	names := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium-browser",
		"chromium",
		"chrome",
	}
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}

	// Fallback to OS-specific well-known paths
	var candidates []string
	switch runtime.GOOS {
	case "darwin":
		candidates = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		}
	case "linux":
		candidates = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
			"/snap/bin/chromium",
		}
	case "windows":
		for _, envVar := range []string{"PROGRAMFILES", "PROGRAMFILES(X86)", "LOCALAPPDATA"} {
			if dir := os.Getenv(envVar); dir != "" {
				candidates = append(candidates, filepath.Join(dir, `Google\Chrome\Application\chrome.exe`))
			}
		}
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// EnsureBrowser ensures that a browser is available for Playwright.
// If system Chrome is found, configures Playwright to use it via Channel.
// If no browser is found, installs Playwright's bundled Chromium.
// Returns the channel string to use ("chrome" or "") and any error.
func EnsureBrowser() (channel string, err error) {
	if FindSystemChrome() != "" {
		// System Chrome found — install only the Playwright driver (no bundled browsers)
		if err := playwright.Install(&playwright.RunOptions{
			SkipInstallBrowsers: true,
		}); err != nil {
			// Driver might already be installed, try to proceed
			return "chrome", nil
		}
		return "chrome", nil
	}

	// No system Chrome — install Playwright's bundled Chromium
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
	}); err != nil {
		return "", fmt.Errorf("sb: could not install browser: %w", err)
	}
	return "", nil
}
