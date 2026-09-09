package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func (s *ClassService) PutStudentAbsence(ctx context.Context, token string, lessonID, studentID uuid.UUID) (result entity.Absence, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.PutStudentAbsence(ctx, tx, sessionTokenHash(token), lessonID, studentID)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func (s *ClassService) DeleteStudentAbsence(ctx context.Context, token string, lessonID, studentID uuid.UUID) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.DeleteStudentAbsence(ctx, tx, sessionTokenHash(token), lessonID, studentID); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
