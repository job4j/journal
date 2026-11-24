package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type ListLessonAttachmentsRequest struct{ LessonID, StudentID uuid.UUID }
type ListLessonAttachmentsResponse struct{ Items []entity.LessonAttachment }

func (d *AttachmentDomain) ListLessonAttachments(ctx context.Context, tx repository.Transaction, hash string, req ListLessonAttachmentsRequest) (ListLessonAttachmentsResponse, error) {
	result := ListLessonAttachmentsResponse{Items: []entity.LessonAttachment{}}
	if _, err := d.authorize(ctx, tx, hash, req.LessonID, req.StudentID, false); err != nil {
		return result, err
	}
	items, err := d.repo.ListLessonAttachments(ctx, tx, req.LessonID)
	if items != nil {
		result.Items = items
	}
	return result, err
}
