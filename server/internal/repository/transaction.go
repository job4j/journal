package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionManager struct{ pool *pgxpool.Pool }

func NewTransactionManager(pool *pgxpool.Pool) *transactionManager {
	return &transactionManager{pool: pool}
}
func (m *transactionManager) Begin(ctx context.Context) (Transaction, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	return tx, nil
}
func pgxTransaction(transaction Transaction) (pgx.Tx, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("postgres repository requires pgx transaction")
	}
	return tx, nil
}

type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
type TransactionManager interface {
	Begin(context.Context) (Transaction, error)
}
