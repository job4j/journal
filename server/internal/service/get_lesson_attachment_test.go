package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"testing"
)

func TestAttachmentDownloadDoesNotReadStorageWithoutPermission(t *testing.T) {
	tx, repo, store := attachmentServiceFixture()
	repo.denied = true
	_, _, err := NewAttachmentService(attachmentManager{tx}, repo, store).GetLessonAttachment(context.Background(), "token", domain.GetLessonAttachmentRequest{AttachmentID: uuid.New()})
	if !errors.Is(err, domain.ErrForbidden) || store.opened {
		t.Fatalf("unauthorized storage read: %v", err)
	}
}
