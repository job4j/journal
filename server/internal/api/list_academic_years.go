package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *AcademicYearHandler) ListAcademicYears(c *fiber.Ctx) error {
	items, err := h.service.ListAcademicYears(c.UserContext(), c.Cookies("journal_session"))
	if errors.Is(err, domain.ErrUnauthenticated) {
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	}
	if errors.Is(err, domain.ErrForbidden) {
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	}
	if err != nil {
		h.logger.Error("list academic years failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
	result := make([]gen.AcademicYear, len(items))
	for i, item := range items {
		result[i] = academicYearResponse(item)
	}
	return c.JSON(gen.ListAcademicYearsResponse{Items: result})
}

func academicYearResponse(item domain.AcademicYearView) gen.AcademicYear {
	quarters := make([]gen.AcademicYearQuarter, len(item.Quarters))
	for j, q := range item.Quarters {
		quarters[j] = gen.AcademicYearQuarter{Id: openapi_types.UUID(q.ID), Number: int(q.Number), Name: q.Name, StartsOn: openapi_types.Date{Time: q.StartsOn}, EndsOn: openapi_types.Date{Time: q.EndsOn}}
	}
	return gen.AcademicYear{Id: openapi_types.UUID(item.Year.ID), Name: item.Year.Name, StartsOn: openapi_types.Date{Time: item.Year.StartsOn}, EndsOn: openapi_types.Date{Time: item.Year.EndsOn}, Status: gen.AcademicYearStatus(item.Year.Status), Quarters: quarters}
}
