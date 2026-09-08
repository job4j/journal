package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *AcademicYearHandler) CreateAcademicYear(c *fiber.Ctx) error {
	var request gen.CreateAcademicYearRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	quarters := make([]domain.AcademicYearQuarterInput, len(request.Quarters))
	for i, q := range request.Quarters {
		quarters[i] = domain.AcademicYearQuarterInput{Number: int16(q.Number), StartsOn: q.StartsOn.Time, EndsOn: q.EndsOn.Time}
	}
	result, err := h.service.CreateAcademicYear(c.UserContext(), c.Cookies("journal_session"), domain.CreateAcademicYearRequest{Name: request.Name, Status: string(request.Status), StartsOn: request.StartsOn.Time, EndsOn: request.EndsOn.Time, Quarters: quarters})
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrInvalidAcademicYear):
		return writeError(c, 400, "validation_failed", "Некорректные даты учебного года или четвертей")
	case errors.Is(err, domain.ErrAcademicYearExists):
		return writeError(c, 409, "academic_year_exists", "Учебный год уже существует")
	case err != nil:
		h.logger.Error("create academic year failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
	return c.Status(201).JSON(gen.CreateAcademicYearResponse{AcademicYear: academicYearResponse(result)})
}
