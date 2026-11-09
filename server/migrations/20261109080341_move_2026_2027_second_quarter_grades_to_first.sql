-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    first_quarter CONSTANT UUID := '55000000-0000-4000-8000-000000000001';
    second_quarter CONSTANT UUID := '55000000-0000-4000-8000-000000000002';
BEGIN
    LOCK TABLE academic_years, academic_year_quarters, quarter_grades IN SHARE ROW EXCLUSIVE MODE;

    IF (SELECT count(*) FROM academic_year_quarters
        WHERE academic_year_id = '50000000-0000-4000-8000-000000000001'
          AND ((id = first_quarter AND number = 1 AND starts_on = DATE '2026-09-01' AND ends_on = DATE '2026-10-30')
            OR (id = second_quarter AND number = 2 AND starts_on = DATE '2026-11-09' AND ends_on = DATE '2026-12-30'))) <> 2 THEN
        RAISE EXCEPTION 'Expected corrected 2026/2027 quarters are missing; apply the calendar correction first';
    END IF;

    -- One-time archive of both quarters, including overwritten rows and audit fields.
    -- Do not remove it until the correction and its rollback window have been reviewed.
    CREATE TABLE migration_20261109080341_quarter_grades_backup (
        grade_id UUID PRIMARY KEY,
        grade_data JSONB NOT NULL
    );
    INSERT INTO migration_20261109080341_quarter_grades_backup (grade_id, grade_data)
    SELECT g.id, to_jsonb(g) FROM quarter_grades g
    WHERE g.quarter_id IN (first_quarter, second_quarter);

    -- The entire record from quarter II wins, including its ID and audit fields.
    DELETE FROM quarter_grades destination
    USING quarter_grades source
    WHERE destination.quarter_id = first_quarter AND source.quarter_id = second_quarter
      AND destination.class_subject_id = source.class_subject_id
      AND destination.user_id = source.user_id;

    UPDATE quarter_grades
    SET quarter_id = first_quarter
    WHERE quarter_id = second_quarter;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
DECLARE
    first_quarter CONSTANT UUID := '55000000-0000-4000-8000-000000000001';
    second_quarter CONSTANT UUID := '55000000-0000-4000-8000-000000000002';
BEGIN
    LOCK TABLE academic_years, academic_year_quarters, quarter_grades,
        migration_20261109080341_quarter_grades_backup IN SHARE ROW EXCLUSIVE MODE;

    -- Refuse rollback if either quarter was edited after the transfer.
    IF EXISTS (
        WITH original AS (
            SELECT grade_data FROM migration_20261109080341_quarter_grades_backup
        ), expected AS (
            SELECT grade_data || jsonb_build_object('quarter_id', first_quarter) AS grade_data
            FROM original WHERE grade_data->>'quarter_id' = second_quarter::TEXT
            UNION ALL
            SELECT a.grade_data FROM original a
            WHERE a.grade_data->>'quarter_id' = first_quarter::TEXT
              AND NOT EXISTS (
                  SELECT 1 FROM original b
                  WHERE b.grade_data->>'quarter_id' = second_quarter::TEXT
                    AND b.grade_data->>'class_subject_id' = a.grade_data->>'class_subject_id'
                    AND b.grade_data->>'user_id' = a.grade_data->>'user_id'
              )
        ), actual AS (
            SELECT to_jsonb(g) AS grade_data FROM quarter_grades g
            WHERE g.quarter_id IN (first_quarter, second_quarter)
        )
        (TABLE expected EXCEPT TABLE actual)
        UNION ALL
        (TABLE actual EXCEPT TABLE expected)
    ) THEN
        RAISE EXCEPTION 'Final grades changed after the transfer; automatic rollback would lose those edits';
    END IF;

    DELETE FROM quarter_grades WHERE quarter_id IN (first_quarter, second_quarter);
    INSERT INTO quarter_grades
    SELECT (jsonb_populate_record(NULL::quarter_grades, grade_data)).*
    FROM migration_20261109080341_quarter_grades_backup;
    DROP TABLE migration_20261109080341_quarter_grades_backup;
END;
$$;
-- +goose StatementEnd
