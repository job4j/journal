package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) UpdateClassSubject(ctx context.Context, token string, assignmentID, teacherID uuid.UUID) (result domain.ClassSubjectView, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.UpdateClassSubject(ctx, tx, sessionTokenHash(token), assignmentID, teacherID)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
