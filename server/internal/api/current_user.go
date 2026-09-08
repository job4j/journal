package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func apiUser(user entity.User) gen.User {
	roles := make([]gen.RoleCode, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = gen.RoleCode(role)
	}
	return gen.User{Id: openapi_types.UUID(user.ID), Email: openapi_types.Email(user.Email), FirstName: user.FirstName, LastName: user.LastName, Status: gen.UserStatus(user.Status), Roles: roles}
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
