package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestDeleteAttachmentCanRetryButRechecksTeacher(t *testing.T) {
	r := attachmentFixture()
	now := time.Now()
	r.attachment.DeletedAt = &now
	d := NewAttachmentDomain(r)
	_, err := d.DeleteLessonAttachment(context.Background(), testTx{}, "hash", DeleteLessonAttachmentRequest{AttachmentID: r.attachment.ID})
	if err != nil || !r.deleted {
		t.Fatalf("retry failed: %v", err)
	}
	r.deleted = false
	r.assignment.ResponsibleTeacherID = uuid.New()
	_, err = d.DeleteLessonAttachment(context.Background(), testTx{}, "hash", DeleteLessonAttachmentRequest{AttachmentID: r.attachment.ID})
	if !errors.Is(err, ErrForbidden) || r.deleted {
		t.Fatalf("previous teacher retained access: %v", err)
	}
}
