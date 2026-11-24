package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/internal/domain"
)

func (h *AttachmentHandler) DeleteLessonAttachment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("attachmentId"))
	if err != nil {
		return writeError(c, 400, "invalid_attachment_id", "Некорректный идентификатор файла")
	}
	if err = h.service.DeleteLessonAttachment(c.UserContext(), c.Cookies("journal_session"), domain.DeleteLessonAttachmentRequest{AttachmentID: id}); err != nil {
		return h.writeAttachmentError(c, err)
	}
	return c.SendStatus(204)
}
