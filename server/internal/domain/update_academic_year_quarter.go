package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
	"time"
)

func (d *AcademicYearDomain) UpdateAcademicYearQuarter(ctx context.Context, tx repository.Transaction, hash string, id uuid.UUID, name string, startsOn, endsOn time.Time) (entity.AcademicYearQuarter, error) {
	if err := Authorize(ctx, tx, d.repo, hash, "can_manage_academic_year"); err != nil {
		return entity.AcademicYearQuarter{}, err
	}
	current, err := d.repo.GetAcademicYearQuarter(ctx, tx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.AcademicYearQuarter{}, ErrQuarterNotFound
	}
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("get period: %w", err)
	}
	years, err := d.repo.ListAcademicYears(ctx, tx)
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("list academic years: %w", err)
	}
	var year entity.AcademicYear
	found := false
	for _, item := range years {
		if item.ID == current.AcademicYearID {
			year, found = item, true
			break
		}
	}
	name = strings.TrimSpace(name)
	if !found || name == "" || endsOn.Before(startsOn) || startsOn.Before(year.StartsOn) || endsOn.After(year.EndsOn) {
		return entity.AcademicYearQuarter{}, ErrInvalidAcademicYear
	}
	periods, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("list periods: %w", err)
	}
	for _, item := range periods {
		if item.ID == current.ID || item.AcademicYearID != current.AcademicYearID {
			continue
		}
		if !endsOn.Before(item.StartsOn) && !startsOn.After(item.EndsOn) {
			return entity.AcademicYearQuarter{}, ErrAcademicYearExists
		}
	}
	current.Name, current.StartsOn, current.EndsOn = name, startsOn, endsOn
	result, err := d.repo.UpdateAcademicYearQuarter(ctx, tx, current)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.AcademicYearQuarter{}, ErrQuarterNotFound
	}
	if errors.Is(err, repository.ErrConflict) {
		return entity.AcademicYearQuarter{}, ErrAcademicYearExists
	}
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("update period: %w", err)
	}
	return result, nil
}
