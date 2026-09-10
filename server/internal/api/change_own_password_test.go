package api

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"journal/server/internal/domain"
)

func TestChangeOwnPassword(t *testing.T) {
	app := fiber.New()
	handler := NewAuthHandler(
		authServiceStub{},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		false,
	)
	app.Put("/me/password", handler.ChangeOwnPassword)
	request := httptest.NewRequest(
		"PUT",
		"/me/password",
		strings.NewReader(`{"currentPassword":"old","newPassword":"x"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d", response.StatusCode)
	}
}

func TestChangeOwnPasswordRejectsWrongCurrentPassword(t *testing.T) {
	app := fiber.New()
	handler := NewAuthHandler(
		authServiceStub{err: domain.ErrInvalidCredentials},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		false,
	)
	app.Put("/me/password", handler.ChangeOwnPassword)
	request := httptest.NewRequest(
		"PUT",
		"/me/password",
		strings.NewReader(`{"currentPassword":"wrong","newPassword":"x"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d", response.StatusCode)
	}
}
