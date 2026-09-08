package repository

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}
