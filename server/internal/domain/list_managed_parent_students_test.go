package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

func TestListManagedParentStudentsUsesParentPermissions(t *testing.T) {
	parentID, studentID, otherID, permissionID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	value := studentID.String()
	repo := &userRepoStub{authenticated: true, allowed: true, users: map[uuid.UUID]entity.User{parentID: {ID: parentID, Roles: []string{"parent"}}}, permissions: []entity.Permission{{ID: permissionID, Code: "can_view_user", Value: &value}}, userPermissions: []entity.UserPermission{{UserID: parentID, PermissionID: permissionID}}, created: entity.User{ID: otherID}}
	repoUsers := []entity.User{{ID: studentID, Roles: []string{"student"}}, {ID: otherID, Roles: []string{"student"}}}
	repo.created = repoUsers[0]
	items, err := NewUserDomain(&managedParentRepoStub{userRepoStub: repo, listed: repoUsers}).ListManagedParentStudents(context.Background(), testTx{}, "hash", parentID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != studentID {
		t.Fatalf("items=%+v", items)
	}
}

type managedParentRepoStub struct {
	*userRepoStub
	listed []entity.User
}

func (r *managedParentRepoStub) ListUsers(context.Context, repository.Transaction) ([]entity.User, error) {
	return r.listed, nil
}
