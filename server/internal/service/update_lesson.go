package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) UpdateLesson(ctx context.Context, token string, id uuid.UUID, input domain.CreateLessonInput) (result domain.LessonView, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.UpdateLesson(ctx, tx, sessionTokenHash(token), id, input)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
