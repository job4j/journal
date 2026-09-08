package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateRole(ctx context.Context, tx Transaction, value entity.Role) (entity.Role, error) {
	return queryOne[entity.Role](ctx, tx, "create role", `INSERT INTO roles (code, name) VALUES ($1, $2) RETURNING id, code, name, created_at`, value.Code, value.Name)
}
func (r *Repository) GetRole(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Role, error) {
	return queryOne[entity.Role](ctx, tx, "get role", `SELECT id, code, name, created_at FROM roles WHERE id = $1`, id)
}
func (r *Repository) ListRoles(ctx context.Context, tx Transaction) ([]entity.Role, error) {
	return queryMany[entity.Role](ctx, tx, "list roles", `SELECT id, code, name, created_at FROM roles ORDER BY created_at, id`)
}
func (r *Repository) UpdateRole(ctx context.Context, tx Transaction, value entity.Role) (entity.Role, error) {
	return queryOne[entity.Role](ctx, tx, "update role", `UPDATE roles SET code = $1, name = $2 WHERE id = $3 RETURNING id, code, name, created_at`, value.Code, value.Name, value.ID)
}
func (r *Repository) DeleteRole(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete role", `DELETE FROM roles WHERE id = $1`, id)
}

func (r *Repository) FindRolesByCodes(ctx context.Context, tx Transaction, codes []string) ([]entity.Role, error) {
	return queryMany[entity.Role](ctx, tx, "find roles by codes", `SELECT id,code,name,created_at FROM roles WHERE code=ANY($1) ORDER BY code`, codes)
}
