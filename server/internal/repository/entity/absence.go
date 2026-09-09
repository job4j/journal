package entity

import (
	"time"

	"github.com/google/uuid"
)

type Absence struct {
	ID         uuid.UUID `db:"id"`
	LessonID   uuid.UUID `db:"lesson_id"`
	UserID     uuid.UUID `db:"user_id"`
	RecordedBy uuid.UUID `db:"recorded_by"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
