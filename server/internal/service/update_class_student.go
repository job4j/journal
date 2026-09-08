package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) UpdateClassStudent(ctx context.Context, token string, classID, studentID uuid.UUID, leftOn *time.Time) (result domain.ClassStudentView, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.UpdateClassStudent(ctx, tx, sessionTokenHash(token), classID, studentID, leftOn)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
