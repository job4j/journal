package service

import (
	"context"
	"fmt"
	"time"

	"journal/server/internal/domain"
	"journal/server/internal/repository"
)

type AuthService struct {
	txManager       repository.TransactionManager
	repo            repository.AuthRepository
	sessionLifetime time.Duration
	now             func() time.Time
}

func NewAuthService(txManager repository.TransactionManager, repo repository.AuthRepository, sessionLifetime time.Duration) *AuthService {
	return &AuthService{txManager: txManager, repo: repo, sessionLifetime: sessionLifetime, now: time.Now}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (result domain.LoginResult, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.LoginResult{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	result, err = domain.Login(ctx, tx, s.repo, email, password, s.now().UTC(), s.sessionLifetime)
	if err != nil {
		return domain.LoginResult{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.LoginResult{}, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
