-- +goose Up
CREATE TABLE academic_year_quarters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_year_id UUID NOT NULL
        REFERENCES academic_years (id) ON DELETE CASCADE,
    number SMALLINT NOT NULL CHECK (number BETWEEN 1 AND 4),
    starts_on DATE NOT NULL,
    ends_on DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (academic_year_id, number),
    CHECK (ends_on >= starts_on)
);

CREATE INDEX academic_year_quarters_academic_year_id_idx
    ON academic_year_quarters (academic_year_id);

-- +goose Down
DROP TABLE IF EXISTS academic_year_quarters;
