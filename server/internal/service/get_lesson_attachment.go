package service

import (
	"context"
	"fmt"
	"io"
	"journal/server/internal/domain"
)

func (s *AttachmentService) GetLessonAttachment(ctx context.Context, token string, req domain.GetLessonAttachmentRequest) (result domain.GetLessonAttachmentResponse, content io.ReadCloser, err error) {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return result, nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err = s.domain.GetLessonAttachment(ctx, tx, sessionTokenHash(token), req)
	if err != nil {
		return result, nil, err
	}
	if s.store == nil {
		return result, nil, domain.ErrAttachmentStorageUnavailable
	}
	if err = tx.Commit(ctx); err != nil {
		return result, nil, err
	}
	content, err = s.store.Open(ctx, result.Attachment.ObjectKey)
	if err != nil {
		return result, nil, fmt.Errorf("%w: %v", domain.ErrAttachmentStorageUnavailable, err)
	}
	return result, content, nil
}
