package domain

import (
	"context"
	"errors"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func TestStudentAbsenceIsIdempotentAndRecordsTeacher(t *testing.T) {
	stored := entity.Absence{}
	repo, _, studentID := scoreFixture("five_point")
	repo.changedAbsence = &stored
	lessonID := repo.lessons[0].ID
	domain := NewClassDomain(repo)
	first, err := domain.PutStudentAbsence(context.Background(), testTx{}, "hash", lessonID, studentID)
	if err != nil {
		t.Fatal(err)
	}
	second, err := domain.PutStudentAbsence(context.Background(), testTx{}, "hash", lessonID, studentID)
	if err != nil {
		t.Fatal(err)
	}
	if first.UserID != second.UserID || stored.RecordedBy != repo.currentUser.ID {
		t.Fatalf("absence = %+v", stored)
	}
	if err = domain.DeleteStudentAbsence(context.Background(), testTx{}, "hash", lessonID, studentID); err != nil {
		t.Fatal(err)
	}
}
func TestStudentAbsenceChecksMembershipDate(t *testing.T) {
	repo, _, studentID := scoreFixture("five_point")
	left := time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC)
	repo.students[0].LeftOn = &left
	_, err := NewClassDomain(repo).PutStudentAbsence(context.Background(), testTx{}, "hash", repo.lessons[0].ID, studentID)
	if !errors.Is(err, ErrClassStudentNotFound) {
		t.Fatalf("error = %v", err)
	}
}
func TestStudentAbsenceChecksResponsibleTeacher(t *testing.T) {
	repo, _, studentID := scoreFixture("five_point")
	repo.classSubjects[0].ResponsibleTeacherID = studentID
	_, err := NewClassDomain(repo).PutStudentAbsence(context.Background(), testTx{}, "hash", repo.lessons[0].ID, studentID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
