package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) AddStudentToClass(c *fiber.Ctx) error {
	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	var request gen.AddStudentToClassRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.AddStudentToClass(c.UserContext(), c.Cookies("journal_session"), classID, request.StudentId, request.EnrolledOn.Time)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.Status(201).JSON(gen.AddStudentToClassResponse{ClassStudent: classStudentResponse(item)})
}
