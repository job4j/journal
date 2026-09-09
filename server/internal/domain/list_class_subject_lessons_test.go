package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func TestListClassSubjectLessonsFiltersSortsAndNests(t *testing.T) {
	teacherID, assignmentID, lessonA, lessonB := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	day := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	from := day.AddDate(0, 0, -1)
	repo := classRepoStub{authenticated: true, allowed: true, currentUser: entity.User{ID: teacherID, Roles: []string{"teacher"}}, classSubjects: []entity.ClassSubject{{ID: assignmentID, ResponsibleTeacherID: teacherID}}, lessons: []entity.Lesson{{ID: lessonA, ClassSubjectID: assignmentID, LessonDate: day, Position: 2}, {ID: lessonB, ClassSubjectID: assignmentID, LessonDate: day, Position: 1}, {ID: uuid.New(), ClassSubjectID: assignmentID, LessonDate: from.AddDate(0, 0, -1)}}, materials: []entity.LessonMaterial{{ID: uuid.New(), LessonID: lessonB, Position: 1}}, gradeItems: []entity.GradeItem{{ID: uuid.New(), LessonID: lessonB}}}
	items, err := NewClassDomain(repo).ListClassSubjectLessons(context.Background(), testTx{}, "hash", assignmentID, &from, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].Lesson.ID != lessonB || len(items[0].Materials) != 1 || len(items[0].GradeItems) != 1 {
		t.Fatalf("items=%+v", items)
	}
}
