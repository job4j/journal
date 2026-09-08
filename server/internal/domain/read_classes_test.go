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
	users                  []entity.User
	getErr                 error
	year                   entity.AcademicYear
	createdStudent         entity.ClassStudent
	createStudentErr       error
	updateStudentErr       error
}

func (s classRepoStub) GetAcademicYear(context.Context, repository.Transaction, uuid.UUID) (entity.AcademicYear, error) {
	return s.year, s.getErr
}

func (s classRepoStub) GetUser(_ context.Context, _ repository.Transaction, id uuid.UUID) (entity.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, s.getErr
		}
	}
	return entity.User{}, repository.ErrNotFound
}
func (s classRepoStub) CreateClassStudent(_ context.Context, _ repository.Transaction, item entity.ClassStudent) (entity.ClassStudent, error) {
	if s.createdStudent.ClassID != uuid.Nil {
		return s.createdStudent, s.createStudentErr
	}
	return item, s.createStudentErr
}
func (s classRepoStub) GetClassStudent(_ context.Context, _ repository.Transaction, classID, userID uuid.UUID) (entity.ClassStudent, error) {
	for _, item := range s.students {
		if item.ClassID == classID && item.UserID == userID {
			return item, s.getErr
		}
	}
	return entity.ClassStudent{}, repository.ErrNotFound
}
func (s classRepoStub) UpdateClassStudent(_ context.Context, _ repository.Transaction, item entity.ClassStudent) (entity.ClassStudent, error) {
	return item, s.updateStudentErr
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
func (s classRepoStub) ListUsers(context.Context, repository.Transaction) ([]entity.User, error) {
	return s.users, nil
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
