package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/domain"
	"testing"
)

func TestAttachmentDeletionCommitsTombstoneBeforeRemovingObject(t *testing.T) {
	tx, repo, store := attachmentServiceFixture()
	store.deleteErr = errors.New("offline")
	err := NewAttachmentService(attachmentManager{tx}, repo, store).DeleteLessonAttachment(context.Background(), "token", domain.DeleteLessonAttachmentRequest{AttachmentID: uuid.New()})
	if !errors.Is(err, domain.ErrAttachmentStorageUnavailable) || !store.deletedAfterCommit {
		t.Fatalf("unsafe deletion: %v %+v", err, store)
	}
}
func TestAttachmentDeletionDoesNotRemoveObjectOnCommitFailure(t *testing.T) {
	tx, repo, store := attachmentServiceFixture()
	tx.commitErr = errors.New("commit failed")
	err := NewAttachmentService(attachmentManager{tx}, repo, store).DeleteLessonAttachment(context.Background(), "token", domain.DeleteLessonAttachmentRequest{AttachmentID: uuid.New()})
	if err == nil || store.removed {
		t.Fatalf("removed object before commit: %v", err)
	}
}
