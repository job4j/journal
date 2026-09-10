package service

import (
	"context"
	"fmt"
	"journal/server/internal/repository/entity"
	"time"
)

func (s *AcademicYearService) CreateAcademicYearQuarter(ctx context.Context, token, yearID, name string, startsOn, endsOn time.Time) (result entity.AcademicYearQuarter, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateAcademicYearQuarter(ctx, tx, sessionTokenHash(token), yearID, name, startsOn, endsOn)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
