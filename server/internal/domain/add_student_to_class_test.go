package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func TestAddStudentToClass(t *testing.T) {
	classID, yearID, studentID := uuid.New(), uuid.New(), uuid.New()
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	repo := classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}}, year: entity.AcademicYear{ID: yearID, StartsOn: start, EndsOn: start.AddDate(1, 0, 0)}, users: []entity.User{{ID: studentID, Roles: []string{"student"}}}}

	result, err := NewClassDomain(repo).AddStudentToClass(context.Background(), testTx{}, "hash", classID, studentID, start)
	if err != nil {
		t.Fatal(err)
	}
	if result.Membership.ClassID != classID || result.Membership.UserID != studentID {
		t.Fatalf("membership = %+v", result.Membership)
	}
}

func TestAddStudentToClassRejectsRoleAndDate(t *testing.T) {
	classID, yearID, userID := uuid.New(), uuid.New(), uuid.New()
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	base := classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}}, year: entity.AcademicYear{ID: yearID, StartsOn: start, EndsOn: start.AddDate(1, 0, 0)}, users: []entity.User{{ID: userID, Roles: []string{"parent"}}}}
	if _, err := NewClassDomain(base).AddStudentToClass(context.Background(), testTx{}, "hash", classID, userID, start); !errors.Is(err, ErrInvalidClassStudent) {
		t.Fatalf("role error = %v", err)
	}
	base.users[0].Roles = []string{"student"}
	if _, err := NewClassDomain(base).AddStudentToClass(context.Background(), testTx{}, "hash", classID, userID, start.AddDate(0, 0, -1)); !errors.Is(err, ErrInvalidClassStudent) {
		t.Fatalf("date error = %v", err)
	}
}

func TestAddStudentToClassMapsConflict(t *testing.T) {
	classID, yearID, studentID := uuid.New(), uuid.New(), uuid.New()
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	repo := classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}}, year: entity.AcademicYear{ID: yearID, StartsOn: start, EndsOn: start.AddDate(1, 0, 0)}, users: []entity.User{{ID: studentID, Roles: []string{"student"}}}, createStudentErr: repository.ErrConflict}
	_, err := NewClassDomain(repo).AddStudentToClass(context.Background(), testTx{}, "hash", classID, studentID, start)
	if !errors.Is(err, ErrClassStudentExists) {
		t.Fatalf("error = %v", err)
	}
}
