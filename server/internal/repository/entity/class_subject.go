package entity

import (
	"github.com/google/uuid"
	"time"
)

type ClassSubject struct {
	ID                   uuid.UUID `db:"id"`
	ClassID              uuid.UUID `db:"class_id"`
	SubjectID            uuid.UUID `db:"subject_id"`
	ResponsibleTeacherID uuid.UUID `db:"responsible_teacher_id"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}
