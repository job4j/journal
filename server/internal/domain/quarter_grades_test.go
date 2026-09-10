package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func quarterGradeFixture() (classRepoStub, uuid.UUID, uuid.UUID, uuid.UUID) {
	teacherID, studentID, classID, yearID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	assignmentID, quarterID := uuid.New(), uuid.New()
	starts := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	ends := time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)
	return classRepoStub{
		authenticated: true,
		allowed:       true,
		currentUser: entity.User{
			ID: teacherID, Roles: []string{"teacher"},
		},
		classes: []entity.Class{{ID: classID, AcademicYearID: yearID}},
		classSubjects: []entity.ClassSubject{{
			ID: assignmentID, ClassID: classID, ResponsibleTeacherID: teacherID,
		}},
		quarters: []entity.AcademicYearQuarter{{
			ID: quarterID, AcademicYearID: yearID, StartsOn: starts, EndsOn: ends,
		}},
		students: []entity.ClassStudent{{
			ClassID: classID, UserID: studentID, EnrolledOn: starts,
		}},
	}, assignmentID, quarterID, studentID
}
func TestListQuarterGradesAllowsAdministrator(t *testing.T) {
	repo, assignmentID, quarterID, _ := quarterGradeFixture()
	repo.currentUser.Roles = []string{"admin"}

	if _, err := NewClassDomain(repo).ListQuarterGrades(
		context.Background(),
		testTx{},
		"hash",
		assignmentID,
		quarterID,
	); err != nil {
		t.Fatal(err)
	}
}

func TestPutQuarterGradeRejectsAdministrator(t *testing.T) {
	repo, assignmentID, quarterID, studentID := quarterGradeFixture()
	repo.currentUser.Roles = []string{"admin"}
	value := 5.0

	_, err := NewClassDomain(repo).PutQuarterGrade(
		context.Background(),
		testTx{},
		"hash",
		assignmentID,
		quarterID,
		studentID,
		QuarterGradeInput{GradingScale: "five_point", NumericValue: &value},
	)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
func TestPutQuarterGradeValidatesAndStoresAudit(t *testing.T) {
	repo, assignmentID, quarterID, studentID := quarterGradeFixture()
	stored := entity.QuarterGrade{}
	repo.changedQuarterGrade = &stored
	value := 9.5
	max := 10.0
	comment := " Хорошо "
	_, err := NewClassDomain(repo).PutQuarterGrade(
		context.Background(), testTx{}, "hash", assignmentID, quarterID, studentID,
		QuarterGradeInput{
			GradingScale: "points", MaxScore: &max, NumericValue: &value,
			TeacherComment: &comment,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	invalidComment := stored.TeacherComment == nil || *stored.TeacherComment != "Хорошо"
	if stored.CreatedBy != repo.currentUser.ID || invalidComment {
		t.Fatalf("stored = %+v", stored)
	}
}
func TestPutQuarterGradeRejectsWrongYear(t *testing.T) {
	repo, assignmentID, quarterID, studentID := quarterGradeFixture()
	repo.quarters[0].AcademicYearID = uuid.New()
	value := 5.0
	_, err := NewClassDomain(repo).PutQuarterGrade(
		context.Background(), testTx{}, "hash", assignmentID, quarterID, studentID,
		QuarterGradeInput{GradingScale: "five_point", NumericValue: &value},
	)
	if !errors.Is(err, ErrInvalidScore) {
		t.Fatalf("error = %v", err)
	}
}
func TestPutQuarterGradeRejectsStudentOutsideQuarter(t *testing.T) {
	repo, assignmentID, quarterID, studentID := quarterGradeFixture()
	left := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	repo.students[0].LeftOn = &left
	value := 5.0
	_, err := NewClassDomain(repo).PutQuarterGrade(
		context.Background(), testTx{}, "hash", assignmentID, quarterID, studentID,
		QuarterGradeInput{GradingScale: "five_point", NumericValue: &value},
	)
	if !errors.Is(err, ErrClassStudentNotFound) {
		t.Fatalf("error = %v", err)
	}
}
