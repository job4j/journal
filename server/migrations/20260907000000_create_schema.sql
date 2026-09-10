-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_login_unique_idx ON users (lower(login));
CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email)) WHERE email IS NOT NULL;

CREATE TABLE roles (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT NOT NULL UNIQUE, name TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE TABLE permissions (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT NOT NULL, value TEXT, description TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE UNIQUE INDEX permissions_code_value_unique_idx ON permissions (code, value) NULLS NOT DISTINCT;
CREATE TABLE user_roles (user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(user_id,role_id));
CREATE TABLE role_permissions (role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE, permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(role_id,permission_id));
CREATE TABLE user_permissions (user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(user_id,permission_id));
CREATE TABLE sessions (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, token_hash TEXT NOT NULL UNIQUE, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), last_used_at TIMESTAMPTZ NOT NULL DEFAULT now(), revoked_at TIMESTAMPTZ);
CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_active_expiry_idx ON sessions(expires_at) WHERE revoked_at IS NULL;

CREATE TABLE academic_years (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT NOT NULL UNIQUE, starts_on DATE NOT NULL, ends_on DATE NOT NULL, status TEXT NOT NULL DEFAULT 'planned', created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE UNIQUE INDEX academic_years_single_active_idx ON academic_years(status) WHERE status='active';
CREATE TABLE academic_year_quarters (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE, number SMALLINT NOT NULL, starts_on DATE NOT NULL, ends_on DATE NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(academic_year_id,number));
CREATE INDEX academic_year_quarters_academic_year_id_idx ON academic_year_quarters(academic_year_id);
CREATE TABLE classes (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT, name TEXT NOT NULL, grade_level SMALLINT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(academic_year_id,name));
CREATE INDEX classes_academic_year_id_idx ON classes(academic_year_id);
CREATE TABLE subjects (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT NOT NULL UNIQUE, name TEXT NOT NULL UNIQUE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE TABLE class_subjects (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE, subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE RESTRICT, responsible_teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(class_id,subject_id));
CREATE INDEX class_subjects_teacher_id_idx ON class_subjects(responsible_teacher_id);
CREATE TABLE class_students (class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, enrolled_on DATE NOT NULL, left_on DATE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(class_id,user_id));
CREATE INDEX class_students_user_id_idx ON class_students(user_id);
CREATE TABLE lessons (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), class_subject_id UUID NOT NULL REFERENCES class_subjects(id) ON DELETE CASCADE, lesson_date DATE NOT NULL, position SMALLINT NOT NULL DEFAULT 1, topic TEXT NOT NULL, homework TEXT, created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(class_subject_id,lesson_date,position));
CREATE INDEX lessons_class_subject_date_idx ON lessons(class_subject_id,lesson_date);
CREATE TABLE lesson_materials (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, title TEXT NOT NULL, url TEXT NOT NULL, position SMALLINT NOT NULL DEFAULT 1, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(lesson_id,position));
CREATE TABLE grade_items (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, title TEXT NOT NULL, kind TEXT NOT NULL, grading_scale TEXT NOT NULL DEFAULT 'five_point', max_score SMALLINT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE INDEX grade_items_lesson_id_idx ON grade_items(lesson_id);
CREATE TABLE scores (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), grade_item_id UUID NOT NULL REFERENCES grade_items(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, numeric_value SMALLINT, text_value TEXT, teacher_comment TEXT, created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(grade_item_id,user_id));
CREATE INDEX scores_user_id_idx ON scores(user_id);
CREATE TABLE absences (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, recorded_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(lesson_id,user_id));
CREATE INDEX absences_user_id_idx ON absences(user_id);
CREATE TABLE quarter_grades (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), quarter_id UUID NOT NULL REFERENCES academic_year_quarters(id) ON DELETE RESTRICT, class_subject_id UUID NOT NULL REFERENCES class_subjects(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, grading_scale TEXT NOT NULL, max_score SMALLINT, numeric_value SMALLINT, text_value TEXT, teacher_comment TEXT, created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(quarter_id,class_subject_id,user_id));
CREATE INDEX quarter_grades_user_id_idx ON quarter_grades(user_id);

-- +goose Down
DROP TABLE quarter_grades;
DROP TABLE absences;
DROP TABLE scores;
DROP TABLE grade_items;
DROP TABLE lesson_materials;
DROP TABLE lessons;
DROP TABLE class_students;
DROP TABLE class_subjects;
DROP TABLE subjects;
DROP TABLE classes;
DROP TABLE academic_year_quarters;
DROP TABLE academic_years;
DROP TABLE sessions;
DROP TABLE user_permissions;
DROP TABLE role_permissions;
DROP TABLE user_roles;
DROP TABLE permissions;
DROP TABLE roles;
DROP TABLE users;