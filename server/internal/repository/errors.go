package repository

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
	ErrReference   = errors.New("referenced entity does not exist")
	ErrInvalidData = errors.New("invalid data")
)

func operationError(operation string, err error) error {
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) {
		switch postgresError.Code {
		case "23505":
			return fmt.Errorf("%s: %w: %v", operation, ErrConflict, err)
		case "23503":
			return fmt.Errorf("%s: %w: %v", operation, ErrReference, err)
		case "23502", "23514":
			return fmt.Errorf("%s: %w: %v", operation, ErrInvalidData, err)
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
