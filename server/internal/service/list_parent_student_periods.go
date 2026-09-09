package service

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) ListParentStudentPeriods(ctx context.Context, token string, studentID uuid.UUID) ([]domain.ParentStudentPeriod, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListParentStudentPeriods(ctx, tx, sessionTokenHash(token), studentID)
}
