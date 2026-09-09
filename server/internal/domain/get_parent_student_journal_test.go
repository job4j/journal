package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
)

func parentJournalFixture() (classRepoStub, uuid.UUID, uuid.UUID) {
	parentID, studentID, otherID, teacherID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	yearID, classID, assignmentID, subjectID, lessonID, itemID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	date := time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC)
	return classRepoStub{authenticated: true, objectAllowed: map[string]bool{"can_view_user": true, "can_view_journal": true}, currentUser: entity.User{ID: parentID, Roles: []string{"parent"}}, year: entity.AcademicYear{ID: yearID}, classes: []entity.Class{{ID: classID, AcademicYearID: yearID}}, students: []entity.ClassStudent{{ClassID: classID, UserID: studentID, EnrolledOn: date}}, users: []entity.User{{ID: studentID, Roles: []string{"student"}}, {ID: teacherID, Roles: []string{"teacher"}}}, classSubjects: []entity.ClassSubject{{ID: assignmentID, ClassID: classID, SubjectID: subjectID, ResponsibleTeacherID: teacherID}}, subjects: []entity.Subject{{ID: subjectID, Name: "Математика"}}, lessons: []entity.Lesson{{ID: lessonID, ClassSubjectID: assignmentID, LessonDate: date}}, gradeItems: []entity.GradeItem{{ID: itemID, LessonID: lessonID}}, scores: []entity.Score{{ID: uuid.New(), GradeItemID: itemID, UserID: studentID}, {ID: uuid.New(), GradeItemID: itemID, UserID: otherID}}, absences: []entity.Absence{{ID: uuid.New(), LessonID: lessonID, UserID: studentID}, {ID: uuid.New(), LessonID: lessonID, UserID: otherID}}}, studentID, yearID
}

func TestGetParentStudentJournalFiltersOtherStudents(t *testing.T) {
	repo, studentID, yearID := parentJournalFixture()
	result, err := NewClassDomain(repo).GetParentStudentJournal(context.Background(), testTx{}, "hash", studentID, yearID)
	if err != nil {
		t.Fatal(err)
	}
	lesson := result.Subjects[0].Lessons[0]
	if len(lesson.Scores[repo.gradeItems[0].ID]) != 1 || len(lesson.Absences) != 1 {
		t.Fatalf("lesson = %+v", lesson)
	}
}
func TestGetParentStudentJournalRequiresMatchingPermissions(t *testing.T) {
	repo, studentID, yearID := parentJournalFixture()
	repo.objectAllowed["can_view_journal"] = false
	_, err := NewClassDomain(repo).GetParentStudentJournal(context.Background(), testTx{}, "hash", studentID, yearID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v", err)
	}
}
