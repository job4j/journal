package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func TestCreateLessonCreatesMaterials(t *testing.T) {
	teacherID, assignmentID, classID, yearID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	day := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	created := []entity.LessonMaterial{}
	repo := classRepoStub{authenticated: true, allowed: true, currentUser: entity.User{ID: teacherID, Roles: []string{"teacher"}}, classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, ResponsibleTeacherID: teacherID}}, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}}, year: entity.AcademicYear{ID: yearID, StartsOn: day.AddDate(0, 0, -1), EndsOn: day.AddDate(1, 0, 0)}, createdMaterials: &created}
	result, err := NewClassDomain(repo).CreateLesson(context.Background(), testTx{}, "hash", assignmentID, CreateLessonInput{LessonDate: day, Position: 1, Topic: " Тема ", Materials: []LessonMaterialInput{{Title: "Ссылка", URL: "https://example.test", Position: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Lesson.Topic != "Тема" || len(created) != 1 || created[0].LessonID != result.Lesson.ID {
		t.Fatalf("result=%+v materials=%+v", result, created)
	}
}
func TestCreateLessonRejectsDateAndURL(t *testing.T) {
	teacherID, assignmentID, classID, yearID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	repo := classRepoStub{authenticated: true, allowed: true, currentUser: entity.User{ID: teacherID, Roles: []string{"teacher"}}, classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, ResponsibleTeacherID: teacherID}}, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}}, year: entity.AcademicYear{StartsOn: start, EndsOn: start.AddDate(1, 0, 0)}}
	_, err := NewClassDomain(repo).CreateLesson(context.Background(), testTx{}, "hash", assignmentID, CreateLessonInput{LessonDate: start.AddDate(0, 0, -1), Position: 1, Topic: "Тема"})
	if !errors.Is(err, ErrInvalidLesson) {
		t.Fatalf("date error=%v", err)
	}
	_, err = NewClassDomain(repo).CreateLesson(context.Background(), testTx{}, "hash", assignmentID, CreateLessonInput{LessonDate: start, Position: 1, Topic: "Тема", Materials: []LessonMaterialInput{{Title: "bad", URL: "javascript:alert(1)", Position: 1}}})
	if !errors.Is(err, ErrInvalidLesson) {
		t.Fatalf("url error=%v", err)
	}
}
