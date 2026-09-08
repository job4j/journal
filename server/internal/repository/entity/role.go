package entity

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	ID        uuid.UUID `db:"id"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}
