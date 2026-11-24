package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"io"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
	"testing"
)

type attachmentTransaction struct {
	committed bool
	commitErr error
}

func (t *attachmentTransaction) Commit(context.Context) error {
	if t.commitErr == nil {
		t.committed = true
	}
	return t.commitErr
}
func (t *attachmentTransaction) Rollback(context.Context) error { return nil }

type attachmentManager struct{ tx *attachmentTransaction }

func (m attachmentManager) Begin(context.Context) (repository.Transaction, error) { return m.tx, nil }

type attachmentServiceRepo struct {
	repository.AttachmentRepository
	teacher, lessonID, assignmentID uuid.UUID
	hash                            string
	denied                          bool
}

func (r *attachmentServiceRepo) FindActiveUserBySessionHash(_ context.Context, _ repository.Transaction, hash string) (entity.User, error) {
	r.hash = hash
	return entity.User{ID: r.teacher, Roles: []string{"teacher"}}, nil
}
func (r *attachmentServiceRepo) GetLesson(context.Context, repository.Transaction, uuid.UUID) (entity.Lesson, error) {
	return entity.Lesson{ID: r.lessonID, ClassSubjectID: r.assignmentID, Topic: "Тема"}, nil
}
func (r *attachmentServiceRepo) GetClassSubject(context.Context, repository.Transaction, uuid.UUID) (entity.ClassSubject, error) {
	return entity.ClassSubject{ID: r.assignmentID, ResponsibleTeacherID: r.teacher}, nil
}
func (r *attachmentServiceRepo) CheckSessionPermissionForValue(context.Context, repository.Transaction, string, string, string) (bool, bool, error) {
	return true, !r.denied, nil
}
func (r *attachmentServiceRepo) CreateLessonAttachment(_ context.Context, _ repository.Transaction, item entity.LessonAttachment) (entity.LessonAttachment, error) {
	return item, nil
}
func (r *attachmentServiceRepo) GetLessonAttachment(context.Context, repository.Transaction, uuid.UUID) (entity.LessonAttachment, error) {
	return entity.LessonAttachment{ID: uuid.New(), LessonID: r.lessonID, ObjectKey: "private-key", SizeBytes: 3}, nil
}
func (r *attachmentServiceRepo) DeleteLessonAttachment(context.Context, repository.Transaction, uuid.UUID) error {
	return nil
}

type attachmentStoreStub struct {
	tx                   *attachmentTransaction
	putErr, deleteErr    error
	put, removed, opened bool
	deletedAfterCommit   bool
	content              string
}

func (s *attachmentStoreStub) Put(_ context.Context, _ string, _ string, r io.Reader, _ int64) error {
	s.put = true
	b, _ := io.ReadAll(r)
	s.content = string(b)
	return s.putErr
}
func (s *attachmentStoreStub) Delete(context.Context, string) error {
	s.removed = true
	s.deletedAfterCommit = s.tx.committed
	return s.deleteErr
}
func (s *attachmentStoreStub) Open(context.Context, string) (io.ReadCloser, error) {
	s.opened = true
	return io.NopCloser(strings.NewReader("abc")), nil
}
func attachmentServiceFixture() (*attachmentTransaction, *attachmentServiceRepo, *attachmentStoreStub) {
	tx := &attachmentTransaction{}
	repo := &attachmentServiceRepo{teacher: uuid.New(), lessonID: uuid.New(), assignmentID: uuid.New()}
	return tx, repo, &attachmentStoreStub{tx: tx}
}

func TestAttachmentUploadTransactionAndCompensation(t *testing.T) {
	for _, test := range []string{"success", "upload failure", "commit failure", "forbidden", "disabled storage"} {
		t.Run(test, func(t *testing.T) {
			tx, repo, store := attachmentServiceFixture()
			if test == "upload failure" {
				store.putErr = errors.New("storage offline")
			}
			if test == "commit failure" {
				tx.commitErr = errors.New("commit failed")
			}
			repo.denied = test == "forbidden"
			var objectStore AttachmentObjectStore = store
			if test == "disabled storage" {
				objectStore = nil
			}
			_, err := NewAttachmentService(attachmentManager{tx}, repo, objectStore).CreateLessonAttachment(context.Background(), "session-secret", domain.CreateLessonAttachmentRequest{LessonID: repo.lessonID, FileName: "file.txt", ContentType: "text/plain", SizeBytes: 3}, strings.NewReader("abc"))
			if test == "success" {
				if err != nil || !tx.committed || !store.put || store.removed || store.content != "abc" {
					t.Fatalf("upload failed: %v %+v", err, store)
				}
			} else {
				if err == nil || tx.committed {
					t.Fatal("failed upload committed")
				}
				if test == "forbidden" || test == "disabled storage" {
					if store.put || store.removed {
						t.Fatal("storage touched without authorization/config")
					}
				} else if !store.removed {
					t.Fatal("failed upload left an object")
				}
			}
			if repo.hash == "session-secret" || repo.hash == "" {
				t.Fatal("session was not hashed")
			}
		})
	}
}
