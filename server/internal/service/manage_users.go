package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func (s *UserService) CreateUser(ctx context.Context, token string, input domain.UserInput) (result entity.User, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateUser(ctx, tx, domain.CreateUserRequest{SessionTokenHash: sessionTokenHash(token), Input: input})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *UserService) ListUsers(ctx context.Context, token, role string, limit, offset int) (result domain.ListUsersResult, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.ListUsers(ctx, tx, domain.ListUsersRequest{SessionTokenHash: sessionTokenHash(token), Role: role, Limit: limit, Offset: offset})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *UserService) GetUser(ctx context.Context, token string, id uuid.UUID) (result entity.User, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.GetUser(ctx, tx, domain.GetUserRequest{SessionTokenHash: sessionTokenHash(token), ID: id})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *UserService) UpdateUser(ctx context.Context, token string, id uuid.UUID, input domain.UserInput) (result entity.User, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.UpdateUser(ctx, tx, domain.UpdateUserRequest{SessionTokenHash: sessionTokenHash(token), ID: id, Input: input})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *UserService) DeleteUser(ctx context.Context, token string, id uuid.UUID) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.DeleteUser(ctx, tx, domain.DeleteUserRequest{SessionTokenHash: sessionTokenHash(token), ID: id}); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
