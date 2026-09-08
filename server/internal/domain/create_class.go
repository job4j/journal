package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
	"strings"
)

func (d *ClassDomain) CreateClass(ctx context.Context, tx repository.Transaction, hash string, yearID uuid.UUID, name string, gradeLevel int16) (ClassView, error) {
	if err := Authorize(ctx, tx, d.repo, hash, "can_create_class"); err != nil {
		return ClassView{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 100 || gradeLevel < 1 || gradeLevel > 11 {
		return ClassView{}, ErrInvalidClass
	}
	if _, err := d.repo.GetAcademicYear(ctx, tx, yearID); errors.Is(err, repository.ErrNotFound) {
		return ClassView{}, ErrAcademicYearNotFound
	} else if err != nil {
		return ClassView{}, fmt.Errorf("get academic year: %w", err)
	}
	item, err := d.repo.CreateClass(ctx, tx, entity.Class{AcademicYearID: yearID, Name: name, GradeLevel: gradeLevel})
	if errors.Is(err, repository.ErrConflict) {
		return ClassView{}, ErrClassExists
	}
	if err != nil {
		return ClassView{}, fmt.Errorf("create class: %w", err)
	}
	return ClassView{Class: item}, nil
}
