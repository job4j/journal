package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func (r *roleRepoStub) ListGlobalPermissions(context.Context, repository.Transaction) ([]entity.Permission, error) {
	return r.permissions, nil
}
func (r *roleRepoStub) GetRole(context.Context, repository.Transaction, uuid.UUID) (entity.Role, error) {
	return r.createdRole, nil
}
func (r *roleRepoStub) ListRoles(context.Context, repository.Transaction) ([]entity.Role, error) {
	return nil, nil
}
func (r *roleRepoStub) UpdateRole(_ context.Context, _ repository.Transaction, v entity.Role) (entity.Role, error) {
	return v, nil
}
func (r *roleRepoStub) DeleteRole(context.Context, repository.Transaction, uuid.UUID) error {
	return nil
}
func (r *roleRepoStub) ListRolePermissions(context.Context, repository.Transaction) ([]entity.RolePermission, error) {
	return r.assigned, nil
}
func (r *roleRepoStub) ListRolePermissionsByRoleID(context.Context, repository.Transaction, uuid.UUID) ([]entity.RolePermission, error) {
	return r.assigned, nil
}
func (r *roleRepoStub) DeleteRolePermission(context.Context, repository.Transaction, uuid.UUID, uuid.UUID) error {
	return nil
}
