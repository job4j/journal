package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
)

func TestRevokeParentStudentAccessHandlesPartialAndRepeatedLinks(t *testing.T) {
	parentID, studentID, profileID, journalID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	value := studentID.String()
	repo := &userRepoStub{authenticated: true, allowed: true, users: map[uuid.UUID]entity.User{parentID: {ID: parentID, Roles: []string{"parent"}}, studentID: {ID: studentID, Roles: []string{"student"}}}, permissions: []entity.Permission{{ID: profileID, Code: "can_view_user", Value: &value}, {ID: journalID, Code: "can_view_journal", Value: &value}}, userPermissions: []entity.UserPermission{{UserID: parentID, PermissionID: profileID}}}
	domain := NewUserDomain(repo)
	for range 2 {
		if _, err := domain.RevokeParentStudentAccess(context.Background(), testTx{}, "hash", parentID, studentID); err != nil {
			t.Fatal(err)
		}
	}
	if len(repo.userPermissions) != 0 {
		t.Fatalf("links = %+v", repo.userPermissions)
	}
}
