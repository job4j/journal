package domain

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGetAttachmentRejectsRemovedAndRevokedAccess(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*attachmentRepoStub)
		want   error
	}{
		{"allowed", func(*attachmentRepoStub) {}, nil},
		{"revoked", func(r *attachmentRepoStub) { r.denied = "can_view_lesson" }, ErrForbidden},
		{"deleted", func(r *attachmentRepoStub) { now := time.Now(); r.attachment.DeletedAt = &now }, ErrAttachmentNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := attachmentFixture()
			test.change(r)
			_, err := NewAttachmentDomain(r).GetLessonAttachment(context.Background(), testTx{}, "hash", GetLessonAttachmentRequest{AttachmentID: r.attachment.ID})
			if !errors.Is(err, test.want) {
				t.Fatalf("got %v want %v", err, test.want)
			}
		})
	}
}
