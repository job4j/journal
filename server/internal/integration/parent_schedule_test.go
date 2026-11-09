//go:build integration

package integration

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func TestParentScheduleProtectsFamiliesMembershipAndHomework(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	repo := repository.NewRepository()
	d := domain.NewScheduleDomain(repo)
	parent, otherParent, limitedParent := user(t, tx, "parent"), user(t, tx, "parent"), user(t, tx, "parent")
	child, otherChild, teacher := user(t, tx, "student"), user(t, tx, "student"), user(t, tx, "teacher")
	for _, code := range []string{"can_view_user", "can_view_journal"} {
		permission(t, tx, parent, code, child)
		permission(t, tx, otherParent, code, otherChild)
	}
	permission(t, tx, limitedParent, "can_view_user", child)
	hash := session(t, tx, parent)
	yearID, classID, otherClass, subjectID, assignmentID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	_, err := tx.Exec(ctx, `INSERT INTO academic_years(id,name,starts_on,ends_on,status) VALUES($1,$2,'2026-09-01','2027-05-31','planned')`, yearID, yearID.String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO classes(id,academic_year_id,name,grade_level) VALUES($1,$2,'A',7),($3,$2,'B',7)`, classID, yearID, otherClass)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO class_students(class_id,user_id,enrolled_on,left_on) VALUES($1,$2,'2026-11-10','2026-11-10')`, classID, child)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO subjects(id,code,name) VALUES($1,$2,$2)`, subjectID, subjectID.String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(ctx, `INSERT INTO class_subjects(id,class_id,subject_id,responsible_teacher_id) VALUES($1,$2,$3,$4)`, assignmentID, classID, subjectID, teacher)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)
	week, err := repo.SaveScheduleWeek(ctx, tx, classID, start)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err = repo.CreateScheduleEntry(ctx, tx, entity.ScheduleEntry{ID: uuid.New(), WeekID: week.ID, LessonDate: start.AddDate(0, 0, i), Position: 1, ClassSubjectID: &assignmentID}); err != nil {
			t.Fatal(err)
		}
	}
	homework := "§ 2, упражнение 7"
	if _, err = repo.CreateLesson(ctx, tx, entity.Lesson{ClassSubjectID: assignmentID, LessonDate: start, Position: 1, Topic: "Тема", Homework: &homework, CreatedBy: teacher}); err != nil {
		t.Fatal(err)
	}
	req := domain.GetParentScheduleRequest{GetClassScheduleRequest: domain.GetClassScheduleRequest{ClassID: classID, WeekStart: start}, StudentID: child}
	result, err := d.GetParentSchedule(ctx, tx, hash, req)
	if err != nil || len(result.Entries) != 1 || !result.Entries[0].Entry.LessonDate.Equal(start.AddDate(0, 0, 1)) {
		t.Fatalf("membership: %+v %v", result, err)
	}
	if result.Entries[0].Homework != nil {
		t.Fatal("homework before enrollment leaked")
	}
	if _, err = d.ReplaceClassSchedule(ctx, tx, hash, domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: req.GetClassScheduleRequest}); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("parent could modify schedule: %v", err)
	}
	_, err = tx.Exec(ctx, `UPDATE class_students SET enrolled_on='2026-09-01' WHERE class_id=$1 AND user_id=$2`, classID, child)
	if err != nil {
		t.Fatal(err)
	}
	result, err = d.GetParentSchedule(ctx, tx, hash, req)
	if err != nil || len(result.Entries) != 2 || result.Entries[1].Homework == nil || *result.Entries[1].Homework != homework {
		t.Fatalf("homework: %+v %v", result, err)
	}
	for _, test := range []struct {
		token   string
		request domain.GetParentScheduleRequest
		want    error
	}{
		{"missing", req, domain.ErrUnauthenticated},
		{session(t, tx, otherParent), req, domain.ErrForbidden},
		{session(t, tx, limitedParent), req, domain.ErrForbidden},
		{session(t, tx, teacher), req, domain.ErrForbidden},
		{hash, domain.GetParentScheduleRequest{GetClassScheduleRequest: domain.GetClassScheduleRequest{ClassID: otherClass, WeekStart: start}, StudentID: child}, domain.ErrForbidden},
	} {
		if _, err = d.GetParentSchedule(ctx, tx, test.token, test.request); !errors.Is(err, test.want) {
			t.Fatalf("access want %v got %v", test.want, err)
		}
	}
}
