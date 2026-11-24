package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"journal/server/gen"
	"journal/server/internal/domain"
	"net/http"
)

func (h *AttachmentHandler) CreateLessonAttachment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("lessonId"))
	if err != nil {
		return writeError(c, 400, "invalid_lesson_id", "Некорректный идентификатор урока")
	}
	form, err := c.MultipartForm()
	if err != nil || len(form.File) != 1 || len(form.File["file"]) != 1 {
		return h.writeAttachmentError(c, domain.ErrInvalidAttachment)
	}
	header := form.File["file"][0]
	if header.Size <= 0 || header.Size > domain.MaxAttachmentBytes {
		return h.writeAttachmentError(c, domain.ErrInvalidAttachment)
	}
	file, err := header.Open()
	if err != nil {
		return h.writeAttachmentError(c, err)
	}
	defer file.Close()
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return h.writeAttachmentError(c, err)
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return h.writeAttachmentError(c, err)
	}
	result, err := h.service.CreateLessonAttachment(c.UserContext(), c.Cookies("journal_session"), domain.CreateLessonAttachmentRequest{LessonID: id, FileName: header.Filename, ContentType: http.DetectContentType(buffer[:n]), SizeBytes: header.Size}, file)
	if err != nil {
		return h.writeAttachmentError(c, err)
	}
	return c.Status(201).JSON(gen.CreateLessonAttachmentResponse{Attachment: attachmentResponse(result.Attachment)})
}
