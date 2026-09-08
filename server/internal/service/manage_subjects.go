package service

import (
	"context"
	"fmt"
	"journal/server/internal/repository/entity"
)

func (s *SubjectService) ListSubjects(ctx context.Context, token string) ([]entity.Subject, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListSubjects(ctx, tx, sessionTokenHash(token))
}
func (s *SubjectService) CreateSubject(ctx context.Context, token, code, name string) (result entity.Subject, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateSubject(ctx, tx, sessionTokenHash(token), code, name)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
