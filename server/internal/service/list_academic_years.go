package service

import (
	"context"
	"journal/server/internal/domain"
)

func (s *AcademicYearService) ListAcademicYears(ctx context.Context, token string) ([]domain.AcademicYearView, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListAcademicYears(ctx, tx, sessionTokenHash(token))
}
