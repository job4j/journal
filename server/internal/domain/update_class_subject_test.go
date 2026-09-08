package domain

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func TestUpdateClassSubjectMovesPermissions(t *testing.T) {
	classID, assignmentID, subjectID, oldTeacherID, newTeacherID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	value := assignmentID.String()
	classValue := classID.String()
	permissions := []entity.Permission{{ID: uuid.New(), Code: "can_view_class_subject", Value: &value}, {ID: uuid.New(), Code: "can_view_class", Value: &classValue}}
	granted := []entity.UserPermission{}
	deleted := []entity.UserPermission{}
	repo := classRepoStub{authenticated: true, allowed: true, classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, SubjectID: subjectID, ResponsibleTeacherID: oldTeacherID}}, subjects: []entity.Subject{{ID: subjectID}}, users: []entity.User{{ID: newTeacherID, Roles: []string{"teacher"}}}, permissions: &permissions, userPermissions: &granted, deletedUserPermissions: &deleted}
	result, err := NewClassDomain(repo).UpdateClassSubject(context.Background(), testTx{}, "hash", assignmentID, newTeacherID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Assignment.ResponsibleTeacherID != newTeacherID || len(granted) != 8 || len(deleted) != 8 {
		t.Fatalf("result=%+v granted=%d deleted=%d", result, len(granted), len(deleted))
	}
}

func TestUpdateClassSubjectKeepsSharedClassPermission(t *testing.T) {
	classID, assignmentID, otherID, subjectID, oldTeacherID, newTeacherID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	classValue := classID.String()
	classPermissionID := uuid.New()
	permissions := []entity.Permission{{ID: classPermissionID, Code: "can_view_class", Value: &classValue}}
	deleted := []entity.UserPermission{}
	repo := classRepoStub{authenticated: true, allowed: true, classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, SubjectID: subjectID, ResponsibleTeacherID: oldTeacherID}, {ID: otherID, ClassID: classID, ResponsibleTeacherID: oldTeacherID}}, subjects: []entity.Subject{{ID: subjectID}}, users: []entity.User{{ID: newTeacherID, Roles: []string{"teacher"}}}, permissions: &permissions, deletedUserPermissions: &deleted}
	_, err := NewClassDomain(repo).UpdateClassSubject(context.Background(), testTx{}, "hash", assignmentID, newTeacherID)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range deleted {
		if item.PermissionID == classPermissionID {
			t.Fatalf("shared class permission was deleted: %+v", deleted)
		}
	}
}
