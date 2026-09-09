package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestRequestLogAddsRequestIDAndStructuredFields(t *testing.T) {
	var output bytes.Buffer
	app := fiber.New()
	app.Use(RequestLog(slog.New(slog.NewJSONHandler(&output, nil))))
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
	response, err := app.Test(httptest.NewRequest("GET", "/health?token=secret", nil))
	if err != nil {
		t.Fatal(err)
	}
	requestID := response.Header.Get(requestIDHeader)
	if _, err := uuid.Parse(requestID); err != nil {
		t.Fatalf("invalid request ID %q", requestID)
	}
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record["request_id"] != requestID || record["method"] != "GET" || record["path"] != "/health" || record["status"] != float64(204) {
		t.Fatalf("unexpected log record: %v", record)
	}
	if strings.Contains(output.String(), "secret") {
		t.Fatal("query secret leaked to log")
	}
}

func TestRequestLogPreservesValidRequestIDWithoutLoggingHeaders(t *testing.T) {
	var output bytes.Buffer
	app := fiber.New()
	app.Use(RequestLog(slog.New(slog.NewJSONHandler(&output, nil))))
	app.Post("/login", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	request := httptest.NewRequest("POST", "/login", strings.NewReader(`{"password":"hidden"}`))
	request.Header.Set(requestIDHeader, "38b042e6-89c0-4aa5-a1a4-6211dbfa0fa6")
	request.Header.Set("Authorization", "Bearer hidden")
	if _, err := app.Test(request); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "hidden") {
		t.Fatal("sensitive request data leaked to log")
	}
}
