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

var assignmentPermissionCodes = []string{
	"can_view_class_subject",
	"can_create_lesson",
	"can_view_lesson",
	"can_update_lesson",
	"can_create_score",
	"can_view_score",
	"can_update_score",
}

func (d *ClassDomain) AssignSubjectToClass(ctx context.Context, tx repository.Transaction, hash string, classID, subjectID, teacherID uuid.UUID) (ClassSubjectView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_manage_class_subjects", classID.String()); err != nil {
		return ClassSubjectView{}, err
	}
	if _, err := d.repo.GetClass(ctx, tx, classID); errors.Is(err, repository.ErrNotFound) {
		return ClassSubjectView{}, ErrClassNotFound
	} else if err != nil {
		return ClassSubjectView{}, fmt.Errorf("get class: %w", err)
	}
	subject, err := d.repo.GetSubject(ctx, tx, subjectID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassSubjectView{}, ErrSubjectNotFound
	} else if err != nil {
		return ClassSubjectView{}, fmt.Errorf("get subject: %w", err)
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
	assignment, err := d.repo.CreateClassSubject(ctx, tx, entity.ClassSubject{ClassID: classID, SubjectID: subjectID, ResponsibleTeacherID: teacherID})
	if errors.Is(err, repository.ErrConflict) {
		return ClassSubjectView{}, ErrClassSubjectExists
	} else if err != nil {
		return ClassSubjectView{}, fmt.Errorf("create class subject: %w", err)
	}
	if err := d.grantAssignmentPermissions(ctx, tx, teacherID, classID, assignment.ID); err != nil {
		return ClassSubjectView{}, err
	}
	return ClassSubjectView{Assignment: assignment, Subject: subject, Teacher: teacher}, nil
}

func (d *ClassDomain) grantAssignmentPermissions(ctx context.Context, tx repository.Transaction, teacherID, classID, assignmentID uuid.UUID) error {
	values := []struct{ code, value string }{{"can_view_class", classID.String()}}
	for _, code := range assignmentPermissionCodes {
		values = append(values, struct{ code, value string }{code, assignmentID.String()})
	}
	for _, item := range values {
		permission, err := d.repo.EnsurePermission(ctx, tx, entity.Permission{Code: item.code, Value: &item.value, Description: "Объектный доступ по назначению учителя"})
		if err != nil {
			return fmt.Errorf("ensure %s permission: %w", item.code, err)
		}
		if _, err = d.repo.EnsureUserPermission(ctx, tx, entity.UserPermission{UserID: teacherID, PermissionID: permission.ID}); err != nil {
			return fmt.Errorf("grant %s permission: %w", item.code, err)
		}
	}
	return nil
}
