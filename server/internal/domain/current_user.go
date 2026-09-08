package domain

import (
	"context"
	"errors"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func (d *AuthDomain) CurrentUser(ctx context.Context, tx repository.Transaction, tokenHash string) (entity.User, error) {
	if tokenHash == "" {
		return entity.User{}, ErrUnauthenticated
	}
	user, err := d.repo.FindUserBySessionTokenHash(ctx, tx, tokenHash, d.now().UTC())
	if errors.Is(err, repository.ErrNotFound) {
		return entity.User{}, ErrUnauthenticated
	}
	return user, err
}
