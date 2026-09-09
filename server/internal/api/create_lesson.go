package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *ClassHandler) CreateLesson(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("classSubjectId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_subject_id", "Некорректный идентификатор назначения")
	}
	var request gen.CreateLessonRequest
	if err = c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	if request.Position < 1 || request.Position > 32767 {
		return writeError(c, 400, "invalid_lesson", "Некорректные данные урока")
	}
	materials := make([]domain.LessonMaterialInput, len(request.Materials))
	for i, item := range request.Materials {
		if item.Position < 1 || item.Position > 32767 {
			return writeError(c, 400, "invalid_material", "Некорректные данные материала")
		}
		materials[i] = domain.LessonMaterialInput{Title: item.Title, URL: item.Url, Position: int16(item.Position)}
	}
	item, err := h.service.CreateLesson(c.UserContext(), c.Cookies("journal_session"), id, domain.CreateLessonInput{LessonDate: request.LessonDate.Time, Position: int16(request.Position), Topic: request.Topic, Homework: request.Homework, Materials: materials})
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.Status(201).JSON(gen.CreateLessonResponse{Lesson: lessonResponse(item)})
}
