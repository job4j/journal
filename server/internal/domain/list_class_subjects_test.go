package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func TestListClassSubjectsBuildsViews(t *testing.T) {
	classID, assignmentID, subjectID, teacherID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	repo := classRepoStub{authenticated: true, allowed: true, classes: []entity.Class{{ID: classID}}, classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, SubjectID: subjectID, ResponsibleTeacherID: teacherID}, {ID: uuid.New(), ClassID: uuid.New()}}, subjects: []entity.Subject{{ID: subjectID, Name: "Математика"}}, users: []entity.User{{ID: teacherID, Roles: []string{"teacher"}}}}
	items, err := NewClassDomain(repo).ListClassSubjects(context.Background(), testTx{}, "hash", classID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Subject.ID != subjectID || items[0].Teacher.ID != teacherID {
		t.Fatalf("items = %+v", items)
	}
}

func TestListClassSubjectsRequiresClassAccess(t *testing.T) {
	_, err := NewClassDomain(classRepoStub{authenticated: true}).ListClassSubjects(context.Background(), testTx{}, "hash", uuid.New())
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
