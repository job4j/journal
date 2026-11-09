package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"time"
)

func scheduleParams(c *fiber.Ctx) (domain.GetClassScheduleRequest, error) {
	id, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return domain.GetClassScheduleRequest{}, err
	}
	day, err := time.Parse(time.DateOnly, c.Params("weekStart"))
	return domain.GetClassScheduleRequest{ClassID: id, WeekStart: day}, err
}
func scheduleResponse(result domain.GetClassScheduleResponse) gen.ScheduleWeek {
	week := gen.ScheduleWeek{WeekStart: result.Week.WeekStart.Format(time.DateOnly), Revision: result.Week.Revision, Entries: []gen.ScheduleEntry{}}
	for _, e := range result.Entries {
		var subject *string
		if e.ClassSubjectID != nil {
			id := e.ClassSubjectID.String()
			subject = &id
		}
		week.Entries = append(week.Entries, gen.ScheduleEntry{Id: e.ID.String(), Day: e.LessonDate.Format(time.DateOnly), Position: e.Position, ClassSubjectId: subject, Title: e.Title, Notes: e.Notes, StartsAt: e.StartsAt, EndsAt: e.EndsAt})
	}
	return week
}
func (h *ScheduleHandler) writeScheduleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthenticated", "Войдите в систему")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Расписание доступно только администратору")
	case errors.Is(err, domain.ErrClassNotFound), errors.Is(err, repository.ErrNotFound):
		return writeError(c, 404, "not_found", "Класс или предмет не найден")
	case errors.Is(err, domain.ErrInvalidSchedule), errors.Is(err, repository.ErrInvalidData), errors.Is(err, repository.ErrReference):
		return writeError(c, 400, "invalid_schedule", "Проверьте даты, время, предметы и номера занятий")
	case errors.Is(err, domain.ErrScheduleConflict), errors.Is(err, repository.ErrConflict):
		return writeError(c, 409, "schedule_conflict", "Расписание изменилось. Загрузите неделю заново перед сохранением")
	default:
		h.logger.Error("schedule request failed", "error", err)
		return writeError(c, 500, "internal_error", "Не удалось обработать расписание")
	}
}
