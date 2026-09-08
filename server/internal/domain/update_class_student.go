package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository"
)

func (d *ClassDomain) UpdateClassStudent(ctx context.Context, tx repository.Transaction, hash string, classID, studentID uuid.UUID, leftOn *time.Time) (ClassStudentView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_manage_class_students", classID.String()); err != nil {
		return ClassStudentView{}, err
	}
	membership, err := d.repo.GetClassStudent(ctx, tx, classID, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassStudentView{}, ErrClassStudentNotFound
	}
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("get class student: %w", err)
	}
	if leftOn != nil && leftOn.Before(membership.EnrolledOn) {
		return ClassStudentView{}, ErrInvalidClassStudent
	}
	membership.LeftOn = leftOn
	membership, err = d.repo.UpdateClassStudent(ctx, tx, membership)
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("update class student: %w", err)
	}
	student, err := d.repo.GetUser(ctx, tx, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassStudentView{}, ErrUserNotFound
	}
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("get student: %w", err)
	}
	return ClassStudentView{Membership: membership, Student: student}, nil
}
