package service

import (
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"time"
)

type AuthService struct {
	txManager repository.TransactionManager
	domain    *domain.AuthDomain
}

func NewAuthService(txManager repository.TransactionManager, repo repository.AuthRepository, sessionLifetime time.Duration) *AuthService {
	return &AuthService{txManager: txManager, domain: domain.NewAuthDomain(repo, sessionLifetime)}
}

type RoleService struct {
	txManager repository.TransactionManager
	domain    *domain.RoleDomain
}

func NewRoleService(txManager repository.TransactionManager, repo repository.RoleRepository) *RoleService {
	return &RoleService{txManager: txManager, domain: domain.NewRoleDomain(repo)}
}

type UserService struct {
	txManager repository.TransactionManager
	domain    *domain.UserDomain
}

type AcademicYearService struct {
	txManager repository.TransactionManager
	domain    *domain.AcademicYearDomain
}

func NewAcademicYearService(txManager repository.TransactionManager, repo repository.AcademicYearRepository) *AcademicYearService {
	return &AcademicYearService{txManager: txManager, domain: domain.NewAcademicYearDomain(repo)}
}

type SubjectService struct {
	txManager repository.TransactionManager
	domain    *domain.SubjectDomain
}

func NewSubjectService(txManager repository.TransactionManager, repo repository.SubjectRepository) *SubjectService {
	return &SubjectService{txManager: txManager, domain: domain.NewSubjectDomain(repo)}
}

func NewUserService(txManager repository.TransactionManager, repo repository.UserManagementRepository) *UserService {
	return &UserService{txManager: txManager, domain: domain.NewUserDomain(repo)}
}
