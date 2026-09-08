package domain

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func (d *ClassDomain) AddStudentToClass(ctx context.Context, tx repository.Transaction, hash string, classID, studentID uuid.UUID, enrolledOn time.Time) (ClassStudentView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_manage_class_students", classID.String()); err != nil {
		return ClassStudentView{}, err
	}
	class, err := d.repo.GetClass(ctx, tx, classID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassStudentView{}, ErrClassNotFound
	}
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("get class: %w", err)
	}
	year, err := d.repo.GetAcademicYear(ctx, tx, class.AcademicYearID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassStudentView{}, ErrAcademicYearNotFound
	}
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("get academic year: %w", err)
	}
	student, err := d.repo.GetUser(ctx, tx, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassStudentView{}, ErrUserNotFound
	}
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("get student: %w", err)
	}
	if !slices.Contains(student.Roles, "student") || enrolledOn.Before(year.StartsOn) || enrolledOn.After(year.EndsOn) {
		return ClassStudentView{}, ErrInvalidClassStudent
	}
	membership, err := d.repo.CreateClassStudent(ctx, tx, entity.ClassStudent{ClassID: classID, UserID: studentID, EnrolledOn: enrolledOn})
	if errors.Is(err, repository.ErrConflict) {
		return ClassStudentView{}, ErrClassStudentExists
	}
	if err != nil {
		return ClassStudentView{}, fmt.Errorf("create class student: %w", err)
	}
	return ClassStudentView{Membership: membership, Student: student}, nil
}
