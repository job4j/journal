package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateLesson(ctx context.Context, tx Transaction, value entity.Lesson) (entity.Lesson, error) {
	return queryOne[entity.Lesson](ctx, tx, "create lesson", `INSERT INTO lessons (class_subject_id, lesson_date, position, topic, homework, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, class_subject_id, lesson_date, position, topic, homework, created_by, created_at, updated_at`, value.ClassSubjectID, value.LessonDate, value.Position, value.Topic, value.Homework, value.CreatedBy)
}
func (r *Repository) GetLesson(ctx context.Context, tx Transaction, id uuid.UUID) (entity.Lesson, error) {
	return queryOne[entity.Lesson](ctx, tx, "get lesson", `SELECT id, class_subject_id, lesson_date, position, topic, homework, created_by, created_at, updated_at FROM lessons WHERE id = $1`, id)
}
func (r *Repository) ListLessons(ctx context.Context, tx Transaction) ([]entity.Lesson, error) {
	return queryMany[entity.Lesson](ctx, tx, "list lessons", `SELECT id, class_subject_id, lesson_date, position, topic, homework, created_by, created_at, updated_at FROM lessons ORDER BY created_at, id`)
}
func (r *Repository) UpdateLesson(ctx context.Context, tx Transaction, value entity.Lesson) (entity.Lesson, error) {
	return queryOne[entity.Lesson](ctx, tx, "update lesson", `UPDATE lessons SET class_subject_id = $1, lesson_date = $2, position = $3, topic = $4, homework = $5, created_by = $6, updated_at = now() WHERE id = $7 RETURNING id, class_subject_id, lesson_date, position, topic, homework, created_by, created_at, updated_at`, value.ClassSubjectID, value.LessonDate, value.Position, value.Topic, value.Homework, value.CreatedBy, value.ID)
}
func (r *Repository) DeleteLesson(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete lesson", `DELETE FROM lessons WHERE id = $1`, id)
}
