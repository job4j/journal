package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type roleRepoStub struct {
	authenticated bool
	allowed       bool
	permissions   []entity.Permission
	createdRole   entity.Role
	createRoleErr error
	assigned      []entity.RolePermission
}

func (r *roleRepoStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return r.authenticated, r.allowed, nil
}
func (r *roleRepoStub) FindGlobalPermissionsByCodes(context.Context, repository.Transaction, []string) ([]entity.Permission, error) {
	return r.permissions, nil
}
func (r *roleRepoStub) CreateRole(_ context.Context, _ repository.Transaction, value entity.Role) (entity.Role, error) {
	if r.createRoleErr != nil {
		return entity.Role{}, r.createRoleErr
	}
	value.ID = uuid.New()
	r.createdRole = value
	return value, nil
}
func (r *roleRepoStub) CreateRolePermission(_ context.Context, _ repository.Transaction, value entity.RolePermission) (entity.RolePermission, error) {
	r.assigned = append(r.assigned, value)
	return value, nil
}

func TestCreateRoleCreatesRoleAndAssignments(t *testing.T) {
	p1 := entity.Permission{ID: uuid.New(), Code: "can_create_user"}
	p2 := entity.Permission{ID: uuid.New(), Code: "can_view_user"}
	repo := &roleRepoStub{authenticated: true, allowed: true, permissions: []entity.Permission{p1, p2}}
	result, err := NewRoleDomain(repo).CreateRole(context.Background(), nil, CreateRoleRequest{SessionTokenHash: "hash", Code: " Editor ", Name: " Редактор ", PermissionCodes: []string{"can_view_user", "can_create_user"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Role.Code != "editor" || result.Role.Name != "Редактор" {
		t.Fatalf("unexpected role: %+v", result.Role)
	}
	if len(repo.assigned) != 2 {
		t.Fatalf("assignments = %d", len(repo.assigned))
	}
	if len(result.PermissionCodes) != 2 || result.PermissionCodes[0] != "can_create_user" {
		t.Fatalf("permissions = %v", result.PermissionCodes)
	}
}
func TestCreateRoleRejectsUnauthenticatedSession(t *testing.T) {
	_, err := NewRoleDomain(&roleRepoStub{}).CreateRole(context.Background(), nil, CreateRoleRequest{})
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("error = %v", err)
	}
}
func TestCreateRoleRejectsMissingPermission(t *testing.T) {
	repo := &roleRepoStub{authenticated: true, allowed: true}
	_, err := NewRoleDomain(repo).CreateRole(context.Background(), nil, CreateRoleRequest{Code: "editor", Name: "Редактор", PermissionCodes: []string{"missing_permission"}})
	if !errors.Is(err, ErrPermissionNotFound) {
		t.Fatalf("error = %v", err)
	}
}
func TestCreateRoleMapsConflict(t *testing.T) {
	permission := entity.Permission{ID: uuid.New(), Code: "can_view_user"}
	repo := &roleRepoStub{authenticated: true, allowed: true, permissions: []entity.Permission{permission}, createRoleErr: repository.ErrConflict}
	_, err := NewRoleDomain(repo).CreateRole(context.Background(), nil, CreateRoleRequest{Code: "editor", Name: "Редактор", PermissionCodes: []string{"can_view_user"}})
	if !errors.Is(err, ErrRoleExists) {
		t.Fatalf("error = %v", err)
	}
}
