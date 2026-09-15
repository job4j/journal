package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *ClassHandler) DeleteStudentScore(c *fiber.Ctx) error {
	gradeItemID, err := uuid.Parse(c.Params("gradeItemId"))
	if err != nil {
		return writeError(c, 400, "invalid_grade_item_id", "Некорректный идентификатор работы")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	if err = h.service.DeleteStudentScore(c.UserContext(), c.Cookies("journal_session"), gradeItemID, studentID); err != nil {
		return h.writeClassError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
