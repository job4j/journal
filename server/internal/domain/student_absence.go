package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

func (d *ClassDomain) PutStudentAbsence(ctx context.Context, tx repository.Transaction, hash string, lessonID, studentID uuid.UUID) (entity.Absence, error) {
	teacherID, err := d.authorizeLessonStudent(ctx, tx, hash, lessonID, studentID)
	if err != nil {
		return entity.Absence{}, err
	}
	result, err := d.repo.EnsureAbsence(ctx, tx, entity.Absence{LessonID: lessonID, UserID: studentID, RecordedBy: teacherID})
	if err != nil {
		return entity.Absence{}, fmt.Errorf("ensure absence: %w", err)
	}
	return result, nil
}

func (d *ClassDomain) DeleteStudentAbsence(ctx context.Context, tx repository.Transaction, hash string, lessonID, studentID uuid.UUID) error {
	if _, err := d.authorizeLessonStudent(ctx, tx, hash, lessonID, studentID); err != nil {
		return err
	}
	if err := d.repo.DeleteAbsence(ctx, tx, lessonID, studentID); err != nil {
		return fmt.Errorf("delete absence: %w", err)
	}
	return nil
}

func (d *ClassDomain) authorizeLessonStudent(ctx context.Context, tx repository.Transaction, hash string, lessonID, studentID uuid.UUID) (uuid.UUID, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return uuid.Nil, err
	}
	lesson, err := d.repo.GetLesson(ctx, tx, lessonID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, ErrLessonNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get lesson: %w", err)
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, lesson.ClassSubjectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get class subject: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return uuid.Nil, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_create_score", assignment.ID.String()); err != nil {
		return uuid.Nil, err
	}
	membership, err := d.repo.GetClassStudent(ctx, tx, assignment.ClassID, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, ErrClassStudentNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get class student: %w", err)
	}
	lessonDay := dayUTC(lesson.LessonDate)
	if lessonDay.Before(dayUTC(membership.EnrolledOn)) || (membership.LeftOn != nil && lessonDay.After(dayUTC(*membership.LeftOn))) {
		return uuid.Nil, ErrClassStudentNotFound
	}
	return teacherID, nil
}
