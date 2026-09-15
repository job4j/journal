-- +goose Up
WITH admin_user AS (
    INSERT INTO users (id, login, password_hash, name, status)
    VALUES (
        '10000000-0000-4000-8000-000000000002'::UUID,
        'direktorenko.alexandra',
        '$argon2id$v=19$m=65536,t=3,p=2$D/i9w/OKhPXHL85ALOUmcw$'
            || 'engTax6n6lAMRsYLrVbzIKfEUg/A1R3pyCD/v/NVUvE',
        'Директоренко Александра',
        'active'
    )
    RETURNING id
)
INSERT INTO user_roles (user_id, role_id)
SELECT admin_user.id, roles.id
FROM admin_user
CROSS JOIN roles
WHERE roles.code = 'admin';

-- +goose Down
DELETE FROM users
WHERE id = '10000000-0000-4000-8000-000000000002';
