package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
)

func absenceIDs(c *fiber.Ctx) (uuid.UUID, uuid.UUID, error) {
	lessonID, err := uuid.Parse(c.Params("lessonId"))
	if err != nil {
		return uuid.Nil, uuid.Nil, writeError(c, 400, "invalid_lesson_id", "Некорректный идентификатор урока")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return uuid.Nil, uuid.Nil, writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	return lessonID, studentID, nil
}
func (h *ClassHandler) PutStudentAbsence(c *fiber.Ctx) error {
	lessonID, studentID, err := absenceIDs(c)
	if err != nil {
		return err
	}
	item, err := h.service.PutStudentAbsence(c.UserContext(), c.Cookies("journal_session"), lessonID, studentID)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.PutStudentAbsenceResponse{Absence: gen.Absence{Id: openapi_types.UUID(item.ID), LessonId: openapi_types.UUID(item.LessonID), StudentId: openapi_types.UUID(item.UserID)}})
}
func (h *ClassHandler) DeleteStudentAbsence(c *fiber.Ctx) error {
	lessonID, studentID, err := absenceIDs(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteStudentAbsence(c.UserContext(), c.Cookies("journal_session"), lessonID, studentID); err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.DeleteStudentAbsenceResponse{Removed: true})
}
