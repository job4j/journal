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

type academicYearServiceStub struct {
	items []domain.AcademicYearView
	err   error
}

func (s academicYearServiceStub) ListAcademicYears(context.Context, string) ([]domain.AcademicYearView, error) {
	return s.items, s.err
}
func (s academicYearServiceStub) CreateAcademicYear(context.Context, string, domain.CreateAcademicYearRequest) (domain.AcademicYearView, error) {
	return domain.AcademicYearView{}, s.err
}

func TestListAcademicYearsReturnsQuarters(t *testing.T) {
	id, quarterID := uuid.New(), uuid.New()
	date := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	service := academicYearServiceStub{items: []domain.AcademicYearView{{Year: entity.AcademicYear{ID: id, Name: "2026/2027", StartsOn: date, EndsOn: date.AddDate(1, 0, 0), Status: "active"}, Quarters: []entity.AcademicYearQuarter{{ID: quarterID, AcademicYearID: id, Number: 1, StartsOn: date, EndsOn: date.AddDate(0, 1, 0)}}}}}
	app := fiber.New()
	NewAcademicYearHandler(service, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	request := httptest.NewRequest("GET", "/academic-years", nil)
	request.Header.Set("Cookie", "journal_session=token")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status = %d", response.StatusCode)
	}
}

func TestListAcademicYearsRejectsMissingSession(t *testing.T) {
	app := fiber.New()
	NewAcademicYearHandler(academicYearServiceStub{err: domain.ErrUnauthenticated}, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	response, err := app.Test(httptest.NewRequest("GET", "/academic-years", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 401 {
		t.Fatalf("status = %d", response.StatusCode)
	}
}

func TestCreateAcademicYearMapsValidationError(t *testing.T) {
	app := fiber.New()
	NewAcademicYearHandler(academicYearServiceStub{err: domain.ErrInvalidAcademicYear}, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	body := `{"name":"2026/2027","startsOn":"2026-09-01","endsOn":"2027-05-31","status":"planned","quarters":[]}`
	request := httptest.NewRequest("POST", "/academic-years", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 400 {
		t.Fatalf("status = %d", response.StatusCode)
	}
}
