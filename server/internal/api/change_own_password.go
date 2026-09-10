package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *AuthHandler) ChangeOwnPassword(c *fiber.Ctx) error {
	var request gen.ChangeOwnPasswordRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid_request", "Некорректное тело запроса")
	}
	err := h.service.ChangeOwnPassword(
		c.UserContext(),
		c.Cookies("journal_session"),
		request.CurrentPassword,
		request.NewPassword,
	)
	if errors.Is(err, domain.ErrUnauthenticated) {
		return writeError(c, fiber.StatusUnauthorized, "unauthenticated", "Требуется вход")
	}
	if errors.Is(err, domain.ErrInvalidCredentials) {
		return writeError(c, fiber.StatusBadRequest, "invalid_password", "Текущий пароль неверен")
	}
	if errors.Is(err, domain.ErrInvalidUser) {
		return writeError(c, fiber.StatusBadRequest, "invalid_password", "Пароль не может быть пустым")
	}
	if errors.Is(err, domain.ErrForbidden) {
		return writeError(c, fiber.StatusForbidden, "forbidden", "Недостаточно прав")
	}
	if err != nil {
		return writeError(
			c,
			fiber.StatusInternalServerError,
			"internal_error",
			"Внутренняя ошибка сервера",
		)
	}
	return c.JSON(gen.ChangeOwnPasswordResponse{Changed: true})
}
