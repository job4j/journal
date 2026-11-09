package domain

import (
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func TestScheduleHomeworkUsesPreviousDateAndOnlyFirstDailyLesson(t *testing.T) {
	subject, other := uuid.New(), uuid.New()
	day := time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)
	text := func(s string) *string { return &s }
	entries := []ParentScheduleEntry{
		{Entry: entity.ScheduleEntry{ID: uuid.New(), ClassSubjectID: &subject, LessonDate: day, Position: 4}},
		{Entry: entity.ScheduleEntry{ID: uuid.New(), ClassSubjectID: &subject, LessonDate: day, Position: 3}},
		{Entry: entity.ScheduleEntry{ID: uuid.New(), ClassSubjectID: &subject, LessonDate: day.AddDate(0, 0, 1), Position: 1}},
		{Entry: entity.ScheduleEntry{ID: uuid.New(), Title: "Классный час", LessonDate: day, Position: 5}},
	}
	lessons := []entity.Lesson{
		{ClassSubjectID: subject, LessonDate: day.AddDate(0, 0, -3), Position: 1, Homework: text("Старое")},
		{ClassSubjectID: subject, LessonDate: day.AddDate(0, 0, -3), Position: 2, Homework: text("№ 2\nПовторить тему")},
		{ClassSubjectID: other, LessonDate: day.AddDate(0, 0, -1), Position: 1, Homework: text("Чужой предмет")},
		{ClassSubjectID: subject, LessonDate: day, Position: 1, Homework: text("Ко вторнику")},
		{ClassSubjectID: subject, LessonDate: day.AddDate(0, 0, 2), Position: 1, Homework: text("Будущее")},
	}
	fillScheduleHomework(entries, lessons)
	if entries[0].Homework == nil || *entries[0].Homework != "№ 2\nПовторить тему" || *entries[0].HomeworkFrom != "2026-11-06" {
		t.Fatalf("first lesson: %+v", entries[0])
	}
	if entries[1].Homework != nil || entries[1].HomeworkFrom != nil {
		t.Fatal("second daily lesson received homework")
	}
	if entries[2].Homework != nil {
		t.Fatal("extracurricular entry received homework")
	}
	if entries[3].Homework == nil || *entries[3].Homework != "Ко вторнику" {
		t.Fatal("next day did not receive Monday homework")
	}
}
func TestScheduleHomeworkDoesNotReuseOlderAssignmentAfterEmptyLesson(t *testing.T) {
	subject := uuid.New()
	day := time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)
	old := "Старое задание"
	for _, empty := range []*string{nil, new(string)} {
		entries := []ParentScheduleEntry{{Entry: entity.ScheduleEntry{ClassSubjectID: &subject, LessonDate: day, Position: 1}}}
		fillScheduleHomework(entries, []entity.Lesson{{ClassSubjectID: subject, LessonDate: day.AddDate(0, 0, -5), Position: 1, Homework: &old}, {ClassSubjectID: subject, LessonDate: day.AddDate(0, 0, -1), Position: 1, Homework: empty}})
		if entries[0].Homework != nil {
			t.Fatal("old assignment was reused")
		}
	}
}
