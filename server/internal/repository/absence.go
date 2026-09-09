package repository

import (
	"context"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) EnsureAbsence(ctx context.Context, tx Transaction, value entity.Absence) (entity.Absence, error) {
	return queryOne[entity.Absence](ctx, tx, "ensure absence", `INSERT INTO absences (lesson_id, user_id, recorded_by) VALUES ($1, $2, $3)
		ON CONFLICT (lesson_id, user_id) DO UPDATE SET recorded_by = EXCLUDED.recorded_by, updated_at = now()
		RETURNING id, lesson_id, user_id, recorded_by, created_at, updated_at`, value.LessonID, value.UserID, value.RecordedBy)
}

func (r *Repository) DeleteAbsence(ctx context.Context, tx Transaction, lessonID, studentID uuid.UUID) error {
	transaction, err := pgxTransaction(tx)
	if err != nil {
		return err
	}
	_, err = transaction.Exec(ctx, `DELETE FROM absences WHERE lesson_id = $1 AND user_id = $2`, lessonID, studentID)
	if err != nil {
		return operationError("delete absence", err)
	}
	return nil
}

func (r *Repository) ListAbsences(ctx context.Context, tx Transaction) ([]entity.Absence, error) {
	return queryMany[entity.Absence](ctx, tx, "list absences", `SELECT id, lesson_id, user_id, recorded_by, created_at, updated_at FROM absences ORDER BY created_at, id`)
}
