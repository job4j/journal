package domain

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"path"
	"strings"
	"unicode"
)

type CreateLessonAttachmentRequest struct {
	LessonID              uuid.UUID
	FileName, ContentType string
	SizeBytes             int64
}
type CreateLessonAttachmentResponse struct{ Attachment entity.LessonAttachment }

func (d *AttachmentDomain) CreateLessonAttachment(ctx context.Context, tx repository.Transaction, hash string, req CreateLessonAttachmentRequest) (CreateLessonAttachmentResponse, error) {
	result := CreateLessonAttachmentResponse{}
	user, err := d.authorize(ctx, tx, hash, req.LessonID, uuid.Nil, true)
	if err != nil {
		return result, err
	}
	name := strings.TrimSpace(path.Base(strings.ReplaceAll(req.FileName, "\\", "/")))
	if name == "" || name == "." || name == ".." || len(name) > 255 || strings.ContainsFunc(name, unicode.IsControl) || req.SizeBytes <= 0 || req.SizeBytes > MaxAttachmentBytes || req.ContentType == "" || len(req.ContentType) > 255 {
		return result, ErrInvalidAttachment
	}
	id := uuid.New()
	result.Attachment, err = d.repo.CreateLessonAttachment(ctx, tx, entity.LessonAttachment{ID: id, LessonID: req.LessonID, ObjectKey: "lessons/" + req.LessonID.String() + "/" + id.String(), FileName: name, ContentType: req.ContentType, SizeBytes: req.SizeBytes, CreatedBy: user.ID})
	return result, err
}
