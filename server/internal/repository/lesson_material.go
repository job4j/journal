package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (r *Repository) CreateLessonMaterial(ctx context.Context, tx Transaction, value entity.LessonMaterial) (entity.LessonMaterial, error) {
	return queryOne[entity.LessonMaterial](ctx, tx, "create lesson_material", `INSERT INTO lesson_materials (lesson_id, title, url, position) VALUES ($1, $2, $3, $4) RETURNING id, lesson_id, title, url, position, created_at`, value.LessonID, value.Title, value.URL, value.Position)
}
func (r *Repository) GetLessonMaterial(ctx context.Context, tx Transaction, id uuid.UUID) (entity.LessonMaterial, error) {
	return queryOne[entity.LessonMaterial](ctx, tx, "get lesson_material", `SELECT id, lesson_id, title, url, position, created_at FROM lesson_materials WHERE id = $1`, id)
}
func (r *Repository) ListLessonMaterials(ctx context.Context, tx Transaction) ([]entity.LessonMaterial, error) {
	return queryMany[entity.LessonMaterial](ctx, tx, "list lesson_materials", `SELECT id, lesson_id, title, url, position, created_at FROM lesson_materials ORDER BY created_at, id`)
}
func (r *Repository) UpdateLessonMaterial(ctx context.Context, tx Transaction, value entity.LessonMaterial) (entity.LessonMaterial, error) {
	return queryOne[entity.LessonMaterial](ctx, tx, "update lesson_material", `UPDATE lesson_materials SET lesson_id = $1, title = $2, url = $3, position = $4 WHERE id = $5 RETURNING id, lesson_id, title, url, position, created_at`, value.LessonID, value.Title, value.URL, value.Position, value.ID)
}
func (r *Repository) DeleteLessonMaterial(ctx context.Context, tx Transaction, id uuid.UUID) error {
	return deleteRows(ctx, tx, "delete lesson_material", `DELETE FROM lesson_materials WHERE id = $1`, id)
}
