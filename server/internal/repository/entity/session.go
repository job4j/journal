package entity

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID         uuid.UUID  `db:"id"`
	UserID     uuid.UUID  `db:"user_id"`
	TokenHash  string     `db:"token_hash"`
	ExpiresAt  time.Time  `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
	LastUsedAt time.Time  `db:"last_used_at"`
	RevokedAt  *time.Time `db:"revoked_at"`
}
