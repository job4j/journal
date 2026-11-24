package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type GetLessonAttachmentRequest struct{ AttachmentID, StudentID uuid.UUID }
type GetLessonAttachmentResponse struct{ Attachment entity.LessonAttachment }

func (d *AttachmentDomain) GetLessonAttachment(ctx context.Context, tx repository.Transaction, hash string, req GetLessonAttachmentRequest) (GetLessonAttachmentResponse, error) {
	result := GetLessonAttachmentResponse{}
	// Reject expired sessions before looking up attachment metadata.
	if _, err := d.repo.FindActiveUserBySessionHash(ctx, tx, hash); errors.Is(err, repository.ErrNotFound) {
		return result, ErrUnauthenticated
	} else if err != nil {
		return result, err
	}
	item, err := d.repo.GetLessonAttachment(ctx, tx, req.AttachmentID)
	if errors.Is(err, repository.ErrNotFound) || (err == nil && item.DeletedAt != nil) {
		return result, ErrAttachmentNotFound
	}
	if err != nil {
		return result, err
	}
	if _, err = d.authorize(ctx, tx, hash, item.LessonID, req.StudentID, false); err != nil {
		return result, err
	}
	result.Attachment = item
	return result, nil
}
