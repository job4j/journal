-- +goose Up
WITH admin_user AS (
    INSERT INTO users (
        email,
        password_hash,
        first_name,
        last_name,
        status
    ) VALUES (
        'parsentev@yandex.ru',
        '$argon2id$v=19$m=65536,t=3,p=2$gCnDHRz1hL5HdvQZHk6aGw$w7RWGkpSszSaYxkAwYjb0YalihF6mrYxHrP93ynF5g0',
        'parsentev',
        'Администратор',
        'active'
    )
    RETURNING id
)
INSERT INTO user_roles (user_id, role_id)
SELECT admin_user.id, (
    SELECT roles.id
    FROM roles
    WHERE roles.code = 'admin'
)
FROM admin_user;

-- +goose Down
DELETE FROM users
WHERE lower(email) = lower('parsentev@yandex.ru');
