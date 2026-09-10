-- +goose Up
CREATE TABLE absences (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, recorded_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(lesson_id,user_id));
CREATE INDEX absences_user_id_idx ON absences(user_id);
-- +goose Down
DROP TABLE absences;