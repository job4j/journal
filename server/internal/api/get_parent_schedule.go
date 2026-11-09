package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
	"time"
)

func (h *ScheduleHandler) GetParentSchedule(c *fiber.Ctx) error {
	params, err := scheduleParams(c)
	if err != nil {
		return writeError(c, 400, "invalid_schedule", "Некорректный класс или дата недели")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный ученик")
	}
	result, err := h.service.GetParentSchedule(c.UserContext(), c.Cookies("journal_session"), domain.GetParentScheduleRequest{GetClassScheduleRequest: params, StudentID: studentID})
	if err != nil {
		return h.writeScheduleError(c, err)
	}
	week := gen.ParentScheduleWeek{WeekStart: result.WeekStart.Format(time.DateOnly), Entries: []gen.ParentScheduleEntry{}}
	for _, item := range result.Entries {
		e := item.Entry
		var subject *string
		if e.ClassSubjectID != nil {
			id := e.ClassSubjectID.String()
			subject = &id
		}
		week.Entries = append(week.Entries, gen.ParentScheduleEntry{Id: e.ID.String(), Day: e.LessonDate.Format(time.DateOnly), Position: e.Position, ClassSubjectId: subject, Title: e.Title, Notes: e.Notes, StartsAt: e.StartsAt, EndsAt: e.EndsAt, Homework: item.Homework, HomeworkFrom: item.HomeworkFrom})
	}
	return c.JSON(gen.GetParentScheduleResponse{Week: week})
}
