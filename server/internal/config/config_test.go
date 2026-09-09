package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("COOKIE_SECURE", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddress != ":8080" || !cfg.CookieSecure || cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://db.example/journal")
	t.Setenv("HTTP_ADDR", "127.0.0.1:9000")
	t.Setenv("COOKIE_SECURE", "false")
	t.Setenv("SHUTDOWN_TIMEOUT", "3s")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddress != "127.0.0.1:9000" || cfg.CookieSecure || cfg.ShutdownTimeout != 3*time.Second {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct{ name, key, value string }{
		{"database URL", "DATABASE_URL", "not-a-url"},
		{"HTTP address", "HTTP_ADDR", "8080"},
		{"cookie flag", "COOKIE_SECURE", "sometimes"},
		{"shutdown timeout", "SHUTDOWN_TIMEOUT", "0s"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("DATABASE_URL", "postgres://localhost/journal")
			t.Setenv("HTTP_ADDR", ":8080")
			t.Setenv("COOKIE_SECURE", "true")
			t.Setenv("SHUTDOWN_TIMEOUT", "10s")
			t.Setenv(test.key, test.value)
			if _, err := Load(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
