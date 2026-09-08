package service

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func (r *serviceRoleRepoStub) ListGlobalPermissions(context.Context, repository.Transaction) ([]entity.Permission, error) {
	return nil, nil
}
func (r *serviceRoleRepoStub) GetRole(context.Context, repository.Transaction, uuid.UUID) (entity.Role, error) {
	return entity.Role{}, nil
}
func (r *serviceRoleRepoStub) ListRoles(context.Context, repository.Transaction) ([]entity.Role, error) {
	return nil, nil
}
func (r *serviceRoleRepoStub) UpdateRole(_ context.Context, _ repository.Transaction, v entity.Role) (entity.Role, error) {
	return v, nil
}
func (r *serviceRoleRepoStub) DeleteRole(context.Context, repository.Transaction, uuid.UUID) error {
	return nil
}
func (r *serviceRoleRepoStub) ListRolePermissions(context.Context, repository.Transaction) ([]entity.RolePermission, error) {
	return nil, nil
}
func (r *serviceRoleRepoStub) ListRolePermissionsByRoleID(context.Context, repository.Transaction, uuid.UUID) ([]entity.RolePermission, error) {
	return nil, nil
}
func (r *serviceRoleRepoStub) DeleteRolePermission(context.Context, repository.Transaction, uuid.UUID, uuid.UUID) error {
	return nil
}
