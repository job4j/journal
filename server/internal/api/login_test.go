package api

import (
	"context"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type authServiceStub struct {
	result domain.LoginResponse
	err    error
}

func (s authServiceStub) Login(context.Context, string, string) (domain.LoginResponse, error) {
	return s.result, s.err
}

func TestLoginSetsSessionCookie(t *testing.T) {
	expires := time.Now().Add(time.Hour)
	service := authServiceStub{result: domain.LoginResponse{
		Token: "secret-token", ExpiresAt: expires,
		User: entity.User{ID: uuid.New(), Email: "user@example.com", FirstName: "Иван", LastName: "Иванов", Status: "active", Roles: []string{"teacher"}},
	}}
	app := fiber.New()
	NewAuthHandler(service, slog.New(slog.NewTextHandler(io.Discard, nil)), true).Register(app)
	request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"secret"}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d", response.StatusCode)
	}
	cookie := response.Header.Get("Set-Cookie")
	for _, part := range []string{"journal_session=secret-token", "HttpOnly", "secure", "SameSite=Lax"} {
		if !strings.Contains(cookie, part) {
			t.Fatalf("cookie %q does not contain %q", cookie, part)
		}
	}
}

func TestLoginMapsInvalidCredentialsToUnauthorized(t *testing.T) {
	app := fiber.New()
	NewAuthHandler(authServiceStub{err: domain.ErrInvalidCredentials}, slog.New(slog.NewTextHandler(io.Discard, nil)), false).Register(app)
	request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"wrong"}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d", response.StatusCode)
	}
}
