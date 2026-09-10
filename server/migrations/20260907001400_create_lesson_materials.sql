-- +goose Up
CREATE TABLE lesson_materials (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, title TEXT NOT NULL, url TEXT NOT NULL, position SMALLINT NOT NULL DEFAULT 1, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(lesson_id,position));
-- +goose Down
DROP TABLE lesson_materials;