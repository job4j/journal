-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    target_year CONSTANT UUID := '50000000-0000-4000-8000-000000000001';
BEGIN
    -- Keep existing quarter IDs and every final grade unchanged.
    LOCK TABLE academic_years, academic_year_quarters, quarter_grades IN SHARE ROW EXCLUSIVE MODE;

    IF NOT EXISTS (
        SELECT 1 FROM academic_years
        WHERE id = target_year
          AND starts_on <= DATE '2026-09-01'
          AND ends_on >= DATE '2027-05-26'
    ) THEN
        RAISE EXCEPTION 'Expected 2026/2027 academic year is missing or has incompatible dates';
    END IF;

    IF (SELECT count(*) FROM academic_year_quarters
        WHERE academic_year_id = target_year AND number BETWEEN 1 AND 4
          AND id = ('55000000-0000-4000-8000-' || lpad(number::TEXT, 12, '0'))::UUID) <> 4
       OR EXISTS (
           SELECT 1 FROM academic_year_quarters
           WHERE academic_year_id = target_year
             AND (number NOT BETWEEN 1 AND 6
                  OR id <> ('55000000-0000-4000-8000-' || lpad(number::TEXT, 12, '0'))::UUID)
       ) THEN
        RAISE EXCEPTION 'Unexpected periods for 2026/2027; review them before correcting the calendar';
    END IF;

    IF EXISTS (
        SELECT 1 FROM quarter_grades g
        JOIN academic_year_quarters q ON q.id = g.quarter_id
        WHERE q.academic_year_id = target_year AND q.number IN (5, 6)
    ) THEN
        RAISE EXCEPTION 'Periods 5 or 6 contain final grades; the calendar correction must not change them';
    END IF;

    UPDATE academic_year_quarters q
    SET name = dates.name, starts_on = dates.starts_on, ends_on = dates.ends_on, updated_at = now()
    FROM (VALUES
        (1, 'I четверть', DATE '2026-09-01', DATE '2026-10-30'),
        (2, 'II четверть', DATE '2026-11-09', DATE '2026-12-30'),
        (3, 'III четверть', DATE '2027-01-11', DATE '2027-03-19'),
        (4, 'IV четверть', DATE '2027-03-29', DATE '2027-05-26')
    ) AS dates(number, name, starts_on, ends_on)
    WHERE q.academic_year_id = target_year AND q.number = dates.number;

    DELETE FROM academic_year_quarters
    WHERE academic_year_id = target_year AND number IN (5, 6);
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- Irreversible: installations have different original dates and five or six periods.
-- +goose StatementBegin
DO $$
BEGIN
    RAISE EXCEPTION 'Calendar correction cannot be safely reversed automatically; restore from a reviewed backup';
END;
$$;
-- +goose StatementEnd
