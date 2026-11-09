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

type GetClassScheduleRequest struct {
	ClassID   uuid.UUID
	WeekStart time.Time
}
type GetClassScheduleResponse struct {
	Week    entity.ScheduleWeek
	Entries []entity.ScheduleEntry
}

func (d *ScheduleDomain) scheduleContext(ctx context.Context, tx repository.Transaction, hash string, req GetClassScheduleRequest) (entity.AcademicYear, error) {
	user, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.AcademicYear{}, ErrUnauthenticated
	}
	if err != nil {
		return entity.AcademicYear{}, err
	}
	if !slices.Contains(user.Roles, "admin") {
		return entity.AcademicYear{}, ErrForbidden
	}
	class, err := d.repo.GetClass(ctx, tx, req.ClassID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.AcademicYear{}, ErrClassNotFound
	}
	if err != nil {
		return entity.AcademicYear{}, err
	}
	year, err := d.repo.GetAcademicYear(ctx, tx, class.AcademicYearID)
	if err != nil {
		return year, err
	}
	start := dayUTC(req.WeekStart)
	if req.ClassID == uuid.Nil || start.Weekday() != time.Monday || start.After(dayUTC(year.EndsOn)) || start.AddDate(0, 0, 6).Before(dayUTC(year.StartsOn)) {
		return year, ErrInvalidSchedule
	}
	return year, nil
}
func (d *ScheduleDomain) readSchedule(ctx context.Context, tx repository.Transaction, req GetClassScheduleRequest) (GetClassScheduleResponse, error) {
	result := GetClassScheduleResponse{Week: entity.ScheduleWeek{ClassID: req.ClassID, WeekStart: dayUTC(req.WeekStart)}, Entries: []entity.ScheduleEntry{}}
	week, err := d.repo.GetScheduleWeek(ctx, tx, req.ClassID, req.WeekStart)
	if errors.Is(err, repository.ErrNotFound) {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	result.Week = week
	result.Entries, err = d.repo.ListScheduleEntries(ctx, tx, week.ID)
	if result.Entries == nil {
		result.Entries = []entity.ScheduleEntry{}
	}
	return result, err
}
func (d *ScheduleDomain) GetClassSchedule(ctx context.Context, tx repository.Transaction, hash string, req GetClassScheduleRequest) (GetClassScheduleResponse, error) {
	if _, err := d.scheduleContext(ctx, tx, hash, req); err != nil {
		return GetClassScheduleResponse{}, err
	}
	return d.readSchedule(ctx, tx, req)
}
