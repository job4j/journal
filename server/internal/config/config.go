package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL     string
	HTTPAddress     string
	CookieSecure    bool
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL: "postgres://postgres:password@127.0.0.1:5433/journal?sslmode=disable",
		HTTPAddress: ":8080", CookieSecure: true, ShutdownTimeout: 10 * time.Second,
	}
	if value := os.Getenv("DATABASE_URL"); value != "" {
		cfg.DatabaseURL = value
	}
	if value := os.Getenv("HTTP_ADDR"); value != "" {
		cfg.HTTPAddress = value
	}
	if value := os.Getenv("COOKIE_SECURE"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("COOKIE_SECURE: %w", err)
		}
		cfg.CookieSecure = parsed
	}
	if value := os.Getenv("SHUTDOWN_TIMEOUT"); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil || parsed <= 0 {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be a positive duration")
		}
		cfg.ShutdownTimeout = parsed
	}
	parsedURL, err := url.ParseRequestURI(cfg.DatabaseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return Config{}, fmt.Errorf("DATABASE_URL must be an absolute URL")
	}
	if _, _, err := net.SplitHostPort(cfg.HTTPAddress); err != nil {
		return Config{}, fmt.Errorf("HTTP_ADDR: %w", err)
	}
	return cfg, nil
}
