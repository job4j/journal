package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *AcademicYearHandler) DeleteAcademicYearQuarter(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("quarterId"))
	if err != nil {
		return writeError(c, 400, "invalid_quarter_id", "Некорректный идентификатор периода")
	}
	if err = h.service.DeleteAcademicYearQuarter(c.UserContext(), c.Cookies("journal_session"), id); err != nil {
		return h.writeQuarterError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
