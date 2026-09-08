package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateGradeItem(ctx context.Context, tx Transaction, value entity.GradeItem) (entity.GradeItem, error) {
	return queryOne[entity.GradeItem](ctx, tx, "create grade_item", `INSERT INTO grade_items (lesson_id, title, kind, grading_scale, max_score) VALUES ($1, $2, $3, $4, $5) RETURNING id, lesson_id, title, kind, grading_scale, max_score, created_at, updated_at`, value.LessonID, value.Title, value.Kind, value.GradingScale, value.MaxScore)
}
func (r *Repository) GetGradeItem(ctx context.Context, tx Transaction, id uuid.UUID) (entity.GradeItem, error) {
	return queryOne[entity.GradeItem](ctx, tx, "get grade_item", `SELECT id, lesson_id, title, kind, grading_scale, max_score, created_at, updated_at FROM grade_items WHERE id = $1`, id)
}
func (r *Repository) ListGradeItems(ctx context.Context, tx Transaction) ([]entity.GradeItem, error) {
	return queryMany[entity.GradeItem](ctx, tx, "list grade_items", `SELECT id, lesson_id, title, kind, grading_scale, max_score, created_at, updated_at FROM grade_items ORDER BY created_at, id`)
}
func (r *Repository) UpdateGradeItem(ctx context.Context, tx Transaction, value entity.GradeItem) (entity.GradeItem, error) {
	return queryOne[entity.GradeItem](ctx, tx, "update grade_item", `UPDATE grade_items SET lesson_id = $1, title = $2, kind = $3, grading_scale = $4, max_score = $5, updated_at = now() WHERE id = $6 RETURNING id, lesson_id, title, kind, grading_scale, max_score, created_at, updated_at`, value.LessonID, value.Title, value.Kind, value.GradingScale, value.MaxScore, value.ID)
}
func (r *Repository) DeleteGradeItem(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete grade_item", `DELETE FROM grade_items WHERE id = $1`, id)
}
