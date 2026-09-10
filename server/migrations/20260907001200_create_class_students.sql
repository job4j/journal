-- +goose Up
CREATE TABLE class_students (class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, enrolled_on DATE NOT NULL, left_on DATE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(class_id,user_id));
CREATE INDEX class_students_user_id_idx ON class_students(user_id);
-- +goose Down
DROP TABLE class_students;