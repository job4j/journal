//go:build integration

package integration

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClassScheduleAccessIsolationAndValidation(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	repo := repository.NewRepository()
	d := domain.NewScheduleDomain(repo)
	admin, teacher, parent := user(t, tx, "admin"), user(t, tx, "teacher"), user(t, tx, "parent")
	hash := session(t, tx, admin)
	yearID, classID, otherClass, subjectID, assignmentID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	_, err := tx.Exec(ctx, `INSERT INTO academic_years(id,name,starts_on,ends_on,status) VALUES($1,$2,'2026-09-01','2027-05-31','planned');`, yearID, yearID.String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO classes(id,academic_year_id,name,grade_level) VALUES($1,$2,'A',7),($3,$2,'B',7)`, classID, yearID, otherClass)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO subjects(id,code,name) VALUES($1,$2,'Математика')`, subjectID, subjectID.String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO class_subjects(id,class_id,subject_id,responsible_teacher_id) VALUES($1,$2,$3,$4)`, assignmentID, classID, subjectID, teacher)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	req := domain.GetClassScheduleRequest{ClassID: classID, WeekStart: start}
	for _, test := range []struct {
		token string
		want  error
	}{{"missing", domain.ErrUnauthenticated}, {session(t, tx, teacher), domain.ErrForbidden}, {session(t, tx, parent), domain.ErrForbidden}} {
		if _, err = d.GetClassSchedule(ctx, tx, test.token, req); !errors.Is(err, test.want) {
			t.Fatalf("read access: %v", err)
		}
		if _, err = d.ReplaceClassSchedule(ctx, tx, test.token, domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: req}); !errors.Is(err, test.want) {
			t.Fatalf("write access: %v", err)
		}
	}
	empty, err := d.GetClassSchedule(ctx, tx, hash, req)
	if err != nil || empty.Week.Revision != 0 || len(empty.Entries) != 0 {
		t.Fatalf("empty week: %+v %v", empty, err)
	}
	entry := entity.ScheduleEntry{LessonDate: start.AddDate(0, 0, 1), Position: 1, ClassSubjectID: &assignmentID}
	create := domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: req, Entries: []entity.ScheduleEntry{entry}}
	first, err := d.ReplaceClassSchedule(ctx, tx, hash, create)
	if err != nil || first.Week.Revision != 1 || len(first.Entries) != 1 {
		t.Fatalf("create: %+v %v", first, err)
	}
	if first.Entries[0].Title != "Математика" {
		t.Fatal("linked subject name not returned")
	}
	if _, err = d.ReplaceClassSchedule(ctx, tx, hash, create); !errors.Is(err, domain.ErrScheduleConflict) {
		t.Fatalf("stale revision: %v", err)
	}
	next := domain.GetClassScheduleRequest{ClassID: classID, WeekStart: start.AddDate(0, 0, 7)}
	extra := entity.ScheduleEntry{LessonDate: next.WeekStart, Position: 1, Title: "Шахматы", Notes: "Классный час раз в месяц"}
	nextSaved, err := d.ReplaceClassSchedule(ctx, tx, hash, domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: next, Entries: []entity.ScheduleEntry{extra}})
	if err != nil {
		t.Fatal(err)
	}
	kept, err := d.GetClassSchedule(ctx, tx, hash, req)
	if err != nil || kept.Entries[0].ID != first.Entries[0].ID {
		t.Fatal("adjacent week changed")
	}
	for _, test := range []struct {
		name   string
		change func(*domain.ReplaceClassScheduleRequest)
	}{
		{"outside year", func(r *domain.ReplaceClassScheduleRequest) { r.Entries[0].LessonDate = start }},
		{"outside week", func(r *domain.ReplaceClassScheduleRequest) { r.Entries[0].LessonDate = next.WeekStart }},
		{"wrong class subject", func(r *domain.ReplaceClassScheduleRequest) {
			r.ClassID = otherClass
			r.Revision = 0
			r.Entries[0].ID = uuid.Nil
		}},
		{"foreign entry ID", func(r *domain.ReplaceClassScheduleRequest) { r.Entries[0].ID = nextSaved.Entries[0].ID }},
		{"bad position", func(r *domain.ReplaceClassScheduleRequest) { r.Entries[0].Position = 0 }},
		{"empty activity", func(r *domain.ReplaceClassScheduleRequest) {
			r.Entries[0].ClassSubjectID = nil
			r.Entries[0].Title = " "
		}},
		{"invalid time", func(r *domain.ReplaceClassScheduleRequest) { v := "25:00"; r.Entries[0].StartsAt = &v }},
		{"reversed times", func(r *domain.ReplaceClassScheduleRequest) {
			a, b := "14:00", "13:00"
			r.Entries[0].StartsAt = &a
			r.Entries[0].EndsAt = &b
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			change := domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: req, Revision: 1, Entries: []entity.ScheduleEntry{first.Entries[0]}}
			test.change(&change)
			if _, err = d.ReplaceClassSchedule(ctx, tx, hash, change); !errors.Is(err, domain.ErrInvalidSchedule) {
				t.Fatalf("validation: %v", err)
			}
		})
	}
	if _, err = d.GetClassSchedule(ctx, tx, hash, domain.GetClassScheduleRequest{ClassID: classID, WeekStart: start.AddDate(0, 0, 1)}); !errors.Is(err, domain.ErrInvalidSchedule) {
		t.Fatalf("Tuesday week: %v", err)
	}
	// Two groups may share the same day and lesson number.
	extras := []entity.ScheduleEntry{nextSaved.Entries[0], {LessonDate: next.WeekStart, Position: 1, Title: "Труд", Notes: "Мальчики"}}
	changed, err := d.ReplaceClassSchedule(ctx, tx, hash, domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: next, Revision: 1, Entries: extras})
	if err != nil || len(changed.Entries) != 2 {
		t.Fatalf("groups: %v", err)
	}
	cleared, err := d.ReplaceClassSchedule(ctx, tx, hash, domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: next, Revision: 2, Entries: []entity.ScheduleEntry{}})
	if err != nil || len(cleared.Entries) != 0 {
		t.Fatalf("clear: %v", err)
	}
}

func TestSeededScheduleIncludesHolidaysAndPreservesEdits(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	var weeks, entries int
	if err := tx.QueryRow(ctx, `SELECT (SELECT count(*) FROM class_schedule_weeks WHERE class_id='60000000-0000-4000-8000-000000000001'),(SELECT count(*) FROM class_schedule_entries e JOIN class_schedule_weeks w ON w.id=e.week_id WHERE w.class_id='60000000-0000-4000-8000-000000000001')`).Scan(&weeks, &entries); err != nil {
		t.Fatal(err)
	}
	if weeks != 40 || entries != 1170 {
		t.Fatalf("weeks=%d entries=%d", weeks, entries)
	}
	var holidays, boundaries bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM class_schedule_entries WHERE lesson_date='2027-01-01'),NOT EXISTS(SELECT 1 FROM class_schedule_entries e JOIN class_schedule_weeks w ON w.id=e.week_id WHERE w.class_id='60000000-0000-4000-8000-000000000001' AND (e.lesson_date<'2026-09-01' OR e.lesson_date>'2027-05-31'))`).Scan(&holidays, &boundaries); err != nil || !holidays || !boundaries {
		t.Fatalf("dates: %v", err)
	}
	_, err := tx.Exec(ctx, `UPDATE class_schedule_entries SET notes='Changed by administrator' WHERE id=(SELECT id FROM class_schedule_entries LIMIT 1); CREATE TEMP TABLE schedule_before AS TABLE class_schedule_entries;`)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("../../migrations/20261109085441_seed_grade_7_weekly_schedule.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.Exec(ctx, strings.SplitN(string(data), "-- +goose Down", 2)[0]); err != nil {
		t.Fatal(err)
	}
	var same bool
	if err = tx.QueryRow(ctx, `SELECT NOT EXISTS ((TABLE class_schedule_entries EXCEPT TABLE schedule_before) UNION ALL (TABLE schedule_before EXCEPT TABLE class_schedule_entries))`).Scan(&same); err != nil || !same {
		t.Fatalf("seed overwrote edits: %v", err)
	}
}
