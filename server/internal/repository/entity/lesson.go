package entity

import (
	"github.com/google/uuid"
	"time"
)

type Lesson struct {
	ID             uuid.UUID `db:"id"`
	ClassSubjectID uuid.UUID `db:"class_subject_id"`
	LessonDate     time.Time `db:"lesson_date"`
	Position       int16     `db:"position"`
	Topic          string    `db:"topic"`
	Homework       *string   `db:"homework"`
	CreatedBy      uuid.UUID `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
