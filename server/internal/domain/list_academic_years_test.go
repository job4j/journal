package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"testing"
	"time"
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

func TestCreateAcademicYearQuarterStoresTrimmedName(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	repo := academicYearRepoStub{
		authenticated: true,
		allowed:       true,
		years: []entity.AcademicYear{{
			ID: uuid.New(), StartsOn: start, EndsOn: start.AddDate(1, 0, 0),
		}},
	}

	result, err := NewAcademicYearDomain(repo).CreateAcademicYearQuarter(
		context.Background(),
		testTx{},
		"hash",
		repo.years[0].ID.String(),
		" Осенний период ",
		start,
		start.AddDate(0, 1, 0),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Осенний период" {
		t.Fatalf("name = %q", result.Name)
	}
}

func TestCreateAcademicYearQuarterRejectsEmptyName(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	repo := academicYearRepoStub{
		authenticated: true,
		allowed:       true,
		years: []entity.AcademicYear{{
			ID: uuid.New(), StartsOn: start, EndsOn: start.AddDate(1, 0, 0),
		}},
	}

	_, err := NewAcademicYearDomain(repo).CreateAcademicYearQuarter(
		context.Background(),
		testTx{},
		"hash",
		repo.years[0].ID.String(),
		"   ",
		start,
		start.AddDate(0, 1, 0),
	)
	if !errors.Is(err, ErrInvalidAcademicYear) {
		t.Fatalf("error = %v", err)
	}
}
