package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"testing"
)

func TestOperationErrorMapsPostgresConstraints(t *testing.T) {
	tests := []struct {
		code   string
		target error
	}{
		{code: "23505", target: ErrConflict},
		{code: "23503", target: ErrReference},
		{code: "23502", target: ErrInvalidData},
		{code: "23514", target: ErrInvalidData},
	}
	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			err := operationError("write", &pgconn.PgError{Code: test.code})
			if !errors.Is(err, test.target) {
				t.Fatalf("expected %v, got %v", test.target, err)
			}
		})
	}
}
