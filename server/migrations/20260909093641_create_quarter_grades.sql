-- +goose Up
CREATE TABLE quarter_grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quarter_id UUID NOT NULL REFERENCES academic_year_quarters (id) ON DELETE RESTRICT,
    class_subject_id UUID NOT NULL REFERENCES class_subjects (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    grading_scale TEXT NOT NULL,
    max_score NUMERIC(8, 2),
    numeric_value NUMERIC(8, 2),
    text_value TEXT,
    teacher_comment TEXT,
    created_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    updated_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (quarter_id, class_subject_id, user_id),
    CHECK ((numeric_value IS NOT NULL AND text_value IS NULL) OR (numeric_value IS NULL AND text_value IS NOT NULL))
);

CREATE INDEX quarter_grades_user_id_idx ON quarter_grades (user_id);

-- +goose Down
DROP TABLE IF EXISTS quarter_grades;
