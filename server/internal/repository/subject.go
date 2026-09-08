package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateSubject(ctx context.Context, tx Transaction, value entity.Subject) (entity.Subject, error) {
	return queryOne[entity.Subject](ctx, tx, "create subject", `INSERT INTO subjects (code, name) VALUES ($1, $2) RETURNING id, code, name, created_at, updated_at`, value.Code, value.Name)
}
func (r *Repository) GetSubject(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Subject, error) {
	return queryOne[entity.Subject](ctx, tx, "get subject", `SELECT id, code, name, created_at, updated_at FROM subjects WHERE id = $1`, id)
}
func (r *Repository) ListSubjects(ctx context.Context, tx Transaction) ([]entity.Subject, error) {
	return queryMany[entity.Subject](ctx, tx, "list subjects", `SELECT id, code, name, created_at, updated_at FROM subjects ORDER BY name, id`)
}
func (r *Repository) UpdateSubject(ctx context.Context, tx Transaction, value entity.Subject) (entity.Subject, error) {
	return queryOne[entity.Subject](ctx, tx, "update subject", `UPDATE subjects SET code = $1, name = $2, updated_at = now() WHERE id = $3 RETURNING id, code, name, created_at, updated_at`, value.Code, value.Name, value.ID)
}
func (r *Repository) DeleteSubject(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete subject", `DELETE FROM subjects WHERE id = $1`, id)
}
