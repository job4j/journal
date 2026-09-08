package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateClassSubject(ctx context.Context, tx Transaction, value entity.ClassSubject) (entity.ClassSubject, error) {
	return queryOne[entity.ClassSubject](ctx, tx, "create class_subject", `INSERT INTO class_subjects (class_id, subject_id, responsible_teacher_id) VALUES ($1, $2, $3) RETURNING id, class_id, subject_id, responsible_teacher_id, created_at, updated_at`, value.ClassID, value.SubjectID, value.ResponsibleTeacherID)
}
func (r *Repository) GetClassSubject(ctx context.Context, tx Transaction, id uuid.UUID) (entity.ClassSubject, error) {
	return queryOne[entity.ClassSubject](ctx, tx, "get class_subject", `SELECT id, class_id, subject_id, responsible_teacher_id, created_at, updated_at FROM class_subjects WHERE id = $1`, id)
}
func (r *Repository) ListClassSubjects(ctx context.Context, tx Transaction) ([]entity.ClassSubject, error) {
	return queryMany[entity.ClassSubject](ctx, tx, "list class_subjects", `SELECT id, class_id, subject_id, responsible_teacher_id, created_at, updated_at FROM class_subjects ORDER BY created_at, id`)
}
func (r *Repository) UpdateClassSubject(ctx context.Context, tx Transaction, value entity.ClassSubject) (entity.ClassSubject, error) {
	return queryOne[entity.ClassSubject](ctx, tx, "update class_subject", `UPDATE class_subjects SET class_id = $1, subject_id = $2, responsible_teacher_id = $3, updated_at = now() WHERE id = $4 RETURNING id, class_id, subject_id, responsible_teacher_id, created_at, updated_at`, value.ClassID, value.SubjectID, value.ResponsibleTeacherID, value.ID)
}
func (r *Repository) DeleteClassSubject(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete class_subject", `DELETE FROM class_subjects WHERE id = $1`, id)
}
