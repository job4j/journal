package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
)

func (d *ClassDomain) DeleteQuarterGrade(ctx context.Context, tx repository.Transaction, hash string, assignmentID, quarterID, studentID uuid.UUID) error {
	if _, _, _, err := d.quarterGradeContext(ctx, tx, hash, assignmentID, quarterID, "can_create_score", true); err != nil {
		return err
	}
	err := d.repo.DeleteQuarterGrade(ctx, tx, quarterID, assignmentID, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrScoreNotFound
	}
	if err != nil {
		return fmt.Errorf("delete quarter grade: %w", err)
	}
	return nil
}
