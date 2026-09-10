-- +goose Up
CREATE TABLE lessons (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), class_subject_id UUID NOT NULL REFERENCES class_subjects(id) ON DELETE CASCADE, lesson_date DATE NOT NULL, position SMALLINT NOT NULL DEFAULT 1, topic TEXT NOT NULL, homework TEXT, created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(class_subject_id,lesson_date,position));
CREATE INDEX lessons_class_subject_date_idx ON lessons(class_subject_id,lesson_date);
-- +goose Down
DROP TABLE lessons;