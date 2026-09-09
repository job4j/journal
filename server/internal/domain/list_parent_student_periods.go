package domain

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type ParentStudentPeriod struct {
	AcademicYear AcademicYearView
	Class        ClassView
	Subjects     []ClassSubjectView
}

func (d *ClassDomain) ListParentStudentPeriods(ctx context.Context, tx repository.Transaction, hash string, studentID uuid.UUID) ([]ParentStudentPeriod, error) {
	viewer, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUnauthenticated
	}
	if err != nil {
		return nil, fmt.Errorf("find parent: %w", err)
	}
	if !slices.Contains(viewer.Roles, "parent") {
		return nil, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_view_user", studentID.String()); err != nil {
		return nil, err
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_view_journal", studentID.String()); err != nil {
		return nil, err
	}
	student, err := d.repo.GetUser(ctx, tx, studentID)
	if err != nil || !slices.Contains(student.Roles, "student") {
		return nil, ErrUserNotFound
	}
	years, err := d.repo.ListAcademicYears(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list years: %w", err)
	}
	quarters, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list quarters: %w", err)
	}
	classes, err := d.repo.ListClasses(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list classes: %w", err)
	}
	memberships, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list memberships: %w", err)
	}
	assignments, err := d.repo.ListClassSubjects(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list assignments: %w", err)
	}
	subjects, err := d.repo.ListSubjects(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list subjects: %w", err)
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	result := []ParentStudentPeriod{}
	for _, membership := range memberships {
		if membership.UserID != studentID {
			continue
		}
		var class entity.Class
		for _, candidate := range classes {
			if candidate.ID == membership.ClassID {
				class = candidate
			}
		}
		if class.ID == uuid.Nil {
			continue
		}
		period := ParentStudentPeriod{Class: ClassView{Class: class}, Subjects: []ClassSubjectView{}}
		for _, year := range years {
			if year.ID == class.AcademicYearID {
				period.AcademicYear.Year = year
			}
		}
		for _, quarter := range quarters {
			if quarter.AcademicYearID == class.AcademicYearID {
				period.AcademicYear.Quarters = append(period.AcademicYear.Quarters, quarter)
			}
		}
		for _, assignment := range assignments {
			if assignment.ClassID != class.ID {
				continue
			}
			view := ClassSubjectView{Assignment: assignment}
			for _, subject := range subjects {
				if subject.ID == assignment.SubjectID {
					view.Subject = subject
				}
			}
			for _, user := range users {
				if user.ID == assignment.ResponsibleTeacherID {
					view.Teacher = user
				}
			}
			period.Subjects = append(period.Subjects, view)
		}
		result = append(result, period)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].AcademicYear.Year.StartsOn.After(result[j].AcademicYear.Year.StartsOn)
	})
	return result, nil
}
