package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) CreateClass(ctx context.Context, token string, yearID uuid.UUID, name string, grade int16) (result domain.ClassView, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateClass(ctx, tx, sessionTokenHash(token), yearID, name, grade)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
