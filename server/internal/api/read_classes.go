package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func classResponse(item domain.ClassView) gen.Class {
	return gen.Class{Id: openapi_types.UUID(item.Class.ID), AcademicYearId: openapi_types.UUID(item.Class.AcademicYearID), Name: item.Class.Name, GradeLevel: int(item.Class.GradeLevel), StudentCount: item.StudentCount}
}
func (h *ClassHandler) ListClasses(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Query("academicYearId"))
	if err != nil {
		return writeError(c, 400, "invalid_academic_year_id", "Некорректный идентификатор учебного года")
	}
	items, err := h.service.ListClasses(c.UserContext(), c.Cookies("journal_session"), id)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.Class, len(items))
	for i, item := range items {
		result[i] = classResponse(item)
	}
	return c.JSON(gen.ListClassesResponse{Items: result})
}
func (h *ClassHandler) GetClass(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	item, err := h.service.GetClass(c.UserContext(), c.Cookies("journal_session"), id)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.GetClassResponse{Class: classResponse(item)})
}
func (h *ClassHandler) writeClassError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrClassNotFound):
		return writeError(c, 404, "class_not_found", "Класс не найден")
	case errors.Is(err, domain.ErrAcademicYearNotFound):
		return writeError(c, 404, "academic_year_not_found", "Учебный год не найден")
	case errors.Is(err, domain.ErrUserNotFound):
		return writeError(c, 404, "student_not_found", "Ученик не найден")
	case errors.Is(err, domain.ErrInvalidClassStudent):
		return writeError(c, 400, "invalid_class_student", "Некорректные данные зачисления")
	case errors.Is(err, domain.ErrClassStudentExists):
		return writeError(c, 409, "class_student_exists", "Ученик уже зачислен в класс")
	case errors.Is(err, domain.ErrClassStudentNotFound):
		return writeError(c, 404, "class_student_not_found", "Ученик не состоит в классе")
	case errors.Is(err, domain.ErrInvalidClass):
		return writeError(c, 400, "validation_failed", "Некорректные данные класса")
	case errors.Is(err, domain.ErrClassExists):
		return writeError(c, 409, "class_exists", "Класс уже существует")
	default:
		h.logger.Error("read classes failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
}
