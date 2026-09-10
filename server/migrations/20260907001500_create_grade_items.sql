-- +goose Up
CREATE TABLE grade_items (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, title TEXT NOT NULL, kind TEXT NOT NULL, grading_scale TEXT NOT NULL DEFAULT 'five_point', max_score SMALLINT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE INDEX grade_items_lesson_id_idx ON grade_items(lesson_id);
-- +goose Down
DROP TABLE grade_items;