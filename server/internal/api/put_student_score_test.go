package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"journal/server/gen"
	"journal/server/internal/domain"
	"journal/server/internal/repository/entity"
)

type scoreServiceStub struct {
	ClassService
	remediated bool
	comment    *string
	err        error
}

func (s *scoreServiceStub) PutStudentScore(_ context.Context, _ string, itemID, studentID uuid.UUID, numeric *float64, text, comment *string, remediated bool) (entity.Score, error) {
	s.remediated, s.comment = remediated, comment
	return entity.Score{ID: uuid.New(), GradeItemID: itemID, UserID: studentID, NumericValue: numeric, TextValue: text, TeacherComment: comment, IsRemediated: remediated}, s.err
}
func TestPutStudentScoreExposesRemediationAndValidationError(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{{"saved", nil, 200}, {"missing explanation", domain.ErrRemediationCommentRequired, 400}} {
		t.Run(tc.name, func(t *testing.T) {
			stub := &scoreServiceStub{err: tc.err}
			app := fiber.New()
			NewClassHandler(stub, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(app)
			request := httptest.NewRequest("PUT", "/grade-items/"+uuid.NewString()+"/scores/"+uuid.NewString(), strings.NewReader(`{"numericValue":2,"isRemediated":true,"teacherComment":"Повторили тему"}`))
			request.Header.Set("Content-Type", "application/json")
			response, err := app.Test(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != tc.status {
				t.Fatalf("status = %d", response.StatusCode)
			}
			if !stub.remediated || stub.comment == nil {
				t.Fatal("remediation was not passed to service")
			}
			if tc.err == nil {
				var body gen.PutStudentScoreResponse
				if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if !body.Score.IsRemediated {
					t.Fatal("remediation missing in response")
				}
			} else {
				var body gen.ErrorResponse
				if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body.Code != "remediation_comment_required" {
					t.Fatalf("code=%s", body.Code)
				}
			}
		})
	}
}
func TestLessonResponsePreservesRemediation(t *testing.T) {
	id := uuid.New()
	lesson := lessonResponse(domain.LessonView{GradeItems: []entity.GradeItem{{ID: id}}, Scores: map[uuid.UUID][]entity.Score{id: {{IsRemediated: true}}}})
	if !lesson.GradeItems[0].Scores[0].IsRemediated {
		t.Fatal("journal response lost remediation")
	}
}
