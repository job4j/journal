package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (s *ClassService) PutStudentScore(ctx context.Context, token string, gradeItemID, studentID uuid.UUID, numeric *float64, text, comment *string) (result entity.Score, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.PutStudentScore(ctx, tx, sessionTokenHash(token), gradeItemID, studentID, numeric, text, comment)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
