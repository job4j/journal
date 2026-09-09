package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *UserHandler) ListManagedParentStudents(c *fiber.Ctx) error {
	parentID, err := uuid.Parse(c.Params("parentId"))
	if err != nil {
		return writeError(c, 400, "invalid_parent_id", "Некорректный идентификатор родителя")
	}
	items, err := h.service.ListManagedParentStudents(c.UserContext(), c.Cookies("journal_session"), parentID)
	if err != nil {
		return h.writeUserError(c, err)
	}
	result := make([]gen.UserSummary, len(items))
	for i, item := range items {
		result[i] = userSummary(item)
	}
	return c.JSON(gen.ListParentStudentsResponse{Items: result})
}
