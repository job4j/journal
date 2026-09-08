package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
)

func (h *UserHandler) RevokeParentStudentAccess(c *fiber.Ctx) error {
	parentID, err := uuid.Parse(c.Params("parentId"))
	if err != nil {
		return writeError(c, 400, "invalid_parent_id", "Некорректный идентификатор родителя")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	result, err := h.service.RevokeParentStudentAccess(c.UserContext(), c.Cookies("journal_session"), parentID, studentID)
	if err != nil {
		return h.writeUserError(c, err)
	}
	return c.JSON(gen.RevokeParentStudentAccessResponse{ParentId: openapi_types.UUID(result.ParentID), StudentId: openapi_types.UUID(result.StudentID), Revoked: gen.RevokeParentStudentAccessResponseRevokedTrue})
}
