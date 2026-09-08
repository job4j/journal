package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, fiber.StatusUnauthorized, "unauthorized", "Требуется авторизация")
	}
	var request gen.CreateRoleRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid_request", "Некорректное тело запроса")
	}
	permissionCodes := make([]string, len(request.Permissions))
	for i, code := range request.Permissions {
		permissionCodes[i] = string(code)
	}
	result, err := h.service.CreateRole(c.UserContext(), token, string(request.Code), request.Name, permissionCodes)
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, fiber.StatusUnauthorized, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, fiber.StatusForbidden, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrInvalidRole):
		return writeError(c, fiber.StatusBadRequest, "validation_failed", "Некорректные данные роли")
	case errors.Is(err, domain.ErrPermissionNotFound):
		return writeError(c, fiber.StatusNotFound, "permission_not_found", "Одно или несколько разрешений не найдены")
	case errors.Is(err, domain.ErrRoleExists):
		return writeError(c, fiber.StatusConflict, "role_exists", "Роль с таким кодом уже существует")
	case err != nil:
		h.logger.Error("create role failed", "error", err)
		return writeError(c, fiber.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера")
	}
	permissions := make([]gen.PermissionCode, len(result.PermissionCodes))
	for i, code := range result.PermissionCodes {
		permissions[i] = gen.PermissionCode(code)
	}
	return c.Status(fiber.StatusCreated).JSON(gen.CreateRoleResponse{Role: gen.Role{Id: openapi_types.UUID(result.Role.ID), Code: gen.RoleCode(result.Role.Code), Name: result.Role.Name, Permissions: permissions}})
}
