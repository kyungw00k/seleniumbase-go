package sb

import (
	"os"
	"strings"
	"testing"
)

func containsArg(args []string, arg string) bool {
	for _, a := range args {
		if a == arg {
			return true
		}
	}
	return false
}

func TestFindChrome(t *testing.T) {
	path, err := findChrome()
	if err != nil {
		t.Skip("Chrome not installed")
	}
	if path == "" {
		t.Fatal("findChrome returned empty path")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Chrome path does not exist: %s", path)
	}
}

func TestFreePort(t *testing.T) {
	port, err := freePort()
	if err != nil {
		t.Fatalf("freePort failed: %v", err)
	}
	if port <= 0 || port >= 65536 {
		t.Fatalf("invalid port: %d", port)
	}
}

func TestBuildStealthArgs_Defaults(t *testing.T) {
	cfg := &Config{ViewportWidth: 1280, ViewportHeight: 720}
	args := buildStealthArgs(cfg, 9222, "/tmp/test-profile")

	checks := []string{
		"--remote-debugging-port=9222",
		"--user-data-dir=/tmp/test-profile",
		"--remote-debugging-host=127.0.0.1",
		"--window-size=1280,720",
	}
	for _, c := range checks {
		if !containsArg(args, c) {
			t.Errorf("expected arg %q not found", c)
		}
	}
	for _, def := range defaultStealthArgs {
		if !containsArg(args, def) {
			t.Errorf("default stealth arg missing: %q", def)
		}
	}
}

func TestBuildStealthArgs_Headless(t *testing.T) {
	cfg := &Config{Headless: true, ViewportWidth: 1280, ViewportHeight: 720}
	args := buildStealthArgs(cfg, 9222, "/tmp/test")
	if !containsArg(args, "--headless=new") {
		t.Error("--headless=new not found")
	}
}

func TestBuildStealthArgs_Proxy(t *testing.T) {
	cfg := &Config{Proxy: "http://proxy:8080", ViewportWidth: 1280, ViewportHeight: 720}
	args := buildStealthArgs(cfg, 9222, "/tmp/test")
	checks := []string{
		"--proxy-server=http://proxy:8080",
		"--ignore-certificate-errors",
		"--ignore-ssl-errors=yes",
	}
	for _, c := range checks {
		if !containsArg(args, c) {
			t.Errorf("expected arg %q not found", c)
		}
	}
}

func TestBuildStealthArgs_UserAgent(t *testing.T) {
	cfg := &Config{UserAgent: "MyBot/1.0", ViewportWidth: 1280, ViewportHeight: 720}
	args := buildStealthArgs(cfg, 9222, "/tmp/test")
	if !containsArg(args, "--user-agent=MyBot/1.0") {
		t.Error("--user-agent=MyBot/1.0 not found")
	}
}

func TestBuildStealthArgs_ExtraArgs(t *testing.T) {
	cfg := &Config{
		ExtraArgs:      []string{"--custom-flag", "--another=val"},
		ViewportWidth:  1280,
		ViewportHeight: 720,
	}
	args := buildStealthArgs(cfg, 9222, "/tmp/test")
	if !containsArg(args, "--custom-flag") {
		t.Error("--custom-flag not found")
	}
	if !containsArg(args, "--another=val") {
		t.Error("--another=val not found")
	}
}

func TestBuildStealthArgs_NoDuplicates(t *testing.T) {
	cfg := &Config{
		ExtraArgs:      []string{"--no-first-run", "--animation-duration-scale=5"},
		ViewportWidth:  1280,
		ViewportHeight: 720,
	}
	args := buildStealthArgs(cfg, 9222, "/tmp/test")

	prefixes := []string{"--no-first-run", "--animation-duration-scale"}
	for _, prefix := range prefixes {
		count := 0
		for _, a := range args {
			argPrefix := a
			if idx := strings.Index(a, "="); idx >= 0 {
				argPrefix = a[:idx]
			}
			if argPrefix == prefix {
				count++
			}
		}
		if count != 1 {
			t.Errorf("prefix %q appears %d times, expected 1", prefix, count)
		}
	}
}

func TestDefaultStealthArgs_NoEnableAutomation(t *testing.T) {
	for _, arg := range defaultStealthArgs {
		if arg == "--enable-automation" {
			t.Error("defaultStealthArgs contains --enable-automation")
		}
		if strings.Contains(arg, "enable-automation") {
			t.Errorf("defaultStealthArgs contains enable-automation substring: %q", arg)
		}
	}
}

func TestBuildHeadlessUA_NoHeadless(t *testing.T) {
	// Empty chromePath falls back to default version — UA must not contain "Headless".
	ua := buildHeadlessUA("")
	if strings.Contains(ua, "Headless") {
		t.Errorf("UA contains Headless token: %q", ua)
	}
	if !strings.Contains(ua, "Chrome/") {
		t.Errorf("UA missing Chrome/ token: %q", ua)
	}
}

func TestBuildHeadlessUA_WithChrome(t *testing.T) {
	path, err := findChrome()
	if err != nil {
		t.Skip("Chrome not installed")
	}
	ua := buildHeadlessUA(path)
	if strings.Contains(ua, "Headless") {
		t.Errorf("UA contains Headless token: %q", ua)
	}
	if !strings.Contains(ua, "Chrome/") {
		t.Errorf("UA missing Chrome/ token: %q", ua)
	}
}

func TestGetChromeVersion(t *testing.T) {
	path, err := findChrome()
	if err != nil {
		t.Skip("Chrome not installed")
	}
	version := getChromeVersion(path)
	if version == "" {
		t.Error("getChromeVersion returned empty string")
	}
	// Version should look like "145.0.7449.84"
	if !strings.Contains(version, ".") {
		t.Errorf("version %q does not look like a version string", version)
	}
}

