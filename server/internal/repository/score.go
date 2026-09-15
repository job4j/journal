package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateScore(ctx context.Context, tx Transaction, value entity.Score) (entity.Score, error) {
	return queryOne[entity.Score](ctx, tx, "create score", `INSERT INTO scores (grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by, created_at, updated_at`, value.GradeItemID, value.UserID, value.NumericValue, value.TextValue, value.TeacherComment, value.CreatedBy, value.UpdatedBy)
}
func (r *Repository) GetScore(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Score, error) {
	return queryOne[entity.Score](ctx, tx, "get score", `SELECT id, grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by, created_at, updated_at FROM scores WHERE id = $1`, id)
}
func (r *Repository) ListScores(ctx context.Context, tx Transaction) ([]entity.Score, error) {
	return queryMany[entity.Score](ctx, tx, "list scores", `SELECT id, grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by, created_at, updated_at FROM scores ORDER BY created_at, id`)
}
func (r *Repository) UpdateScore(ctx context.Context, tx Transaction, value entity.Score) (entity.Score, error) {
	return queryOne[entity.Score](ctx, tx, "update score", `UPDATE scores SET grade_item_id = $1, user_id = $2, numeric_value = $3, text_value = $4, teacher_comment = $5, updated_by = $6, updated_at = now() WHERE id = $7 RETURNING id, grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by, created_at, updated_at`, value.GradeItemID, value.UserID, value.NumericValue, value.TextValue, value.TeacherComment, value.UpdatedBy, value.ID)
}

func (r *Repository) UpsertScore(ctx context.Context, tx Transaction, value entity.Score) (entity.Score, error) {
	return queryOne[entity.Score](ctx, tx, "upsert score", `INSERT INTO scores (grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		ON CONFLICT (grade_item_id, user_id) DO UPDATE SET numeric_value = EXCLUDED.numeric_value, text_value = EXCLUDED.text_value,
			teacher_comment = EXCLUDED.teacher_comment, updated_by = EXCLUDED.updated_by, updated_at = now()
		RETURNING id, grade_item_id, user_id, numeric_value, text_value, teacher_comment, created_by, updated_by, created_at, updated_at`,
		value.GradeItemID, value.UserID, value.NumericValue, value.TextValue, value.TeacherComment, value.CreatedBy)
}
func (r *Repository) DeleteScore(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete score", `DELETE FROM scores WHERE id = $1`, id)
}

func (r *Repository) DeleteStudentScore(ctx context.Context, tx Transaction, gradeItemID, studentID uuid.UUID) error {
	return deleteRows(ctx, tx, "delete student score", `DELETE FROM scores WHERE grade_item_id = $1 AND user_id = $2`, gradeItemID, studentID)
}
