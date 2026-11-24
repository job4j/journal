package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateLessonAttachment(ctx context.Context, tx Transaction, item entity.LessonAttachment) (entity.LessonAttachment, error) {
	return queryOne[entity.LessonAttachment](ctx, tx, "create lesson attachment", `INSERT INTO lesson_attachments(id,lesson_id,object_key,file_name,content_type,size_bytes,created_by) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING *`, item.ID, item.LessonID, item.ObjectKey, item.FileName, item.ContentType, item.SizeBytes, item.CreatedBy)
}
func (r *Repository) ListLessonAttachments(ctx context.Context, tx Transaction, lessonID uuid.UUID) ([]entity.LessonAttachment, error) {
	return queryMany[entity.LessonAttachment](ctx, tx, "list lesson attachments", `SELECT * FROM lesson_attachments WHERE lesson_id=$1 AND deleted_at IS NULL ORDER BY created_at,id`, lessonID)
}
func (r *Repository) GetLessonAttachment(ctx context.Context, tx Transaction, id uuid.UUID) (entity.LessonAttachment, error) {
	return queryOne[entity.LessonAttachment](ctx, tx, "get lesson attachment", `SELECT * FROM lesson_attachments WHERE id=$1`, id)
}
func (r *Repository) DeleteLessonAttachment(ctx context.Context, tx Transaction, id uuid.UUID) error {
	_, err := queryOne[entity.LessonAttachment](ctx, tx, "delete lesson attachment", `UPDATE lesson_attachments SET deleted_at=coalesce(deleted_at,now()) WHERE id=$1 RETURNING *`, id)
	return err
}
