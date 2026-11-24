package service

import (
	"context"
	"fmt"
	"journal/server/internal/domain"
)

func (s *AttachmentService) DeleteLessonAttachment(ctx context.Context, token string, req domain.DeleteLessonAttachmentRequest) error {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := s.domain.DeleteLessonAttachment(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return err
	}
	if s.store == nil {
		return domain.ErrAttachmentStorageUnavailable
	}
	// Revoke downloads before removing bytes. The tombstone permits a safe retry.
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	if err = s.store.Delete(ctx, result.Attachment.ObjectKey); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrAttachmentStorageUnavailable, err)
	}
	return nil
}
