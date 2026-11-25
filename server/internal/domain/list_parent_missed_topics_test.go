package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type missedTopicsScopedRepo struct {
	classRepoStub
	child uuid.UUID
}

func (s missedTopicsScopedRepo) CheckSessionPermissionForValue(_ context.Context, _ repository.Transaction, _ string, code, value string) (bool, bool, error) {
	return s.authenticated, value == s.child.String() && s.objectAllowed[code], nil
}
func TestParentMissedTopicsPermissionsMustReferToTheRequestedChild(t *testing.T) {
	repo, req := missedTopicsFixture()
	allowedChild := req.StudentID
	req.StudentID = uuid.New()
	repo.users = append(repo.users, entity.User{ID: req.StudentID, Roles: []string{"student"}})
	_, err := NewClassDomain(missedTopicsScopedRepo{classRepoStub: repo, child: allowedChild}).ListParentMissedTopics(context.Background(), testTx{}, "hash", req)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("another family's child accessible: %v", err)
	}
}

func missedTopicsFixture() (classRepoStub, ListParentMissedTopicsRequest) {
	repo, student, year := parentJournalFixture()
	day := repo.lessons[0].LessonDate
	repo.year.StartsOn = day.AddDate(0, 0, -7)
	repo.year.EndsOn = day.AddDate(1, 0, 0)
	repo.lessons[0].Topic = "Дроби"
	repo.lessons[0].Position = 1
	repo.classes[0].Name = "7А"
	return repo, ListParentMissedTopicsRequest{StudentID: student, AcademicYearID: year}
}

func TestParentMissedTopicsRequireRoleAndBothChildPermissions(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*classRepoStub, *ListParentMissedTopicsRequest)
		want   error
	}{
		{"allowed", func(*classRepoStub, *ListParentMissedTopicsRequest) {}, nil},
		{"expired", func(r *classRepoStub, _ *ListParentMissedTopicsRequest) { r.currentUser.ID = uuid.Nil }, ErrUnauthenticated},
		{"teacher", func(r *classRepoStub, _ *ListParentMissedTopicsRequest) { r.currentUser.Roles = []string{"teacher"} }, ErrForbidden},
		{"missing user permission", func(r *classRepoStub, _ *ListParentMissedTopicsRequest) { r.objectAllowed["can_view_user"] = false }, ErrForbidden},
		{"missing journal permission", func(r *classRepoStub, _ *ListParentMissedTopicsRequest) { r.objectAllowed["can_view_journal"] = false }, ErrForbidden},
		{"unknown child", func(_ *classRepoStub, q *ListParentMissedTopicsRequest) { q.StudentID = uuid.New() }, ErrUserNotFound},
		{"not enrolled", func(r *classRepoStub, _ *ListParentMissedTopicsRequest) { r.students = nil }, ErrClassStudentNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, q := missedTopicsFixture()
			test.change(&r, &q)
			result, err := NewClassDomain(r).ListParentMissedTopics(context.Background(), testTx{}, "hash", q)
			if !errors.Is(err, test.want) {
				t.Fatalf("got %v want %v", err, test.want)
			}
			if err != nil && len(result.Items) > 0 {
				t.Fatal("unauthorized data returned")
			}
		})
	}
}

