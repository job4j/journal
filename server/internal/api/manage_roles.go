package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func roleResponse(view domain.RoleView) gen.Role {
	permissions := make([]gen.PermissionCode, len(view.PermissionCodes))
	for i, code := range view.PermissionCodes {
		permissions[i] = gen.PermissionCode(code)
	}
	return gen.Role{Id: openapi_types.UUID(view.Role.ID), Code: gen.RoleCode(view.Role.Code), Name: view.Role.Name, Permissions: permissions}
}
func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	result, err := h.service.ListRoles(c.UserContext(), token)
	if err != nil {
		return h.writeManagementError(c, err)
	}
	roles := make([]gen.Role, len(result.Roles))
	for i, role := range result.Roles {
		roles[i] = roleResponse(role)
	}
	return c.JSON(gen.ListRolesResponse{Roles: roles})
}
func (h *RoleHandler) ListPermissions(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	result, err := h.service.ListPermissions(c.UserContext(), token)
	if err != nil {
		return h.writeManagementError(c, err)
	}
	permissions := make([]gen.Permission, len(result.Permissions))
	for i, p := range result.Permissions {
		permissions[i] = gen.Permission{Code: gen.PermissionCode(p.Code), Description: p.Description}
	}
	return c.JSON(gen.ListPermissionsResponse{Permissions: permissions})
}
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	id, err := uuid.Parse(c.Params("roleId"))
	if err != nil {
		return writeError(c, 400, "invalid_role_id", "Некорректный идентификатор роли")
	}
	var request gen.UpdateRoleRequest
	if err = c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	codes := make([]string, len(request.Permissions))
	for i, code := range request.Permissions {
		codes[i] = string(code)
	}
	result, err := h.service.UpdateRole(c.UserContext(), token, id, string(request.Code), request.Name, codes)
	if err != nil {
		return h.writeManagementError(c, err)
	}
	return c.JSON(gen.UpdateRoleResponse{Role: roleResponse(result)})
}
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	id, err := uuid.Parse(c.Params("roleId"))
	if err != nil {
		return writeError(c, 400, "invalid_role_id", "Некорректный идентификатор роли")
	}
	if err = h.service.DeleteRole(c.UserContext(), token, id); err != nil {
		return h.writeManagementError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *RoleHandler) writeManagementError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrInvalidRole):
		return writeError(c, 400, "validation_failed", "Некорректные данные роли")
	case errors.Is(err, domain.ErrPermissionNotFound):
		return writeError(c, 404, "permission_not_found", "Одно или несколько разрешений не найдены")
	case errors.Is(err, domain.ErrRoleNotFound):
		return writeError(c, 404, "role_not_found", "Роль не найдена")
	case errors.Is(err, domain.ErrRoleExists):
		return writeError(c, 409, "role_exists", "Роль с таким кодом уже существует")
	case errors.Is(err, domain.ErrRoleInUse):
		return writeError(c, 409, "role_in_use", "Роль назначена пользователям")
	case errors.Is(err, domain.ErrProtectedRole):
		return writeError(c, 409, "protected_role", "Системную роль нельзя изменить таким способом")
	default:
		h.logger.Error("manage roles failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
}
