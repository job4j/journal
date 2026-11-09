package service

import (
	"context"
	"journal/server/internal/domain"
)

func (s *ScheduleService) GetClassSchedule(ctx context.Context, token string, req domain.GetClassScheduleRequest) (domain.GetClassScheduleResponse, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.GetClassScheduleResponse{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := s.domain.GetClassSchedule(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}
