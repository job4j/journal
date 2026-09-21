package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
)

func (d *AcademicYearDomain) DeleteAcademicYearQuarter(ctx context.Context, tx repository.Transaction, hash string, id uuid.UUID) error {
	if err := Authorize(ctx, tx, d.repo, hash, "can_manage_academic_year"); err != nil {
		return err
	}
	if _, err := d.repo.GetAcademicYearQuarter(ctx, tx, id); errors.Is(err, repository.ErrNotFound) {
		return ErrQuarterNotFound
	} else if err != nil {
		return fmt.Errorf("get period: %w", err)
	}
	err := d.repo.DeleteAcademicYearQuarter(ctx, tx, id)
	if errors.Is(err, repository.ErrReference) {
		return ErrQuarterInUse
	}
	if errors.Is(err, repository.ErrNotFound) {
		return ErrQuarterNotFound
	}
	if err != nil {
		return fmt.Errorf("delete period: %w", err)
	}
	return nil
}
