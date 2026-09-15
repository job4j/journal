-- +goose Up
INSERT INTO permissions (code, value, description)
VALUES (
    'can_view_class',
    '60000000-0000-4000-8000-000000000001',
    'Объектный доступ по назначению учителя'
);

WITH permission_codes (code) AS (
    VALUES
        ('can_view_class_subject'),
        ('can_create_lesson'),
        ('can_view_lesson'),
        ('can_update_lesson'),
        ('can_create_score'),
        ('can_view_score'),
        ('can_update_score')
)
INSERT INTO permissions (code, value, description)
SELECT
    permission_codes.code,
    class_subjects.id::TEXT,
    'Объектный доступ по назначению учителя'
FROM class_subjects
CROSS JOIN permission_codes
WHERE class_subjects.id::TEXT LIKE '65000000-0000-4000-8000-%';

INSERT INTO user_permissions (user_id, permission_id)
SELECT DISTINCT
    class_subjects.responsible_teacher_id,
    permissions.id
FROM class_subjects
JOIN permissions
    ON (
        permissions.code = 'can_view_class'
        AND permissions.value = class_subjects.class_id::TEXT
    )
    OR (
        permissions.code IN (
            'can_view_class_subject',
            'can_create_lesson',
            'can_view_lesson',
            'can_update_lesson',
            'can_create_score',
            'can_view_score',
            'can_update_score'
        )
        AND permissions.value = class_subjects.id::TEXT
    )
WHERE class_subjects.id::TEXT LIKE '65000000-0000-4000-8000-%';

-- +goose Down
DELETE FROM permissions
WHERE description = 'Объектный доступ по назначению учителя'
    AND (
        value = '60000000-0000-4000-8000-000000000001'
        OR value LIKE '65000000-0000-4000-8000-%'
    );
