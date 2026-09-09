package service

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) GetParentStudentJournal(ctx context.Context, token string, studentID, yearID uuid.UUID) (domain.ParentJournal, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.ParentJournal{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.GetParentStudentJournal(ctx, tx, sessionTokenHash(token), studentID, yearID)
}
