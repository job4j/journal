package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateAcademicYear(ctx context.Context, tx Transaction, value entity.AcademicYear) (entity.AcademicYear, error) {
	return queryOne[entity.AcademicYear](ctx, tx, "create academic_year", `INSERT INTO academic_years (name, starts_on, ends_on, status) VALUES ($1, $2, $3, $4) RETURNING id, name, starts_on, ends_on, status, created_at, updated_at`, value.Name, value.StartsOn, value.EndsOn, value.Status)
}
func (r *Repository) GetAcademicYear(ctx context.Context, tx Transaction, id uuid.UUID) (entity.AcademicYear, error) {
	return queryOne[entity.AcademicYear](ctx, tx, "get academic_year", `SELECT id, name, starts_on, ends_on, status, created_at, updated_at FROM academic_years WHERE id = $1`, id)
}
func (r *Repository) ListAcademicYears(ctx context.Context, tx Transaction) ([]entity.AcademicYear, error) {
	return queryMany[entity.AcademicYear](ctx, tx, "list academic_years", `SELECT id, name, starts_on, ends_on, status, created_at, updated_at FROM academic_years ORDER BY starts_on DESC, id`)
}
func (r *Repository) UpdateAcademicYear(ctx context.Context, tx Transaction, value entity.AcademicYear) (entity.AcademicYear, error) {
	return queryOne[entity.AcademicYear](ctx, tx, "update academic_year", `UPDATE academic_years SET name = $1, starts_on = $2, ends_on = $3, status = $4, updated_at = now() WHERE id = $5 RETURNING id, name, starts_on, ends_on, status, created_at, updated_at`, value.Name, value.StartsOn, value.EndsOn, value.Status, value.ID)
}
func (r *Repository) DeleteAcademicYear(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete academic_year", `DELETE FROM academic_years WHERE id = $1`, id)
}
