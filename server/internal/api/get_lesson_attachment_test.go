package api

import (
	"github.com/google/uuid"
	"io"
	"journal/server/internal/domain"
	"mime"
	"net/http/httptest"
	"testing"
)

func TestDownloadParentAttachmentUsesSafeHeadersAndChildContext(t *testing.T) {
	s := &attachmentServiceStub{}
	student := uuid.New()
	res, err := attachmentApp(s).Test(httptest.NewRequest("GET", "/parent/students/"+student.String()+"/lesson-attachments/"+attachmentID.String()+"/download", nil))
	if err != nil {
		t.Fatal(err)
	}
	content, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	disposition, params, err := mime.ParseMediaType(res.Header.Get("Content-Disposition"))
	if err != nil || disposition != "attachment" || params["filename"] != "Задание.txt" {
		t.Fatalf("invalid filename: %s", res.Header.Get("Content-Disposition"))
	}
	if res.StatusCode != 200 || string(content) != "abc" || res.Header.Get("X-Content-Type-Options") != "nosniff" || res.Header.Get("Cache-Control") != "private, no-store" || res.Header.Get("Content-Type") != "application/octet-stream" {
		t.Fatalf("unsafe download: %+v %s", res.Header, content)
	}
	if s.get.StudentID != student || s.get.AttachmentID != attachmentID {
		t.Fatal("lost parent context")
	}
}
func TestDownloadAttachmentRejectsMissingAndForbidden(t *testing.T) {
	for _, test := range []struct {
		err    error
		status int
	}{{domain.ErrAttachmentNotFound, 404}, {domain.ErrForbidden, 403}, {domain.ErrUnauthenticated, 401}, {domain.ErrAttachmentStorageUnavailable, 503}} {
		res, err := attachmentApp(&attachmentServiceStub{err: test.err}).Test(httptest.NewRequest("GET", "/lesson-attachments/"+attachmentID.String()+"/download", nil))
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != test.status {
			t.Fatalf("status=%d", res.StatusCode)
		}
	}
}
