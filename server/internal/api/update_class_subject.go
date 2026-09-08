package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) UpdateClassSubject(c *fiber.Ctx) error {
	assignmentID, err := uuid.Parse(c.Params("classSubjectId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_subject_id", "Некорректный идентификатор назначения")
	}
	var request gen.UpdateClassSubjectRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.UpdateClassSubject(c.UserContext(), c.Cookies("journal_session"), assignmentID, request.ResponsibleTeacherId)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.UpdateClassSubjectResponse{ClassSubject: classSubjectResponse(item)})
}
