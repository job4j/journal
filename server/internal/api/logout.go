package api

import (
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
)

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.service.Logout(c.UserContext(), c.Cookies("journal_session")); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
	}
	c.Cookie(&fiber.Cookie{Name: "journal_session", Value: "", MaxAge: -1, HTTPOnly: true, Secure: h.secureCookie, SameSite: fiber.CookieSameSiteLaxMode, Path: "/"})
	return c.JSON(gen.LogoutResponse{LoggedOut: true})
}
