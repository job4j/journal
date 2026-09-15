package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ClassService) DeleteStudentScore(ctx context.Context, token string, gradeItemID, studentID uuid.UUID) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.DeleteStudentScore(ctx, tx, sessionTokenHash(token), gradeItemID, studentID); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
