package domain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func (d *UserDomain) ListParentStudents(ctx context.Context, tx repository.Transaction, hash string) ([]entity.User, error) {
	parent, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUnauthenticated
	}
	if err != nil {
		return nil, fmt.Errorf("find parent session: %w", err)
	}
	if !slices.Contains(parent.Roles, "parent") {
		return nil, ErrForbidden
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	_, global, err := d.repo.CheckSessionPermission(ctx, tx, hash, "can_view_user")
	if err != nil {
		return nil, fmt.Errorf("check global permission: %w", err)
	}
	allowed := map[uuid.UUID]struct{}{}
	if !global {
		permissions, err := d.repo.ListPermissions(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("list permissions: %w", err)
		}
		links, err := d.repo.ListUserPermissions(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("list user permissions: %w", err)
		}
		values := map[uuid.UUID]uuid.UUID{}
		for _, permission := range permissions {
			if permission.Code == "can_view_user" && permission.Value != nil {
				if id, parseErr := uuid.Parse(*permission.Value); parseErr == nil {
					values[permission.ID] = id
				}
			}
		}
		for _, link := range links {
			if link.UserID == parent.ID {
				if id, ok := values[link.PermissionID]; ok {
					allowed[id] = struct{}{}
				}
			}
		}
	}
	result := []entity.User{}
	for _, user := range users {
		_, personal := allowed[user.ID]
		if slices.Contains(user.Roles, "student") && (global || personal) {
			result = append(result, user)
		}
	}
	return result, nil
}
