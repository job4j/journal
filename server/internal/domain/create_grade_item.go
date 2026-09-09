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

func (d *ClassDomain) CreateGradeItem(ctx context.Context, tx repository.Transaction, hash string, lessonID uuid.UUID, title, kind, scale string, maxScore *float64) (entity.GradeItem, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return entity.GradeItem{}, err
	}
	lesson, err := d.repo.GetLesson(ctx, tx, lessonID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.GradeItem{}, ErrLessonNotFound
	}
	if err != nil {
		return entity.GradeItem{}, fmt.Errorf("get lesson: %w", err)
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, lesson.ClassSubjectID)
	if err != nil {
		return entity.GradeItem{}, fmt.Errorf("get class subject: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return entity.GradeItem{}, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_create_score", assignment.ID.String()); err != nil {
		return entity.GradeItem{}, err
	}
	title = strings.TrimSpace(title)
	validKind := kind == "homework" || kind == "classwork" || kind == "knowledge_check" || kind == "other"
	validScale := scale == "five_point" || scale == "points" || scale == "pass_fail"
	if title == "" || !validKind || !validScale || (scale == "points" && (maxScore == nil || *maxScore <= 0)) || (scale != "points" && maxScore != nil) {
		return entity.GradeItem{}, ErrInvalidGradeItem
	}
	item, err := d.repo.CreateGradeItem(ctx, tx, entity.GradeItem{LessonID: lessonID, Title: title, Kind: kind, GradingScale: scale, MaxScore: maxScore})
	if err != nil {
		return entity.GradeItem{}, fmt.Errorf("create grade item: %w", err)
	}
	return item, nil
}
