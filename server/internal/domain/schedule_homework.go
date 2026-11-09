package domain

import (
	"journal/server/internal/repository/entity"
	"sort"
	"strings"
	"time"
)

func fillScheduleHomework(entries []ParentScheduleEntry, lessons []entity.Lesson) {
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i].Entry, entries[j].Entry
		if !a.LessonDate.Equal(b.LessonDate) {
			return a.LessonDate.Before(b.LessonDate)
		}
		if a.Position != b.Position {
			return a.Position < b.Position
		}
		startA, startB := "", ""
		if a.StartsAt != nil {
			startA = *a.StartsAt
		}
		if b.StartsAt != nil {
			startB = *b.StartsAt
		}
		if startA != startB {
			return startA < startB
		}
		return a.ID.String() < b.ID.String()
	})
	seen := map[string]bool{}
	for i := range entries {
		entry := &entries[i]
		if entry.Entry.ClassSubjectID == nil {
			continue
		}
		key := entry.Entry.ClassSubjectID.String() + entry.Entry.LessonDate.Format(time.DateOnly)
		if seen[key] {
			continue
		}
		seen[key] = true
		var previous *entity.Lesson
		for j := range lessons {
			candidate := &lessons[j]
			if candidate.ClassSubjectID != *entry.Entry.ClassSubjectID || !dayUTC(candidate.LessonDate).Before(dayUTC(entry.Entry.LessonDate)) {
				continue
			}
			if previous == nil || candidate.LessonDate.After(previous.LessonDate) || (candidate.LessonDate.Equal(previous.LessonDate) && candidate.Position > previous.Position) {
				previous = candidate
			}
		}
		if previous != nil && previous.Homework != nil && strings.TrimSpace(*previous.Homework) != "" {
			entry.Homework = previous.Homework
			from := previous.LessonDate.Format(time.DateOnly)
			entry.HomeworkFrom = &from
		}
	}
}
