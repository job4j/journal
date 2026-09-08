package entity

import (
	"github.com/google/uuid"
	"time"
)

type Permission struct {
	ID          uuid.UUID `db:"id"`
	Code        string    `db:"code"`
	Value       *string   `db:"value"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
