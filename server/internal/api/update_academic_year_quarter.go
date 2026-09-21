package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *AcademicYearHandler) UpdateAcademicYearQuarter(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("quarterId"))
	if err != nil {
		return writeError(c, 400, "invalid_quarter_id", "Некорректный идентификатор периода")
	}
	var body gen.UpdateAcademicYearQuarterRequest
	if err = c.BodyParser(&body); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.UpdateAcademicYearQuarter(c.UserContext(), c.Cookies("journal_session"), id, body.Name, body.StartsOn.Time, body.EndsOn.Time)
	if err != nil {
		return h.writeQuarterError(c, err)
	}
	return c.JSON(gen.UpdateAcademicYearQuarterResponse{Quarter: academicYearQuarterResponse(item)})
}
