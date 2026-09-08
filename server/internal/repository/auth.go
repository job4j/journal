package repository

import (
	"context"
	"errors"
	"fmt"

	"journal/server/internal/repository/entity"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r *authRepository) FindUserByEmail(ctx context.Context, transaction Transaction, email string) (entity.User, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User
	err = tx.QueryRow(ctx, `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.status,
		       COALESCE(array_agg(r.code ORDER BY r.code) FILTER (WHERE r.code IS NOT NULL), '{}')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		WHERE lower(u.email) = lower($1)
		GROUP BY u.id
	`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Status, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, ErrNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

func (r *authRepository) FindUserBySessionTokenHash(ctx context.Context, transaction Transaction, tokenHash string, now time.Time) (entity.User, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return entity.User{}, err
	}
	var user entity.User
	err = tx.QueryRow(ctx, `SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.status,
		COALESCE(array_agg(r.code ORDER BY r.code) FILTER (WHERE r.code IS NOT NULL), '{}')
		FROM sessions s JOIN users u ON u.id=s.user_id LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id
		WHERE s.token_hash=$1 AND s.expires_at>$2 AND s.revoked_at IS NULL AND u.status='active' GROUP BY u.id`, tokenHash, now).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Status, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, ErrNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("find user by session: %w", err)
	}
	return user, nil
}

func (r *authRepository) RevokeSession(ctx context.Context, transaction Transaction, tokenHash string, now time.Time) error {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE sessions SET revoked_at=$2 WHERE token_hash=$1 AND revoked_at IS NULL`, tokenHash, now)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (r *authRepository) InsertSession(ctx context.Context, transaction Transaction, session entity.Session) error {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO sessions (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, session.ID, session.UserID, session.TokenHash, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}
