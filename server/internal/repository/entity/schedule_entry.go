package entity

import (
	"github.com/google/uuid"
	"time"
)

type ScheduleEntry struct {
	ID             uuid.UUID  `db:"id"`
	WeekID         uuid.UUID  `db:"week_id"`
	LessonDate     time.Time  `db:"lesson_date"`
	Position       int        `db:"position"`
	ClassSubjectID *uuid.UUID `db:"class_subject_id"`
	Title          string     `db:"title"`
	Notes          string     `db:"notes"`
	StartsAt       *string    `db:"starts_at"`
	EndsAt         *string    `db:"ends_at"`
}
