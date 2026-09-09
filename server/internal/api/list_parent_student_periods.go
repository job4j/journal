package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) ListParentStudentPeriods(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	items, err := h.service.ListParentStudentPeriods(c.UserContext(), c.Cookies("journal_session"), studentID)
	if err != nil {
		return h.writeClassError(c, err)
	}
	result := make([]gen.ParentStudentPeriod, len(items))
	for i, item := range items {
		subjects := make([]gen.ClassSubject, len(item.Subjects))
		for j, subject := range item.Subjects {
			subjects[j] = classSubjectResponse(subject)
		}
		result[i] = gen.ParentStudentPeriod{AcademicYear: academicYearResponse(item.AcademicYear), Class: classResponse(item.Class), Subjects: subjects}
	}
	return c.JSON(gen.ListParentStudentPeriodsResponse{Items: result})
}
