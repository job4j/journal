package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func subjectResponse(item entity.Subject) gen.Subject {
	return gen.Subject{Id: openapi_types.UUID(item.ID), Code: item.Code, Name: item.Name}
}
func (h *SubjectHandler) ListSubjects(c *fiber.Ctx) error {
	items, err := h.service.ListSubjects(c.UserContext(), c.Cookies("journal_session"))
	if err != nil {
		return h.writeSubjectError(c, err)
	}
	result := make([]gen.Subject, len(items))
	for i, item := range items {
		result[i] = subjectResponse(item)
	}
	return c.JSON(gen.ListSubjectsResponse{Items: result})
}
func (h *SubjectHandler) CreateSubject(c *fiber.Ctx) error {
	var request gen.CreateSubjectRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.CreateSubject(c.UserContext(), c.Cookies("journal_session"), request.Code, request.Name)
	if err != nil {
		return h.writeSubjectError(c, err)
	}
	return c.Status(201).JSON(gen.CreateSubjectResponse{Subject: subjectResponse(item)})
}
func (h *SubjectHandler) writeSubjectError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrInvalidSubject):
		return writeError(c, 400, "validation_failed", "Некорректные данные предмета")
	case errors.Is(err, domain.ErrSubjectExists):
		return writeError(c, 409, "subject_exists", "Предмет уже существует")
	default:
		h.logger.Error("manage subjects failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
}
