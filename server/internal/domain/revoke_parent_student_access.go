package domain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"journal/server/internal/repository"
)

func (d *UserDomain) RevokeParentStudentAccess(ctx context.Context, tx repository.Transaction, hash string, parentID, studentID uuid.UUID) (ParentStudentAccess, error) {
	if err := d.authorize(ctx, tx, hash, "can_update_user"); err != nil {
		return ParentStudentAccess{}, err
	}
	parent, err := d.repo.GetUser(ctx, tx, parentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentStudentAccess{}, ErrUserNotFound
	}
	if err != nil {
		return ParentStudentAccess{}, fmt.Errorf("get parent: %w", err)
	}
	student, err := d.repo.GetUser(ctx, tx, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentStudentAccess{}, ErrUserNotFound
	}
	if err != nil {
		return ParentStudentAccess{}, fmt.Errorf("get student: %w", err)
	}
	if parentID == studentID || !slices.Contains(parent.Roles, "parent") || !slices.Contains(student.Roles, "student") {
		return ParentStudentAccess{}, ErrInvalidParentStudent
	}
	permissions, err := d.repo.ListPermissions(ctx, tx)
	if err != nil {
		return ParentStudentAccess{}, fmt.Errorf("list permissions: %w", err)
	}
	studentValue := studentID.String()
	for _, permission := range permissions {
		if permission.Value == nil || *permission.Value != studentValue || (permission.Code != "can_view_user" && permission.Code != "can_view_journal") {
			continue
		}
		if err := d.repo.DeleteUserPermission(ctx, tx, parentID, permission.ID); err != nil && !errors.Is(err, repository.ErrNotFound) {
			return ParentStudentAccess{}, fmt.Errorf("revoke %s permission: %w", permission.Code, err)
		}
	}
	return ParentStudentAccess{ParentID: parentID, StudentID: studentID}, nil
}
