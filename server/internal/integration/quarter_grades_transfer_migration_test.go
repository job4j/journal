//go:build integration

package integration

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func gradeTransferMigration(t *testing.T) (string, string) {
	t.Helper()
	data, err := os.ReadFile("../../migrations/20261109080341_move_2026_2027_second_quarter_grades_to_first.sql")
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(data), "-- +goose Down", 2)
	if len(parts) != 2 {
		t.Fatal("missing Down section")
	}
	return parts[0], parts[1]
}

func gradeTransferFixture(t *testing.T) pgx.Tx {
	t.Helper()
	tx := calendarFixture(t, 6)
	calendarUp, _ := calendarMigration(t)
	if _, err := tx.Exec(context.Background(), calendarUp); err != nil {
		t.Fatal(err)
	}
	// The archive created by the migration must also stay in the temporary schema.
	_, err := tx.Exec(context.Background(), `SET LOCAL search_path = pg_temp, public;
		INSERT INTO quarter_grades
		SELECT (jsonb_populate_record(NULL::quarter_grades, to_jsonb(g) || jsonb_build_object(
		'id',gen_random_uuid(),'quarter_id',q.id))).*
		FROM (SELECT * FROM quarter_grades LIMIT 1) g
		CROSS JOIN academic_year_quarters q WHERE q.number=3 OR q.number=1 AND q.name='Unchanged';`)
	if err != nil {
		t.Fatal(err)
	}
	return tx
}

func TestGradeTransferMigrationAndRollback(t *testing.T) {
	up, down := gradeTransferMigration(t)
	for _, name := range []string{"without_collision", "second_quarter_wins"} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tx := gradeTransferFixture(t)
			if name == "second_quarter_wins" {
				_, err := tx.Exec(ctx, `INSERT INTO quarter_grades
				SELECT (jsonb_populate_record(NULL::quarter_grades, to_jsonb(g) || jsonb_build_object(
				'id',gen_random_uuid(),'quarter_id','55000000-0000-4000-8000-000000000001',
				'numeric_value',2,'teacher_comment','Old first-quarter comment',
				'created_by',gen_random_uuid(),'updated_by',gen_random_uuid()))).*
				FROM quarter_grades g WHERE quarter_id='55000000-0000-4000-8000-000000000002' LIMIT 1;`)
				if err != nil {
					t.Fatal(err)
				}
			}
			if _, err := tx.Exec(ctx, `CREATE TEMP TABLE transfer_before AS TABLE quarter_grades`); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.Exec(ctx, up); err != nil {
				t.Fatal(err)
			}
			var archived, total, second int
			if err := tx.QueryRow(ctx, `SELECT
			(SELECT count(*) FROM migration_20261109080341_quarter_grades_backup),
			(SELECT count(*) FROM quarter_grades),
			(SELECT count(*) FROM quarter_grades WHERE quarter_id='55000000-0000-4000-8000-000000000002')`).Scan(&archived, &total, &second); err != nil {
				t.Fatal(err)
			}
			wantArchived := 97
			if name == "second_quarter_wins" {
				wantArchived++
			}
			if archived != wantArchived || total != 99 || second != 0 {
				t.Fatalf("archive=%d, total=%d, second=%d", archived, total, second)
			}
			var preserved bool
			// Every original quarter-II record survives with only its quarter changed.
			err := tx.QueryRow(ctx, `SELECT NOT EXISTS (
			SELECT 1 FROM transfer_before b LEFT JOIN quarter_grades a ON a.id=b.id
			WHERE b.quarter_id='55000000-0000-4000-8000-000000000002'
			AND (to_jsonb(a) IS DISTINCT FROM (to_jsonb(b) || jsonb_build_object('quarter_id','55000000-0000-4000-8000-000000000001'))))
			AND NOT EXISTS (
			SELECT 1 FROM transfer_before b LEFT JOIN quarter_grades a ON a.id=b.id
			WHERE b.quarter_id NOT IN ('55000000-0000-4000-8000-000000000001','55000000-0000-4000-8000-000000000002')
			AND to_jsonb(a) IS DISTINCT FROM to_jsonb(b))`).Scan(&preserved)
			if err != nil || !preserved {
				t.Fatalf("source record or unrelated grade changed: %v", err)
			}
			if _, err := tx.Exec(ctx, down); err != nil {
				t.Fatal(err)
			}
			err = tx.QueryRow(ctx, `SELECT NOT EXISTS ((TABLE quarter_grades EXCEPT TABLE transfer_before)
			UNION ALL (TABLE transfer_before EXCEPT TABLE quarter_grades))
			AND to_regclass('migration_20261109080341_quarter_grades_backup') IS NULL`).Scan(&preserved)
			if err != nil || !preserved {
				t.Fatalf("rollback did not restore all original records: %v", err)
			}
		})
	}
}

func TestGradeTransferRollbackRejectsLaterChanges(t *testing.T) {
	up, down := gradeTransferMigration(t)
	for _, test := range []struct{ name, sql string }{
		{"edit", `UPDATE quarter_grades SET teacher_comment='New teacher edit' WHERE quarter_id='55000000-0000-4000-8000-000000000001'`},
		{"delete", `DELETE FROM quarter_grades WHERE id=(SELECT id FROM quarter_grades WHERE quarter_id='55000000-0000-4000-8000-000000000001' LIMIT 1)`},
		{"new_second_grade", `INSERT INTO quarter_grades SELECT (jsonb_populate_record(NULL::quarter_grades,to_jsonb(g) || jsonb_build_object('id',gen_random_uuid(),'class_subject_id',gen_random_uuid(),'quarter_id','55000000-0000-4000-8000-000000000002'))).* FROM quarter_grades g LIMIT 1`},
	} {
		t.Run(test.name, func(t *testing.T) {
			tx := gradeTransferFixture(t)
			ctx := context.Background()
			if _, err := tx.Exec(ctx, up); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.Exec(ctx, test.sql+`; CREATE TEMP TABLE after_teacher_change AS TABLE quarter_grades; SAVEPOINT before_down`); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.Exec(ctx, down); err == nil || !strings.Contains(err.Error(), "Final grades changed") {
				t.Fatalf("expected refusal to discard teacher changes: %v", err)
			}
			if _, err := tx.Exec(ctx, `ROLLBACK TO SAVEPOINT before_down`); err != nil {
				t.Fatal(err)
			}
			var preserved bool
			if err := tx.QueryRow(ctx, `SELECT NOT EXISTS ((TABLE quarter_grades EXCEPT TABLE after_teacher_change)
			UNION ALL (TABLE after_teacher_change EXCEPT TABLE quarter_grades))`).Scan(&preserved); err != nil || !preserved {
				t.Fatalf("teacher changes were lost: %v", err)
			}
		})
	}
}

func TestGradeTransferRequiresCorrectedCalendar(t *testing.T) {
	up, _ := gradeTransferMigration(t)
	tx := calendarFixture(t, 6)
	if _, err := tx.Exec(context.Background(), up); err == nil || !strings.Contains(err.Error(), "calendar correction first") {
		t.Fatalf("expected rejection of uncorrected calendar: %v", err)
	}
}
