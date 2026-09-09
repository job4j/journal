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

func (d *UserDomain) ListManagedParentStudents(ctx context.Context, tx repository.Transaction, hash string, parentID uuid.UUID) ([]entity.User, error) {
	if err := d.authorize(ctx, tx, hash, "can_update_user"); err != nil {
		return nil, err
	}
	parent, err := d.repo.GetUser(ctx, tx, parentID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get parent: %w", err)
	}
	if !slices.Contains(parent.Roles, "parent") {
		return nil, ErrInvalidParentStudent
	}
	permissions, err := d.repo.ListPermissions(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}
	links, err := d.repo.ListUserPermissions(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list user permissions: %w", err)
	}
	studentIDs := map[uuid.UUID]struct{}{}
	permissionValues := map[uuid.UUID]string{}
	for _, permission := range permissions {
		if permission.Code == "can_view_user" && permission.Value != nil {
			permissionValues[permission.ID] = *permission.Value
		}
	}
	for _, link := range links {
		if link.UserID != parentID {
			continue
		}
		if raw, ok := permissionValues[link.PermissionID]; ok {
			if id, parseErr := uuid.Parse(raw); parseErr == nil {
				studentIDs[id] = struct{}{}
			}
		}
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	result := []entity.User{}
	for _, user := range users {
		if _, ok := studentIDs[user.ID]; ok && slices.Contains(user.Roles, "student") {
			result = append(result, user)
		}
	}
	return result, nil
}
