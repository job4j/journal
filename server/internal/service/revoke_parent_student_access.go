package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *UserService) RevokeParentStudentAccess(ctx context.Context, token string, parentID, studentID uuid.UUID) (result domain.ParentStudentAccess, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.RevokeParentStudentAccess(ctx, tx, sessionTokenHash(token), parentID, studentID)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
