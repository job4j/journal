package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"slices"
	"sort"
)

type ListParentMissedTopicsRequest struct{ StudentID, AcademicYearID uuid.UUID }
type ParentMissedTopic struct {
	Lesson  entity.Lesson
	Subject entity.Subject
	Class   entity.Class
}
type ListParentMissedTopicsResponse struct{ Items []ParentMissedTopic }

func (d *ClassDomain) ListParentMissedTopics(ctx context.Context, tx repository.Transaction, hash string, req ListParentMissedTopicsRequest) (ListParentMissedTopicsResponse, error) {
	result := ListParentMissedTopicsResponse{Items: []ParentMissedTopic{}}
	viewer, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return result, ErrUnauthenticated
	}
	if err != nil {
		return result, err
	}
	if !slices.Contains(viewer.Roles, "parent") {
		return result, ErrForbidden
	}
	for _, code := range []string{"can_view_user", "can_view_journal"} {
		if err = AuthorizeObject(ctx, tx, d.repo, hash, code, req.StudentID.String()); err != nil {
			return result, err
		}
	}
	student, err := d.repo.GetUser(ctx, tx, req.StudentID)
	if errors.Is(err, repository.ErrNotFound) {
		return result, ErrUserNotFound
	}
	if err != nil {
		return result, err
	}
	if !slices.Contains(student.Roles, "student") {
		return result, ErrUserNotFound
	}
	year, err := d.repo.GetAcademicYear(ctx, tx, req.AcademicYearID)
	if errors.Is(err, repository.ErrNotFound) {
		return result, ErrAcademicYearNotFound
	}
	if err != nil {
		return result, err
	}
	classes, err := d.repo.ListClasses(ctx, tx)
	if err != nil {
		return result, err
	}
	memberships, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return result, err
	}
	classByID := map[uuid.UUID]entity.Class{}
	for _, item := range classes {
		if item.AcademicYearID == req.AcademicYearID {
			classByID[item.ID] = item
		}
	}
	membershipByClass := map[uuid.UUID]entity.ClassStudent{}
	for _, item := range memberships {
		if _, ok := classByID[item.ClassID]; ok && item.UserID == req.StudentID {
			membershipByClass[item.ClassID] = item
		}
	}
	if len(membershipByClass) == 0 {
		return result, ErrClassStudentNotFound
	}
	assignments, err := d.repo.ListClassSubjects(ctx, tx)
	if err != nil {
		return result, err
	}
	subjects, err := d.repo.ListSubjects(ctx, tx)
	if err != nil {
		return result, err
	}
	lessons, err := d.repo.ListLessons(ctx, tx)
	if err != nil {
		return result, err
	}
	absences, err := d.repo.ListAbsences(ctx, tx)
	if err != nil {
		return result, err
	}
	subjectByID := map[uuid.UUID]entity.Subject{}
	for _, item := range subjects {
		subjectByID[item.ID] = item
	}
	assignmentByID := map[uuid.UUID]entity.ClassSubject{}
	for _, item := range assignments {
		if _, ok := membershipByClass[item.ClassID]; ok {
			assignmentByID[item.ID] = item
		}
	}
	missed := map[uuid.UUID]bool{}
	for _, item := range absences {
		if item.UserID == req.StudentID {
			missed[item.LessonID] = true
		}
	}
	for _, lesson := range lessons {
		assignment, ok := assignmentByID[lesson.ClassSubjectID]
		if !ok || !missed[lesson.ID] {
			continue
		}
		membership := membershipByClass[assignment.ClassID]
		date := dayUTC(lesson.LessonDate)
		if date.Before(dayUTC(year.StartsOn)) || date.After(dayUTC(year.EndsOn)) || date.Before(dayUTC(membership.EnrolledOn)) || (membership.LeftOn != nil && date.After(dayUTC(*membership.LeftOn))) {
			continue
		}
		result.Items = append(result.Items, ParentMissedTopic{Lesson: lesson, Subject: subjectByID[assignment.SubjectID], Class: classByID[assignment.ClassID]})
	}
	sort.Slice(result.Items, func(i, j int) bool {
		a, b := result.Items[i], result.Items[j]
		if !a.Lesson.LessonDate.Equal(b.Lesson.LessonDate) {
			return a.Lesson.LessonDate.After(b.Lesson.LessonDate)
		}
		if a.Lesson.Position != b.Lesson.Position {
			return a.Lesson.Position > b.Lesson.Position
		}
		if a.Subject.Name != b.Subject.Name {
			return a.Subject.Name < b.Subject.Name
		}
		return a.Lesson.ID.String() < b.Lesson.ID.String()
	})
	return result, nil
}