func TestParentMissedTopicsFilterYearClassMembershipAndOtherChildren(t *testing.T) {
	repo, req := missedTopicsFixture()
	day := repo.lessons[0].LessonDate
	classID := repo.classes[0].ID
	left := day.AddDate(0, 0, 2)
	repo.students[0].LeftOn = &left
	for _, test := range []struct {
		name       string
		dateOffset int
		student    uuid.UUID
		assignment uuid.UUID
		want       bool
	}{
		{"enrollment", 0, req.StudentID, repo.classSubjects[0].ID, true},
		{"leaving", 2, req.StudentID, repo.classSubjects[0].ID, true},
		{"before enrollment", -1, req.StudentID, repo.classSubjects[0].ID, false},
		{"after leaving", 3, req.StudentID, repo.classSubjects[0].ID, false},
		{"other student", 1, uuid.New(), repo.classSubjects[0].ID, false},
		{"unknown class assignment", 1, req.StudentID, uuid.New(), false},
		{"outside year", 400, req.StudentID, repo.classSubjects[0].ID, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := repo
			r.lessons = []entity.Lesson{{ID: uuid.New(), ClassSubjectID: test.assignment, LessonDate: day.AddDate(0, 0, test.dateOffset), Topic: test.name}}
			r.absences = []entity.Absence{{LessonID: r.lessons[0].ID, UserID: test.student}, {LessonID: r.lessons[0].ID, UserID: test.student}}
			result, err := NewClassDomain(r).ListParentMissedTopics(context.Background(), testTx{}, "hash", req)
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			if test.want {
				want = 1
			}
			if len(result.Items) != want {
				t.Fatalf("got %d want %d", len(result.Items), want)
			}
			if want == 1 && result.Items[0].Class.ID != classID {
				t.Fatal("wrong class")
			}
		})
	}
	otherYearClass, otherAssignment := uuid.New(), uuid.New()
	repo.classes = append(repo.classes, entity.Class{ID: otherYearClass, AcademicYearID: uuid.New()})
	repo.students = append(repo.students, entity.ClassStudent{ClassID: otherYearClass, UserID: req.StudentID, EnrolledOn: day})
	repo.classSubjects = append(repo.classSubjects, entity.ClassSubject{ID: otherAssignment, ClassID: otherYearClass, SubjectID: repo.subjects[0].ID})
	repo.lessons = append(repo.lessons, entity.Lesson{ID: uuid.New(), ClassSubjectID: otherAssignment, LessonDate: day})
	repo.absences = append(repo.absences, entity.Absence{LessonID: repo.lessons[1].ID, UserID: req.StudentID})
	result, err := NewClassDomain(repo).ListParentMissedTopics(context.Background(), testTx{}, "hash", req)
	if err != nil || len(result.Items) != 1 {
		t.Fatalf("another year leaked: %+v %v", result, err)
	}
}

func TestParentMissedTopicsIncludeTransferredClassesAndKeepLessonOrder(t *testing.T) {
	repo, req := missedTopicsFixture()
	day := repo.lessons[0].LessonDate
	classID, assignmentID := uuid.New(), uuid.New()
	repo.classes = append(repo.classes, entity.Class{ID: classID, AcademicYearID: req.AcademicYearID, Name: "7Б"})
	repo.students = append(repo.students, entity.ClassStudent{ClassID: classID, UserID: req.StudentID, EnrolledOn: day.AddDate(0, 0, 1)})
	repo.classSubjects = append(repo.classSubjects, entity.ClassSubject{ID: assignmentID, ClassID: classID, SubjectID: repo.subjects[0].ID})
	for position := int16(1); position <= 2; position++ {
		id := uuid.New()
		repo.lessons = append(repo.lessons, entity.Lesson{ID: id, ClassSubjectID: assignmentID, LessonDate: day.AddDate(0, 0, 1), Position: position, Topic: "Новая тема"})
		repo.absences = append(repo.absences, entity.Absence{LessonID: id, UserID: req.StudentID})
	}
	repo.lessons = append(repo.lessons, entity.Lesson{ID: uuid.New(), ClassSubjectID: assignmentID, LessonDate: day.AddDate(0, 0, 2), Topic: "Нет оценки, но нет и пропуска"})
	result, err := NewClassDomain(repo).ListParentMissedTopics(context.Background(), testTx{}, "hash", req)
	if err != nil || len(result.Items) != 3 || result.Items[0].Lesson.Position != 2 || result.Items[1].Class.Name != "7Б" || result.Items[2].Class.Name != "7А" {
		t.Fatalf("wrong history/order: %+v %v", result, err)
	}
	repo.absences = nil
	result, err = NewClassDomain(repo).ListParentMissedTopics(context.Background(), testTx{}, "hash", req)
	if err != nil || result.Items == nil || len(result.Items) != 0 {
		t.Fatalf("empty result: %+v %v", result, err)
	}
}
