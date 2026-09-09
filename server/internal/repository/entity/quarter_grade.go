package entity

import (
	"github.com/google/uuid"
	"time"
)

type QuarterGrade struct {
	ID             uuid.UUID `db:"id"`
	QuarterID      uuid.UUID `db:"quarter_id"`
	ClassSubjectID uuid.UUID `db:"class_subject_id"`
	UserID         uuid.UUID `db:"user_id"`
	GradingScale   string    `db:"grading_scale"`
	MaxScore       *float64  `db:"max_score"`
	NumericValue   *float64  `db:"numeric_value"`
	TextValue      *string   `db:"text_value"`
	TeacherComment *string   `db:"teacher_comment"`
	CreatedBy      uuid.UUID `db:"created_by"`
	UpdatedBy      uuid.UUID `db:"updated_by"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
