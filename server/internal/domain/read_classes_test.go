package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type classRepoStub struct {
	authenticated, allowed bool
	classes                []entity.Class
	students               []entity.ClassStudent
	getErr                 error
}

func (s classRepoStub) GetAcademicYear(context.Context, repository.Transaction, uuid.UUID) (entity.AcademicYear, error) {
	return entity.AcademicYear{ID: uuid.New()}, s.getErr
}
func (s classRepoStub) CreateClass(_ context.Context, _ repository.Transaction, item entity.Class) (entity.Class, error) {
	item.ID = uuid.New()
	return item, s.getErr
}

func (s classRepoStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return s.authenticated, s.allowed, nil
}
func (s classRepoStub) CheckSessionPermissionForValue(context.Context, repository.Transaction, string, string, string) (bool, bool, error) {
	return s.authenticated, s.allowed, nil
}
func (s classRepoStub) ListClasses(context.Context, repository.Transaction) ([]entity.Class, error) {
	return s.classes, nil
}
func (s classRepoStub) ListClassStudents(context.Context, repository.Transaction) ([]entity.ClassStudent, error) {
	return s.students, nil
}
func (s classRepoStub) GetClass(context.Context, repository.Transaction, uuid.UUID) (entity.Class, error) {
	if len(s.classes) == 0 {
		return entity.Class{}, s.getErr
	}
	return s.classes[0], s.getErr
}
func TestListClassesFiltersYearAndCountsActiveStudents(t *testing.T) {
	yearID, classID := uuid.New(), uuid.New()
	left := entity.ClassStudent{ClassID: classID}
	now := left.EnrolledOn
	left.LeftOn = &now
	repo := classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}, {ID: uuid.New(), AcademicYearID: uuid.New()}}, students: []entity.ClassStudent{{ClassID: classID}, left}}
	items, err := NewClassDomain(repo).ListClasses(context.Background(), testTx{}, "hash", yearID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].StudentCount != 1 {
		t.Fatalf("items = %+v", items)
	}
}
func TestGetClassChecksObjectAccess(t *testing.T) {
	_, err := NewClassDomain(classRepoStub{authenticated: true}).GetClass(context.Background(), testTx{}, "hash", uuid.New())
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
