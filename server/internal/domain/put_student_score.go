package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"math"
	"strings"
	"time"
)

func (d *ClassDomain) PutStudentScore(ctx context.Context, tx repository.Transaction, hash string, gradeItemID, studentID uuid.UUID, numeric *float64, text, comment *string) (entity.Score, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return entity.Score{}, err
	}
	item, err := d.repo.GetGradeItem(ctx, tx, gradeItemID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.Score{}, ErrGradeItemNotFound
	}
	if err != nil {
		return entity.Score{}, fmt.Errorf("get grade item: %w", err)
	}
	lesson, err := d.repo.GetLesson(ctx, tx, item.LessonID)
	if err != nil {
		return entity.Score{}, fmt.Errorf("get lesson: %w", err)
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, lesson.ClassSubjectID)
	if err != nil {
		return entity.Score{}, fmt.Errorf("get class subject: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return entity.Score{}, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_create_score", assignment.ID.String()); err != nil {
		return entity.Score{}, err
	}
	membership, err := d.repo.GetClassStudent(ctx, tx, assignment.ClassID, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.Score{}, ErrClassStudentNotFound
	}
	if err != nil {
		return entity.Score{}, fmt.Errorf("get class student: %w", err)
	}
	lessonDay := dayUTC(lesson.LessonDate)
	if lessonDay.Before(dayUTC(membership.EnrolledOn)) || (membership.LeftOn != nil && lessonDay.After(dayUTC(*membership.LeftOn))) {
		return entity.Score{}, ErrClassStudentNotFound
	}
	if !validScore(item, numeric, text) {
		return entity.Score{}, ErrInvalidScore
	}
	comment = trimmedOptional(comment)
	if text != nil {
		value := strings.TrimSpace(*text)
		text = &value
	}
	result, err := d.repo.UpsertScore(ctx, tx, entity.Score{GradeItemID: gradeItemID, UserID: studentID, NumericValue: numeric, TextValue: text, TeacherComment: comment, CreatedBy: teacherID, UpdatedBy: teacherID})
	if err != nil {
		return entity.Score{}, fmt.Errorf("upsert score: %w", err)
	}
	return result, nil
}

func dayUTC(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func validScore(item entity.GradeItem, numeric *float64, text *string) bool {
	if numeric != nil && (math.IsNaN(*numeric) || math.IsInf(*numeric, 0)) {
		return false
	}
	switch item.GradingScale {
	case "five_point":
		return numeric != nil && text == nil && *numeric >= 2 && *numeric <= 5 && math.Trunc(*numeric) == *numeric
	case "points":
		return numeric != nil && text == nil && item.MaxScore != nil && *numeric >= 0 && *numeric <= *item.MaxScore
	case "pass_fail":
		return numeric == nil && text != nil && (*text == "pass" || *text == "fail")
	default:
		return false
	}
}

func trimmedOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
