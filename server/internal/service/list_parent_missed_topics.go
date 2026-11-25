package service

import (
	"context"
	"journal/server/internal/domain"
)

func (s *ClassService) ListParentMissedTopics(ctx context.Context, token string, req domain.ListParentMissedTopicsRequest) (result domain.ListParentMissedTopicsResponse, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.ListParentMissedTopics(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return result, err
	}
	err = tx.Commit(ctx)
	return result, err
}
