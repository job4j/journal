-- +goose Up
ALTER TABLE academic_year_quarters
ADD COLUMN name TEXT;

UPDATE academic_year_quarters
SET name = number::TEXT || ' период';

ALTER TABLE academic_year_quarters
ALTER COLUMN name SET NOT NULL;

-- +goose Down
ALTER TABLE academic_year_quarters
DROP COLUMN name;