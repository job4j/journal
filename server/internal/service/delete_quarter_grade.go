package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (s *ClassService) DeleteQuarterGrade(ctx context.Context, token string, assignmentID, quarterID, studentID uuid.UUID) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.DeleteQuarterGrade(ctx, tx, sessionTokenHash(token), assignmentID, quarterID, studentID); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
