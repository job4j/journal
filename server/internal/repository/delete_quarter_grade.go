package repository

import (
	"context"
	"github.com/google/uuid"
)

func (r *Repository) DeleteQuarterGrade(ctx context.Context, tx Transaction, quarterID, classSubjectID, studentID uuid.UUID) error {
	return deleteRows(ctx, tx, "delete quarter grade", `DELETE FROM quarter_grades WHERE quarter_id=$1 AND class_subject_id=$2 AND user_id=$3`, quarterID, classSubjectID, studentID)
}
