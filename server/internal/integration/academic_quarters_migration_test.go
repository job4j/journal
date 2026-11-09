//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func calendarMigration(t *testing.T) (string, string) {
	t.Helper()
	data, err := os.ReadFile("../../migrations/20261109075354_correct_2026_2027_academic_quarters.sql")
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(data), "-- +goose Down", 2)
	if len(parts) != 2 {
		t.Fatal("missing Down section")
	}
	return parts[0], parts[1]
}

func calendarFixture(t *testing.T, periods int) pgx.Tx {
	t.Helper()
	tx := transaction(t)
	// Temporary copies keep this migration test away from application data.
	_, err := tx.Exec(context.Background(), `
		CREATE TEMP TABLE academic_years (LIKE public.academic_years INCLUDING ALL) ON COMMIT DROP;
		CREATE TEMP TABLE academic_year_quarters (LIKE public.academic_year_quarters INCLUDING ALL) ON COMMIT DROP;
		CREATE TEMP TABLE quarter_grades (LIKE public.quarter_grades INCLUDING ALL) ON COMMIT DROP;
		INSERT INTO academic_years(id,name,starts_on,ends_on,status)
		VALUES ('50000000-0000-4000-8000-000000000001','2026/2027','2026-09-01','2027-05-31','active'),
		       ('50000000-0000-4000-8000-000000000002','Other year','2027-09-01','2028-05-31','planned');
		INSERT INTO academic_year_quarters(id,academic_year_id,number,name,starts_on,ends_on)
		VALUES ('55000000-0000-4000-8000-000000000099','50000000-0000-4000-8000-000000000002',1,'Unchanged','2027-09-01','2027-10-30');
	`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(context.Background(), `
		INSERT INTO academic_year_quarters(id,academic_year_id,number,name,starts_on,ends_on)
		SELECT ('55000000-0000-4000-8000-' || lpad(n::TEXT,12,'0'))::UUID,
		'50000000-0000-4000-8000-000000000001', n, n || ' period',
		DATE '2026-09-01' + n, DATE '2026-09-10' + n FROM generate_series(1,$1::INT) n;
		`, periods)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(context.Background(), `
		INSERT INTO quarter_grades(quarter_id,class_subject_id,user_id,grading_scale,numeric_value,teacher_comment,created_by,updated_by)
		SELECT ('55000000-0000-4000-8000-' || lpad((1 + n % 2)::TEXT,12,'0'))::UUID,
		gen_random_uuid(),gen_random_uuid(),'five_point',4,'Keep comment',gen_random_uuid(),gen_random_uuid()
		FROM generate_series(1,97) n;
		CREATE TEMP TABLE grades_before AS TABLE quarter_grades;
		CREATE TEMP TABLE other_year_before AS SELECT * FROM academic_year_quarters
		WHERE academic_year_id='50000000-0000-4000-8000-000000000002';
	`)
	if err != nil {
		t.Fatal(err)
	}
	return tx
}

func TestCalendarMigrationPreservesGrades(t *testing.T) {
	up, down := calendarMigration(t)
	for _, periods := range []int{5, 6} {
		t.Run(fmt.Sprintf("%d_periods", periods), func(t *testing.T) {
			tx := calendarFixture(t, periods)
			ctx := context.Background()
			if _, err := tx.Exec(ctx, up); err != nil {
				t.Fatal(err)
			}
			var unchanged bool
			err := tx.QueryRow(ctx, `SELECT
			NOT EXISTS ((TABLE quarter_grades EXCEPT TABLE grades_before) UNION ALL (TABLE grades_before EXCEPT TABLE quarter_grades))
			AND NOT EXISTS ((SELECT * FROM academic_year_quarters WHERE academic_year_id='50000000-0000-4000-8000-000000000002' EXCEPT TABLE other_year_before)
			UNION ALL (TABLE other_year_before EXCEPT SELECT * FROM academic_year_quarters WHERE academic_year_id='50000000-0000-4000-8000-000000000002'))`).Scan(&unchanged)
			if err != nil || !unchanged {
				t.Fatalf("grades or other year changed: %v", err)
			}
			var calendar string
			err = tx.QueryRow(ctx, `SELECT string_agg(number || ':' || name || ':' || starts_on || ':' || ends_on, ',' ORDER BY number)
			FROM academic_year_quarters WHERE academic_year_id='50000000-0000-4000-8000-000000000001'`).Scan(&calendar)
			want := "1:I четверть:2026-09-01:2026-10-30,2:II четверть:2026-11-09:2026-12-30,3:III четверть:2027-01-11:2027-03-19,4:IV четверть:2027-03-29:2027-05-26"
			if err != nil || calendar != want {
				t.Fatalf("calendar = %q: %v", calendar, err)
			}
			if _, err = tx.Exec(ctx, down); err == nil || !strings.Contains(err.Error(), "cannot be safely reversed") {
				t.Fatalf("unsafe rollback allowed: %v", err)
			}
		})
	}
}

func TestCalendarMigrationRejectsUnexpectedData(t *testing.T) {
	up, _ := calendarMigration(t)
	for _, test := range []struct{ name, change, message string }{
		{"graded_fifth", `UPDATE quarter_grades SET quarter_id='55000000-0000-4000-8000-000000000005'`, "contain final grades"},
		{"graded_sixth", `UPDATE quarter_grades SET quarter_id='55000000-0000-4000-8000-000000000006'`, "contain final grades"},
		{"unexpected_id", `UPDATE academic_year_quarters SET id=gen_random_uuid() WHERE number=3`, "Unexpected periods"},
		{"missing_period", `DELETE FROM academic_year_quarters WHERE number=4`, "Unexpected periods"},
		{"short_year", `UPDATE academic_years SET ends_on='2027-05-01'`, "incompatible dates"},
	} {
		t.Run(test.name, func(t *testing.T) {
			tx := calendarFixture(t, 6)
			ctx := context.Background()
			if _, err := tx.Exec(ctx, test.change+`; CREATE TEMP TABLE periods_before AS TABLE academic_year_quarters; SAVEPOINT before_migration`); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.Exec(ctx, up); err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("expected rejection %q: %v", test.message, err)
			}
			if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT before_migration"); err != nil {
				t.Fatal(err)
			}
			var unchanged bool
			if err := tx.QueryRow(ctx, `SELECT NOT EXISTS ((TABLE academic_year_quarters EXCEPT TABLE periods_before) UNION ALL (TABLE periods_before EXCEPT TABLE academic_year_quarters))`).Scan(&unchanged); err != nil || !unchanged {
				t.Fatalf("failed migration changed periods: %v", err)
			}
		})
	}
}
