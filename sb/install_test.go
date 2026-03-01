package sb

import "testing"

func TestFindSystemChrome(t *testing.T) {
	// This test just verifies the function doesn't panic
	// The result depends on the system's Chrome installation
	path := FindSystemChrome()
	t.Logf("System Chrome path: %q", path)
}

func TestWithChannel(t *testing.T) {
	cfg := newDefaultConfig()
	WithChannel("chrome")(cfg)
	if cfg.Channel != "chrome" {
		t.Errorf("expected chrome, got %s", cfg.Channel)
	}
}

func TestWithChannel_MSEdge(t *testing.T) {
	cfg := newDefaultConfig()
	WithChannel("msedge")(cfg)
	if cfg.Channel != "msedge" {
		t.Errorf("expected msedge, got %s", cfg.Channel)
	}
}
