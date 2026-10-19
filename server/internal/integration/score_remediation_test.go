//go:build integration

package integration

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

func TestScoreRemediationRoundTrip(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	teacher, student := user(t, tx, "teacher"), user(t, tx, "student")
	year, class, subject, assignment, lesson, item := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	statements := []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO academic_years(id,name,starts_on,ends_on,status)VALUES($1,$2,'2026-09-01','2027-05-31','active')`, []any{year, year.String()}},
		{`INSERT INTO classes(id,academic_year_id,name,grade_level)VALUES($1,$2,'7A',7)`, []any{class, year}},
		{`INSERT INTO subjects(id,code,name)VALUES($1,$2,'Math')`, []any{subject, subject.String()}},
		{`INSERT INTO class_subjects(id,class_id,subject_id,responsible_teacher_id)VALUES($1,$2,$3,$4)`, []any{assignment, class, subject, teacher}},
		{`INSERT INTO lessons(id,class_subject_id,lesson_date,position,topic,created_by)VALUES($1,$2,'2026-09-10',1,'Topic',$3)`, []any{lesson, assignment, teacher}},
		{`INSERT INTO grade_items(id,lesson_id,title,kind,grading_scale)VALUES($1,$2,'Work','classwork','five_point')`, []any{item, lesson}},
	}
	for _, s := range statements {
		if _, err := tx.Exec(ctx, s.sql, s.args...); err != nil {
			t.Fatal(err)
		}
	}
	repo := repository.NewRepository()
	value := 2.0
	score, err := repo.CreateScore(ctx, tx, entity.Score{GradeItemID: item, UserID: student, NumericValue: &value, CreatedBy: teacher, UpdatedBy: teacher})
	if err != nil {
		t.Fatal(err)
	}
	if score.IsRemediated {
		t.Fatal("old grades must default to pending")
	}
	explanation := "Повторили тему и решили задачи"
	score.IsRemediated, score.TeacherComment = true, &explanation
	if _, err = repo.UpdateScore(ctx, tx, score); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.GetScore(ctx, tx, score.ID)
	if err != nil || !loaded.IsRemediated || loaded.TeacherComment == nil || *loaded.TeacherComment != explanation {
		t.Fatalf("read: %+v %v", loaded, err)
	}
	listed, err := repo.ListScores(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, s := range listed {
		if s.ID == score.ID {
			found = s.IsRemediated
		}
	}
	if !found {
		t.Fatal("list lost remediation")
	}
	score.IsRemediated, score.TeacherComment = false, nil
	cleared, err := repo.UpsertScore(ctx, tx, score)
	if err != nil || cleared.IsRemediated {
		t.Fatalf("clear: %+v %v", cleared, err)
	}
	nested, err := tx.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	score.IsRemediated = true
	_, err = repo.UpsertScore(ctx, nested, score)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23514" {
		t.Fatalf("missing DB comment constraint: %v", err)
	}
	_ = nested.Rollback(ctx)
}
