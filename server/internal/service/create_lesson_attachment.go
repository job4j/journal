package service

import (
	"context"
	"fmt"
	"io"
	"journal/server/internal/domain"
	"time"
)

func (s *AttachmentService) CreateLessonAttachment(ctx context.Context, token string, req domain.CreateLessonAttachmentRequest, content io.Reader) (result domain.CreateLessonAttachmentResponse, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.CreateLessonAttachment(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return result, err
	}
	if s.store == nil {
		return result, domain.ErrAttachmentStorageUnavailable
	}
	item := result.Attachment
	// Remove a partial/uncommitted object even when the upload request is cancelled.
	committed := false
	defer func() {
		if !committed {
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
			defer cancel()
			if cleanupErr := s.store.Delete(cleanup, item.ObjectKey); cleanupErr != nil {
				err = fmt.Errorf("%w; cleanup upload: %v", err, cleanupErr)
			}
		}
	}()
	if err = s.store.Put(ctx, item.ObjectKey, item.ContentType, content, item.SizeBytes); err != nil {
		return result, fmt.Errorf("%w: %v", domain.ErrAttachmentStorageUnavailable, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit attachment: %w", err)
	}
	committed = true
	return result, nil
}
