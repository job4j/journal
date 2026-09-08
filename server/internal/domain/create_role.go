package domain

import (
	"context"
	"errors"
	"fmt"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"regexp"
	"sort"
	"strings"
)

var roleCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type CreateRoleRequest struct {
	SessionTokenHash string
	Code             string
	Name             string
	PermissionCodes  []string
}
type CreateRoleResponse struct {
	Role            entity.Role
	PermissionCodes []string
}

func (d *RoleDomain) authorize(ctx context.Context, tx repository.Transaction, tokenHash string) error {
	return Authorize(ctx, tx, d.repo, tokenHash, "can_manage_roles")
}
func normalizeRole(code, name string, rawCodes []string) (string, string, []string, error) {
	code = strings.ToLower(strings.TrimSpace(code))
	name = strings.TrimSpace(name)
	if len(code) < 2 || len(code) > 50 || !roleCodePattern.MatchString(code) || name == "" || len([]rune(name)) > 100 {
		return "", "", nil, ErrInvalidRole
	}
	codes := make([]string, 0, len(rawCodes))
	seen := map[string]struct{}{}
	for _, raw := range rawCodes {
		permissionCode := strings.ToLower(strings.TrimSpace(raw))
		if len(permissionCode) < 2 || len(permissionCode) > 100 || !roleCodePattern.MatchString(permissionCode) {
			return "", "", nil, ErrInvalidRole
		}
		if _, ok := seen[permissionCode]; ok {
			return "", "", nil, ErrInvalidRole
		}
		seen[permissionCode] = struct{}{}
		codes = append(codes, permissionCode)
	}
	sort.Strings(codes)
	return code, name, codes, nil
}
func (d *RoleDomain) resolvePermissions(ctx context.Context, tx repository.Transaction, codes []string) ([]entity.Permission, error) {
	permissions, err := d.repo.FindGlobalPermissionsByCodes(ctx, tx, codes)
	if err != nil {
		return nil, fmt.Errorf("load permissions: %w", err)
	}
	if len(permissions) != len(codes) {
		return nil, ErrPermissionNotFound
	}
	return permissions, nil
}
func (d *RoleDomain) CreateRole(ctx context.Context, tx repository.Transaction, request CreateRoleRequest) (CreateRoleResponse, error) {
	if err := d.authorize(ctx, tx, request.SessionTokenHash); err != nil {
		return CreateRoleResponse{}, err
	}
	code, name, codes, err := normalizeRole(request.Code, request.Name, request.PermissionCodes)
	if err != nil {
		return CreateRoleResponse{}, err
	}
	permissions, err := d.resolvePermissions(ctx, tx, codes)
	if err != nil {
		return CreateRoleResponse{}, err
	}
	role, err := d.repo.CreateRole(ctx, tx, entity.Role{Code: code, Name: name})
	if errors.Is(err, repository.ErrConflict) {
		return CreateRoleResponse{}, ErrRoleExists
	}
	if err != nil {
		return CreateRoleResponse{}, fmt.Errorf("create role: %w", err)
	}
	for _, permission := range permissions {
		if _, err = d.repo.CreateRolePermission(ctx, tx, entity.RolePermission{RoleID: role.ID, PermissionID: permission.ID}); err != nil {
			return CreateRoleResponse{}, fmt.Errorf("assign permission %s: %w", permission.Code, err)
		}
	}
	return CreateRoleResponse{Role: role, PermissionCodes: codes}, nil
}
