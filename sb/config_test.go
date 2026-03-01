package sb

import (
	"testing"
	"time"
)

func TestWithBrowser(t *testing.T) {
	cfg := newDefaultConfig()
	WithBrowser("firefox")(cfg)
	if cfg.Browser != "firefox" {
		t.Errorf("expected firefox, got %s", cfg.Browser)
	}
}

func TestWithHeadless(t *testing.T) {
	cfg := newDefaultConfig()
	WithHeadless(false)(cfg)
	if cfg.Headless {
		t.Error("expected Headless to be false")
	}
}

func TestWithScreenshotOnFailure(t *testing.T) {
	cfg := newDefaultConfig()
	WithScreenshotOnFailure(true)(cfg)
	if !cfg.ScreenshotOnFailure {
		t.Error("expected ScreenshotOnFailure to be true")
	}
}

func TestWithMobile(t *testing.T) {
	cfg := newDefaultConfig()
	WithMobile(true)(cfg)
	if !cfg.Mobile {
		t.Error("expected Mobile to be true")
	}
}

func TestWithDevice(t *testing.T) {
	cfg := newDefaultConfig()
	WithDevice("iPhone 13")(cfg)
	if cfg.DeviceName != "iPhone 13" {
		t.Errorf("expected iPhone 13, got %s", cfg.DeviceName)
	}
	if !cfg.Mobile {
		t.Error("expected Mobile to be true when device is set")
	}
}

func TestWithRecordVideo(t *testing.T) {
	cfg := newDefaultConfig()
	WithRecordVideo("/tmp/videos")(cfg)
	if cfg.RecordVideo != "/tmp/videos" {
		t.Errorf("expected /tmp/videos, got %s", cfg.RecordVideo)
	}
}

func TestWithRecordHAR(t *testing.T) {
	cfg := newDefaultConfig()
	WithRecordHAR("/tmp/test.har")(cfg)
	if cfg.RecordHAR != "/tmp/test.har" {
		t.Errorf("expected /tmp/test.har, got %s", cfg.RecordHAR)
	}
}

func TestWithDisableCSP(t *testing.T) {
	cfg := newDefaultConfig()
	WithDisableCSP(true)(cfg)
	if !cfg.DisableCSP {
		t.Error("expected DisableCSP to be true")
	}
}

func TestWithIncognito(t *testing.T) {
	cfg := newDefaultConfig()
	WithIncognito(true)(cfg)
	if !cfg.Incognito {
		t.Error("expected Incognito to be true")
	}
}

func TestWithProxy(t *testing.T) {
	cfg := newDefaultConfig()
	WithProxy("http://proxy:8080")(cfg)
	if cfg.Proxy != "http://proxy:8080" {
		t.Errorf("expected http://proxy:8080, got %s", cfg.Proxy)
	}
}

func TestWithUserAgent(t *testing.T) {
	cfg := newDefaultConfig()
	WithUserAgent("TestAgent/1.0")(cfg)
	if cfg.UserAgent != "TestAgent/1.0" {
		t.Errorf("expected TestAgent/1.0, got %s", cfg.UserAgent)
	}
}

func TestWithViewportSize(t *testing.T) {
	cfg := newDefaultConfig()
	WithViewportSize(1920, 1080)(cfg)
	if cfg.ViewportWidth != 1920 || cfg.ViewportHeight != 1080 {
		t.Errorf("expected 1920x1080, got %dx%d", cfg.ViewportWidth, cfg.ViewportHeight)
	}
}

func TestWithSlowMo(t *testing.T) {
	cfg := newDefaultConfig()
	WithSlowMo(100)(cfg)
	if cfg.SlowMo != 100 {
		t.Errorf("expected 100, got %f", cfg.SlowMo)
	}
}

func TestWithTimeout(t *testing.T) {
	cfg := newDefaultConfig()
	WithTimeout(30 * time.Second)(cfg)
	if cfg.Timeout != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.Timeout)
	}
}

func TestWithLocale(t *testing.T) {
	cfg := newDefaultConfig()
	WithLocale("ko-KR")(cfg)
	if cfg.Locale != "ko-KR" {
		t.Errorf("expected ko-KR, got %s", cfg.Locale)
	}
}

func TestWithColorScheme(t *testing.T) {
	cfg := newDefaultConfig()
	WithColorScheme("dark")(cfg)
	if cfg.ColorScheme != "dark" {
		t.Errorf("expected dark, got %s", cfg.ColorScheme)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := newDefaultConfig()
	if cfg.Browser != "chromium" {
		t.Errorf("expected chromium, got %s", cfg.Browser)
	}
	if !cfg.Headless {
		t.Error("expected Headless to be true by default")
	}
	if cfg.ViewportWidth != 1280 || cfg.ViewportHeight != 720 {
		t.Errorf("expected 1280x720, got %dx%d", cfg.ViewportWidth, cfg.ViewportHeight)
	}
	if cfg.Timeout != LargeTimeout {
		t.Errorf("expected LargeTimeout, got %v", cfg.Timeout)
	}
}
