package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func (s *ClassService) PutQuarterGrade(ctx context.Context, token string, assignmentID, quarterID, studentID uuid.UUID, input domain.QuarterGradeInput) (result entity.QuarterGrade, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.PutQuarterGrade(ctx, tx, sessionTokenHash(token), assignmentID, quarterID, studentID, input)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *ClassService) ListQuarterGrades(ctx context.Context, token string, assignmentID, quarterID uuid.UUID) ([]entity.QuarterGrade, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListQuarterGrades(ctx, tx, sessionTokenHash(token), assignmentID, quarterID)
}
