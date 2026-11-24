package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type DeleteLessonAttachmentRequest struct{ AttachmentID uuid.UUID }
type DeleteLessonAttachmentResponse struct{ Attachment entity.LessonAttachment }

func (d *AttachmentDomain) DeleteLessonAttachment(ctx context.Context, tx repository.Transaction, hash string, req DeleteLessonAttachmentRequest) (DeleteLessonAttachmentResponse, error) {
	result := DeleteLessonAttachmentResponse{}
	if _, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash); errors.Is(err, repository.ErrNotFound) {
		return result, ErrUnauthenticated
	} else if err != nil {
		return result, err
	}
	item, err := d.repo.GetLessonAttachment(ctx, tx, req.AttachmentID)
	if errors.Is(err, repository.ErrNotFound) {
		return result, ErrAttachmentNotFound
	}
	if err != nil {
		return result, err
	}
	if _, err = d.authorize(ctx, tx, hash, item.LessonID, uuid.Nil, true); err != nil {
		return result, err
	}
	result.Attachment = item
	err = d.repo.DeleteLessonAttachment(ctx, tx, item.ID)
	return result, err
}
