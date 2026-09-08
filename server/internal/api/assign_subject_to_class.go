package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) AssignSubjectToClass(c *fiber.Ctx) error {
	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	var request gen.AssignSubjectToClassRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.AssignSubjectToClass(c.UserContext(), c.Cookies("journal_session"), classID, request.SubjectId, request.ResponsibleTeacherId)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.Status(201).JSON(gen.AssignSubjectToClassResponse{ClassSubject: classSubjectResponse(item)})
}
