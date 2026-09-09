package domain

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func TestListParentStudentsUsesOnlyOwnStudentPermissions(t *testing.T) {
	parentID, studentID, teacherID := uuid.New(), uuid.New(), uuid.New()
	permissionID, foreignPermissionID := uuid.New(), uuid.New()
	studentValue, teacherValue := studentID.String(), teacherID.String()
	repo := &userRepoStub{authenticated: true, currentUser: entity.User{ID: parentID, Roles: []string{"parent"}}, users: map[uuid.UUID]entity.User{studentID: {ID: studentID, Roles: []string{"student"}}, teacherID: {ID: teacherID, Roles: []string{"teacher"}}}, permissions: []entity.Permission{{ID: permissionID, Code: "can_view_user", Value: &studentValue}, {ID: foreignPermissionID, Code: "can_view_user", Value: &teacherValue}}, userPermissions: []entity.UserPermission{{UserID: parentID, PermissionID: permissionID}, {UserID: parentID, PermissionID: foreignPermissionID}}}
	items, err := NewUserDomain(repo).ListParentStudents(context.Background(), testTx{}, "hash")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != studentID {
		t.Fatalf("items = %+v", items)
	}
}
func TestListParentStudentsSupportsGlobalPermission(t *testing.T) {
	parentID, studentID := uuid.New(), uuid.New()
	repo := &userRepoStub{authenticated: true, allowed: true, currentUser: entity.User{ID: parentID, Roles: []string{"parent"}}, users: map[uuid.UUID]entity.User{studentID: {ID: studentID, Roles: []string{"student"}}}}
	items, err := NewUserDomain(repo).ListParentStudents(context.Background(), testTx{}, "hash")
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%+v err=%v", items, err)
	}
}
