package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
)

func TestListTeacherClassesUsesSessionTeacher(t *testing.T) {
	teacherID, classID := uuid.New(), uuid.New()
	repo := classRepoStub{authenticated: true, allowed: true, currentUser: entity.User{ID: teacherID, Roles: []string{"teacher"}}, classes: []entity.Class{{ID: classID}, {ID: uuid.New()}}, classSubjects: []entity.ClassSubject{{ID: uuid.New(), ClassID: classID, ResponsibleTeacherID: teacherID}, {ID: uuid.New(), ClassID: uuid.New(), ResponsibleTeacherID: uuid.New()}}}
	items, err := NewClassDomain(repo).ListTeacherClasses(context.Background(), testTx{}, "hash")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Class.ID != classID {
		t.Fatalf("items=%+v", items)
	}
}
func TestListTeacherClassesRejectsNonTeacher(t *testing.T) {
	repo := classRepoStub{currentUser: entity.User{ID: uuid.New(), Roles: []string{"parent"}}}
	_, err := NewClassDomain(repo).ListTeacherClasses(context.Background(), testTx{}, "hash")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error=%v", err)
	}
}
