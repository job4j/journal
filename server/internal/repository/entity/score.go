package entity

import (
	"github.com/google/uuid"
	"time"
)

type Score struct {
	ID             uuid.UUID `db:"id"`
	GradeItemID    uuid.UUID `db:"grade_item_id"`
	UserID         uuid.UUID `db:"user_id"`
	NumericValue   *int8     `db:"numeric_value"`
	TextValue      *string   `db:"text_value"`
	TeacherComment *string   `db:"teacher_comment"`
	CreatedBy      uuid.UUID `db:"created_by"`
	UpdatedBy      uuid.UUID `db:"updated_by"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
