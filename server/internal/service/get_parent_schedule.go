package service

import (
	"context"
	"journal/server/internal/domain"
)

func (s *ScheduleService) GetParentSchedule(ctx context.Context, token string, req domain.GetParentScheduleRequest) (domain.GetParentScheduleResponse, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.GetParentScheduleResponse{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := s.domain.GetParentSchedule(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}
