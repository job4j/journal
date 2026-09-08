package postgres

import (
	"context"
	"errors"
	"fmt"

	"journal/server/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository { return &AuthRepository{} }

func (r *AuthRepository) FindUserByEmail(ctx context.Context, transaction repository.Transaction, email string) (repository.User, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return repository.User{}, err
	}

	var user repository.User
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
		return repository.User{}, repository.ErrNotFound
	}
	if err != nil {
		return repository.User{}, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

func (r *AuthRepository) InsertSession(ctx context.Context, transaction repository.Transaction, session repository.Session) error {
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

type TransactionManager struct{ pool *pgxpool.Pool }

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

func (m *TransactionManager) Begin(ctx context.Context) (repository.Transaction, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	return tx, nil
}

func pgxTransaction(transaction repository.Transaction) (pgx.Tx, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("postgres repository requires pgx transaction")
	}
	return tx, nil
}
