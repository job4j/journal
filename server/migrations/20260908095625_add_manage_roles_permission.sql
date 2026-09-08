-- +goose Up
WITH permission AS (
    INSERT INTO permissions (code, value, description)
    VALUES ('can_manage_roles', NULL, 'Управление ролями и их правами')
    RETURNING id
)
INSERT INTO role_permissions (role_id, permission_id)
SELECT roles.id, permission.id
FROM roles
CROSS JOIN permission
WHERE roles.code = 'admin';

-- +goose Down
DELETE FROM permissions
WHERE code = 'can_manage_roles' AND value IS NULL;
