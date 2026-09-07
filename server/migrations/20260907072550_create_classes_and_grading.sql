-- +goose Up
CREATE TABLE academic_years (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    starts_on DATE NOT NULL,
    ends_on DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'planned'
        CHECK (status IN ('planned', 'active', 'completed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (ends_on > starts_on)
);

CREATE UNIQUE INDEX academic_years_single_active_idx
    ON academic_years (status)
    WHERE status = 'active';

CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    academic_year_id UUID NOT NULL
        REFERENCES academic_years (id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    grade_level SMALLINT NOT NULL CHECK (grade_level BETWEEN 1 AND 11),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (academic_year_id, name)
);

CREATE INDEX classes_academic_year_id_idx ON classes (academic_year_id);

CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE, 
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE class_subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects (id) ON DELETE RESTRICT,
    responsible_teacher_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (class_id, subject_id)
);

CREATE INDEX class_subjects_teacher_id_idx
    ON class_subjects (responsible_teacher_id);

CREATE TABLE class_students (
    class_id UUID NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    enrolled_on DATE NOT NULL,
    left_on DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (class_id, user_id),
    CHECK (left_on IS NULL OR left_on >= enrolled_on)
);

CREATE INDEX class_students_user_id_idx ON class_students (user_id);

CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_subject_id UUID NOT NULL
        REFERENCES class_subjects (id) ON DELETE CASCADE,
    lesson_date DATE NOT NULL,
    position SMALLINT NOT NULL DEFAULT 1 CHECK (position > 0),
    topic TEXT NOT NULL,
    homework TEXT,
    created_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (class_subject_id, lesson_date, position)
);

CREATE INDEX lessons_class_subject_date_idx
    ON lessons (class_subject_id, lesson_date);

CREATE TABLE lesson_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    position SMALLINT NOT NULL DEFAULT 1 CHECK (position > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (lesson_id, position)
);

CREATE TABLE grade_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    kind TEXT NOT NULL
        CHECK (kind IN ('homework', 'classwork', 'knowledge_check', 'other')),
    grading_scale TEXT NOT NULL DEFAULT 'five_point'
        CHECK (grading_scale IN ('five_point', 'points', 'pass_fail')),
    max_score NUMERIC(8, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (
        (grading_scale = 'points' AND max_score IS NOT NULL AND max_score > 0)
        OR (grading_scale <> 'points' AND max_score IS NULL)
    )
);

CREATE INDEX grade_items_lesson_id_idx ON grade_items (lesson_id);

CREATE TABLE scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grade_item_id UUID NOT NULL REFERENCES grade_items (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    numeric_value NUMERIC(8, 2),
    text_value TEXT,
    teacher_comment TEXT,
    created_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    updated_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (grade_item_id, user_id),
    CHECK (
        (numeric_value IS NOT NULL AND text_value IS NULL)
        OR (numeric_value IS NULL AND text_value IS NOT NULL)
    )
);

CREATE INDEX scores_user_id_idx ON scores (user_id);

-- +goose Down
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS grade_items;
DROP TABLE IF EXISTS lesson_materials;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS class_students;
DROP TABLE IF EXISTS class_subjects;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS academic_years;
