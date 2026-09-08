package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
)

func TestDeleteRoleProtectsAdmin(t *testing.T) {
	repo := &roleRepoStub{authenticated: true, allowed: true, createdRole: entity.Role{ID: uuid.New(), Code: "admin"}}
	err := NewRoleDomain(repo).DeleteRole(context.Background(), nil, DeleteRoleRequest{SessionTokenHash: "hash", ID: repo.createdRole.ID})
	if !errors.Is(err, ErrProtectedRole) {
		t.Fatalf("error = %v", err)
	}
}
func TestUpdateAdminKeepsManageRolesPermission(t *testing.T) {
	repo := &roleRepoStub{authenticated: true, allowed: true, createdRole: entity.Role{ID: uuid.New(), Code: "admin"}}
	_, err := NewRoleDomain(repo).UpdateRole(context.Background(), nil, UpdateRoleRequest{SessionTokenHash: "hash", ID: repo.createdRole.ID, Code: "admin", Name: "Администратор", PermissionCodes: []string{"can_view_user"}})
	if !errors.Is(err, ErrProtectedRole) {
		t.Fatalf("error = %v", err)
	}
}
