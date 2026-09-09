package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
)

func (h *ClassHandler) GetParentStudentJournal(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("studentId"))
	if err != nil {
		return writeError(c, 400, "invalid_student_id", "Некорректный идентификатор ученика")
	}
	yearID, err := uuid.Parse(c.Query("academicYearId"))
	if err != nil {
		return writeError(c, 400, "invalid_academic_year_id", "Некорректный идентификатор учебного года")
	}
	journal, err := h.service.GetParentStudentJournal(c.UserContext(), c.Cookies("journal_session"), studentID, yearID)
	if err != nil {
		return h.writeClassError(c, err)
	}
	subjects := make([]gen.JournalSubject, len(journal.Subjects))
	for i, subject := range journal.Subjects {
		lessons := make([]gen.JournalLesson, len(subject.Lessons))
		for j, lesson := range subject.Lessons {
			response := lessonResponse(lesson)
			scores := []gen.JournalScore{}
			for _, item := range response.GradeItems {
				for _, score := range item.Scores {
					scores = append(scores, gen.JournalScore{GradeItem: item, Score: score})
				}
			}
			lessons[j] = gen.JournalLesson{Lesson: response, Scores: scores}
		}
		subjects[i] = gen.JournalSubject{ClassSubject: classSubjectResponse(subject.ClassSubjectView), Lessons: lessons}
	}
	return c.JSON(gen.GetParentStudentJournalResponse{Student: userSummary(journal.Student), AcademicYear: academicYearResponse(journal.AcademicYear), Class: classResponse(journal.Class), Subjects: subjects})
}
