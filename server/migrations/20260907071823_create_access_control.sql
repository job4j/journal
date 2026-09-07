-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email));

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    scope TEXT NOT NULL
        CHECK (scope IN ('all', 'own', 'children', 'assigned', 'granted')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id, scope)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ,
    CHECK (expires_at > created_at)
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_active_expiry_idx ON sessions (expires_at)
    WHERE revoked_at IS NULL;

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users (id) ON DELETE SET NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    middle_name TEXT,
    birth_date DATE,
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE guardian_students (
    guardian_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students (id) ON DELETE CASCADE,
    relationship TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (guardian_user_id, student_id)
);

CREATE TABLE access_grants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    resource_type TEXT NOT NULL,
    resource_id UUID NOT NULL,
    granted_by UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    CHECK (expires_at IS NULL OR expires_at > created_at)
);

CREATE UNIQUE INDEX access_grants_active_unique_idx
    ON access_grants (user_id, permission_id, resource_type, resource_id)
    WHERE revoked_at IS NULL;

CREATE INDEX access_grants_lookup_idx
    ON access_grants (user_id, permission_id, resource_type, resource_id);

INSERT INTO roles (code, name) VALUES
    ('admin', 'Администратор'),
    ('teacher', 'Учитель'),
    ('parent', 'Родитель');

INSERT INTO permissions (code, description) VALUES
    ('users.create', 'Создание пользователей'),
    ('users.view', 'Просмотр пользователей'),
    ('users.update', 'Изменение пользователей'),
    ('students.view', 'Просмотр учеников'),
    ('students.update', 'Изменение учеников'),
    ('lessons.create', 'Создание уроков'),
    ('lessons.update', 'Изменение уроков'),
    ('grades.create', 'Создание оценок'),
    ('grades.update', 'Изменение оценок'),
    ('journal.view', 'Просмотр журнала');

INSERT INTO role_permissions (role_id, permission_id, scope)
SELECT r.id, p.id, 'all'
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'admin';

INSERT INTO role_permissions (role_id, permission_id, scope)
SELECT r.id, p.id, value.scope
FROM (
    VALUES
        ('users.view', 'own'),
        ('students.view', 'children'),
        ('journal.view', 'children')
) AS value(permission_code, scope)
JOIN roles r ON r.code = 'parent'
JOIN permissions p ON p.code = value.permission_code;

INSERT INTO role_permissions (role_id, permission_id, scope)
SELECT r.id, p.id, value.scope
FROM (
    VALUES
        ('users.view', 'own'),
        ('students.view', 'assigned'),
        ('lessons.create', 'assigned'),
        ('lessons.update', 'assigned'),
        ('grades.create', 'assigned'),
        ('grades.update', 'assigned'),
        ('journal.view', 'assigned')
) AS value(permission_code, scope)
JOIN roles r ON r.code = 'teacher'
JOIN permissions p ON p.code = value.permission_code;

-- +goose Down
DROP TABLE IF EXISTS access_grants;
DROP TABLE IF EXISTS guardian_students;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
