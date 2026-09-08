package domain

import (
	"errors"
	"journal/server/internal/repository"
	"time"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserBlocked          = errors.New("user blocked")
	ErrUnauthenticated      = errors.New("unauthenticated")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidRole          = errors.New("invalid role")
	ErrPermissionNotFound   = errors.New("permission not found")
	ErrRoleExists           = errors.New("role already exists")
	ErrRoleNotFound         = errors.New("role not found")
	ErrRoleInUse            = errors.New("role is in use")
	ErrProtectedRole        = errors.New("protected role")
	ErrInvalidUser          = errors.New("invalid user")
	ErrUserExists           = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserInUse            = errors.New("user is in use")
	ErrInvalidAcademicYear  = errors.New("invalid academic year")
	ErrAcademicYearExists   = errors.New("academic year already exists")
	ErrInvalidSubject       = errors.New("invalid subject")
	ErrSubjectExists        = errors.New("subject already exists")
	ErrClassNotFound        = errors.New("class not found")
	ErrInvalidClass         = errors.New("invalid class")
	ErrClassExists          = errors.New("class already exists")
	ErrAcademicYearNotFound = errors.New("academic year not found")
)

type AuthDomain struct {
	repo            repository.AuthRepository
	sessionLifetime time.Duration
	now             func() time.Time
}

type SubjectDomain struct{ repo repository.SubjectRepository }

func NewSubjectDomain(repo repository.SubjectRepository) *SubjectDomain {
	return &SubjectDomain{repo: repo}
}

type ClassDomain struct {
	repo repository.ClassReadRepository
}

func NewClassDomain(repo repository.ClassReadRepository) *ClassDomain {
	return &ClassDomain{repo: repo}
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
