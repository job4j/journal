package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func academicYearQuarterResponse(item entity.AcademicYearQuarter) gen.AcademicYearQuarter {
	return gen.AcademicYearQuarter{Id: openapi_types.UUID(item.ID), Number: int(item.Number), Name: item.Name, StartsOn: openapi_types.Date{Time: item.StartsOn}, EndsOn: openapi_types.Date{Time: item.EndsOn}}
}
func (h *AcademicYearHandler) writeQuarterError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrQuarterNotFound):
		return writeError(c, 404, "quarter_not_found", "Период не найден")
	case errors.Is(err, domain.ErrInvalidAcademicYear):
		return writeError(c, 400, "validation_failed", "Период должен находиться внутри учебного года")
	case errors.Is(err, domain.ErrAcademicYearExists):
		return writeError(c, 409, "period_overlap", "Период пересекается с существующим")
	case errors.Is(err, domain.ErrQuarterInUse):
		return writeError(c, 409, "quarter_in_use", "Период содержит итоговые оценки и не может быть удалён")
	default:
		h.logger.Error("manage academic period failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
}
