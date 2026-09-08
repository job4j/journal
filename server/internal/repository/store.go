package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type Repository struct{}

func NewRepository() *Repository { return &Repository{} }

func queryOne[T any](ctx context.Context, transaction Transaction, operation, sql string, args ...any) (T, error) {
	var zero T
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return zero, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return zero, operationError(operation, err)
	}
	value, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if errors.Is(err, pgx.ErrNoRows) {
		return zero, ErrNotFound
	}
	if err != nil {
		return zero, operationError(operation, err)
	}
	return value, nil
}
func queryMany[T any](ctx context.Context, transaction Transaction, operation, sql string, args ...any) ([]T, error) {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, operationError(operation, err)
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, operationError(operation, err)
	}
	return values, nil
}
func deleteRows(ctx context.Context, transaction Transaction, operation, sql string, args ...any) error {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return operationError(operation, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
