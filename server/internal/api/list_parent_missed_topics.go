package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
)

func (h *ClassHandler) ListParentMissedTopics(c *fiber.Ctx) error {
	student, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	year, err := uuid.Parse(c.Query("academicYearId"))
	if err != nil {
		return writeError(c, 400, "invalid_academic_year_id", "Некорректный идентификатор учебного года")
	}
	result, err := h.service.ListParentMissedTopics(c.UserContext(), c.Cookies("journal_session"), domain.ListParentMissedTopicsRequest{StudentID: student, AcademicYearID: year})
	if err != nil {
		return h.writeClassError(c, err)
	}
	items := make([]gen.ParentMissedTopic, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, gen.ParentMissedTopic{LessonId: item.Lesson.ID, LessonDate: openapi_types.Date{Time: item.Lesson.LessonDate}, Position: int(item.Lesson.Position), Topic: item.Lesson.Topic, SubjectId: item.Subject.ID, SubjectName: item.Subject.Name, ClassId: item.Class.ID, ClassName: item.Class.Name})
	}
	return c.JSON(gen.ListParentMissedTopicsResponse{Items: items})
}
