package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"net/url"
	"strings"
	"time"
)

type LessonMaterialInput struct {
	Title, URL string
	Position   int16
}
type CreateLessonInput struct {
	LessonDate time.Time
	Position   int16
	Topic      string
	Homework   *string
	Materials  []LessonMaterialInput
}

func (d *ClassDomain) CreateLesson(ctx context.Context, tx repository.Transaction, hash string, assignmentID uuid.UUID, input CreateLessonInput) (LessonView, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return LessonView{}, err
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, assignmentID)
	if errors.Is(err, repository.ErrNotFound) {
		return LessonView{}, ErrClassSubjectNotFound
	}
	if err != nil {
		return LessonView{}, fmt.Errorf("get class subject: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return LessonView{}, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_create_lesson", assignmentID.String()); err != nil {
		return LessonView{}, err
	}
	class, err := d.repo.GetClass(ctx, tx, assignment.ClassID)
	if err != nil {
		return LessonView{}, fmt.Errorf("get class: %w", err)
	}
	year, err := d.repo.GetAcademicYear(ctx, tx, class.AcademicYearID)
	if err != nil {
		return LessonView{}, fmt.Errorf("get academic year: %w", err)
	}
	input.Topic = strings.TrimSpace(input.Topic)
	if input.Topic == "" || input.Position < 1 || input.LessonDate.Before(year.StartsOn) || input.LessonDate.After(year.EndsOn) {
		return LessonView{}, ErrInvalidLesson
	}
	seen := map[int16]struct{}{}
	for i := range input.Materials {
		material := &input.Materials[i]
		material.Title = strings.TrimSpace(material.Title)
		material.URL = strings.TrimSpace(material.URL)
		parsed, parseErr := url.ParseRequestURI(material.URL)
		if material.Title == "" || material.Position < 1 || parseErr != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return LessonView{}, ErrInvalidLesson
		}
		if _, ok := seen[material.Position]; ok {
			return LessonView{}, ErrInvalidLesson
		}
		seen[material.Position] = struct{}{}
	}
	lesson, err := d.repo.CreateLesson(ctx, tx, entity.Lesson{ClassSubjectID: assignmentID, LessonDate: input.LessonDate, Position: input.Position, Topic: input.Topic, Homework: input.Homework, CreatedBy: teacherID})
	if errors.Is(err, repository.ErrConflict) {
		return LessonView{}, ErrLessonExists
	}
	if err != nil {
		return LessonView{}, fmt.Errorf("create lesson: %w", err)
	}
	materials := make([]entity.LessonMaterial, 0, len(input.Materials))
	for _, value := range input.Materials {
		item, createErr := d.repo.CreateLessonMaterial(ctx, tx, entity.LessonMaterial{LessonID: lesson.ID, Title: value.Title, URL: value.URL, Position: value.Position})
		if createErr != nil {
			return LessonView{}, fmt.Errorf("create lesson material: %w", createErr)
		}
		materials = append(materials, item)
	}
	return LessonView{Lesson: lesson, Materials: materials, GradeItems: []entity.GradeItem{}}, nil
}
