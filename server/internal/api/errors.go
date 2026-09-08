package api

import (
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
)

func writeError(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(gen.ErrorResponse{Code: code, Message: message})
}
