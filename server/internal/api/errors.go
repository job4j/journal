package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *AttachmentHandler) writeAttachmentError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrUnauthenticated):
		return writeError(c, 401, "unauthorized", "Сессия недействительна")
	case errors.Is(err, domain.ErrForbidden):
		return writeError(c, 403, "forbidden", "Недостаточно прав")
	case errors.Is(err, domain.ErrAttachmentNotFound):
		return writeError(c, 404, "attachment_not_found", "Файл не найден")
	case errors.Is(err, domain.ErrLessonNotFound):
		return writeError(c, 404, "lesson_not_found", "Урок не найден")
	case errors.Is(err, domain.ErrInvalidLesson):
		return writeError(c, 400, "invalid_lesson", "Сначала сохраните тему урока")
	case errors.Is(err, domain.ErrInvalidAttachment):
		return writeError(c, 400, "invalid_attachment", "Выберите непустой файл размером не более 20 МБ с именем до 255 байт")
	case errors.Is(err, domain.ErrAttachmentStorageUnavailable):
		h.logger.Error("attachment storage unavailable", "error", err)
		return writeError(c, 503, "storage_unavailable", "Хранилище файлов временно недоступно. Повторите попытку")
	default:
		h.logger.Error("attachment operation failed", "error", err)
		return writeError(c, 500, "internal_error", "Внутренняя ошибка сервера")
	}
}

func writeError(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(gen.ErrorResponse{Code: code, Message: message})
}
