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

type ParentJournal struct {
	Student      entity.User
	AcademicYear AcademicYearView
	Class        ClassView
	Subjects     []ParentJournalSubject
}
type ParentJournalSubject struct {
	ClassSubjectView
	Lessons []LessonView
}

func (d *ClassDomain) GetParentStudentJournal(ctx context.Context, tx repository.Transaction, hash string, studentID, yearID uuid.UUID) (ParentJournal, error) {
	viewer, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentJournal{}, ErrUnauthenticated
	}
	if err != nil {
		return ParentJournal{}, fmt.Errorf("find parent: %w", err)
	}
	if !slices.Contains(viewer.Roles, "parent") {
		return ParentJournal{}, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_view_user", studentID.String()); err != nil {
		return ParentJournal{}, err
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, "can_view_journal", studentID.String()); err != nil {
		return ParentJournal{}, err
	}
	student, err := d.repo.GetUser(ctx, tx, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentJournal{}, ErrUserNotFound
	}
	if err != nil {
		return ParentJournal{}, fmt.Errorf("get student: %w", err)
	}
	if !slices.Contains(student.Roles, "student") {
		return ParentJournal{}, ErrUserNotFound
	}
	year, err := d.repo.GetAcademicYear(ctx, tx, yearID)
	if errors.Is(err, repository.ErrNotFound) {
		return ParentJournal{}, ErrAcademicYearNotFound
	}
	if err != nil {
		return ParentJournal{}, fmt.Errorf("get year: %w", err)
	}
	classes, err := d.repo.ListClasses(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list classes: %w", err)
	}
	memberships, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list memberships: %w", err)
	}
	var selected entity.Class
	var enrolled entity.ClassStudent
	for _, membership := range memberships {
		if membership.UserID != studentID {
			continue
		}
		for _, class := range classes {
			if class.ID == membership.ClassID && class.AcademicYearID == yearID && (selected.ID == uuid.Nil || membership.EnrolledOn.After(enrolled.EnrolledOn)) {
				selected = class
				enrolled = membership
			}
		}
	}
	if selected.ID == uuid.Nil {
		return ParentJournal{}, ErrClassStudentNotFound
	}
	quarters, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list quarters: %w", err)
	}
	yearView := AcademicYearView{Year: year}
	for _, q := range quarters {
		if q.AcademicYearID == yearID {
			yearView.Quarters = append(yearView.Quarters, q)
		}
	}
	assignments, err := d.repo.ListClassSubjects(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list class subjects: %w", err)
	}
	subjects, err := d.repo.ListSubjects(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list subjects: %w", err)
	}
	users, err := d.repo.ListUsers(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list users: %w", err)
	}
	lessons, err := d.repo.ListLessons(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list lessons: %w", err)
	}
	materials, err := d.repo.ListLessonMaterials(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list materials: %w", err)
	}
	items, err := d.repo.ListGradeItems(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list grade items: %w", err)
	}
	scores, err := d.repo.ListScores(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list scores: %w", err)
	}
	absences, err := d.repo.ListAbsences(ctx, tx)
	if err != nil {
		return ParentJournal{}, fmt.Errorf("list absences: %w", err)
	}
	result := ParentJournal{Student: student, AcademicYear: yearView, Class: ClassView{Class: selected}, Subjects: []ParentJournalSubject{}}
	for _, assignment := range assignments {
		if assignment.ClassID != selected.ID {
			continue
		}
		view := ParentJournalSubject{ClassSubjectView: ClassSubjectView{Assignment: assignment}, Lessons: []LessonView{}}
		for _, s := range subjects {
			if s.ID == assignment.SubjectID {
				view.Subject = s
			}
		}
		for _, u := range users {
			if u.ID == assignment.ResponsibleTeacherID {
				view.Teacher = u
			}
		}
		for _, lesson := range lessons {
			if lesson.ClassSubjectID != assignment.ID {
				continue
			}
			lessonView := LessonView{Lesson: lesson, Materials: []entity.LessonMaterial{}, GradeItems: []entity.GradeItem{}, Scores: map[uuid.UUID][]entity.Score{}, Absences: []entity.Absence{}}
			for _, m := range materials {
				if m.LessonID == lesson.ID {
					lessonView.Materials = append(lessonView.Materials, m)
				}
			}
			for _, item := range items {
				if item.LessonID == lesson.ID {
					lessonView.GradeItems = append(lessonView.GradeItems, item)
					for _, score := range scores {
						if score.GradeItemID == item.ID && score.UserID == studentID {
							lessonView.Scores[item.ID] = append(lessonView.Scores[item.ID], score)
						}
					}
				}
			}
			for _, absence := range absences {
				if absence.LessonID == lesson.ID && absence.UserID == studentID {
					lessonView.Absences = append(lessonView.Absences, absence)
				}
			}
			view.Lessons = append(view.Lessons, lessonView)
		}
		sort.Slice(view.Lessons, func(i, j int) bool {
			return view.Lessons[i].Lesson.LessonDate.Before(view.Lessons[j].Lesson.LessonDate)
		})
		result.Subjects = append(result.Subjects, view)
	}
	return result, nil
}
