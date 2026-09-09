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
	if request.Name == "" || !request.EndsOn.After(request.StartsOn) || (request.Status != "planned" && request.Status != "active" && request.Status != "completed") {
		return AcademicYearView{}, ErrInvalidAcademicYear
	}
	year, err := d.repo.CreateAcademicYear(ctx, tx, entity.AcademicYear{Name: request.Name, StartsOn: request.StartsOn, EndsOn: request.EndsOn, Status: request.Status})
	if errors.Is(err, repository.ErrConflict) {
		return AcademicYearView{}, ErrAcademicYearExists
	}
	if err != nil {
		return AcademicYearView{}, fmt.Errorf("create academic year: %w", err)
	}
	return AcademicYearView{Year: year, Quarters: []entity.AcademicYearQuarter{}}, nil
}
