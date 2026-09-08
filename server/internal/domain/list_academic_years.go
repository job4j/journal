package domain

import (
	"context"
	"fmt"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type AcademicYearView struct {
	Year     entity.AcademicYear
	Quarters []entity.AcademicYearQuarter
}

func (d *AcademicYearDomain) ListAcademicYears(ctx context.Context, tx repository.Transaction, tokenHash string) ([]AcademicYearView, error) {
	if err := Authorize(ctx, tx, d.repo, tokenHash, "can_manage_academic_year"); err != nil {
		return nil, err
	}
	years, err := d.repo.ListAcademicYears(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list academic years: %w", err)
	}
	quarters, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list academic year quarters: %w", err)
	}
	result := make([]AcademicYearView, len(years))
	for i, year := range years {
		result[i] = AcademicYearView{Year: year, Quarters: []entity.AcademicYearQuarter{}}
		for _, quarter := range quarters {
			if quarter.AcademicYearID == year.ID {
				result[i].Quarters = append(result[i].Quarters, quarter)
			}
		}
	}
	return result, nil
}
