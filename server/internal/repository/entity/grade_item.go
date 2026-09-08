package entity

import (
	"github.com/google/uuid"
	"time"
)

type GradeItem struct {
	ID           uuid.UUID `db:"id"`
	LessonID     uuid.UUID `db:"lesson_id"`
	Title        string    `db:"title"`
	Kind         string    `db:"kind"`
	GradingScale string    `db:"grading_scale"`
	MaxScore     *int8     `db:"max_score"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
