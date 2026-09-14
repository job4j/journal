-- +goose Up
INSERT INTO class_students (class_id, user_id, enrolled_on)
SELECT
    '60000000-0000-4000-8000-000000000001',
    users.id,
    DATE '2026-09-01'
FROM users
JOIN user_roles ON user_roles.user_id = users.id
JOIN roles ON roles.id = user_roles.role_id
WHERE roles.code = 'student'
    AND users.id::TEXT LIKE '20000000-0000-4000-8000-%';

-- +goose Down
DELETE FROM class_students
WHERE class_id = '60000000-0000-4000-8000-000000000001';
