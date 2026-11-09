package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"slices"
	"time"
)

type GetParentScheduleRequest struct {
	GetClassScheduleRequest
	StudentID uuid.UUID
}
type ParentScheduleEntry struct {
	Entry                  entity.ScheduleEntry
	Homework, HomeworkFrom *string
}
type GetParentScheduleResponse struct {
	WeekStart time.Time
	Entries   []ParentScheduleEntry
}

func (d *ScheduleDomain) GetParentSchedule(ctx context.Context, tx repository.Transaction, hash string, req GetParentScheduleRequest) (GetParentScheduleResponse, error) {
	result := GetParentScheduleResponse{WeekStart: dayUTC(req.WeekStart), Entries: []ParentScheduleEntry{}}
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
	for _, permission := range []string{"can_view_user", "can_view_journal"} {
		if err = AuthorizeObject(ctx, tx, d.repo, hash, permission, req.StudentID.String()); err != nil {
			return result, err
		}
	}
	student, err := d.repo.GetUser(ctx, tx, req.StudentID)
	if err != nil {
		return result, err
	}
	if !slices.Contains(student.Roles, "student") {
		return result, ErrUserNotFound
	}
	class, err := d.repo.GetClass(ctx, tx, req.ClassID)
	if err != nil {
		return result, err
	}
	membership, err := d.repo.GetClassStudent(ctx, tx, req.ClassID, req.StudentID)
	if errors.Is(err, repository.ErrNotFound) {
		return result, ErrForbidden
	}
	if err != nil {
		return result, err
	}
	year, err := d.repo.GetAcademicYear(ctx, tx, class.AcademicYearID)
	if err != nil {
		return result, err
	}
	start := dayUTC(req.WeekStart)
	if start.Weekday() != time.Monday || start.After(dayUTC(year.EndsOn)) || start.AddDate(0, 0, 6).Before(dayUTC(year.StartsOn)) {
		return result, ErrInvalidSchedule
	}
	schedule, err := d.readSchedule(ctx, tx, req.GetClassScheduleRequest)
	if err != nil {
		return result, err
	}
	for _, entry := range schedule.Entries {
		day := dayUTC(entry.LessonDate)
		if day.Before(dayUTC(year.StartsOn)) || day.After(dayUTC(year.EndsOn)) || day.Before(dayUTC(membership.EnrolledOn)) || (membership.LeftOn != nil && day.After(dayUTC(*membership.LeftOn))) {
			continue
		}
		result.Entries = append(result.Entries, ParentScheduleEntry{Entry: entry})
	}
	from := dayUTC(year.StartsOn)
	if dayUTC(membership.EnrolledOn).After(from) {
		from = dayUTC(membership.EnrolledOn)
	}
	to := start.AddDate(0, 0, 6)
	if dayUTC(year.EndsOn).Before(to) {
		to = dayUTC(year.EndsOn)
	}
	if membership.LeftOn != nil && dayUTC(*membership.LeftOn).Before(to) {
		to = dayUTC(*membership.LeftOn)
	}
	if len(result.Entries) > 0 {
		lessons, readErr := d.repo.ListScheduleLessons(ctx, tx, req.ClassID, from, to)
		if readErr != nil {
			return result, readErr
		}
		fillScheduleHomework(result.Entries, lessons)
	}
	return result, nil
}
