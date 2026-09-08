package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type subjectRepoStub struct {
	authenticated, allowed bool
	items                  []entity.Subject
	createErr              error
}

func (s subjectRepoStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return s.authenticated, s.allowed, nil
}
func (s subjectRepoStub) ListSubjects(context.Context, repository.Transaction) ([]entity.Subject, error) {
	return s.items, nil
}
func (s subjectRepoStub) CreateSubject(_ context.Context, _ repository.Transaction, item entity.Subject) (entity.Subject, error) {
	item.ID = uuid.New()
	return item, s.createErr
}
func TestCreateSubjectNormalizesInput(t *testing.T) {
	item, err := NewSubjectDomain(subjectRepoStub{authenticated: true, allowed: true}).CreateSubject(context.Background(), testTx{}, "hash", " MATH ", " Математика ")
	if err != nil {
		t.Fatal(err)
	}
	if item.Code != "math" || item.Name != "Математика" {
		t.Fatalf("item = %+v", item)
	}
}
func TestCreateSubjectRejectsInvalidCode(t *testing.T) {
	_, err := NewSubjectDomain(subjectRepoStub{authenticated: true, allowed: true}).CreateSubject(context.Background(), testTx{}, "hash", "Bad Code", "Предмет")
	if !errors.Is(err, ErrInvalidSubject) {
		t.Fatalf("error = %v", err)
	}
}
func TestListSubjectsChecksAccess(t *testing.T) {
	_, err := NewSubjectDomain(subjectRepoStub{}).ListSubjects(context.Background(), testTx{}, "hash")
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("error = %v", err)
	}
}
