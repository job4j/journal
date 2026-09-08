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

type ParentStudentAccess struct {
	ParentID, StudentID uuid.UUID
	Permissions         []string
}

func (d *UserDomain) GrantParentStudentAccess(ctx context.Context, tx repository.Transaction, hash string, parentID, studentID uuid.UUID) (ParentStudentAccess, error) {
	if err := d.authorize(ctx, tx, hash, "can_update_user"); err != nil {
		return ParentStudentAccess{}, err
	}
	parent, err := d.repo.GetUser(ctx, tx, parentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentStudentAccess{}, ErrUserNotFound
	} else if err != nil {
		return ParentStudentAccess{}, fmt.Errorf("get parent: %w", err)
	}
	student, err := d.repo.GetUser(ctx, tx, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentStudentAccess{}, ErrUserNotFound
	} else if err != nil {
		return ParentStudentAccess{}, fmt.Errorf("get student: %w", err)
	}
	if parentID == studentID || !slices.Contains(parent.Roles, "parent") || !slices.Contains(student.Roles, "student") {
		return ParentStudentAccess{}, ErrInvalidParentStudent
	}
	codes := []string{"can_view_user", "can_view_journal"}
	value := studentID.String()
	for _, code := range codes {
		permission, err := d.repo.EnsurePermission(ctx, tx, entity.Permission{Code: code, Value: &value, Description: "Доступ родителя к данным ученика"})
		if err != nil {
			return ParentStudentAccess{}, fmt.Errorf("ensure %s permission: %w", code, err)
		}
		if _, err = d.repo.EnsureUserPermission(ctx, tx, entity.UserPermission{UserID: parentID, PermissionID: permission.ID}); err != nil {
			return ParentStudentAccess{}, fmt.Errorf("grant %s permission: %w", code, err)
		}
	}
	return ParentStudentAccess{ParentID: parentID, StudentID: studentID, Permissions: codes}, nil
}
