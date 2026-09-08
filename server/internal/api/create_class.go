package api

import (
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
)

func (h *ClassHandler) CreateClass(c *fiber.Ctx) error {
	var request gen.CreateClassRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.CreateClass(c.UserContext(), c.Cookies("journal_session"), request.AcademicYearId, request.Name, int16(request.GradeLevel))
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.Status(201).JSON(gen.CreateClassResponse{Class: classResponse(item)})
}
