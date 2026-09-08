package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func (s *AuthService) Logout(ctx context.Context, token string) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	hash := sha256.Sum256([]byte(token))
	if err = s.domain.Logout(ctx, tx, hex.EncodeToString(hash[:])); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
