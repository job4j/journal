package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"journal/server/internal/repository"
)

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
