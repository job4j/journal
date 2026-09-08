package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type RoleView struct {
	Role            entity.Role
	PermissionCodes []string
}
type ListRolesRequest struct{ SessionTokenHash string }
type ListRolesResponse struct{ Roles []RoleView }
type ListPermissionsResponse struct{ Permissions []entity.Permission }
type UpdateRoleRequest struct {
	SessionTokenHash string
	ID               uuid.UUID
	Code             string
	Name             string
	PermissionCodes  []string
}
type DeleteRoleRequest struct {
	SessionTokenHash string
	ID               uuid.UUID
}

func (d *RoleDomain) ListRoles(ctx context.Context, tx repository.Transaction, request ListRolesRequest) (ListRolesResponse, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash); err != nil {
		return ListRolesResponse{}, err
	}
	roles, err := d.repo.ListRoles(ctx, tx)
	if err != nil {
		return ListRolesResponse{}, fmt.Errorf("list roles: %w", err)
	}
	permissions, err := d.repo.ListGlobalPermissions(ctx, tx)
	if err != nil {
		return ListRolesResponse{}, fmt.Errorf("list permissions: %w", err)
	}
	links, err := d.repo.ListRolePermissions(ctx, tx)
	if err != nil {
		return ListRolesResponse{}, fmt.Errorf("list role permissions: %w", err)
	}
	codes := map[uuid.UUID]string{}
	for _, p := range permissions {
		codes[p.ID] = p.Code
	}
	result := make([]RoleView, len(roles))
	for i, role := range roles {
		result[i].Role = role
		result[i].PermissionCodes = []string{}
		for _, link := range links {
			if link.RoleID == role.ID {
				if code, ok := codes[link.PermissionID]; ok {
					result[i].PermissionCodes = append(result[i].PermissionCodes, code)
				}
			}
		}
	}
	return ListRolesResponse{Roles: result}, nil
}
func (d *RoleDomain) ListPermissions(ctx context.Context, tx repository.Transaction, request ListRolesRequest) (ListPermissionsResponse, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash); err != nil {
		return ListPermissionsResponse{}, err
	}
	permissions, err := d.repo.ListGlobalPermissions(ctx, tx)
	if err != nil {
		return ListPermissionsResponse{}, fmt.Errorf("list permissions: %w", err)
	}
	return ListPermissionsResponse{Permissions: permissions}, nil
}
func (d *RoleDomain) UpdateRole(ctx context.Context, tx repository.Transaction, request UpdateRoleRequest) (RoleView, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash); err != nil {
		return RoleView{}, err
	}
	code, name, codes, err := normalizeRole(request.Code, request.Name, request.PermissionCodes)
	if err != nil {
		return RoleView{}, err
	}
	current, err := d.repo.GetRole(ctx, tx, request.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return RoleView{}, ErrRoleNotFound
	}
	if err != nil {
		return RoleView{}, fmt.Errorf("get role: %w", err)
	}
	if current.Code == "admin" && code != "admin" {
		return RoleView{}, ErrProtectedRole
	}
	if current.Code == "admin" && !containsCode(codes, "can_manage_roles") {
		return RoleView{}, ErrProtectedRole
	}
	permissions, err := d.resolvePermissions(ctx, tx, codes)
	if err != nil {
		return RoleView{}, err
	}
	currentLinks, err := d.repo.ListRolePermissionsByRoleID(ctx, tx, current.ID)
	if err != nil {
		return RoleView{}, fmt.Errorf("list current permissions: %w", err)
	}
	role, err := d.repo.UpdateRole(ctx, tx, entity.Role{ID: current.ID, Code: code, Name: name})
	if errors.Is(err, repository.ErrConflict) || errors.Is(err, repository.ErrReference) {
		return RoleView{}, ErrRoleExists
	}
	if err != nil {
		return RoleView{}, fmt.Errorf("update role: %w", err)
	}
	for _, link := range currentLinks {
		if err = d.repo.DeleteRolePermission(ctx, tx, link.RoleID, link.PermissionID); err != nil {
			return RoleView{}, fmt.Errorf("remove permission: %w", err)
		}
	}
	for _, permission := range permissions {
		if _, err = d.repo.CreateRolePermission(ctx, tx, entity.RolePermission{RoleID: role.ID, PermissionID: permission.ID}); err != nil {
			return RoleView{}, fmt.Errorf("assign permission: %w", err)
		}
	}
	return RoleView{Role: role, PermissionCodes: codes}, nil
}
func (d *RoleDomain) DeleteRole(ctx context.Context, tx repository.Transaction, request DeleteRoleRequest) error {
	if err := d.authorize(ctx, tx, request.SessionTokenHash); err != nil {
		return err
	}
	role, err := d.repo.GetRole(ctx, tx, request.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrRoleNotFound
	}
	if err != nil {
		return fmt.Errorf("get role: %w", err)
	}
	if role.Code == "admin" {
		return ErrProtectedRole
	}
	err = d.repo.DeleteRole(ctx, tx, role.ID)
	if errors.Is(err, repository.ErrConflict) || errors.Is(err, repository.ErrReference) {
		return ErrRoleInUse
	}
	if errors.Is(err, repository.ErrNotFound) {
		return ErrRoleNotFound
	}
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}
func containsCode(codes []string, target string) bool {
	for _, code := range codes {
		if code == target {
			return true
		}
	}
	return false
}
