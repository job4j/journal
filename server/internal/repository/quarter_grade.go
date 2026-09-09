package repository

import (
	"context"
	"journal/server/internal/repository/entity"
)

func (r *Repository) UpsertQuarterGrade(ctx context.Context, tx Transaction, value entity.QuarterGrade) (entity.QuarterGrade, error) {
	return queryOne[entity.QuarterGrade](ctx, tx, "upsert quarter grade", `INSERT INTO quarter_grades(quarter_id,class_subject_id,user_id,grading_scale,max_score,numeric_value,text_value,teacher_comment,created_by,updated_by) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$9) ON CONFLICT(quarter_id,class_subject_id,user_id) DO UPDATE SET grading_scale=EXCLUDED.grading_scale,max_score=EXCLUDED.max_score,numeric_value=EXCLUDED.numeric_value,text_value=EXCLUDED.text_value,teacher_comment=EXCLUDED.teacher_comment,updated_by=EXCLUDED.updated_by,updated_at=now() RETURNING id,quarter_id,class_subject_id,user_id,grading_scale,max_score,numeric_value,text_value,teacher_comment,created_by,updated_by,created_at,updated_at`, value.QuarterID, value.ClassSubjectID, value.UserID, value.GradingScale, value.MaxScore, value.NumericValue, value.TextValue, value.TeacherComment, value.CreatedBy)
}
func (r *Repository) ListQuarterGrades(ctx context.Context, tx Transaction) ([]entity.QuarterGrade, error) {
	return queryMany[entity.QuarterGrade](ctx, tx, "list quarter grades", `SELECT id,quarter_id,class_subject_id,user_id,grading_scale,max_score,numeric_value,text_value,teacher_comment,created_by,updated_by,created_at,updated_at FROM quarter_grades ORDER BY created_at,id`)
}
