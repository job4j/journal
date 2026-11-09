-- +goose Up
CREATE TABLE class_schedule_weeks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE RESTRICT,
    week_start DATE NOT NULL CHECK (extract(isodow FROM week_start) = 1),
    revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (class_id, week_start)
);

-- +goose Down
DROP TABLE class_schedule_weeks;
