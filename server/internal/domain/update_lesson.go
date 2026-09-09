package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
)

func (d *ClassDomain) UpdateLesson(ctx context.Context, tx repository.Transaction, hash string, id uuid.UUID, input CreateLessonInput) (LessonView, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return LessonView{}, err
	}
	lesson, err := d.repo.GetLesson(ctx, tx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return LessonView{}, ErrLessonNotFound
	}
	if err != nil {
		return LessonView{}, fmt.Errorf("get lesson: %w", err)
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, lesson.ClassSubjectID)
	if err != nil {
		return LessonView{}, fmt.Errorf("get assignment: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return LessonView{}, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_update_lesson", assignment.ID.String()); err != nil {
		return LessonView{}, err
	}
	input.Topic = strings.TrimSpace(input.Topic)
	if input.Topic == "" || input.Position < 1 {
		return LessonView{}, ErrInvalidLesson
	}
	lesson.LessonDate = input.LessonDate
	lesson.Position = input.Position
	lesson.Topic = input.Topic
	lesson.Homework = input.Homework
	lesson, err = d.repo.UpdateLesson(ctx, tx, lesson)
	if err != nil {
		return LessonView{}, fmt.Errorf("update lesson: %w", err)
	}
	return LessonView{Lesson: lesson, Materials: []entity.LessonMaterial{}, GradeItems: []entity.GradeItem{}}, nil
}
