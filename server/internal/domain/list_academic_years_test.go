package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
)

type academicYearRepoStub struct {
	authenticated, allowed bool
	years                  []entity.AcademicYear
	quarters               []entity.AcademicYearQuarter
}

func (s academicYearRepoStub) CheckSessionPermission(context.Context, repository.Transaction, string, string) (bool, bool, error) {
	return s.authenticated, s.allowed, nil
}
func (s academicYearRepoStub) CreateAcademicYear(_ context.Context, _ repository.Transaction, value entity.AcademicYear) (entity.AcademicYear, error) {
	value.ID = uuid.New()
	return value, nil
}
func (s academicYearRepoStub) CreateAcademicYearQuarter(_ context.Context, _ repository.Transaction, value entity.AcademicYearQuarter) (entity.AcademicYearQuarter, error) {
	value.ID = uuid.New()
	return value, nil
}
func (s academicYearRepoStub) ListAcademicYears(context.Context, repository.Transaction) ([]entity.AcademicYear, error) {
	return s.years, nil
}
func (s academicYearRepoStub) ListAcademicYearQuarters(context.Context, repository.Transaction) ([]entity.AcademicYearQuarter, error) {
	return s.quarters, nil
}

func TestListAcademicYearsGroupsQuarters(t *testing.T) {
	yearID, otherID := uuid.New(), uuid.New()
	repo := academicYearRepoStub{authenticated: true, allowed: true, years: []entity.AcademicYear{{ID: yearID, Name: "2026/2027"}}, quarters: []entity.AcademicYearQuarter{{AcademicYearID: otherID, Number: 1}, {AcademicYearID: yearID, Number: 1}, {AcademicYearID: yearID, Number: 2}}}
	result, err := NewAcademicYearDomain(repo).ListAcademicYears(context.Background(), testTx{}, "hash")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || len(result[0].Quarters) != 2 {
		t.Fatalf("result = %+v", result)
	}
}

func TestListAcademicYearsChecksAccess(t *testing.T) {
	_, err := NewAcademicYearDomain(academicYearRepoStub{}).ListAcademicYears(context.Background(), testTx{}, "hash")
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("error = %v", err)
	}
}
