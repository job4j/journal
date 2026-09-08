package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"journal/server/internal/repository/entity"
)

func (s *AuthService) CurrentUser(ctx context.Context, token string) (entity.User, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	hash := sha256.Sum256([]byte(token))
	return s.domain.CurrentUser(ctx, tx, hex.EncodeToString(hash[:]))
}
