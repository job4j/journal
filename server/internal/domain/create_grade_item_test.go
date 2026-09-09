package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
)

func gradeItemRepo() (classRepoStub, uuid.UUID) {
	teacherID, assignmentID, lessonID := uuid.New(), uuid.New(), uuid.New()
	return classRepoStub{authenticated: true, allowed: true, currentUser: entity.User{ID: teacherID, Roles: []string{"teacher"}}, classSubjects: []entity.ClassSubject{{ID: assignmentID, ResponsibleTeacherID: teacherID}}, lessons: []entity.Lesson{{ID: lessonID, ClassSubjectID: assignmentID}}}, lessonID
}
func TestCreateGradeItemValidatesScales(t *testing.T) {
	repo, lessonID := gradeItemRepo()
	max := float64(20)
	item, err := NewClassDomain(repo).CreateGradeItem(context.Background(), testTx{}, "hash", lessonID, "Баллы", "homework", "points", &max)
	if err != nil || item.MaxScore == nil || *item.MaxScore != 20 {
		t.Fatalf("item=%+v error=%v", item, err)
	}
	_, err = NewClassDomain(repo).CreateGradeItem(context.Background(), testTx{}, "hash", lessonID, "Зачёт", "other", "pass_fail", &max)
	if !errors.Is(err, ErrInvalidGradeItem) {
		t.Fatalf("error=%v", err)
	}
}
