package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *AcademicYearHandler) CreateAcademicYearQuarter(c *fiber.Ctx) error {
	yearID, err := uuid.Parse(c.Params("academicYearId"))
	if err != nil {
		return writeError(c, 400, "invalid_academic_year_id", "Некорректный идентификатор учебного года")
	}
	var body gen.CreateAcademicYearQuarterRequest
	if err = c.BodyParser(&body); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.CreateAcademicYearQuarter(c.UserContext(), c.Cookies("journal_session"), yearID.String(), body.Name, body.StartsOn.Time, body.EndsOn.Time)
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrAcademicYearNotFound):
		return writeError(c, 404, "academic_year_not_found", "Учебный год не найден")
	case errors.Is(err, domain.ErrInvalidAcademicYear):
		return writeError(c, 400, "validation_failed", "Период должен находиться внутри учебного года")
	case errors.Is(err, domain.ErrAcademicYearExists):
		return writeError(c, 409, "period_overlap", "Период пересекается с существующим")
	case err != nil:
		h.logger.Error("create academic period failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
	return c.Status(201).JSON(gen.CreateAcademicYearQuarterResponse{Quarter: gen.AcademicYearQuarter{Id: openapi_types.UUID(item.ID), Number: int(item.Number), Name: item.Name, StartsOn: openapi_types.Date{Time: item.StartsOn}, EndsOn: openapi_types.Date{Time: item.EndsOn}}})
}
