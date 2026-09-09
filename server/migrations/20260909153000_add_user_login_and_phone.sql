-- +goose Up
ALTER TABLE users ADD COLUMN login TEXT;
UPDATE users SET login = lower(email);
ALTER TABLE users ALTER COLUMN login SET NOT NULL;
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
ALTER TABLE users ADD COLUMN phone TEXT;
CREATE UNIQUE INDEX users_login_unique_idx ON users (lower(login));
DROP INDEX users_email_unique_idx;
CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email)) WHERE email IS NOT NULL;

-- +goose Down
DROP INDEX users_email_unique_idx;
CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email));
ALTER TABLE users DROP COLUMN phone;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
DROP INDEX users_login_unique_idx;
ALTER TABLE users DROP COLUMN login;
