//go:build integration

package integration

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
)

const defaultDatabaseURL = "postgres://postgres:password@127.0.0.1:5433/journal?sslmode=disable"

func transaction(t *testing.T) pgx.Tx {
	t.Helper()
	url := os.Getenv("INTEGRATION_DATABASE_URL")
	if url == "" {
		url = defaultDatabaseURL
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Skipf("postgres is unavailable: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	return tx
}
func user(t *testing.T, tx pgx.Tx, roles ...string) uuid.UUID {
	t.Helper()
	id := uuid.New()
	_, err := tx.Exec(context.Background(), `INSERT INTO users(id,login,email,password_hash,name,status)VALUES($1,$2,$2,'hash','Test User','active')`, id, id.String()+"@test.local")
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range roles {
		_, err = tx.Exec(context.Background(), `INSERT INTO user_roles(user_id,role_id)SELECT $1,id FROM roles WHERE code=$2`, id, code)
		if err != nil {
			t.Fatal(err)
		}
	}
	return id
}
func session(t *testing.T, tx pgx.Tx, userID uuid.UUID) string {
	t.Helper()
	hash := uuid.NewString()
	_, err := tx.Exec(context.Background(), `INSERT INTO sessions(user_id,token_hash,expires_at)VALUES($1,$2,now()+interval '1 hour')`, userID, hash)
	if err != nil {
		t.Fatal(err)
	}
	return hash
}
func permission(t *testing.T, tx pgx.Tx, userID uuid.UUID, code string, value uuid.UUID) {
	t.Helper()
	permissionID := uuid.Nil
	err := tx.QueryRow(context.Background(), `SELECT id FROM permissions WHERE code=$1 AND value=$2`, code, value.String()).Scan(&permissionID)
	if errors.Is(err, pgx.ErrNoRows) {
		permissionID = uuid.New()
		_, err = tx.Exec(context.Background(), `INSERT INTO permissions(id,code,value,description)VALUES($1,$2,$3,$2)`, permissionID, code, value.String())
	}
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO user_permissions(user_id,permission_id)VALUES($1,$2)`, userID, permissionID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParentsCannotReadAnotherFamily(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	firstParent, secondParent := user(t, tx, "parent"), user(t, tx, "parent")
	firstChild, secondChild := user(t, tx, "student"), user(t, tx, "student")
	firstHash, secondHash := session(t, tx, firstParent), session(t, tx, secondParent)
	for _, code := range []string{"can_view_user", "can_view_journal"} {
		permission(t, tx, firstParent, code, firstChild)
		permission(t, tx, secondParent, code, secondChild)
	}
	repo := repository.NewRepository()
	for _, check := range []struct {
		hash, code string
		value      uuid.UUID
		want       bool
	}{{firstHash, "can_view_user", firstChild, true}, {firstHash, "can_view_journal", firstChild, true}, {firstHash, "can_view_user", secondChild, false}, {secondHash, "can_view_journal", firstChild, false}} {
		authenticated, allowed, err := repo.CheckSessionPermissionForValue(ctx, tx, check.hash, check.code, check.value.String())
		if err != nil || !authenticated || allowed != check.want {
			t.Fatalf("check %+v: authenticated=%v allowed=%v err=%v", check, authenticated, allowed, err)
		}
	}
}

func TestTeacherAssignmentAndFormerStudentBoundaries(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	teacher1, teacher2, studentID := user(t, tx, "teacher"), user(t, tx, "teacher"), user(t, tx, "student")
	hash1, hash2 := session(t, tx, teacher1), session(t, tx, teacher2)
	yearID, classID, subjectID, assignmentID, lessonID, itemID := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	var err error
	statements := []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO academic_years(id,name,starts_on,ends_on,status)VALUES($1,'2026',DATE '2026-09-01',DATE '2027-05-31','active')`, []any{yearID}},
		{`INSERT INTO classes(id,academic_year_id,name,grade_level)VALUES($1,$2,'7A',7)`, []any{classID, yearID}},
		{`INSERT INTO subjects(id,code,name)VALUES($1,$2,'Math')`, []any{subjectID, "math_" + subjectID.String()[:8]}},
		{`INSERT INTO class_subjects(id,class_id,subject_id,responsible_teacher_id)VALUES($1,$2,$3,$4)`, []any{assignmentID, classID, subjectID, teacher1}},
		{`INSERT INTO class_students(class_id,user_id,enrolled_on)VALUES($1,$2,DATE '2026-09-01')`, []any{classID, studentID}},
		{`INSERT INTO lessons(id,class_subject_id,lesson_date,position,topic,created_by)VALUES($1,$2,DATE '2026-09-10',1,'Topic',$3)`, []any{lessonID, assignmentID, teacher1}},
		{`INSERT INTO grade_items(id,lesson_id,title,kind,grading_scale)VALUES($1,$2,'Work','classwork','five_point')`, []any{itemID, lessonID}},
	}
	for _, statement := range statements {
		if _, err := tx.Exec(ctx, statement.sql, statement.args...); err != nil {
			t.Fatal(err)
		}
	}
	for _, teacherID := range []uuid.UUID{teacher1, teacher2} {
		permission(t, tx, teacherID, "can_create_score", assignmentID)
	}
	repo := repository.NewRepository()
	journal := domain.NewClassDomain(repo)
	value := 5.0
	if _, err = journal.PutStudentScore(ctx, tx, hash2, itemID, studentID, &value, nil, nil); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("other teacher: %v", err)
	}
	left := time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	_, err = tx.Exec(ctx, `UPDATE class_students SET left_on=$1 WHERE class_id=$2 AND user_id=$3`, left, classID, studentID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = journal.PutStudentScore(ctx, tx, hash1, itemID, studentID, &value, nil, nil); !errors.Is(err, domain.ErrClassStudentNotFound) {
		t.Fatalf("former student: %v", err)
	}
	if _, err = tx.Exec(ctx, `UPDATE class_students SET left_on=NULL WHERE class_id=$1 AND user_id=$2`, classID, studentID); err != nil {
		t.Fatal(err)
	}
	if _, err = tx.Exec(ctx, `UPDATE class_subjects SET responsible_teacher_id=$1 WHERE id=$2`, teacher2, assignmentID); err != nil {
		t.Fatal(err)
	}
	if _, err = journal.PutStudentScore(ctx, tx, hash1, itemID, studentID, &value, nil, nil); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("replaced teacher: %v", err)
	}
	if _, err = journal.PutStudentScore(ctx, tx, hash2, itemID, studentID, &value, nil, nil); err != nil {
		t.Fatalf("new teacher: %v", err)
	}
}
