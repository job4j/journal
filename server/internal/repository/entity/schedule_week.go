package entity

import (
	"github.com/google/uuid"
	"time"
)

type ScheduleWeek struct {
	ID        uuid.UUID `db:"id"`
	ClassID   uuid.UUID `db:"class_id"`
	WeekStart time.Time `db:"week_start"`
	Revision  int       `db:"revision"`
	UpdatedAt time.Time `db:"updated_at"`
}
