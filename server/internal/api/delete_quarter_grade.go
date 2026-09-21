package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *ClassHandler) DeleteQuarterGrade(c *fiber.Ctx) error {
	assignmentID, quarterID, err := quarterGradeIDs(c)
	if err != nil {
		return err
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	if err = h.service.DeleteQuarterGrade(c.UserContext(), c.Cookies("journal_session"), assignmentID, quarterID, studentID); err != nil {
		return h.writeClassError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
