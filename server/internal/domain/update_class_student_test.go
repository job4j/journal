package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func TestUpdateClassStudentSetsAndClearsLeftOn(t *testing.T) {
	classID, studentID := uuid.New(), uuid.New()
	enrolled := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	left := enrolled.AddDate(0, 1, 0)
	repo := classRepoStub{authenticated: true, allowed: true, students: []entity.ClassStudent{{ClassID: classID, UserID: studentID, EnrolledOn: enrolled}}, users: []entity.User{{ID: studentID}}}
	domain := NewClassDomain(repo)

	result, err := domain.UpdateClassStudent(context.Background(), testTx{}, "hash", classID, studentID, &left)
	if err != nil || result.Membership.LeftOn == nil || !result.Membership.LeftOn.Equal(left) {
		t.Fatalf("result = %+v, error = %v", result, err)
	}
	result, err = domain.UpdateClassStudent(context.Background(), testTx{}, "hash", classID, studentID, nil)
	if err != nil || result.Membership.LeftOn != nil {
		t.Fatalf("result = %+v, error = %v", result, err)
	}
}

func TestUpdateClassStudentRejectsEarlyDate(t *testing.T) {
	classID, studentID := uuid.New(), uuid.New()
	enrolled := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	left := enrolled.AddDate(0, 0, -1)
	repo := classRepoStub{authenticated: true, allowed: true, students: []entity.ClassStudent{{ClassID: classID, UserID: studentID, EnrolledOn: enrolled}}}
	_, err := NewClassDomain(repo).UpdateClassStudent(context.Background(), testTx{}, "hash", classID, studentID, &left)
	if !errors.Is(err, ErrInvalidClassStudent) {
		t.Fatalf("error = %v", err)
	}
}

func TestUpdateClassStudentRequiresMembershipAndAccess(t *testing.T) {
	classID, studentID := uuid.New(), uuid.New()
	_, err := NewClassDomain(classRepoStub{authenticated: true, allowed: true}).UpdateClassStudent(context.Background(), testTx{}, "hash", classID, studentID, nil)
	if !errors.Is(err, ErrClassStudentNotFound) {
		t.Fatalf("missing error = %v", err)
	}
	_, err = NewClassDomain(classRepoStub{authenticated: true}).UpdateClassStudent(context.Background(), testTx{}, "hash", classID, studentID, nil)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("access error = %v", err)
	}
}
