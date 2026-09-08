package entity

import (
	"github.com/google/uuid"
	"time"
)

type RolePermission struct {
	RoleID       uuid.UUID `db:"role_id"`
	PermissionID uuid.UUID `db:"permission_id"`
	CreatedAt    time.Time `db:"created_at"`
}
