package config

import "testing"

func TestStorageConfiguration(t *testing.T) {
	for _, test := range []struct {
		name, endpoint, key, secret string
		valid                       bool
	}{
		{"disabled", "", "", "", true},
		{"http", "http://localhost:9000", "test-key", "test-secret", true},
		{"https", "https://storage.example.test", "test-key", "test-secret", true},
		{"missing endpoint", "", "key", "secret", false},
		{"missing credentials", "http://localhost:9000", "", "", false},
		{"credentials in URL", "http://key:secret@localhost:9000", "key", "secret", false},
		{"not an origin", "http://localhost:9000/bucket", "key", "secret", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("S3_ENDPOINT", test.endpoint)
			t.Setenv("S3_ACCESS_KEY", test.key)
			t.Setenv("S3_SECRET_KEY", test.secret)
			t.Setenv("S3_BUCKET", "")
			cfg, err := loadStorage()
			if (err == nil) != test.valid {
				t.Fatalf("unexpected config error: %v", err)
			}
			if err == nil && (cfg.Bucket != "journal-lesson-attachments" || cfg.Secure != (test.name == "https")) {
				t.Fatal("incorrect defaults or TLS")
			}
		})
	}
}
