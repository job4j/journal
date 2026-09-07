-- +goose Up
INSERT INTO roles (code, name) VALUES
    ('admin', 'Администратор'),
    ('teacher', 'Учитель'),
    ('parent', 'Родитель'),
    ('student', 'Ученик');

INSERT INTO permissions (code, value, description) VALUES
    ('can_create_user', NULL, 'Создание пользователей'),
    ('can_view_user', NULL, 'Просмотр всех пользователей'),
    ('can_update_user', NULL, 'Изменение пользователей'),
    ('can_manage_academic_year', NULL, 'Управление учебными годами'),
    ('can_create_class', NULL, 'Создание классов'),
    ('can_view_class', NULL, 'Просмотр всех классов'),
    ('can_update_class', NULL, 'Изменение классов'),
    ('can_manage_class_students', NULL, 'Управление составом класса'),
    ('can_create_subject', NULL, 'Создание предметов'),
    ('can_view_subject', NULL, 'Просмотр всех предметов'),
    ('can_update_subject', NULL, 'Изменение предметов'),
    ('can_manage_class_subjects', NULL, 'Назначение предметов и учителей'),
    ('can_create_lesson', NULL, 'Создание уроков'),
    ('can_view_lesson', NULL, 'Просмотр всех уроков'),
    ('can_update_lesson', NULL, 'Изменение уроков'),
    ('can_create_score', NULL, 'Создание оценок'),
    ('can_view_score', NULL, 'Просмотр всех оценок'),
    ('can_update_score', NULL, 'Изменение оценок'),
    ('can_view_journal', NULL, 'Просмотр всего журнала');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'admin';

-- Permissions with a concrete resource identifier in value are created when
-- access is assigned and linked directly to a user through user_permissions.

-- +goose Down
DELETE FROM roles
WHERE code IN ('admin', 'teacher', 'parent', 'student');

DELETE FROM permissions
WHERE code IN (
    'can_create_user',
    'can_view_user',
    'can_update_user',
    'can_manage_academic_year',
    'can_create_class',
    'can_view_class',
    'can_update_class',
    'can_manage_class_students',
    'can_create_subject',
    'can_view_subject',
    'can_update_subject',
    'can_manage_class_subjects',
    'can_create_lesson',
    'can_view_lesson',
    'can_update_lesson',
    'can_create_score',
    'can_view_score',
    'can_update_score',
    'can_view_journal'
);
