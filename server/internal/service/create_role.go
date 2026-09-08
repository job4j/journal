package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func sessionTokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
func (s *RoleService) CreateRole(ctx context.Context, token, code, name string, permissions []string) (result domain.CreateRoleResponse, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateRole(ctx, tx, domain.CreateRoleRequest{SessionTokenHash: sessionTokenHash(token), Code: code, Name: name, PermissionCodes: permissions})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *RoleService) ListRoles(ctx context.Context, token string) (result domain.ListRolesResponse, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.ListRoles(ctx, tx, domain.ListRolesRequest{SessionTokenHash: sessionTokenHash(token)})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *RoleService) ListPermissions(ctx context.Context, token string) (result domain.ListPermissionsResponse, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.ListPermissions(ctx, tx, domain.ListRolesRequest{SessionTokenHash: sessionTokenHash(token)})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, token string, id uuid.UUID, code, name string, permissions []string) (result domain.RoleView, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.UpdateRole(ctx, tx, domain.UpdateRoleRequest{SessionTokenHash: sessionTokenHash(token), ID: id, Code: code, Name: name, PermissionCodes: permissions})
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}
func (s *RoleService) DeleteRole(ctx context.Context, token string, id uuid.UUID) (err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = s.domain.DeleteRole(ctx, tx, domain.DeleteRoleRequest{SessionTokenHash: sessionTokenHash(token), ID: id}); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
