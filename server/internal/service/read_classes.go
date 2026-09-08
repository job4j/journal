package service

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (s *ClassService) ListClasses(ctx context.Context, token string, yearID uuid.UUID) ([]domain.ClassView, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListClasses(ctx, tx, sessionTokenHash(token), yearID)
}
func (s *ClassService) GetClass(ctx context.Context, token string, id uuid.UUID) (domain.ClassView, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return domain.ClassView{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.GetClass(ctx, tx, sessionTokenHash(token), id)
}
func (s *ClassService) ListClassStudents(ctx context.Context, token string, id uuid.UUID) ([]domain.ClassStudentView, error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	return s.domain.ListClassStudents(ctx, tx, sessionTokenHash(token), id)
}
