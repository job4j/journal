package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type ClassView struct {
	Class        entity.Class
	StudentCount int
	Quarters     []entity.AcademicYearQuarter
}

func quartersForYear(
	quarters []entity.AcademicYearQuarter,
	yearID uuid.UUID,
) []entity.AcademicYearQuarter {
	result := []entity.AcademicYearQuarter{}
	for _, quarter := range quarters {
		if quarter.AcademicYearID == yearID {
			result = append(result, quarter)
		}
	}
	return result
}

func (d *ClassDomain) ListClasses(
	ctx context.Context,
	tx repository.Transaction,
	hash string,
	yearID uuid.UUID,
) ([]ClassView, error) {
	if err := Authorize(ctx, tx, d.repo, hash, "can_view_class"); err != nil {
		return nil, err
	}
	classes, err := d.repo.ListClasses(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list classes: %w", err)
	}
	students, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list class students: %w", err)
	}
	quarters, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list academic year quarters: %w", err)
	}
	result := []ClassView{}
	for _, class := range classes {
		if class.AcademicYearID != yearID {
			continue
		}
		view := ClassView{
			Class:    class,
			Quarters: quartersForYear(quarters, class.AcademicYearID),
		}
		for _, student := range students {
			if student.ClassID == class.ID && student.LeftOn == nil {
				view.StudentCount++
			}
		}
		result = append(result, view)
	}
	return result, nil
}

func (d *ClassDomain) GetClass(
	ctx context.Context,
	tx repository.Transaction,
	hash string,
	id uuid.UUID,
) (ClassView, error) {
	if err := AuthorizeObject(ctx, tx, d.repo, hash, "can_view_class", id.String()); err != nil {
		return ClassView{}, err
	}
	class, err := d.repo.GetClass(ctx, tx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ClassView{}, ErrClassNotFound
	}
	if err != nil {
		return ClassView{}, fmt.Errorf("get class: %w", err)
	}
	students, err := d.repo.ListClassStudents(ctx, tx)
	if err != nil {
		return ClassView{}, fmt.Errorf("list class students: %w", err)
	}
	quarters, err := d.repo.ListAcademicYearQuarters(ctx, tx)
	if err != nil {
		return ClassView{}, fmt.Errorf("list academic year quarters: %w", err)
	}
	view := ClassView{
		Class:    class,
		Quarters: quartersForYear(quarters, class.AcademicYearID),
	}
	for _, student := range students {
		if student.ClassID == id && student.LeftOn == nil {
			view.StudentCount++
		}
	}
	return view, nil
}
