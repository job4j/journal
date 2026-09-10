-- +goose Up
CREATE TABLE class_subjects (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE, subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE RESTRICT, responsible_teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(class_id,subject_id));
CREATE INDEX class_subjects_teacher_id_idx ON class_subjects(responsible_teacher_id);
-- +goose Down
DROP TABLE class_subjects;