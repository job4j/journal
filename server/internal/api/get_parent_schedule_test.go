package api

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"journal/server/internal/domain"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParentScheduleHTTPContract(t *testing.T) {
	for _, test := range []struct {
		err    error
		status int
	}{{nil, 200}, {domain.ErrUnauthenticated, 401}, {domain.ErrForbidden, 403}, {domain.ErrInvalidSchedule, 400}} {
		app := fiber.New()
		NewScheduleHandler(&scheduleServiceStub{err: test.err}, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
		res, err := app.Test(httptest.NewRequest("GET", "/parent/students/10000000-0000-4000-8000-000000000001/classes/60000000-0000-4000-8000-000000000001/schedule/2026-11-09", nil))
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if res.StatusCode != test.status {
			t.Fatalf("status=%d", res.StatusCode)
		}
		if test.status == 200 && (!strings.Contains(string(body), `"entries":[]`) || strings.Contains(string(body), "revision")) {
			t.Fatalf("response=%s", body)
		}
	}
}
