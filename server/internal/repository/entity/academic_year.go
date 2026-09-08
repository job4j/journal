package entity

import (
	"github.com/google/uuid"
	"time"
)

type AcademicYear struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	StartsOn  time.Time `db:"starts_on"`
	EndsOn    time.Time `db:"ends_on"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
