package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"slices"
	"sort"
	"time"
)

type LessonView struct {
	Lesson     entity.Lesson
	Materials  []entity.LessonMaterial
	GradeItems []entity.GradeItem
	Scores     map[uuid.UUID][]entity.Score
	Absences   []entity.Absence
}

func (d *ClassDomain) ListClassSubjectLessons(ctx context.Context, tx repository.Transaction, hash string, assignmentID uuid.UUID, dateFrom, dateTo *time.Time) ([]LessonView, error) {
	user, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUnauthenticated
	}
	if err != nil {
		return nil, fmt.Errorf("find session user: %w", err)
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, assignmentID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrClassSubjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get class subject: %w", err)
	}
	if slices.Contains(user.Roles, "teacher") && assignment.ResponsibleTeacherID != user.ID {
		return nil, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_view_lesson", assignmentID.String()); err != nil {
		return nil, err
	}
	if dateFrom != nil && dateTo != nil && dateFrom.After(*dateTo) {
		return nil, ErrInvalidClass
	}
	lessons, err := d.repo.ListLessons(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list lessons: %w", err)
	}
	materials, err := d.repo.ListLessonMaterials(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list lesson materials: %w", err)
	}
	gradeItems, err := d.repo.ListGradeItems(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list grade items: %w", err)
	}
	scores, err := d.repo.ListScores(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list scores: %w", err)
	}
	absences, err := d.repo.ListAbsences(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list absences: %w", err)
	}
	result := []LessonView{}
	for _, lesson := range lessons {
		if lesson.ClassSubjectID != assignmentID || (dateFrom != nil && lesson.LessonDate.Before(*dateFrom)) || (dateTo != nil && lesson.LessonDate.After(*dateTo)) {
			continue
		}
		view := LessonView{Lesson: lesson, Materials: []entity.LessonMaterial{}, GradeItems: []entity.GradeItem{}, Scores: map[uuid.UUID][]entity.Score{}}
		for _, item := range materials {
			if item.LessonID == lesson.ID {
				view.Materials = append(view.Materials, item)
			}
		}
		for _, item := range gradeItems {
			if item.LessonID == lesson.ID {
				view.GradeItems = append(view.GradeItems, item)
				for _, score := range scores {
					if score.GradeItemID == item.ID {
						view.Scores[item.ID] = append(view.Scores[item.ID], score)
					}
				}
			}
		}
		for _, absence := range absences {
			if absence.LessonID == lesson.ID {
				view.Absences = append(view.Absences, absence)
			}
		}
		sort.Slice(view.Materials, func(i, j int) bool { return view.Materials[i].Position < view.Materials[j].Position })
		result = append(result, view)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Lesson.LessonDate.Equal(result[j].Lesson.LessonDate) {
			return result[i].Lesson.Position < result[j].Lesson.Position
		}
		return result[i].Lesson.LessonDate.Before(result[j].Lesson.LessonDate)
	})
	return result, nil
}
