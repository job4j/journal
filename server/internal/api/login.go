package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request gen.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid_request", "Некорректное тело запроса")
	}
	result, err := h.service.Login(c.UserContext(), request.Login, request.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		return writeError(c, fiber.StatusUnauthorized, "invalid_credentials", "Неверный логин или пароль")
	}
	if errors.Is(err, domain.ErrUserBlocked) {
		return writeError(c, fiber.StatusUnauthorized, "user_blocked", "Пользователь заблокирован")
	}
	if err != nil {
		h.logger.Error("login failed", "error", err)
		return writeError(c, fiber.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
	}
	c.Cookie(&fiber.Cookie{Name: "journal_session", Value: result.Token, Expires: result.ExpiresAt, HTTPOnly: true, Secure: h.secureCookie, SameSite: fiber.CookieSameSiteLaxMode, Path: "/"})
	return c.Status(fiber.StatusOK).JSON(gen.LoginResponse{User: apiUser(result.User)})
}
