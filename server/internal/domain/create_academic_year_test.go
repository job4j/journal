package domain

import (
	"context"
	"testing"
	"time"
)

func validAcademicYearRequest() CreateAcademicYearRequest {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	return CreateAcademicYearRequest{SessionTokenHash: "hash", Name: "2026/2027", Status: "planned", StartsOn: start, EndsOn: start.AddDate(1, 0, -1), Quarters: []AcademicYearQuarterInput{
		{1, start, start.AddDate(0, 1, 0)}, {2, start.AddDate(0, 1, 1), start.AddDate(0, 3, 0)}, {3, start.AddDate(0, 3, 1), start.AddDate(0, 5, 0)}, {4, start.AddDate(0, 5, 1), start.AddDate(0, 8, 0)},
	}}
}

func TestCreateAcademicYearWithoutPeriods(t *testing.T) {
	result, err := NewAcademicYearDomain(academicYearRepoStub{authenticated: true, allowed: true}).CreateAcademicYear(context.Background(), testTx{}, validAcademicYearRequest())
	if err != nil {
		t.Fatal(err)
	}
	if result.Year.ID.String() == "00000000-0000-0000-0000-000000000000" || len(result.Quarters) != 0 {
		t.Fatalf("result = %+v", result)
	}
}
