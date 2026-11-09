package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"regexp"
	"strings"
	"unicode/utf8"
)

type ReplaceClassScheduleRequest struct {
	GetClassScheduleRequest
	Revision int
	Entries  []entity.ScheduleEntry
}
type ReplaceClassScheduleResponse struct{ GetClassScheduleResponse }

var scheduleTime = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

func (d *ScheduleDomain) ReplaceClassSchedule(ctx context.Context, tx repository.Transaction, hash string, req ReplaceClassScheduleRequest) (ReplaceClassScheduleResponse, error) {
	var result ReplaceClassScheduleResponse
	year, err := d.scheduleContext(ctx, tx, hash, req.GetClassScheduleRequest)
	if err != nil {
		return result, err
	}
	if len(req.Entries) > 140 || req.Revision < 0 {
		return result, ErrInvalidSchedule
	}
	// Serialize creation and updates, including a week which does not exist yet.
	if err = d.repo.LockScheduleClass(ctx, tx, req.ClassID); err != nil {
		return result, err
	}
	current, err := d.readSchedule(ctx, tx, req.GetClassScheduleRequest)
	if err != nil {
		return result, err
	}
	if current.Week.Revision != req.Revision {
		return result, ErrScheduleConflict
	}
	owned := map[uuid.UUID]bool{}
	seen := map[uuid.UUID]bool{}
	for _, e := range current.Entries {
		owned[e.ID] = true
	}
	for i := range req.Entries {
		e := &req.Entries[i]
		e.LessonDate = dayUTC(e.LessonDate)
		e.Title = strings.TrimSpace(e.Title)
		e.Notes = strings.TrimSpace(e.Notes)
		if e.ID != uuid.Nil && (!owned[e.ID] || seen[e.ID]) {
			return result, ErrInvalidSchedule
		}
		if e.ID == uuid.Nil {
			e.ID = uuid.New()
		}
		seen[e.ID] = true
		if e.LessonDate.Before(dayUTC(req.WeekStart)) || e.LessonDate.After(dayUTC(req.WeekStart).AddDate(0, 0, 6)) || e.LessonDate.Before(dayUTC(year.StartsOn)) || e.LessonDate.After(dayUTC(year.EndsOn)) || e.Position < 1 || e.Position > 20 || utf8.RuneCountInString(e.Title) > 200 || utf8.RuneCountInString(e.Notes) > 1000 {
			return result, ErrInvalidSchedule
		}
		if (e.StartsAt != nil && !scheduleTime.MatchString(*e.StartsAt)) || (e.EndsAt != nil && (!scheduleTime.MatchString(*e.EndsAt) || e.StartsAt == nil || *e.EndsAt <= *e.StartsAt)) {
			return result, ErrInvalidSchedule
		}
		if e.ClassSubjectID == nil {
			if e.Title == "" {
				return result, ErrInvalidSchedule
			}
		} else {
			subject, readErr := d.repo.GetClassSubject(ctx, tx, *e.ClassSubjectID)
			if readErr != nil {
				return result, readErr
			}
			if subject.ClassID != req.ClassID {
				return result, ErrInvalidSchedule
			}
		}
	}
	week, err := d.repo.SaveScheduleWeek(ctx, tx, req.ClassID, req.WeekStart)
	if err != nil {
		return result, err
	}
	if err = d.repo.DeleteScheduleEntries(ctx, tx, week.ID); err != nil {
		return result, err
	}
	for _, e := range req.Entries {
		e.WeekID = week.ID
		if err = d.repo.CreateScheduleEntry(ctx, tx, e); err != nil {
			return result, err
		}
	}
	result.GetClassScheduleResponse, err = d.readSchedule(ctx, tx, req.GetClassScheduleRequest)
	return result, err
}
