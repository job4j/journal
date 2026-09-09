package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) ListTeacherClasses(c *fiber.Ctx) error {
	items, err := h.service.ListTeacherClasses(c.UserContext(), c.Cookies("journal_session"))
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.Class, len(items))
	for i, item := range items {
		result[i] = classResponse(item)
	}
	return c.JSON(gen.ListTeacherClassesResponse{Items: result})
}
func (h *ClassHandler) ListTeacherClassSubjects(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	items, err := h.service.ListTeacherClassSubjects(c.UserContext(), c.Cookies("journal_session"), id)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.ClassSubject, len(items))
	for i, item := range items {
		result[i] = classSubjectResponse(item)
	}
	return c.JSON(gen.ListTeacherClassSubjectsResponse{Items: result})
}
