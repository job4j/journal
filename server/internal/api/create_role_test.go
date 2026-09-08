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
)

type roleServiceStub struct {
	result domain.CreateRoleResponse
	err    error
	token  string
}

func (s *roleServiceStub) CreateRole(_ context.Context, token, code, name string, permissions []string) (domain.CreateRoleResponse, error) {
	s.token = token
	return s.result, s.err
}

func TestCreateRoleReturnsCreatedRole(t *testing.T) {
	service := &roleServiceStub{result: domain.CreateRoleResponse{Role: entity.Role{ID: uuid.New(), Code: "editor", Name: "Редактор"}, PermissionCodes: []string{"can_view_user"}}}
	app := fiber.New()
	NewRoleHandler(service, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	request := httptest.NewRequest("POST", "/roles", strings.NewReader(`{"code":"editor","name":"Редактор","permissions":["can_view_user"]}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "journal_session=token")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d", response.StatusCode)
	}
	if service.token != "token" {
		t.Fatalf("token = %q", service.token)
	}
}
func TestCreateRoleRequiresSessionCookie(t *testing.T) {
	app := fiber.New()
	NewRoleHandler(&roleServiceStub{}, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	request := httptest.NewRequest("POST", "/roles", strings.NewReader(`{"code":"editor","name":"Редактор","permissions":[]}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d", response.StatusCode)
	}
}
func TestCreateRoleMapsForbidden(t *testing.T) {
	app := fiber.New()
	NewRoleHandler(&roleServiceStub{err: domain.ErrForbidden}, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
	request := httptest.NewRequest("POST", "/roles", strings.NewReader(`{"code":"editor","name":"Редактор","permissions":[]}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "journal_session=token")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusForbidden {
		t.Fatalf("status = %d", response.StatusCode)
	}
}
