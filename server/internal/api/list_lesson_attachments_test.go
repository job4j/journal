package api

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

var attachmentLessonID = uuid.MustParse("10000000-0000-4000-8000-000000000001")
var attachmentID = uuid.MustParse("20000000-0000-4000-8000-000000000001")

type attachmentServiceStub struct {
	err     error
	upload  domain.CreateLessonAttachmentRequest
	list    domain.ListLessonAttachmentsRequest
	get     domain.GetLessonAttachmentRequest
	deleted uuid.UUID
	content string
}

func attachmentItem() entity.LessonAttachment {
	return entity.LessonAttachment{ID: attachmentID, LessonID: attachmentLessonID, FileName: "Задание.txt", ContentType: "text/plain", SizeBytes: 3, ObjectKey: "secret-storage-key"}
}
func (s *attachmentServiceStub) CreateLessonAttachment(_ context.Context, _ string, q domain.CreateLessonAttachmentRequest, r io.Reader) (domain.CreateLessonAttachmentResponse, error) {
	s.upload = q
	b, _ := io.ReadAll(r)
	s.content = string(b)
	return domain.CreateLessonAttachmentResponse{Attachment: attachmentItem()}, s.err
}
func (s *attachmentServiceStub) ListLessonAttachments(_ context.Context, _ string, q domain.ListLessonAttachmentsRequest) (domain.ListLessonAttachmentsResponse, error) {
	s.list = q
	return domain.ListLessonAttachmentsResponse{Items: []entity.LessonAttachment{attachmentItem()}}, s.err
}
func (s *attachmentServiceStub) GetLessonAttachment(_ context.Context, _ string, q domain.GetLessonAttachmentRequest) (domain.GetLessonAttachmentResponse, io.ReadCloser, error) {
	s.get = q
	if s.err != nil {
		return domain.GetLessonAttachmentResponse{}, nil, s.err
	}
	return domain.GetLessonAttachmentResponse{Attachment: attachmentItem()}, io.NopCloser(strings.NewReader("abc")), nil
}
func (s *attachmentServiceStub) DeleteLessonAttachment(_ context.Context, _ string, q domain.DeleteLessonAttachmentRequest) error {
	s.deleted = q.AttachmentID
	return s.err
}
func attachmentApp(s *attachmentServiceStub) *fiber.App {
	app := fiber.New(fiber.Config{BodyLimit: 24 << 20})
	NewAttachmentHandler(s, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	return app
}
func TestListParentLessonAttachmentsContract(t *testing.T) {
	student := uuid.New()
	for _, test := range []struct {
		err    error
		status int
	}{{nil, 200}, {domain.ErrUnauthenticated, 401}, {domain.ErrForbidden, 403}, {domain.ErrLessonNotFound, 404}, {errors.New("database error"), 500}} {
		s := &attachmentServiceStub{err: test.err}
		res, err := attachmentApp(s).Test(httptest.NewRequest("GET", "/parent/students/"+student.String()+"/lessons/"+attachmentLessonID.String()+"/attachments", nil))
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if res.StatusCode != test.status {
			t.Fatalf("status=%d body=%s", res.StatusCode, body)
		}
		if s.list.StudentID != student || s.list.LessonID != attachmentLessonID {
			t.Fatal("lost parent context")
		}
		if strings.Contains(string(body), "secret-storage-key") || strings.Contains(string(body), "ObjectKey") {
			t.Fatal("storage key leaked")
		}
	}
}
