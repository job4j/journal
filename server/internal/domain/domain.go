package domain

import (
	"errors"
	"journal/server/internal/repository"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBlocked        = errors.New("user blocked")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidRole        = errors.New("invalid role")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrRoleExists         = errors.New("role already exists")
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleInUse          = errors.New("role is in use")
	ErrProtectedRole      = errors.New("protected role")
	ErrInvalidUser        = errors.New("invalid user")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInUse          = errors.New("user is in use")
)

type AuthDomain struct {
	repo            repository.AuthRepository
	sessionLifetime time.Duration
	now             func() time.Time
}

func NewAuthDomain(repo repository.AuthRepository, sessionLifetime time.Duration) *AuthDomain {
	return &AuthDomain{repo: repo, sessionLifetime: sessionLifetime, now: time.Now}
}

type RoleDomain struct{ repo repository.RoleRepository }

func NewRoleDomain(repo repository.RoleRepository) *RoleDomain { return &RoleDomain{repo: repo} }

type AcademicYearDomain struct {
	repo repository.AcademicYearRepository
}

func NewAcademicYearDomain(repo repository.AcademicYearRepository) *AcademicYearDomain {
	return &AcademicYearDomain{repo: repo}
}
