package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateClass(ctx context.Context, tx Transaction, value entity.Class) (entity.Class, error) {
	return queryOne[entity.Class](ctx, tx, "create class", `INSERT INTO classes (academic_year_id, name, grade_level) VALUES ($1, $2, $3) RETURNING id, academic_year_id, name, grade_level, created_at, updated_at`, value.AcademicYearID, value.Name, value.GradeLevel)
}
func (r *Repository) GetClass(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Class, error) {
	return queryOne[entity.Class](ctx, tx, "get class", `SELECT id, academic_year_id, name, grade_level, created_at, updated_at FROM classes WHERE id = $1`, id)
}
func (r *Repository) ListClasses(ctx context.Context, tx Transaction) ([]entity.Class, error) {
	return queryMany[entity.Class](ctx, tx, "list classes", `SELECT id, academic_year_id, name, grade_level, created_at, updated_at FROM classes ORDER BY grade_level, name, id`)
}
func (r *Repository) UpdateClass(ctx context.Context, tx Transaction, value entity.Class) (entity.Class, error) {
	return queryOne[entity.Class](ctx, tx, "update class", `UPDATE classes SET academic_year_id = $1, name = $2, grade_level = $3, updated_at = now() WHERE id = $4 RETURNING id, academic_year_id, name, grade_level, created_at, updated_at`, value.AcademicYearID, value.Name, value.GradeLevel, value.ID)
}
func (r *Repository) DeleteClass(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete class", `DELETE FROM classes WHERE id = $1`, id)
}
