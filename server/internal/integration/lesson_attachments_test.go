//go:build integration

package integration

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
)

func TestLessonAttachmentsMetadataAndParentAccess(t *testing.T) {
	ctx := context.Background()
	tx := transaction(t)
	repo := repository.NewRepository()
	teacher, parent, child := user(t, tx, "teacher"), user(t, tx, "parent"), user(t, tx, "student")
	yearID, classID, subjectID, assignmentID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	for _, query := range []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO academic_years(id,name,starts_on,ends_on,status) VALUES($1,$2,'2026-09-01','2027-05-31','planned')`, []any{yearID, yearID.String()}},
		{`INSERT INTO classes(id,academic_year_id,name,grade_level) VALUES($1,$2,'7А',7)`, []any{classID, yearID}},
		{`INSERT INTO subjects(id,code,name) VALUES($1,$2,'Предмет')`, []any{subjectID, subjectID.String()}},
		{`INSERT INTO class_subjects(id,class_id,subject_id,responsible_teacher_id) VALUES($1,$2,$3,$4)`, []any{assignmentID, classID, subjectID, teacher}},
		{`INSERT INTO class_students(class_id,user_id,enrolled_on) VALUES($1,$2,'2026-09-01')`, []any{classID, child}},
	} {
		if _, err := tx.Exec(ctx, query.sql, query.args...); err != nil {
			t.Fatal(err)
		}
	}
	for _, code := range []string{"can_view_lesson", "can_update_lesson"} {
		permission(t, tx, teacher, code, assignmentID)
	}
	for _, code := range []string{"can_view_user", "can_view_journal"} {
		permission(t, tx, parent, code, child)
	}
	lesson, err := repo.CreateLesson(ctx, tx, entity.Lesson{ClassSubjectID: assignmentID, LessonDate: time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC), Position: 1, Topic: "Тема", CreatedBy: teacher})
	if err != nil {
		t.Fatal(err)
	}
	d := domain.NewAttachmentDomain(repo)
	teacherHash, parentHash := session(t, tx, teacher), session(t, tx, parent)
	saved, err := d.CreateLessonAttachment(ctx, tx, teacherHash, domain.CreateLessonAttachmentRequest{LessonID: lesson.ID, FileName: "Задание.pdf", ContentType: "application/pdf", SizeBytes: 123})
	if err != nil {
		t.Fatal(err)
	}
	list, err := d.ListLessonAttachments(ctx, tx, parentHash, domain.ListLessonAttachmentsRequest{LessonID: lesson.ID, StudentID: child})
	if err != nil || len(list.Items) != 1 || list.Items[0].ID != saved.Attachment.ID {
		t.Fatalf("list: %+v %v", list, err)
	}
	if _, err = d.GetLessonAttachment(ctx, tx, parentHash, domain.GetLessonAttachmentRequest{AttachmentID: saved.Attachment.ID, StudentID: uuid.New()}); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("another child: %v", err)
	}
	_, err = d.DeleteLessonAttachment(ctx, tx, teacherHash, domain.DeleteLessonAttachmentRequest{AttachmentID: saved.Attachment.ID})
	if err != nil {
		t.Fatal(err)
	}
	list, err = d.ListLessonAttachments(ctx, tx, parentHash, domain.ListLessonAttachmentsRequest{LessonID: lesson.ID, StudentID: child})
	if err != nil || len(list.Items) != 0 {
		t.Fatalf("deleted file listed: %+v %v", list, err)
	}
	if _, err = d.GetLessonAttachment(ctx, tx, parentHash, domain.GetLessonAttachmentRequest{AttachmentID: saved.Attachment.ID, StudentID: child}); !errors.Is(err, domain.ErrAttachmentNotFound) {
		t.Fatalf("deleted file downloadable: %v", err)
	}
	_, err = d.DeleteLessonAttachment(ctx, tx, teacherHash, domain.DeleteLessonAttachmentRequest{AttachmentID: saved.Attachment.ID})
	if err != nil {
		t.Fatalf("deletion retry: %v", err)
	}
}
