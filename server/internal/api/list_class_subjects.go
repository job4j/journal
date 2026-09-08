package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func userSummary(user entity.User) gen.UserSummary {
	roles := make([]gen.RoleCode, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = gen.RoleCode(role)
	}
	return gen.UserSummary{Id: openapi_types.UUID(user.ID), FirstName: user.FirstName, LastName: user.LastName, Roles: roles}
}

func classSubjectResponse(item domain.ClassSubjectView) gen.ClassSubject {
	return gen.ClassSubject{Id: openapi_types.UUID(item.Assignment.ID), ClassId: openapi_types.UUID(item.Assignment.ClassID), Subject: gen.Subject{Id: openapi_types.UUID(item.Subject.ID), Code: item.Subject.Code, Name: item.Subject.Name}, ResponsibleTeacher: userSummary(item.Teacher)}
}

func (h *ClassHandler) ListClassSubjects(c *fiber.Ctx) error {
	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	items, err := h.service.ListClassSubjects(c.UserContext(), c.Cookies("journal_session"), classID)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.ClassSubject, len(items))
	for i, item := range items {
		result[i] = classSubjectResponse(item)
	}
	return c.JSON(gen.ListClassSubjectsResponse{Items: result})
}
