package storage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMinioObjectLifecycleThroughSignedS3Requests(t *testing.T) {
	var mu sync.Mutex
	bucketExists := false
	var data []byte
	signed := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 ") {
			signed = false
		}
		if r.URL.Query().Has("location") {
			w.Header().Set("Content-Type", "application/xml")
			_, _ = io.WriteString(w, `<LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-east-1</LocationConstraint>`)
			return
		}
		if r.URL.Path == "/journal-test/" || r.URL.Path == "/journal-test" {
			if r.Method == "PUT" {
				bucketExists = true
				w.WriteHeader(200)
				return
			}
			if !bucketExists {
				w.WriteHeader(404)
				return
			}
			w.WriteHeader(200)
			return
		}
		switch r.Method {
		case "PUT":
			var reader io.Reader = r.Body
			if strings.Contains(r.Header.Get("Content-Encoding"), "aws-chunked") || strings.HasPrefix(r.Header.Get("X-Amz-Content-Sha256"), "STREAMING-") {
				reader = httputil.NewChunkedReader(reader)
			}
			data, _ = io.ReadAll(reader)
			w.Header().Set("ETag", `"900150983cd24fb0d6963f7d28e17f72"`)
		case "HEAD", "GET":
			if data == nil {
				w.WriteHeader(404)
				return
			}
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			w.Header().Set("ETag", `"900150983cd24fb0d6963f7d28e17f72"`)
			if r.Method == "GET" {
				_, _ = w.Write(data)
			}
		case "DELETE":
			data = nil
			w.WriteHeader(204)
		default:
			w.WriteHeader(405)
		}
	}))
	defer server.Close()
	store, err := NewMinio(strings.TrimPrefix(server.URL, "http://"), "test-key", "test-secret", "journal-test", false)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = store.EnsureBucket(ctx); err != nil {
		t.Fatal(err)
	}
	if err = store.EnsureBucket(ctx); err != nil {
		t.Fatal(err)
	}
	if err = store.Put(ctx, "lessons/test/file", "text/plain", bytes.NewReader([]byte("abc")), 3); err != nil {
		t.Fatal(err)
	}
	content, err := store.Open(ctx, "lessons/test/file")
	if err != nil {
		t.Fatal(err)
	}
	got, readErr := io.ReadAll(content)
	_ = content.Close()
	if readErr != nil || string(got) != "abc" {
		t.Fatalf("download=%q err=%v", got, readErr)
	}
	if err = store.Delete(ctx, "lessons/test/file"); err != nil {
		t.Fatal(err)
	}
	if content, err = store.Open(ctx, "lessons/test/file"); err == nil {
		_ = content.Close()
		t.Fatal("deleted object opened")
	}
	mu.Lock()
	defer mu.Unlock()
	if !signed || !bucketExists {
		t.Fatal("S3 authentication or bucket creation was skipped")
	}
}
