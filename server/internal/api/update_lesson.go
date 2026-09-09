package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *ClassHandler) UpdateLesson(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("lessonId"))
	if err != nil {
		return writeError(c, 400, "invalid_lesson_id", "Некорректный идентификатор урока")
	}
	var body gen.CreateLessonRequest
	if err = c.BodyParser(&body); err != nil || body.Position < 1 || body.Position > 32767 {
		return writeError(c, 400, "invalid_request", "Некорректные данные урока")
	}
	materials := make([]domain.LessonMaterialInput, len(body.Materials))
	for i, m := range body.Materials {
		materials[i] = domain.LessonMaterialInput{Title: m.Title, URL: m.Url, Position: int16(m.Position)}
	}
	item, err := h.service.UpdateLesson(c.UserContext(), c.Cookies("journal_session"), id, domain.CreateLessonInput{LessonDate: body.LessonDate.Time, Position: int16(body.Position), Topic: body.Topic, Homework: body.Homework, Materials: materials})
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.CreateLessonResponse{Lesson: lessonResponse(item)})
}
