package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type missedTopicsAPIStub struct {
	ClassService
	request domain.ListParentMissedTopicsRequest
	err     error
	items   []domain.ParentMissedTopic
	calls   int
}

func (s *missedTopicsAPIStub) ListParentMissedTopics(_ context.Context, _ string, req domain.ListParentMissedTopicsRequest) (domain.ListParentMissedTopicsResponse, error) {
	s.request = req
	s.calls++
	return domain.ListParentMissedTopicsResponse{Items: s.items}, s.err
}
func TestParentMissedTopicsHTTPContract(t *testing.T) {
	student, year := uuid.New(), uuid.New()
	for _, test := range []struct {
		err    error
		status int
	}{{nil, 200}, {domain.ErrUnauthenticated, 401}, {domain.ErrForbidden, 403}, {domain.ErrClassStudentNotFound, 404}} {
		stub := &missedTopicsAPIStub{err: test.err}
		app := fiber.New()
		NewClassHandler(stub, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
		res, err := app.Test(httptest.NewRequest("GET", "/parent/students/"+student.String()+"/missed-topics?academicYearId="+year.String(), nil))
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if res.StatusCode != test.status || stub.request.StudentID != student || stub.request.AcademicYearID != year {
			t.Fatalf("status=%d body=%s", res.StatusCode, body)
		}
		if test.status == 200 && !strings.Contains(string(body), `"items":[]`) {
			t.Fatalf("empty result=%s", body)
		}
	}
}
func TestParentMissedTopicsReturnsOnlyLessonSummaryAndValidatesIDs(t *testing.T) {
	stub := &missedTopicsAPIStub{items: []domain.ParentMissedTopic{{Lesson: entity.Lesson{ID: uuid.New(), Topic: "Дроби", LessonDate: time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC), Position: 1}, Subject: entity.Subject{ID: uuid.New(), Name: "Математика"}, Class: entity.Class{ID: uuid.New(), Name: "7А"}}}}
	app := fiber.New()
	NewClassHandler(stub, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	url := "/parent/students/" + uuid.New().String() + "/missed-topics?academicYearId=" + uuid.New().String()
	res, err := app.Test(httptest.NewRequest("GET", url, nil))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != 200 || !strings.Contains(string(body), `"topic":"Дроби"`) || strings.Contains(string(body), "scores") || strings.Contains(string(body), "createdBy") {
		t.Fatalf("body=%s", body)
	}
	for _, url := range []string{"/parent/students/not-a-uuid/missed-topics?academicYearId=bad", "/parent/students/" + uuid.New().String() + "/missed-topics"} {
		res, err = app.Test(httptest.NewRequest("GET", url, nil))
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("invalid IDs status=%d", res.StatusCode)
		}
	}
	if stub.calls != 1 {
		t.Fatal("invalid IDs reached service")
	}
}
