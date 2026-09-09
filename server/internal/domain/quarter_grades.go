package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"journal/server/internal/repository"
	"journal/server/internal/repository/entity"
)

type QuarterGradeInput struct {
	GradingScale              string
	MaxScore, NumericValue    *float64
	TextValue, TeacherComment *string
}

func (d *ClassDomain) quarterGradeContext(ctx context.Context, tx repository.Transaction, hash string, assignmentID, quarterID uuid.UUID, permission string) (uuid.UUID, entity.ClassSubject, entity.AcademicYearQuarter, error) {
	teacherID, err := d.currentTeacher(ctx, tx, hash)
	if err != nil {
		return uuid.Nil, entity.ClassSubject{}, entity.AcademicYearQuarter{}, err
	}
	assignment, err := d.repo.GetClassSubject(ctx, tx, assignmentID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, assignment, entity.AcademicYearQuarter{}, ErrClassSubjectNotFound
	}
	if err != nil {
		return uuid.Nil, assignment, entity.AcademicYearQuarter{}, fmt.Errorf("get assignment: %w", err)
	}
	if assignment.ResponsibleTeacherID != teacherID {
		return uuid.Nil, assignment, entity.AcademicYearQuarter{}, ErrForbidden
	}
	if err = AuthorizeObject(ctx, tx, d.repo, hash, permission, assignmentID.String()); err != nil {
		return uuid.Nil, assignment, entity.AcademicYearQuarter{}, err
	}
	quarter, err := d.repo.GetAcademicYearQuarter(ctx, tx, quarterID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, assignment, quarter, ErrQuarterNotFound
	}
	if err != nil {
		return uuid.Nil, assignment, quarter, fmt.Errorf("get quarter: %w", err)
	}
	class, err := d.repo.GetClass(ctx, tx, assignment.ClassID)
	if err != nil {
		return uuid.Nil, assignment, quarter, fmt.Errorf("get class: %w", err)
	}
	if quarter.AcademicYearID != class.AcademicYearID {
		return uuid.Nil, assignment, quarter, ErrInvalidScore
	}
	return teacherID, assignment, quarter, nil
}

func (d *ClassDomain) PutQuarterGrade(ctx context.Context, tx repository.Transaction, hash string, assignmentID, quarterID, studentID uuid.UUID, input QuarterGradeInput) (entity.QuarterGrade, error) {
	teacherID, assignment, quarter, err := d.quarterGradeContext(ctx, tx, hash, assignmentID, quarterID, "can_create_score")
	if err != nil {
		return entity.QuarterGrade{}, err
	}
	membership, err := d.repo.GetClassStudent(ctx, tx, assignment.ClassID, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		return entity.QuarterGrade{}, ErrClassStudentNotFound
	}
	if err != nil {
		return entity.QuarterGrade{}, fmt.Errorf("get membership: %w", err)
	}
	if dayUTC(membership.EnrolledOn).After(dayUTC(quarter.EndsOn)) || (membership.LeftOn != nil && dayUTC(*membership.LeftOn).Before(dayUTC(quarter.StartsOn))) {
		return entity.QuarterGrade{}, ErrClassStudentNotFound
	}
	item := entity.GradeItem{GradingScale: input.GradingScale, MaxScore: input.MaxScore}
	if !validScore(item, input.NumericValue, input.TextValue) {
		return entity.QuarterGrade{}, ErrInvalidScore
	}
	input.TeacherComment = trimmedOptional(input.TeacherComment)
	result, err := d.repo.UpsertQuarterGrade(ctx, tx, entity.QuarterGrade{QuarterID: quarterID, ClassSubjectID: assignmentID, UserID: studentID, GradingScale: input.GradingScale, MaxScore: input.MaxScore, NumericValue: input.NumericValue, TextValue: input.TextValue, TeacherComment: input.TeacherComment, CreatedBy: teacherID, UpdatedBy: teacherID})
	if err != nil {
		return entity.QuarterGrade{}, fmt.Errorf("upsert quarter grade: %w", err)
	}
	return result, nil
}
func (d *ClassDomain) ListQuarterGrades(ctx context.Context, tx repository.Transaction, hash string, assignmentID, quarterID uuid.UUID) ([]entity.QuarterGrade, error) {
	_, _, _, err := d.quarterGradeContext(ctx, tx, hash, assignmentID, quarterID, "can_view_score")
	if err != nil {
		return nil, err
	}
	items, err := d.repo.ListQuarterGrades(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("list quarter grades: %w", err)
	}
	result := []entity.QuarterGrade{}
	for _, item := range items {
		if item.ClassSubjectID == assignmentID && item.QuarterID == quarterID {
			result = append(result, item)
		}
	}
	return result, nil
}
