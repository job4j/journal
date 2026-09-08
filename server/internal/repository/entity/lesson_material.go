package entity

import (
	"github.com/google/uuid"
	"time"
)

type LessonMaterial struct {
	ID        uuid.UUID `db:"id"`
	LessonID  uuid.UUID `db:"lesson_id"`
	Title     string    `db:"title"`
	URL       string    `db:"url"`
	Position  int16     `db:"position"`
	CreatedAt time.Time `db:"created_at"`
}
