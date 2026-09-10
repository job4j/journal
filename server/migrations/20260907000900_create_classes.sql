-- +goose Up
CREATE TABLE classes (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT, name TEXT NOT NULL, grade_level SMALLINT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(academic_year_id,name));
CREATE INDEX classes_academic_year_id_idx ON classes(academic_year_id);
-- +goose Down
DROP TABLE classes;