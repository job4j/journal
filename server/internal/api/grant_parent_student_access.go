package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
)

func (h *UserHandler) GrantParentStudentAccess(c *fiber.Ctx) error {
	parentID, err := uuid.Parse(c.Params("parentId"))
	if err != nil {
		return writeError(c, 400, "invalid_parent_id", "Некорректный идентификатор родителя")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	var request gen.GrantParentStudentAccessRequest
	if err := c.BodyParser(&request); err != nil || !request.CanViewProfile.Valid() || !request.CanViewJournal.Valid() {
		return writeError(c, 400, "invalid_request", "Для связи необходим доступ к профилю и журналу")
	}
	result, err := h.service.GrantParentStudentAccess(c.UserContext(), c.Cookies("journal_session"), parentID, studentID)
	if err != nil {
		return h.writeUserError(c, err)
	}
	return c.JSON(gen.GrantParentStudentAccessResponse{ParentId: openapi_types.UUID(result.ParentID), StudentId: openapi_types.UUID(result.StudentID), Permissions: result.Permissions})
}
