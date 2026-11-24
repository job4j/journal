package entity

import (
	"github.com/google/uuid"
	"time"
)

type LessonAttachment struct {
	ID          uuid.UUID  `db:"id"`
	LessonID    uuid.UUID  `db:"lesson_id"`
	ObjectKey   string     `db:"object_key"`
	FileName    string     `db:"file_name"`
	ContentType string     `db:"content_type"`
	SizeBytes   int64      `db:"size_bytes"`
	CreatedBy   uuid.UUID  `db:"created_by"`
	CreatedAt   time.Time  `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
