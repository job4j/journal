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

func TestUpdateAcademicYearQuarterValidatesAndTrims(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	yearID, quarterID := uuid.New(), uuid.New()
	repo := academicYearRepoStub{authenticated: true, allowed: true, years: []entity.AcademicYear{{ID: yearID, StartsOn: start, EndsOn: start.AddDate(1, 0, -1)}}, quarters: []entity.AcademicYearQuarter{{ID: quarterID, AcademicYearID: yearID, Number: 2, Name: "Осень", StartsOn: start, EndsOn: start.AddDate(0, 1, 0)}}}
	result, err := NewAcademicYearDomain(repo).UpdateAcademicYearQuarter(context.Background(), testTx{}, "hash", quarterID, " Первый период ", start, start.AddDate(0, 2, 0))
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Первый период" || result.Number != 2 || result.AcademicYearID != yearID {
		t.Fatalf("result = %+v", result)
	}
}

func TestUpdateAcademicYearQuarterRejectsOverlap(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	yearID, quarterID := uuid.New(), uuid.New()
	repo := academicYearRepoStub{authenticated: true, allowed: true, years: []entity.AcademicYear{{ID: yearID, StartsOn: start, EndsOn: start.AddDate(1, 0, -1)}}, quarters: []entity.AcademicYearQuarter{{ID: quarterID, AcademicYearID: yearID, Number: 1}, {ID: uuid.New(), AcademicYearID: yearID, Number: 2, StartsOn: start.AddDate(0, 2, 0), EndsOn: start.AddDate(0, 3, 0)}}}
	_, err := NewAcademicYearDomain(repo).UpdateAcademicYearQuarter(context.Background(), testTx{}, "hash", quarterID, "Осень", start, start.AddDate(0, 2, 0))
	if !errors.Is(err, ErrAcademicYearExists) {
		t.Fatalf("error = %v", err)
	}
}

type quarterDeleteRepo struct {
	academicYearRepoStub
	deleteErr error
}

func (s quarterDeleteRepo) DeleteAcademicYearQuarter(context.Context, repository.Transaction, uuid.UUID) error {
	return s.deleteErr
}

func TestDeleteAcademicYearQuarterProtectsGrades(t *testing.T) {
	id := uuid.New()
	repo := quarterDeleteRepo{academicYearRepoStub: academicYearRepoStub{authenticated: true, allowed: true, quarters: []entity.AcademicYearQuarter{{ID: id}}}, deleteErr: repository.ErrReference}
	err := NewAcademicYearDomain(repo).DeleteAcademicYearQuarter(context.Background(), testTx{}, "hash", id)
	if !errors.Is(err, ErrQuarterInUse) {
		t.Fatalf("error = %v", err)
	}
}
