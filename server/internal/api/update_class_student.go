package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) UpdateClassStudent(c *fiber.Ctx) error {
	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	var request gen.UpdateClassStudentRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	var value = request.LeftOn
	var leftOn = (*time.Time)(nil)
	if value != nil {
		date := value.Time
		leftOn = &date
	}
	item, err := h.service.UpdateClassStudent(c.UserContext(), c.Cookies("journal_session"), classID, studentID, leftOn)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.UpdateClassStudentResponse{ClassStudent: classStudentResponse(item)})
}
