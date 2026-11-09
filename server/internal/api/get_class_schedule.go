package api

import (
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
)

func (h *ScheduleHandler) GetClassSchedule(c *fiber.Ctx) error {
	req, err := scheduleParams(c)
	if err != nil {
		return writeError(c, 400, "invalid_schedule", "Некорректный класс или дата недели")
	}
	result, err := h.service.GetClassSchedule(c.UserContext(), c.Cookies("journal_session"), req)
	if err != nil {
		return h.writeScheduleError(c, err)
	}
	return c.JSON(gen.GetClassScheduleResponse{Week: scheduleResponse(result)})
}
