package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserPermission struct {
	UserID       uuid.UUID `db:"user_id"`
	PermissionID uuid.UUID `db:"permission_id"`
	CreatedAt    time.Time `db:"created_at"`
}
