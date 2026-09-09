package service

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (s *UserService) ListManagedParentStudents(ctx context.Context, token string, parentID uuid.UUID) ([]entity.User, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListManagedParentStudents(ctx, tx, sessionTokenHash(token), parentID)
}
