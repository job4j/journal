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
