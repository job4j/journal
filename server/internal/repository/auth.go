package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

type TransactionManager interface {
	Begin(context.Context) (Transaction, error)
}

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Status       string
	Roles        []string
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

// AuthRepository contains table-level persistence operations only.
type AuthRepository interface {
	FindUserByEmail(context.Context, Transaction, string) (User, error)
	InsertSession(context.Context, Transaction, Session) error
}
