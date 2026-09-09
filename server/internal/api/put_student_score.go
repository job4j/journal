package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"journal/server/gen"
)

type putScoreBody struct {
	NumericValue   *float64 `json:"numericValue"`
	TextValue      *string  `json:"textValue"`
	TeacherComment *string  `json:"teacherComment"`
}

func (h *ClassHandler) PutStudentScore(c *fiber.Ctx) error {
	gradeItemID, err := uuid.Parse(c.Params("gradeItemId"))
	if err != nil {
		return writeError(c, 400, "invalid_grade_item_id", "Некорректный идентификатор работы")
	}
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	var body putScoreBody
	if err = c.BodyParser(&body); err != nil || (body.NumericValue == nil) == (body.TextValue == nil) {
		return writeError(c, 400, "invalid_request", "Нужно передать один вид оценки")
	}
	score, err := h.service.PutStudentScore(c.UserContext(), c.Cookies("journal_session"), gradeItemID, studentID, body.NumericValue, body.TextValue, body.TeacherComment)
	if err != nil {
		return h.writeClassError(c, err)
	}
	return c.JSON(gen.PutStudentScoreResponse{Score: gen.Score{Id: openapi_types.UUID(score.ID), GradeItemId: openapi_types.UUID(score.GradeItemID), StudentId: openapi_types.UUID(score.UserID), NumericValue: score.NumericValue, TextValue: score.TextValue, TeacherComment: score.TeacherComment}})
}
