package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func TestDeleteStudentScoreChecksUpdatePermissionAndDeletesByStudent(t *testing.T) {
	repo, itemID, studentID := scoreFixture("five_point")
	repo.objectAllowed = map[string]bool{"can_update_score": true}
	deleted := entity.Score{}
	repo.deletedScore = &deleted
	if err := NewClassDomain(repo).DeleteStudentScore(context.Background(), testTx{}, "hash", itemID, studentID); err != nil {
		t.Fatal(err)
	}
	if deleted.GradeItemID != itemID || deleted.UserID != studentID {
		t.Fatalf("deleted = %+v", deleted)
	}
}

func TestDeleteStudentScoreRejectsOtherTeacher(t *testing.T) {
	repo, itemID, studentID := scoreFixture("five_point")
	repo.classSubjects[0].ResponsibleTeacherID = uuid.New()
	err := NewClassDomain(repo).DeleteStudentScore(context.Background(), testTx{}, "hash", itemID, studentID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
