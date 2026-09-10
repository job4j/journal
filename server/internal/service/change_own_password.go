package service

import (
	"context"
	"fmt"
)

func (s *AuthService) ChangeOwnPassword(
	ctx context.Context,
	token string,
	currentPassword string,
	newPassword string,
) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.ChangeOwnPassword(
		ctx,
		tx,
		sessionTokenHash(token),
		currentPassword,
		newPassword,
	); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
