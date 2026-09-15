package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/repository"
)

func (d *ClassDomain) DeleteStudentScore(ctx context.Context, tx repository.Transaction, hash string, gradeItemID, studentID uuid.UUID) error {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return err
	}
	item, err := d.repo.GetGradeItem(ctx, tx, gradeItemID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrGradeItemNotFound
	}
	if err != nil {
		return fmt.Errorf("get grade item: %w", err)
	}
	lesson, err := d.repo.GetLesson(ctx, tx, item.LessonID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrLessonNotFound
	}
	if err != nil {
		return fmt.Errorf("get lesson: %w", err)
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, lesson.ClassSubjectID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrClassSubjectNotFound
	}
	if err != nil {
		return fmt.Errorf("get class subject: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_update_score", assignment.ID.String()); err != nil {
		return err
	}
	if err = d.repo.DeleteStudentScore(ctx, tx, gradeItemID, studentID); errors.Is(err, repository.ErrNotFound) {
		return ErrScoreNotFound
	}
	if err != nil {
		return fmt.Errorf("delete student score: %w", err)
	}
	return nil
}
