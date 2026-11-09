-- +goose Up
CREATE TABLE class_schedule_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_id UUID NOT NULL REFERENCES class_schedule_weeks(id) ON DELETE CASCADE,
    lesson_date DATE NOT NULL,
    position SMALLINT NOT NULL CHECK (position BETWEEN 1 AND 20),
    class_subject_id UUID REFERENCES class_subjects(id) ON DELETE RESTRICT,
    title TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    starts_at TEXT CHECK (starts_at ~ '^([01][0-9]|2[0-3]):[0-5][0-9]$'),
    ends_at TEXT CHECK (ends_at ~ '^([01][0-9]|2[0-3]):[0-5][0-9]$'),
    CHECK (class_subject_id IS NOT NULL OR length(btrim(title)) BETWEEN 1 AND 200),
    CHECK (length(title) <= 200 AND length(notes) <= 1000),
    CHECK (ends_at IS NULL OR (starts_at IS NOT NULL AND ends_at > starts_at))
);
CREATE INDEX class_schedule_entries_week_idx ON class_schedule_entries(week_id, lesson_date, position);

-- +goose Down
DROP TABLE class_schedule_entries;
