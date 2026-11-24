package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func attachmentResponse(item entity.LessonAttachment) gen.LessonAttachment {
	return gen.LessonAttachment{Id: item.ID, LessonId: item.LessonID, FileName: item.FileName, ContentType: item.ContentType, SizeBytes: item.SizeBytes, CreatedAt: item.CreatedAt}
}
func attachmentStudent(c *fiber.Ctx) (uuid.UUID, error) {
	if c.Params("studentId") == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(c.Params("studentId"))
}
func (h *AttachmentHandler) ListLessonAttachments(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("lessonId"))
	if err != nil {
		return writeError(c, 400, "invalid_lesson_id", "Некорректный идентификатор урока")
	}
	student, err := attachmentStudent(c)
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	result, err := h.service.ListLessonAttachments(c.UserContext(), c.Cookies("journal_session"), domain.ListLessonAttachmentsRequest{LessonID: id, StudentID: student})
	if err != nil {
		return h.writeAttachmentError(c, err)
	}
	items := make([]gen.LessonAttachment, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, attachmentResponse(item))
	}
	return c.JSON(gen.ListLessonAttachmentsResponse{Items: items})
}
