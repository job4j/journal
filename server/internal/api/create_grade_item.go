package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
)

func (h *ClassHandler) CreateGradeItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("lessonId"))
	if err != nil {
		return writeError(c, 400, "invalid_lesson_id", "Некорректный идентификатор урока")
	}
	var request gen.CreateGradeItemRequest
	if err = c.BodyParser(&request); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.CreateGradeItem(c.UserContext(), c.Cookies("journal_session"), id, request.Title, string(request.Kind), string(request.GradingScale), request.MaxScore)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.Status(201).JSON(gen.CreateGradeItemResponse{GradeItem: gen.GradeItem{Id: openapi_types.UUID(item.ID), LessonId: openapi_types.UUID(item.LessonID), Title: item.Title, Kind: gen.GradeItemKind(item.Kind), GradingScale: gen.GradingScale(item.GradingScale), MaxScore: item.MaxScore, Scores: []gen.Score{}}})
}
