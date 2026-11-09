package repository

import (
	"context"
	"github.com/google/uuid"
	"journal/server/internal/repository/entity"
	"time"
)

func (r *Repository) LockScheduleClass(ctx context.Context, tx Transaction, id uuid.UUID) error {
	_, err := queryOne[entity.Class](ctx, tx, "lock schedule class", `SELECT id,academic_year_id,name,grade_level,created_at,updated_at FROM classes WHERE id=$1 FOR UPDATE`, id)
	return err
}
func (r *Repository) GetScheduleWeek(ctx context.Context, tx Transaction, classID uuid.UUID, start time.Time) (entity.ScheduleWeek, error) {
	return queryOne[entity.ScheduleWeek](ctx, tx, "get schedule week", `SELECT id,class_id,week_start,revision,updated_at FROM class_schedule_weeks WHERE class_id=$1 AND week_start=$2`, classID, start)
}
func (r *Repository) SaveScheduleWeek(ctx context.Context, tx Transaction, classID uuid.UUID, start time.Time) (entity.ScheduleWeek, error) {
	return queryOne[entity.ScheduleWeek](ctx, tx, "save schedule week", `INSERT INTO class_schedule_weeks(class_id,week_start) VALUES($1,$2) ON CONFLICT(class_id,week_start) DO UPDATE SET revision=class_schedule_weeks.revision+1,updated_at=now() RETURNING id,class_id,week_start,revision,updated_at`, classID, start)
}
func (r *Repository) ListScheduleEntries(ctx context.Context, tx Transaction, weekID uuid.UUID) ([]entity.ScheduleEntry, error) {
	return queryMany[entity.ScheduleEntry](ctx, tx, "list schedule entries", `SELECT e.id,e.week_id,e.lesson_date,e.position,e.class_subject_id,coalesce(s.name,e.title) AS title,e.notes,e.starts_at,e.ends_at FROM class_schedule_entries e LEFT JOIN class_subjects cs ON cs.id=e.class_subject_id LEFT JOIN subjects s ON s.id=cs.subject_id WHERE e.week_id=$1 ORDER BY e.lesson_date,e.position,e.starts_at,e.id`, weekID)
}
func (r *Repository) DeleteScheduleEntries(ctx context.Context, transaction Transaction, weekID uuid.UUID) error {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM class_schedule_entries WHERE week_id=$1`, weekID)
	if err != nil {
		return operationError("delete schedule entries", err)
	}
	return nil
}
func (r *Repository) CreateScheduleEntry(ctx context.Context, transaction Transaction, e entity.ScheduleEntry) error {
	tx, err := pgxTransaction(transaction)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO class_schedule_entries(id,week_id,lesson_date,position,class_subject_id,title,notes,starts_at,ends_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, e.ID, e.WeekID, e.LessonDate, e.Position, e.ClassSubjectID, e.Title, e.Notes, e.StartsAt, e.EndsAt)
	if err != nil {
		return operationError("create schedule entry", err)
	}
	return nil
}
