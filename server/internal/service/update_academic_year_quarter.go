package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"time"
)

func (s *AcademicYearService) UpdateAcademicYearQuarter(ctx context.Context, token string, id uuid.UUID, name string, startsOn, endsOn time.Time) (result entity.AcademicYearQuarter, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.UpdateAcademicYearQuarter(ctx, tx, sessionTokenHash(token), id, name, startsOn, endsOn)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
