package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"mime"
)

func (h *AttachmentHandler) GetLessonAttachment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("attachmentId"))
	if err != nil {
		return writeError(c, 400, "invalid_attachment_id", "Некорректный идентификатор файла")
	}
	student, err := attachmentStudent(c)
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	result, content, err := h.service.GetLessonAttachment(c.UserContext(), c.Cookies("journal_session"), domain.GetLessonAttachmentRequest{AttachmentID: id, StudentID: student})
	if err != nil {
		return h.writeAttachmentError(c, err)
	}
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": result.Attachment.FileName}))
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Cache-Control", "private, no-store")
	// Fiber closes the stream after sending the response.
	return c.SendStream(content, int(result.Attachment.SizeBytes))
}
