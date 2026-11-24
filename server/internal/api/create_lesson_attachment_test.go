package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadAttachmentMultipartAndDetectedContentType(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "Задание.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(part, "abc")
	_ = writer.Close()
	s := &attachmentServiceStub{}
	req := httptest.NewRequest("POST", "/lessons/"+attachmentLessonID.String()+"/attachments", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := attachmentApp(s).Test(req)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != 201 || s.content != "abc" || s.upload.SizeBytes != 3 || !strings.HasPrefix(s.upload.ContentType, "text/plain") {
		t.Fatalf("upload failed: %d %s %+v", res.StatusCode, data, s.upload)
	}
	if strings.Contains(string(data), "secret-storage-key") {
		t.Fatal("object key leaked")
	}
}
func TestUploadAttachmentRejectsEmptyAndMissingFiles(t *testing.T) {
	for _, field := range []string{"file", "wrong-field"} {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_, _ = writer.CreateFormFile(field, "empty.txt")
		_ = writer.Close()
		s := &attachmentServiceStub{}
		req := httptest.NewRequest("POST", "/lessons/"+attachmentLessonID.String()+"/attachments", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := attachmentApp(s).Test(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != 400 || s.upload.SizeBytes != 0 {
			t.Fatalf("invalid upload accepted: %d", res.StatusCode)
		}
	}
}
