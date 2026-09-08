package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserRole struct {
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
}
