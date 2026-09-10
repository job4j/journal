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
	var nullableEmail, nullablePhone *string
	err = tx.QueryRow(ctx, `
		SELECT u.id, u.login, u.email, u.phone, u.password_hash, u.name, u.status,
		       ARRAY(SELECT role.code FROM user_roles user_role JOIN roles role ON role.id=user_role.role_id WHERE user_role.user_id=u.id ORDER BY role.code)
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		WHERE lower(u.login) = lower($1)
		GROUP BY u.id
	`, email).Scan(&user.ID, &user.Login, &nullableEmail, &nullablePhone, &user.PasswordHash, &user.Name, &user.Status, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, ErrNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("find user by email: %w", err)
	}
	if nullableEmail != nil {
		user.Email = *nullableEmail
	}
	if nullablePhone != nil {
		user.Phone = *nullablePhone
	}
	return user, nil
}

func (r *authRepository) FindUserBySessionTokenHash(ctx context.Context, transaction Transaction, tokenHash string, now time.Time) (entity.User, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return entity.User{}, err
	}
	var user entity.User
	var email, phone *string
	err = tx.QueryRow(ctx, `SELECT u.id, u.login, u.email, u.phone, u.password_hash, u.name, u.status,
		ARRAY(SELECT role.code FROM user_roles user_role JOIN roles role ON role.id=user_role.role_id WHERE user_role.user_id=u.id ORDER BY role.code)
		FROM sessions s JOIN users u ON u.id=s.user_id LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id
		WHERE s.token_hash=$1 AND s.expires_at>$2 AND s.revoked_at IS NULL AND u.status='active' GROUP BY u.id`, tokenHash, now).
		Scan(&user.ID, &user.Login, &email, &phone, &user.PasswordHash, &user.Name, &user.Status, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, ErrNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("find user by session: %w", err)
	}
	if email != nil {
		user.Email = *email
	}
	if phone != nil {
		user.Phone = *phone
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
