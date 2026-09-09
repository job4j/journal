package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	Email        string    `db:"email"`
	Phone        string    `db:"phone"`
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Status       string    `db:"status"`
	Roles        []string  `db:"roles"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
