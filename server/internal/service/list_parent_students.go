package service

import (
	"context"
	"journal/server/internal/repository/entity"
)

func (s *UserService) ListParentStudents(ctx context.Context, token string) ([]entity.User, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListParentStudents(ctx, tx, sessionTokenHash(token))
}
