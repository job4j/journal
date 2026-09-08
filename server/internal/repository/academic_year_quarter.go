package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateAcademicYearQuarter(ctx context.Context, tx Transaction, value entity.AcademicYearQuarter) (entity.AcademicYearQuarter, error) {
	return queryOne[entity.AcademicYearQuarter](ctx, tx, "create academic_year_quarter", `INSERT INTO academic_year_quarters (academic_year_id, number, starts_on, ends_on) VALUES ($1, $2, $3, $4) RETURNING id, academic_year_id, number, starts_on, ends_on, created_at, updated_at`, value.AcademicYearID, value.Number, value.StartsOn, value.EndsOn)
}
func (r *Repository) GetAcademicYearQuarter(ctx context.Context, tx Transaction, id uuid.UUID) (entity.AcademicYearQuarter, error) {
	return queryOne[entity.AcademicYearQuarter](ctx, tx, "get academic_year_quarter", `SELECT id, academic_year_id, number, starts_on, ends_on, created_at, updated_at FROM academic_year_quarters WHERE id = $1`, id)
}
func (r *Repository) ListAcademicYearQuarters(ctx context.Context, tx Transaction) ([]entity.AcademicYearQuarter, error) {
	return queryMany[entity.AcademicYearQuarter](ctx, tx, "list academic_year_quarters", `SELECT id, academic_year_id, number, starts_on, ends_on, created_at, updated_at FROM academic_year_quarters ORDER BY academic_year_id, number`)
}
func (r *Repository) UpdateAcademicYearQuarter(ctx context.Context, tx Transaction, value entity.AcademicYearQuarter) (entity.AcademicYearQuarter, error) {
	return queryOne[entity.AcademicYearQuarter](ctx, tx, "update academic_year_quarter", `UPDATE academic_year_quarters SET academic_year_id = $1, number = $2, starts_on = $3, ends_on = $4, updated_at = now() WHERE id = $5 RETURNING id, academic_year_id, number, starts_on, ends_on, created_at, updated_at`, value.AcademicYearID, value.Number, value.StartsOn, value.EndsOn, value.ID)
}
func (r *Repository) DeleteAcademicYearQuarter(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete academic_year_quarter", `DELETE FROM academic_year_quarters WHERE id = $1`, id)
}
