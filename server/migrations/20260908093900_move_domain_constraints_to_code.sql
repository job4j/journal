-- +goose Up
-- Domain validation is intentionally performed in the domain layer. Foreign
-- keys and uniqueness constraints remain as relational integrity guarantees.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS permissions_check;
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS sessions_check;
ALTER TABLE academic_years DROP CONSTRAINT IF EXISTS academic_years_status_check;
ALTER TABLE academic_years DROP CONSTRAINT IF EXISTS academic_years_check;
ALTER TABLE classes DROP CONSTRAINT IF EXISTS classes_grade_level_check;
ALTER TABLE class_students DROP CONSTRAINT IF EXISTS class_students_check;
ALTER TABLE lessons DROP CONSTRAINT IF EXISTS lessons_position_check;
ALTER TABLE lesson_materials DROP CONSTRAINT IF EXISTS lesson_materials_position_check;
ALTER TABLE grade_items DROP CONSTRAINT IF EXISTS grade_items_kind_check;
ALTER TABLE grade_items DROP CONSTRAINT IF EXISTS grade_items_grading_scale_check;
ALTER TABLE grade_items DROP CONSTRAINT IF EXISTS grade_items_check;
ALTER TABLE scores DROP CONSTRAINT IF EXISTS scores_check;
ALTER TABLE academic_year_quarters DROP CONSTRAINT IF EXISTS academic_year_quarters_number_check;
ALTER TABLE academic_year_quarters DROP CONSTRAINT IF EXISTS academic_year_quarters_check;

ALTER TABLE grade_items
    ALTER COLUMN max_score TYPE SMALLINT USING max_score::SMALLINT;
ALTER TABLE scores
    ALTER COLUMN numeric_value TYPE SMALLINT USING numeric_value::SMALLINT;

-- +goose Down
-- Reverting the numeric representation is safe. Removed domain constraints are
-- not recreated because rows written after this migration may not satisfy them.
ALTER TABLE grade_items
    ALTER COLUMN max_score TYPE NUMERIC(8, 2) USING max_score::NUMERIC(8, 2);
ALTER TABLE scores
    ALTER COLUMN numeric_value TYPE NUMERIC(8, 2) USING numeric_value::NUMERIC(8, 2);
