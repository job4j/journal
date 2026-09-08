package service

import (
	"context"
	"fmt"
	"journal/server/internal/domain"
)

func (s *AuthService) Login(ctx context.Context, email, password string) (result domain.LoginResponse, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.LoginResponse{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.Login(ctx, tx, domain.LoginRequest{Email: email, Password: password})
	if err != nil {
		return domain.LoginResponse{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.LoginResponse{}, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
