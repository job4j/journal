package api

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestReplaceScheduleParsesEntry(t *testing.T) {
	stub := &scheduleServiceStub{}
	app := fiber.New()
	NewScheduleHandler(stub, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	body := `{"revision":2,"entries":[{"id":"","day":"2026-11-10","position":3,"classSubjectId":null,"title":"Шахматы","notes":"Группа","startsAt":"13:00","endsAt":null}]}`
	req := httptest.NewRequest("PUT", "/classes/60000000-0000-4000-8000-000000000001/schedule/2026-11-09", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil || res.StatusCode != 200 {
		t.Fatalf("status: %v %+v", err, res)
	}
	if stub.request.Revision != 2 || len(stub.request.Entries) != 1 || stub.request.Entries[0].Title != "Шахматы" || stub.request.Entries[0].LessonDate.Format(time.DateOnly) != "2026-11-10" {
		t.Fatalf("request: %+v", stub.request)
	}
}
