-- +goose Up
CREATE TABLE academic_year_quarters (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE, number SMALLINT NOT NULL, starts_on DATE NOT NULL, ends_on DATE NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(academic_year_id,number));
CREATE INDEX academic_year_quarters_academic_year_id_idx ON academic_year_quarters(academic_year_id);
-- +goose Down
DROP TABLE academic_year_quarters;