package entity

import (
	"github.com/google/uuid"
	"time"
)

type Class struct {
	ID             uuid.UUID `db:"id"`
	AcademicYearID uuid.UUID `db:"academic_year_id"`
	Name           string    `db:"name"`
	GradeLevel     int16     `db:"grade_level"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
