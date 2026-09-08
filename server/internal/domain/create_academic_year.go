package domain

import (
	"context"
	"errors"
	"fmt"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"sort"
	"strings"
	"time"
)

type AcademicYearQuarterInput struct {
	Number           int16
	StartsOn, EndsOn time.Time
}
type CreateAcademicYearRequest struct {
	SessionTokenHash string
	Name, Status     string
	StartsOn, EndsOn time.Time
	Quarters         []AcademicYearQuarterInput
}

func (d *AcademicYearDomain) CreateAcademicYear(ctx context.Context, tx repository.Transaction, request CreateAcademicYearRequest) (AcademicYearView, error) {
	if err := Authorize(ctx, tx, d.repo, request.SessionTokenHash, "can_manage_academic_year"); err != nil {
		return AcademicYearView{}, err
	}
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || !request.EndsOn.After(request.StartsOn) || (request.Status != "planned" && request.Status != "active" && request.Status != "completed") || len(request.Quarters) != 4 {
		return AcademicYearView{}, ErrInvalidAcademicYear
	}
	quarters := append([]AcademicYearQuarterInput(nil), request.Quarters...)
	sort.Slice(quarters, func(i, j int) bool { return quarters[i].Number < quarters[j].Number })
	for i, q := range quarters {
		if q.Number != int16(i+1) || q.StartsOn.Before(request.StartsOn) || q.EndsOn.After(request.EndsOn) || q.EndsOn.Before(q.StartsOn) {
			return AcademicYearView{}, ErrInvalidAcademicYear
		}
		if i > 0 && !q.StartsOn.After(quarters[i-1].EndsOn) {
			return AcademicYearView{}, ErrInvalidAcademicYear
		}
	}
	year, err := d.repo.CreateAcademicYear(ctx, tx, entity.AcademicYear{Name: request.Name, StartsOn: request.StartsOn, EndsOn: request.EndsOn, Status: request.Status})
	if errors.Is(err, repository.ErrConflict) {
		return AcademicYearView{}, ErrAcademicYearExists
	}
	if err != nil {
		return AcademicYearView{}, fmt.Errorf("create academic year: %w", err)
	}
	created := make([]entity.AcademicYearQuarter, 0, 4)
	for _, q := range quarters {
		quarter, err := d.repo.CreateAcademicYearQuarter(ctx, tx, entity.AcademicYearQuarter{AcademicYearID: year.ID, Number: q.Number, StartsOn: q.StartsOn, EndsOn: q.EndsOn})
		if err != nil {
			return AcademicYearView{}, fmt.Errorf("create academic year quarter: %w", err)
		}
		created = append(created, quarter)
	}
	return AcademicYearView{Year: year, Quarters: created}, nil
}
