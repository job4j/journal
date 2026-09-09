package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

func quarterGradeResponse(item entity.QuarterGrade) gen.QuarterGrade {
	return gen.QuarterGrade{Id: openapi_types.UUID(item.ID), QuarterId: openapi_types.UUID(item.QuarterID), ClassSubjectId: openapi_types.UUID(item.ClassSubjectID), StudentId: openapi_types.UUID(item.UserID), GradingScale: gen.GradingScale(item.GradingScale), MaxScore: item.MaxScore, NumericValue: item.NumericValue, TextValue: item.TextValue, TeacherComment: item.TeacherComment}
}
func quarterGradeIDs(c *fiber.Ctx) (uuid.UUID, uuid.UUID, error) {
	assignmentID, err := uuid.Parse(c.Params("classSubjectId"))
	if err != nil {
		return uuid.Nil, uuid.Nil, writeError(c, 400, "invalid_class_subject_id", "Некорректный идентификатор предмета")
	}
	quarterID, err := uuid.Parse(c.Params("quarterId"))
	if err != nil {
		return uuid.Nil, uuid.Nil, writeError(c, 400, "invalid_quarter_id", "Некорректный идентификатор четверти")
	}
	return assignmentID, quarterID, nil
}
func (h *ClassHandler) PutQuarterGrade(c *fiber.Ctx) error {
	assignmentID, quarterID, err := quarterGradeIDs(c)
	if err != nil {
		return err
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	var body gen.PutQuarterGradeRequest
	if err = c.BodyParser(&body); err != nil {
		return writeError(c, 400, "invalid_request", "Некорректное тело запроса")
	}
	item, err := h.service.PutQuarterGrade(c.UserContext(), c.Cookies("journal_session"), assignmentID, quarterID, studentID, domain.QuarterGradeInput{GradingScale: string(body.GradingScale), MaxScore: body.MaxScore, NumericValue: body.NumericValue, TextValue: body.TextValue, TeacherComment: body.TeacherComment})
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.PutQuarterGradeResponse{QuarterGrade: quarterGradeResponse(item)})
}
func (h *ClassHandler) ListQuarterGrades(c *fiber.Ctx) error {
	assignmentID, quarterID, err := quarterGradeIDs(c)
	if err != nil {
		return err
	}
	items, err := h.service.ListQuarterGrades(c.UserContext(), c.Cookies("journal_session"), assignmentID, quarterID)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.QuarterGrade, len(items))
	for i, item := range items {
		result[i] = quarterGradeResponse(item)
	}
	return c.JSON(gen.ListQuarterGradesResponse{Items: result})
}
