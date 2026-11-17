-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    LOCK TABLE academic_year_quarters IN SHARE ROW EXCLUSIVE MODE;

    IF (SELECT count(*) FROM academic_year_quarters
        WHERE academic_year_id = '50000000-0000-4000-8000-000000000001'
          AND ((id = '55000000-0000-4000-8000-000000000001' AND number = 1
                AND starts_on = DATE '2026-09-01' AND ends_on = DATE '2026-10-30')
            OR (id = '55000000-0000-4000-8000-000000000002' AND number = 2
                AND starts_on = DATE '2026-11-09' AND ends_on = DATE '2026-12-30')
            OR (id = '55000000-0000-4000-8000-000000000003' AND number = 3
                AND starts_on = DATE '2027-01-11' AND ends_on = DATE '2027-03-19')
            OR (id = '55000000-0000-4000-8000-000000000004' AND number = 4
                AND starts_on = DATE '2027-03-29' AND ends_on = DATE '2027-05-26'))) <> 4 THEN
        RAISE EXCEPTION 'Expected corrected 2026/2027 quarters are missing; review the calendar before closing the gaps';
    END IF;

    -- Scores belong to lessons: changing the ranges makes them visible
    -- without moving or rewriting lessons or grades.
    UPDATE academic_year_quarters q
    SET starts_on = dates.starts_on, updated_at = now()
    FROM (VALUES
        ('55000000-0000-4000-8000-000000000002'::UUID, DATE '2026-10-31'),
        ('55000000-0000-4000-8000-000000000003'::UUID, DATE '2026-12-31'),
        ('55000000-0000-4000-8000-000000000004'::UUID, DATE '2027-03-20')
    ) AS dates(id, starts_on)
    WHERE q.id = dates.id;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- Reopening the gap would hide existing grades again. Require a reviewed plan.
-- +goose StatementBegin
DO $$
BEGIN
    RAISE EXCEPTION 'Closing the quarter gap cannot be safely reversed automatically; reopening it would hide grades';
END;
$$;
-- +goose StatementEnd
