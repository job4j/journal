package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func apiUser(user entity.User) gen.User {
	roles := make([]gen.RoleCode, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = gen.RoleCode(role)
	}
	return userResponse(user) // user.Name: user.Name, Status: gen.UserStatus(user.Status), Roles: roles}
}

func (h *AuthHandler) CurrentUser(c *fiber.Ctx) error {
	user, err := h.service.CurrentUser(c.UserContext(), c.Cookies("journal_session"))
	if errors.Is(err, domain.ErrUnauthenticated) {
		return writeError(c, fiber.StatusUnauthorized, "unauthenticated", "Требуется вход")
	}
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
	}
	return c.JSON(gen.GetCurrentUserResponse{User: apiUser(user)})
}
