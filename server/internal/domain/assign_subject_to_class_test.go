package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func TestAssignSubjectToClassCreatesAssignmentAndPermissions(t *testing.T) {
	classID, subjectID, teacherID := uuid.New(), uuid.New(), uuid.New()
	permissions := []entity.Permission{}
	links := []entity.UserPermission{}
	repo := &classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID}}, subjects: []entity.Subject{{ID: subjectID}}, users: []entity.User{{ID: teacherID, Roles: []string{"teacher"}}}, permissions: &permissions, userPermissions: &links}
	result, err := NewClassDomain(repo).AssignSubjectToClass(context.Background(), testTx{}, "hash", classID, subjectID, teacherID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Assignment.ID == uuid.Nil || len(permissions) != 8 || len(links) != 8 {
		t.Fatalf("assignment=%+v permissions=%d links=%d", result.Assignment, len(permissions), len(links))
	}
	if permissions[0].Code != "can_view_class" || *permissions[0].Value != classID.String() {
		t.Fatalf("class permission = %+v", permissions[0])
	}
}

func TestAssignSubjectToClassValidatesTeacherAndConflict(t *testing.T) {
	classID, subjectID, teacherID := uuid.New(), uuid.New(), uuid.New()
	repo := &classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID}}, subjects: []entity.Subject{{ID: subjectID}}, users: []entity.User{{ID: teacherID, Roles: []string{"parent"}}}}
	_, err := NewClassDomain(repo).AssignSubjectToClass(context.Background(), testTx{}, "hash", classID, subjectID, teacherID)
	if !errors.Is(err, ErrInvalidTeacher) {
		t.Fatalf("teacher error = %v", err)
	}
	repo.users[0].Roles = []string{"teacher"}
	repo.createClassSubjectErr = repository.ErrConflict
	_, err = NewClassDomain(repo).AssignSubjectToClass(context.Background(), testTx{}, "hash", classID, subjectID, teacherID)
	if !errors.Is(err, ErrClassSubjectExists) {
		t.Fatalf("conflict error = %v", err)
	}
}
