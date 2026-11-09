package api

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"io"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

type scheduleServiceStub struct {
	err     error
	request domain.ReplaceClassScheduleRequest
}

func (s *scheduleServiceStub) GetClassSchedule(_ context.Context, _ string, r domain.GetClassScheduleRequest) (domain.GetClassScheduleResponse, error) {
	return domain.GetClassScheduleResponse{Week: entity.ScheduleWeek{WeekStart: r.WeekStart}, Entries: []entity.ScheduleEntry{}}, s.err
}
func (s *scheduleServiceStub) ReplaceClassSchedule(_ context.Context, _ string, r domain.ReplaceClassScheduleRequest) (domain.ReplaceClassScheduleResponse, error) {
	s.request = r
	return domain.ReplaceClassScheduleResponse{GetClassScheduleResponse: domain.GetClassScheduleResponse{Week: entity.ScheduleWeek{WeekStart: r.WeekStart, Revision: r.Revision + 1}, Entries: r.Entries}}, s.err
}
func TestScheduleHTTPResponses(t *testing.T) {
	for _, test := range []struct {
		err    error
		status int
	}{{nil, 200}, {domain.ErrUnauthenticated, 401}, {domain.ErrForbidden, 403}, {domain.ErrInvalidSchedule, 400}, {domain.ErrClassNotFound, 404}, {domain.ErrScheduleConflict, 409}, {errors.New("secret database details"), 500}} {
		app := fiber.New()
		stub := &scheduleServiceStub{err: test.err}
		NewScheduleHandler(stub, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
		res, err := app.Test(httptest.NewRequest("GET", "/classes/60000000-0000-4000-8000-000000000001/schedule/2026-11-09", nil))
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if res.StatusCode != test.status || strings.Contains(string(body), "secret database") {
			t.Fatalf("status=%d body=%s", res.StatusCode, body)
		}
	}
}
