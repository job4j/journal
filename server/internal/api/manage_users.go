package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
	"strconv"
)

func userResponse(user entity.User) gen.User {
	roles := make([]gen.RoleCode, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = gen.RoleCode(role)
	}
	var email *openapi_types.Email
	if user.Email != "" {
		value := openapi_types.Email(user.Email)
		email = &value
	}
	var phone *string
	if user.Phone != "" {
		phone = &user.Phone
	}
	return gen.User{Id: openapi_types.UUID(user.ID), Login: user.Login, Email: email, Phone: phone, FirstName: user.FirstName, LastName: user.LastName, Status: gen.UserStatus(user.Status), Roles: roles}
}
func optionalEmail(value *openapi_types.Email) string {
	if value == nil {
		return ""
	}
	return string(*value)
}
func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
func userInput(login, email, phone, firstName, lastName, password, status string, roles []gen.RoleCode) domain.UserInput {
	codes := make([]string, len(roles))
	for i, role := range roles {
		codes[i] = string(role)
	}
	return domain.UserInput{Login: login, Email: email, Phone: phone, FirstName: firstName, LastName: lastName, Password: password, Status: status, Roles: codes}
}
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	var request gen.CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	result, err := h.service.CreateUser(c.UserContext(), token, userInput(request.Login, optionalEmail(request.Email), optionalString(request.Phone), request.FirstName, request.LastName, request.Password, "active", request.Roles))
	if err != nil {
		return h.writeUserError(c, err)
	}
	return c.Status(201).JSON(gen.CreateUserResponse{User: userResponse(result)})
}
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	limit := 50
	offset := 0
	var err error
	if raw := c.Query("limit"); raw != "" {
		limit, err = strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > 100 {
			return writeError(c, 400, "invalid_pagination", "Некорректный limit")
		}
	}
	if raw := c.Query("offset"); raw != "" {
		offset, err = strconv.Atoi(raw)
		if err != nil || offset < 0 {
			return writeError(c, 400, "invalid_pagination", "Некорректный offset")
		}
	}
	result, err := h.service.ListUsers(c.UserContext(), token, c.Query("role"), limit, offset)
	if err != nil {
		return h.writeUserError(c, err)
	}
	items := make([]gen.User, len(result.Items))
	for i, user := range result.Items {
		items[i] = userResponse(user)
	}
	return c.JSON(gen.ListUsersResponse{Items: items, Total: result.Total})
}
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	id, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return writeError(c, 400, "invalid_user_id", "Некорректный идентификатор пользователя")
	}
	result, err := h.service.GetUser(c.UserContext(), token, id)
	if err != nil {
		return h.writeUserError(c, err)
	}
	return c.JSON(gen.GetUserResponse{User: userResponse(result)})
}
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	id, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return writeError(c, 400, "invalid_user_id", "Некорректный идентификатор пользователя")
	}
	var request gen.UpdateUserRequest
	if err = c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	password := ""
	if request.Password != nil {
		password = *request.Password
	}
	result, err := h.service.UpdateUser(c.UserContext(), token, id, userInput(request.Login, optionalEmail(request.Email), optionalString(request.Phone), request.FirstName, request.LastName, password, string(request.Status), request.Roles))
	if err != nil {
		return h.writeUserError(c, err)
	}
	return c.JSON(gen.UpdateUserResponse{User: userResponse(result)})
}
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	token := c.Cookies("journal_session")
	if token == "" {
		return writeError(c, 401, "unauthorized", "Требуется авторизация")
	}
	id, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return writeError(c, 400, "invalid_user_id", "Некорректный идентификатор пользователя")
	}
	if err = h.service.DeleteUser(c.UserContext(), token, id); err != nil {
		return h.writeUserError(c, err)
	}
	return c.SendStatus(204)
}
func (h *UserHandler) writeUserError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrInvalidUser):
		return writeError(c, 400, "validation_failed", "Некорректные данные пользователя")
	case errors.Is(err, domain.ErrInvalidParentStudent):
		return writeError(c, 400, "invalid_parent_student", "Требуются пользователи с ролями родителя и ученика")
	case errors.Is(err, domain.ErrRoleNotFound):
		return writeError(c, 404, "role_not_found", "Одна или несколько ролей не найдены")
	case errors.Is(err, domain.ErrUserNotFound):
		return writeError(c, 404, "user_not_found", "Пользователь не найден")
	case errors.Is(err, domain.ErrUserExists):
		return writeError(c, 409, "user_exists", "Пользователь с таким email уже существует")
	case errors.Is(err, domain.ErrUserInUse):
		return writeError(c, 409, "user_in_use", "Пользователь связан с учебными данными")
	default:
		h.logger.Error("manage users failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
}
