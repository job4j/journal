package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
	"time"
)

func (h *ScheduleHandler) ReplaceClassSchedule(c *fiber.Ctx) error {
	params, err := scheduleParams(c)
	if err != nil {
		return writeError(c, 400, "invalid_schedule", "Некорректный класс или дата недели")
	}
	var body gen.ReplaceClassScheduleRequest
	if err = c.BodyParser(&body); err != nil || body.Entries == nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	req := domain.ReplaceClassScheduleRequest{GetClassScheduleRequest: params, Revision: body.Revision, Entries: []entity.ScheduleEntry{}}
	for _, item := range body.Entries {
		e := entity.ScheduleEntry{Position: item.Position, Title: item.Title, Notes: item.Notes, StartsAt: item.StartsAt, EndsAt: item.EndsAt}
		e.LessonDate, err = time.Parse(time.DateOnly, item.Day)
		if err != nil {
			return writeError(c, 400, "invalid_schedule", "Некорректная дата занятия")
		}
		if item.Id != "" {
			e.ID, err = uuid.Parse(item.Id)
			if err != nil {
				return writeError(c, 400, "invalid_schedule", "Некорректный ID занятия")
			}
		}
		if item.ClassSubjectId != nil {
			id, parseErr := uuid.Parse(*item.ClassSubjectId)
			if parseErr != nil {
				return writeError(c, 400, "invalid_schedule", "Некорректный ID предмета")
			}
			e.ClassSubjectID = &id
		}
		req.Entries = append(req.Entries, e)
	}
	result, err := h.service.ReplaceClassSchedule(c.UserContext(), c.Cookies("journal_session"), req)
	if err != nil {
		return h.writeScheduleError(c, err)
	}
	return c.JSON(gen.ReplaceClassScheduleResponse{Week: scheduleResponse(result.GetClassScheduleResponse)})
}
