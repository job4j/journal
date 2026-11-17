//go:build integration

package integration

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func secondQuarterBoundaryMigration(t *testing.T) (string, string) {
	t.Helper()
	data, err := os.ReadFile("../../migrations/20261117074138_start_2026_2027_second_quarter_after_first.sql")
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(data), "-- +goose Down", 2)
	if len(parts) != 2 {
		t.Fatal("missing Down section")
	}
	return parts[0], parts[1]
}

func correctedCalendarFixture(t *testing.T) pgx.Tx {
	t.Helper()
	tx := calendarFixture(t, 5)
	up, _ := calendarMigration(t)
	if _, err := tx.Exec(context.Background(), up); err != nil {
		t.Fatal(err)
	}
	return tx
}

func TestQuarterBoundariesMigration(t *testing.T) {
	tx := correctedCalendarFixture(t)
	ctx := context.Background()
	_, err := tx.Exec(ctx, `
		CREATE TEMP TABLE lessons (LIKE public.lessons INCLUDING ALL) ON COMMIT DROP;
		CREATE TEMP TABLE grade_items (LIKE public.grade_items INCLUDING ALL) ON COMMIT DROP;
		CREATE TEMP TABLE scores (LIKE public.scores INCLUDING ALL) ON COMMIT DROP;
		INSERT INTO lessons(class_subject_id,lesson_date,topic,created_by)
		SELECT gen_random_uuid(), d::date, 'Boundary lesson', gen_random_uuid()
		FROM unnest(ARRAY['2026-10-30','2026-10-31','2026-11-08','2026-11-09',
		    '2026-12-30','2026-12-31','2027-01-10','2027-01-11',
		    '2027-03-19','2027-03-20','2027-03-28','2027-03-29','2027-05-26','2027-05-27']) d;
		INSERT INTO grade_items(lesson_id,title,kind)
		SELECT id,'Boundary work','classwork' FROM lessons;
		INSERT INTO scores(grade_item_id,user_id,numeric_value,created_by,updated_by)
		SELECT id,gen_random_uuid(),4,gen_random_uuid(),gen_random_uuid() FROM grade_items;
		CREATE TEMP TABLE lessons_before AS TABLE lessons;
		CREATE TEMP TABLE items_before AS TABLE grade_items;
		CREATE TEMP TABLE scores_before AS TABLE scores;
		CREATE TEMP TABLE quarters_before AS TABLE academic_year_quarters;
	`)
	if err != nil {
		t.Fatal(err)
	}
	up, down := secondQuarterBoundaryMigration(t)
	if _, err := tx.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	var dates string
	err = tx.QueryRow(ctx, `SELECT string_agg(q.number || ':' || l.lesson_date, ',' ORDER BY l.lesson_date)
		FROM scores s JOIN grade_items gi ON gi.id=s.grade_item_id JOIN lessons l ON l.id=gi.lesson_id
		JOIN academic_year_quarters q ON l.lesson_date BETWEEN q.starts_on AND q.ends_on
		WHERE q.academic_year_id='50000000-0000-4000-8000-000000000001'`).Scan(&dates)
	want := "1:2026-10-30,2:2026-10-31,2:2026-11-08,2:2026-11-09,2:2026-12-30," +
		"3:2026-12-31,3:2027-01-10,3:2027-01-11,3:2027-03-19," +
		"4:2027-03-20,4:2027-03-28,4:2027-03-29,4:2027-05-26"
	if err != nil || dates != want {
		t.Fatalf("quarter scores = %q, want %q: %v", dates, want, err)
	}
	var contiguous bool
	err = tx.QueryRow(ctx, `SELECT count(*)=3 AND bool_and(q.starts_on=previous.ends_on+1)
		FROM academic_year_quarters q JOIN academic_year_quarters previous
		ON previous.academic_year_id=q.academic_year_id AND previous.number=q.number-1
		WHERE q.academic_year_id='50000000-0000-4000-8000-000000000001'`).Scan(&contiguous)
	if err != nil || !contiguous {
		t.Fatalf("quarters have gaps or overlaps: %v", err)
	}
	var unchanged bool
	err = tx.QueryRow(ctx, `SELECT
		NOT EXISTS ((TABLE lessons EXCEPT TABLE lessons_before) UNION ALL (TABLE lessons_before EXCEPT TABLE lessons))
		AND NOT EXISTS ((TABLE grade_items EXCEPT TABLE items_before) UNION ALL (TABLE items_before EXCEPT TABLE grade_items))
		AND NOT EXISTS ((TABLE scores EXCEPT TABLE scores_before) UNION ALL (TABLE scores_before EXCEPT TABLE scores))
		AND NOT EXISTS ((TABLE quarter_grades EXCEPT TABLE grades_before) UNION ALL (TABLE grades_before EXCEPT TABLE quarter_grades))
		AND NOT EXISTS (
			SELECT 1 FROM academic_year_quarters q FULL JOIN quarters_before b USING(id)
			WHERE CASE WHEN q.academic_year_id='50000000-0000-4000-8000-000000000001' AND q.number BETWEEN 2 AND 4
			THEN q.starts_on <> CASE q.number WHEN 2 THEN DATE '2026-10-31'
			     WHEN 3 THEN DATE '2026-12-31' WHEN 4 THEN DATE '2027-03-20' END OR
			     (to_jsonb(q)-'starts_on'-'updated_at') IS DISTINCT FROM (to_jsonb(b)-'starts_on'-'updated_at')
			ELSE to_jsonb(q) IS DISTINCT FROM to_jsonb(b) END
		)`).Scan(&unchanged)
	if err != nil || !unchanged {
		t.Fatalf("unexpected data changes: %v", err)
	}
	if _, err := tx.Exec(ctx, down); err == nil || !strings.Contains(err.Error(), "would hide grades") {
		t.Fatalf("unsafe rollback allowed: %v", err)
	}
}

func TestQuarterBoundariesMigrationRejectsUnexpectedCalendar(t *testing.T) {
	up, _ := secondQuarterBoundaryMigration(t)
	for _, test := range []struct{ name, change string }{
		{"first_end", `UPDATE academic_year_quarters SET ends_on='2026-10-29' WHERE number=1`},
		{"second_start", `UPDATE academic_year_quarters SET starts_on='2026-11-10' WHERE number=2`},
		{"second_end", `UPDATE academic_year_quarters SET ends_on='2026-12-29' WHERE number=2`},
		{"missing_second", `DELETE FROM quarter_grades; DELETE FROM academic_year_quarters WHERE number=2`},
		{"third_start", `UPDATE academic_year_quarters SET starts_on='2027-01-12' WHERE number=3`},
		{"third_end", `UPDATE academic_year_quarters SET ends_on='2027-03-18' WHERE number=3`},
		{"fourth_start", `UPDATE academic_year_quarters SET starts_on='2027-03-30' WHERE number=4`},
		{"fourth_end", `UPDATE academic_year_quarters SET ends_on='2027-05-25' WHERE number=4`},
		{"missing_third", `DELETE FROM academic_year_quarters WHERE number=3`},
		{"missing_fourth", `DELETE FROM academic_year_quarters WHERE number=4`},
	} {
		t.Run(test.name, func(t *testing.T) {
			tx := correctedCalendarFixture(t)
			ctx := context.Background()
			if _, err := tx.Exec(ctx, test.change); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.Exec(ctx, up); err == nil || !strings.Contains(err.Error(), "Expected corrected") {
				t.Fatalf("unexpected calendar accepted: %v", err)
			}
		})
	}
}
