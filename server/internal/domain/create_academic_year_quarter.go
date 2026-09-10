package domain

import (
	"context"
	"errors"
	"fmt"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
	"time"
)

func (d *AcademicYearDomain) CreateAcademicYearQuarter(ctx context.Context, tx repository.Transaction, hash string, yearID, name string, startsOn, endsOn time.Time) (entity.AcademicYearQuarter, error) {
	if err := Authorize(ctx, tx, d.repo, hash, "can_manage_academic_year"); err != nil {
		return entity.AcademicYearQuarter{}, err
	}
	years, err := d.repo.ListAcademicYears(ctx, tx)
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("list academic years: %w", err)
	}
	var year entity.AcademicYear
	found := false
	for _, item := range years {
		if item.ID.String() == yearID {
			year = item
			found = true
			break
		}
	}
	if !found {
		return entity.AcademicYearQuarter{}, ErrAcademicYearNotFound
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return entity.AcademicYearQuarter{}, ErrInvalidAcademicYear
	}
	if endsOn.Before(startsOn) || startsOn.Before(year.StartsOn) || endsOn.After(year.EndsOn) {
		return entity.AcademicYearQuarter{}, ErrInvalidAcademicYear
	}
	periods, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("list periods: %w", err)
	}
	number := int16(1)
	for _, item := range periods {
		if item.AcademicYearID != year.ID {
			continue
		}
		if !endsOn.Before(item.StartsOn) && !startsOn.After(item.EndsOn) {
			return entity.AcademicYearQuarter{}, ErrAcademicYearExists
		}
		if item.Number >= number {
			number = item.Number + 1
		}
	}
	result, err := d.repo.CreateAcademicYearQuarter(ctx, tx, entity.AcademicYearQuarter{AcademicYearID: year.ID, Number: number, Name: name, StartsOn: startsOn, EndsOn: endsOn})
	if errors.Is(err, repository.ErrConflict) {
		return entity.AcademicYearQuarter{}, ErrAcademicYearExists
	}
	if err != nil {
		return entity.AcademicYearQuarter{}, fmt.Errorf("create period: %w", err)
	}
	return result, nil
}
