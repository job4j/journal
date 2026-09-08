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

func (d *ClassDomain) UpdateClassSubject(ctx context.Context, tx repository.Transaction, hash string, assignmentID, teacherID uuid.UUID) (ClassSubjectView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_manage_class_subjects", assignmentID.String()); err != nil {
		return ClassSubjectView{}, err
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, assignmentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassSubjectView{}, ErrClassSubjectNotFound
	} else if err != nil {
		return ClassSubjectView{}, fmt.Errorf("get class subject: %w", err)
	}
	teacher, err := d.repo.GetUser(ctx, tx, teacherID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassSubjectView{}, ErrUserNotFound
	} else if err != nil {
		return ClassSubjectView{}, fmt.Errorf("get teacher: %w", err)
	}
	if !slices.Contains(teacher.Roles, "teacher") {
		return ClassSubjectView{}, ErrInvalidTeacher
	}
	subject, err := d.repo.GetSubject(ctx, tx, assignment.SubjectID)
	if err != nil {
		return ClassSubjectView{}, fmt.Errorf("get subject: %w", err)
	}
	previousTeacherID := assignment.ResponsibleTeacherID
	if previousTeacherID == teacherID {
		return ClassSubjectView{Assignment: assignment, Subject: subject, Teacher: teacher}, nil
	}
	assignment.ResponsibleTeacherID = teacherID
	assignment, err = d.repo.UpdateClassSubject(ctx, tx, assignment)
	if err != nil {
		return ClassSubjectView{}, fmt.Errorf("update class subject: %w", err)
	}
	if err := d.grantAssignmentPermissions(ctx, tx, teacherID, assignment.ClassID, assignment.ID); err != nil {
		return ClassSubjectView{}, err
	}
	if err := d.revokeAssignmentPermissions(ctx, tx, previousTeacherID, assignment); err != nil {
		return ClassSubjectView{}, err
	}
	return ClassSubjectView{Assignment: assignment, Subject: subject, Teacher: teacher}, nil
}

func (d *ClassDomain) revokeAssignmentPermissions(ctx context.Context, tx repository.Transaction, teacherID uuid.UUID, assignment entity.ClassSubject) error {
	permissions, err := d.repo.ListPermissions(ctx, tx)
	if err != nil {
		return fmt.Errorf("list permissions: %w", err)
	}
	keepClassAccess := false
	assignments, err := d.repo.ListClassSubjects(ctx, tx)
	if err != nil {
		return fmt.Errorf("list class subjects: %w", err)
	}
	for _, item := range assignments {
		if item.ID != assignment.ID && item.ClassID == assignment.ClassID && item.ResponsibleTeacherID == teacherID {
			keepClassAccess = true
			break
		}
	}
	for _, permission := range permissions {
		if permission.Value == nil {
			continue
		}
		remove := *permission.Value == assignment.ID.String() && slices.Contains(assignmentPermissionCodes, permission.Code)
		remove = remove || (!keepClassAccess && *permission.Value == assignment.ClassID.String() && permission.Code == "can_view_class")
		if remove {
			if err := d.repo.DeleteUserPermission(ctx, tx, teacherID, permission.ID); err != nil && !errors.Is(err, repository.ErrNotFound) {
				return fmt.Errorf("revoke %s permission: %w", permission.Code, err)
			}
		}
	}
	return nil
}
