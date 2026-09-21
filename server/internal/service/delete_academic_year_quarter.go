package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (s *AcademicYearService) DeleteAcademicYearQuarter(ctx context.Context, token string, id uuid.UUID) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.DeleteAcademicYearQuarter(ctx, tx, sessionTokenHash(token), id); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
