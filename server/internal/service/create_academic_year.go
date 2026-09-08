package service

import (
	"context"
	"fmt"
	"journal/server/internal/domain"
)

func (s *AcademicYearService) CreateAcademicYear(ctx context.Context, token string, request domain.CreateAcademicYearRequest) (result domain.AcademicYearView, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	request.SessionTokenHash = sessionTokenHash(token)
	result, err = s.domain.CreateAcademicYear(ctx, tx, request)
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
