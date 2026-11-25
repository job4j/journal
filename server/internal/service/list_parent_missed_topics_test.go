package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type missedTopicsRepoStub struct {
	repository.ClassReadRepository
	child, class, year uuid.UUID
	hash               string
	allowed            bool
}

func (s *missedTopicsRepoStub) FindActiveUserBySessionHash(_ context.Context, _ repository.Transaction, hash string) (entity.User, error) {
	s.hash = hash
	return entity.User{ID: uuid.New(), Roles: []string{"parent"}}, nil
}
func (s *missedTopicsRepoStub) CheckSessionPermissionForValue(_ context.Context, _ repository.Transaction, _ string, _ string, value string) (bool, bool, error) {
	return true, s.allowed && value == s.child.String(), nil
}
func (s *missedTopicsRepoStub) GetUser(context.Context, repository.Transaction, uuid.UUID) (entity.User, error) {
	return entity.User{ID: s.child, Roles: []string{"student"}}, nil
}
func (s *missedTopicsRepoStub) GetAcademicYear(context.Context, repository.Transaction, uuid.UUID) (entity.AcademicYear, error) {
	return entity.AcademicYear{ID: s.year}, nil
}
func (s *missedTopicsRepoStub) ListClasses(context.Context, repository.Transaction) ([]entity.Class, error) {
	return []entity.Class{{ID: s.class, AcademicYearID: s.year}}, nil
}
func (s *missedTopicsRepoStub) ListClassStudents(context.Context, repository.Transaction) ([]entity.ClassStudent, error) {
	return []entity.ClassStudent{{ClassID: s.class, UserID: s.child}}, nil
}
func (s *missedTopicsRepoStub) ListClassSubjects(context.Context, repository.Transaction) ([]entity.ClassSubject, error) {
	return nil, nil
}
func (s *missedTopicsRepoStub) ListSubjects(context.Context, repository.Transaction) ([]entity.Subject, error) {
	return nil, nil
}
func (s *missedTopicsRepoStub) ListLessons(context.Context, repository.Transaction) ([]entity.Lesson, error) {
	return nil, nil
}
func (s *missedTopicsRepoStub) ListAbsences(context.Context, repository.Transaction) ([]entity.Absence, error) {
	return nil, nil
}
func TestParentMissedTopicsTransactionAndTokenPrivacy(t *testing.T) {
	for _, allowed := range []bool{true, false} {
		tx := &roleTransaction{}
		repo := &missedTopicsRepoStub{child: uuid.New(), class: uuid.New(), year: uuid.New(), allowed: allowed}
		result, err := NewClassService(roleTransactionManager{tx: tx}, repo).ListParentMissedTopics(context.Background(), "private-session", domain.ListParentMissedTopicsRequest{StudentID: repo.child, AcademicYearID: repo.year})
		if allowed {
			if err != nil || !tx.committed || result.Items == nil {
				t.Fatalf("read did not commit: %+v %v", result, err)
			}
		} else if !errors.Is(err, domain.ErrForbidden) || tx.committed {
			t.Fatalf("forbidden read committed: %v", err)
		}
		if repo.hash != sessionTokenHash("private-session") || !tx.rolledBack {
			t.Fatal("session hashing or transaction cleanup missing")
		}
	}
}
