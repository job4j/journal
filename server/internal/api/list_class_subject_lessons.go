package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"time"
)

func lessonResponse(item domain.LessonView) gen.Lesson {
	materials := make([]gen.LessonMaterial, len(item.Materials))
	for i, value := range item.Materials {
		materials[i] = gen.LessonMaterial{Id: openapi_types.UUID(value.ID), Title: value.Title, Url: value.URL, Position: int(value.Position)}
	}
	grades := make([]gen.GradeItem, len(item.GradeItems))
	for i, value := range item.GradeItems {
		var max *float64
		if value.MaxScore != nil {
			number := float64(*value.MaxScore)
			max = &number
		}
		grades[i] = gen.GradeItem{Id: openapi_types.UUID(value.ID), LessonId: openapi_types.UUID(value.LessonID), Title: value.Title, Kind: gen.GradeItemKind(value.Kind), GradingScale: gen.GradingScale(value.GradingScale), MaxScore: max}
	}
	return gen.Lesson{Id: openapi_types.UUID(item.Lesson.ID), ClassSubjectId: openapi_types.UUID(item.Lesson.ClassSubjectID), LessonDate: openapi_types.Date{Time: item.Lesson.LessonDate}, Position: int(item.Lesson.Position), Topic: item.Lesson.Topic, Homework: item.Lesson.Homework, Materials: materials, GradeItems: grades}
}
func parseOptionalDate(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}
	value, err := time.Parse("2006-01-02", raw)
	return &value, err
}
func (h *ClassHandler) ListClassSubjectLessons(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("classSubjectId"))
	if err != nil {
		return writeError(c, 400, "invalid_class_subject_id", "Некорректный идентификатор назначения")
	}
	from, err := parseOptionalDate(c.Query("dateFrom"))
	if err != nil {
		return writeError(c, 400, "invalid_date_from", "Некорректная начальная дата")
	}
	to, err := parseOptionalDate(c.Query("dateTo"))
	if err != nil {
		return writeError(c, 400, "invalid_date_to", "Некорректная конечная дата")
	}
	items, err := h.service.ListClassSubjectLessons(c.UserContext(), c.Cookies("journal_session"), id, from, to)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.Lesson, len(items))
	for i, item := range items {
		result[i] = lessonResponse(item)
	}
	return c.JSON(gen.ListClassSubjectLessonsResponse{Items: result})
}
