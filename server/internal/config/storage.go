package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Storage struct {
	Endpoint, AccessKey, SecretKey, Bucket string
	Secure                                 bool
}

func loadStorage() (Storage, error) {
	cfg := Storage{AccessKey: os.Getenv("S3_ACCESS_KEY"), SecretKey: os.Getenv("S3_SECRET_KEY"), Bucket: os.Getenv("S3_BUCKET")}
	endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	if cfg.Bucket == "" {
		cfg.Bucket = "journal-lesson-attachments"
	}
	if endpoint == "" {
		if cfg.AccessKey != "" || cfg.SecretKey != "" {
			return cfg, fmt.Errorf("S3_ENDPOINT is required with S3 credentials")
		}
		return cfg, nil
	}
	u, err := url.Parse(endpoint)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return cfg, fmt.Errorf("S3_ENDPOINT must be an http(s) origin")
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return cfg, fmt.Errorf("S3_ACCESS_KEY and S3_SECRET_KEY are required")
	}
	cfg.Endpoint, cfg.Secure = u.Host, u.Scheme == "https"
	return cfg, nil
}
