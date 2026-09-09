package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (s *ClassService) CreateGradeItem(ctx context.Context, token string, lessonID uuid.UUID, title, kind, scale string, max *float64) (result entity.GradeItem, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateGradeItem(ctx, tx, sessionTokenHash(token), lessonID, title, kind, scale, max)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
