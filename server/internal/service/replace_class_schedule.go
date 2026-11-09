package service

import (
	"context"
	"journal/server/internal/domain"
)

func (s *ScheduleService) ReplaceClassSchedule(ctx context.Context, token string, req domain.ReplaceClassScheduleRequest) (domain.ReplaceClassScheduleResponse, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.ReplaceClassScheduleResponse{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := s.domain.ReplaceClassSchedule(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}
