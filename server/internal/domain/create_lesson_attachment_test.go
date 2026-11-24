package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

type attachmentRepoStub struct {
	repository.AttachmentRepository
	viewer            entity.User
	lesson            entity.Lesson
	assignment        entity.ClassSubject
	membership        entity.ClassStudent
	attachment        entity.LessonAttachment
	denied            string
	missingMembership bool
	created, deleted  bool
}

func attachmentFixture() *attachmentRepoStub {
	teacher, lessonID, assignmentID, classID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	date := time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)
	return &attachmentRepoStub{viewer: entity.User{ID: teacher, Roles: []string{"teacher"}}, lesson: entity.Lesson{ID: lessonID, ClassSubjectID: assignmentID, LessonDate: date, Topic: "Тема"}, assignment: entity.ClassSubject{ID: assignmentID, ClassID: classID, ResponsibleTeacherID: teacher}, membership: entity.ClassStudent{ClassID: classID, EnrolledOn: date}, attachment: entity.LessonAttachment{ID: uuid.New(), LessonID: lessonID, ObjectKey: "private-key", FileName: "Задание.pdf", SizeBytes: 3}}
}
func (r *attachmentRepoStub) FindActiveUserBySessionHash(context.Context, repository.Transaction, string) (entity.User, error) {
	if r.viewer.ID == uuid.Nil {
		return entity.User{}, repository.ErrNotFound
	}
	return r.viewer, nil
}
func (r *attachmentRepoStub) CheckSessionPermissionForValue(_ context.Context, _ repository.Transaction, _ string, code, value string) (bool, bool, error) {
	return true, code != r.denied, nil
}
func (r *attachmentRepoStub) GetLesson(context.Context, repository.Transaction, uuid.UUID) (entity.Lesson, error) {
	return r.lesson, nil
}
func (r *attachmentRepoStub) GetClassSubject(context.Context, repository.Transaction, uuid.UUID) (entity.ClassSubject, error) {
	return r.assignment, nil
}
func (r *attachmentRepoStub) GetUser(context.Context, repository.Transaction, uuid.UUID) (entity.User, error) {
	return entity.User{Roles: []string{"student"}}, nil
}
func (r *attachmentRepoStub) GetClassStudent(context.Context, repository.Transaction, uuid.UUID, uuid.UUID) (entity.ClassStudent, error) {
	if r.missingMembership {
		return entity.ClassStudent{}, repository.ErrNotFound
	}
	return r.membership, nil
}
func (r *attachmentRepoStub) CreateLessonAttachment(_ context.Context, _ repository.Transaction, item entity.LessonAttachment) (entity.LessonAttachment, error) {
	r.created = true
	r.attachment = item
	return item, nil
}
func (r *attachmentRepoStub) ListLessonAttachments(context.Context, repository.Transaction, uuid.UUID) ([]entity.LessonAttachment, error) {
	return []entity.LessonAttachment{r.attachment}, nil
}
func (r *attachmentRepoStub) GetLessonAttachment(context.Context, repository.Transaction, uuid.UUID) (entity.LessonAttachment, error) {
	return r.attachment, nil
}
func (r *attachmentRepoStub) DeleteLessonAttachment(context.Context, repository.Transaction, uuid.UUID) error {
	r.deleted = true
	return nil
}

func TestCreateLessonAttachmentValidatesAccessAndMetadata(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*attachmentRepoStub, *CreateLessonAttachmentRequest)
		want   error
	}{
		{"teacher", func(*attachmentRepoStub, *CreateLessonAttachmentRequest) {}, nil},
		{"administrator", func(r *attachmentRepoStub, _ *CreateLessonAttachmentRequest) {
			r.viewer = entity.User{ID: uuid.New(), Roles: []string{"admin"}}
		}, nil},
		{"expired session", func(r *attachmentRepoStub, _ *CreateLessonAttachmentRequest) { r.viewer.ID = uuid.Nil }, ErrUnauthenticated},
		{"another teacher", func(r *attachmentRepoStub, _ *CreateLessonAttachmentRequest) { r.viewer.ID = uuid.New() }, ErrForbidden},
		{"parent cannot upload", func(r *attachmentRepoStub, _ *CreateLessonAttachmentRequest) { r.viewer.Roles = []string{"parent"} }, ErrForbidden},
		{"missing permission", func(r *attachmentRepoStub, _ *CreateLessonAttachmentRequest) { r.denied = "can_update_lesson" }, ErrForbidden},
		{"missing topic", func(r *attachmentRepoStub, _ *CreateLessonAttachmentRequest) { r.lesson.Topic = "  " }, ErrInvalidLesson},
		{"empty file", func(_ *attachmentRepoStub, q *CreateLessonAttachmentRequest) { q.SizeBytes = 0 }, ErrInvalidAttachment},
		{"too large", func(_ *attachmentRepoStub, q *CreateLessonAttachmentRequest) { q.SizeBytes = MaxAttachmentBytes + 1 }, ErrInvalidAttachment},
		{"control character in filename", func(_ *attachmentRepoStub, q *CreateLessonAttachmentRequest) { q.FileName = "a\nb.pdf" }, ErrInvalidAttachment},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := attachmentFixture()
			req := CreateLessonAttachmentRequest{LessonID: r.lesson.ID, FileName: "C:\\fakepath\\Задание.pdf", ContentType: "application/pdf", SizeBytes: 3}
			test.change(r, &req)
			result, err := NewAttachmentDomain(r).CreateLessonAttachment(context.Background(), testTx{}, "hash", req)
			if !errors.Is(err, test.want) {
				t.Fatalf("got %v want %v", err, test.want)
			}
			if test.want != nil && r.created {
				t.Fatal("unauthorized or invalid attachment was created")
			}
			if err == nil && (result.Attachment.FileName != "Задание.pdf" || result.Attachment.CreatedBy != r.viewer.ID || result.Attachment.ObjectKey == req.FileName) {
				t.Fatalf("unsafe metadata: %+v", result)
			}
		})
	}
}
