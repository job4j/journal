package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func scoreFixture(scale string) (classRepoStub, uuid.UUID, uuid.UUID) {
	teacherID, studentID, classID := uuid.New(), uuid.New(), uuid.New()
	assignmentID, lessonID, itemID := uuid.New(), uuid.New(), uuid.New()
	max := 10.0
	return classRepoStub{
		authenticated: true, allowed: true, currentUser: entity.User{ID: teacherID, Roles: []string{"teacher"}},
		classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, ResponsibleTeacherID: teacherID}},
		lessons:       []entity.Lesson{{ID: lessonID, ClassSubjectID: assignmentID, LessonDate: time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC)}},
		gradeItems:    []entity.GradeItem{{ID: itemID, LessonID: lessonID, GradingScale: scale, MaxScore: &max}},
		students:      []entity.ClassStudent{{ClassID: classID, UserID: studentID, EnrolledOn: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)}},
	}, itemID, studentID
}

func TestPutStudentScoreAcceptsScalesAndRecordsTeacher(t *testing.T) {
	cases := []struct {
		name, scale string
		numeric     *float64
		text        *string
	}{
		{name: "five point", scale: "five_point", numeric: floatPointer(5)},
		{name: "points", scale: "points", numeric: floatPointer(9.5)},
		{name: "pass fail", scale: "pass_fail", text: stringPointer("pass")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stored := entity.Score{}
			repo, itemID, studentID := scoreFixture(tc.scale)
			repo.upsertedScore = &stored
			_, err := NewClassDomain(repo).PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, tc.numeric, tc.text, stringPointer("  good  "), false)
			if err != nil {
				t.Fatal(err)
			}
			if stored.UserID != studentID || stored.CreatedBy != repo.currentUser.ID || stored.TeacherComment == nil || *stored.TeacherComment != "good" {
				t.Fatalf("stored = %+v", stored)
			}
		})
	}
}

func TestPutStudentScoreRejectsInvalidValues(t *testing.T) {
	cases := []struct {
		scale   string
		numeric *float64
		text    *string
	}{
		{scale: "five_point", numeric: floatPointer(1)},
		{scale: "five_point", numeric: floatPointer(4.5)},
		{scale: "points", numeric: floatPointer(10.1)},
		{scale: "pass_fail", text: stringPointer("passed")},
	}
	for _, tc := range cases {
		repo, itemID, studentID := scoreFixture(tc.scale)
		_, err := NewClassDomain(repo).PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, tc.numeric, tc.text, nil, false)
		if !errors.Is(err, ErrInvalidScore) {
			t.Fatalf("scale %s: %v", tc.scale, err)
		}
	}
}

func TestPutStudentScoreRejectsFormerStudentAndOtherTeacher(t *testing.T) {
	repo, itemID, studentID := scoreFixture("five_point")
	left := time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC)
	repo.students[0].LeftOn = &left
	_, err := NewClassDomain(repo).PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, floatPointer(5), nil, nil, false)
	if !errors.Is(err, ErrClassStudentNotFound) {
		t.Fatalf("former student: %v", err)
	}
	repo, itemID, studentID = scoreFixture("five_point")
	repo.classSubjects[0].ResponsibleTeacherID = uuid.New()
	_, err = NewClassDomain(repo).PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, floatPointer(5), nil, nil, false)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("other teacher: %v", err)
	}
}

func floatPointer(value float64) *float64 { return &value }
func stringPointer(value string) *string  { return &value }

func TestPutStudentScoreRemediationRequiresExplanation(t *testing.T) {
	for _, comment := range []*string{nil, stringPointer(""), stringPointer(" \t\n "), stringPointer("\u00a0")} {
		repo, itemID, studentID := scoreFixture("five_point")
		stored := entity.Score{}
		repo.upsertedScore = &stored
		_, err := NewClassDomain(repo).PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, floatPointer(2), nil, comment, true)
		if !errors.Is(err, ErrRemediationCommentRequired) {
			t.Fatalf("comment %v: %v", comment, err)
		}
		if stored.UserID != uuid.Nil {
			t.Fatal("invalid remediation was persisted")
		}
	}
}

func TestPutStudentScoreStoresAndClearsRemediation(t *testing.T) {
	repo, itemID, studentID := scoreFixture("five_point")
	stored := entity.Score{}
	repo.upsertedScore = &stored
	d := NewClassDomain(repo)
	_, err := d.PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, floatPointer(2), nil, stringPointer("  Повторили дроби и решили задачи  "), true)
	if err != nil {
		t.Fatal(err)
	}
	if !stored.IsRemediated || stored.TeacherComment == nil || *stored.TeacherComment != "Повторили дроби и решили задачи" || *stored.NumericValue != 2 {
		t.Fatalf("stored = %+v", stored)
	}
	_, err = d.PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, floatPointer(2), nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if stored.IsRemediated || stored.TeacherComment != nil {
		t.Fatalf("cleared = %+v", stored)
	}
}

func TestPutStudentScoreRemediationCannotBypassTeacherAccess(t *testing.T) {
	repo, itemID, studentID := scoreFixture("five_point")
	repo.classSubjects[0].ResponsibleTeacherID = uuid.New()
	_, err := NewClassDomain(repo).PutStudentScore(context.Background(), testTx{}, "hash", itemID, studentID, floatPointer(2), nil, stringPointer("Отработано"), true)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("got %v", err)
	}
}
