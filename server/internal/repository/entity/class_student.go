package entity

import (
	"github.com/google/uuid"
	"time"
)

type ClassStudent struct {
	ClassID    uuid.UUID  `db:"class_id"`
	UserID     uuid.UUID  `db:"user_id"`
	EnrolledOn time.Time  `db:"enrolled_on"`
	LeftOn     *time.Time `db:"left_on"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}
