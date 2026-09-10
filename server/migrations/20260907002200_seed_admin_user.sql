-- +goose Up
WITH admin_user AS (INSERT INTO users(login,email,phone,password_hash,name,status) VALUES ('parsentev','parsentev@yandex.ru',NULL,'$argon2id$v=19$m=65536,t=3,p=2$gCnDHRz1hL5HdvQZHk6aGw$w7RWGkpSszSaYxkAwYjb0YalihF6mrYxHrP93ynF5g0','Администратор','active') RETURNING id)
INSERT INTO user_roles(user_id,role_id) SELECT admin_user.id,roles.id FROM admin_user CROSS JOIN roles WHERE roles.code='admin';
-- +goose Down
DELETE FROM users WHERE login='parsentev';