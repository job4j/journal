package domain

import (
	"errors"
	"journal/server/internal/repository"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBlocked        = errors.New("user blocked")
)

type AuthDomain struct {
	repo            repository.AuthRepository
	sessionLifetime time.Duration
	now             func() time.Time
}

func NewAuthDomain(repo repository.AuthRepository, sessionLifetime time.Duration) *AuthDomain {
	return &AuthDomain{repo: repo, sessionLifetime: sessionLifetime, now: time.Now}
}
