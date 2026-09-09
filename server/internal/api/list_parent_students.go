package api

import (
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
)

func (h *UserHandler) ListParentStudents(c *fiber.Ctx) error {
	items, err := h.service.ListParentStudents(c.UserContext(), c.Cookies("journal_session"))
	if err != nil {
		return h.writeUserError(c, err)
	}
	result := make([]gen.UserSummary, len(items))
	for i, item := range items {
		result[i] = userSummary(item)
	}
	return c.JSON(gen.ListParentStudentsResponse{Items: result})
}
