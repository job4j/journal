package api

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func RequestLog(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(requestIDHeader)
		if _, err := uuid.Parse(requestID); err != nil {
			requestID = uuid.NewString()
		}
		c.Set(requestIDHeader, requestID)
		started := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()
		if err != nil && status < 400 {
			status = fiber.StatusInternalServerError
		}
		logger.Info("http request",
			"request_id", requestID, "method", c.Method(), "path", c.Path(),
			"status", status, "duration_ms", time.Since(started).Milliseconds(),
		)
		return err
	}
}
