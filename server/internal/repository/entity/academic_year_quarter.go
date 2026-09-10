package entity

import (
	"github.com/google/uuid"
	"time"
)

type AcademicYearQuarter struct {
	ID             uuid.UUID `db:"id"`
	AcademicYearID uuid.UUID `db:"academic_year_id"`
	Number         int16     `db:"number"`
	Name           string    `db:"name"`
	StartsOn       time.Time `db:"starts_on"`
	EndsOn         time.Time `db:"ends_on"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
