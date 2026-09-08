package domain

import (
	"context"
	"journal/server/internal/repository"
)

func (d *AuthDomain) Logout(ctx context.Context, tx repository.Transaction, tokenHash string) error {
	if tokenHash == "" {
		return nil
	}
	return d.repo.RevokeSession(ctx, tx, tokenHash, d.now().UTC())
}
