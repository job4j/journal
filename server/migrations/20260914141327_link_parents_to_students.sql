-- +goose Up
WITH permission_codes (code) AS (
    VALUES ('can_view_user'), ('can_view_journal')
), students AS (
    SELECT id
    FROM users
    WHERE id::TEXT LIKE '20000000-0000-4000-8000-%'
)
INSERT INTO permissions (code, value, description)
SELECT
    permission_codes.code,
    students.id::TEXT,
    'Доступ родителя к данным ученика'
FROM permission_codes
CROSS JOIN students;

WITH family (parent_login, student_login) AS (
    VALUES
        ('arsenteva.svetlana', 'arsentev.matvey'),
        ('arsentev.petr', 'arsentev.matvey'),
        ('tatyana.bulgakova', 'bulgakova.darina'),
        ('bokova.ludmila', 'bokova.zlata'),
        ('volkova.elena', 'volkova.zoya'),
        ('gertsik.galina', 'gertsik.andrey'),
        ('pershina.alexandra', 'direktorenko.arseniy'),
        ('krivopishina.olga', 'krivopishina.vera'),
        ('kovtun.svetlana', 'kovtun.maria'),
        ('pomogaeva.natalia', 'pomogaeva.inna'),
        ('strumenskaya.irina', 'strumenskoy.vladislav'),
        ('shkurko.anna', 'shkurko.alexander')
), permission_codes (code) AS (
    VALUES ('can_view_user'), ('can_view_journal')
)
INSERT INTO user_permissions (user_id, permission_id)
SELECT parents.id, permissions.id
FROM family
JOIN users AS parents ON parents.login = family.parent_login
JOIN users AS students ON students.login = family.student_login
CROSS JOIN permission_codes
JOIN permissions
    ON permissions.code = permission_codes.code
    AND permissions.value = students.id::TEXT;

-- +goose Down
DELETE FROM permissions
WHERE description = 'Доступ родителя к данным ученика'
    AND value IN (
        SELECT id::TEXT
        FROM users
        WHERE id::TEXT LIKE '20000000-0000-4000-8000-%'
    );
