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

func NewUserService(txManager repository.TransactionManager, repo repository.UserManagementRepository) *UserService {
	return &UserService{txManager: txManager, domain: domain.NewUserDomain(repo)}
}
