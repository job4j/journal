package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *ClassHandler) ListClassStudents(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_id", "Некорректный идентификатор класса")
	}
	items, err := h.service.ListClassStudents(c.UserContext(), c.Cookies("journal_session"), id)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.ClassStudent, len(items))
	for i, item := range items {
		result[i] = classStudentResponse(item)
	}
	return c.JSON(gen.ListClassStudentsResponse{Items: result})
}

func classStudentResponse(item domain.ClassStudentView) gen.ClassStudent {
	roles := make([]gen.RoleCode, len(item.Student.Roles))
	for i, role := range item.Student.Roles {
		roles[i] = gen.RoleCode(role)
	}
	var leftOn *openapi_types.Date
	if item.Membership.LeftOn != nil {
		leftOn = &openapi_types.Date{Time: *item.Membership.LeftOn}
	}
	return gen.ClassStudent{Student: gen.UserSummary{Id: openapi_types.UUID(item.Student.ID), FirstName: item.Student.FirstName, LastName: item.Student.LastName, Roles: roles}, EnrolledOn: openapi_types.Date{Time: item.Membership.EnrolledOn}, LeftOn: leftOn}
}
