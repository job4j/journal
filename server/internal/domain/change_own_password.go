package domain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"journal/server/internal/repository"
)

func (d *AuthDomain) ChangeOwnPassword(
	ctx context.Context,
	tx repository.Transaction,
	tokenHash string,
	currentPassword string,
	newPassword string,
) error {
	if currentPassword == "" || newPassword == "" {
		return ErrInvalidUser
	}
	user, err := d.CurrentUser(ctx, tx, tokenHash)
	if err != nil {
		return err
	}
	if !slices.Contains(user.Roles, "teacher") && !slices.Contains(user.Roles, "parent") {
		return ErrForbidden
	}
	valid, err := verifyArgon2ID(currentPassword, user.PasswordHash)
	if err != nil || !valid {
		return ErrInvalidCredentials
	}
	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err = d.repo.UpdatePasswordHash(ctx, tx, user.ID, passwordHash); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
