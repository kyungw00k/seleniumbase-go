package sb

import "testing"

func TestWithRemoteCDPURL(t *testing.T) {
	cfg := newDefaultConfig()
	WithRemoteCDPURL("http://localhost:9222")(cfg)
	if cfg.RemoteCDPURL != "http://localhost:9222" {
		t.Errorf("expected http://localhost:9222, got %s", cfg.RemoteCDPURL)
	}
}

func TestWithRemoteWSURL(t *testing.T) {
	cfg := newDefaultConfig()
	WithRemoteWSURL("ws://localhost:3000")(cfg)
	if cfg.RemoteWSURL != "ws://localhost:3000" {
		t.Errorf("expected ws://localhost:3000, got %s", cfg.RemoteWSURL)
	}
}
