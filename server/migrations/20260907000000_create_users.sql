-- +goose Up
CREATE TABLE users (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), login TEXT NOT NULL, email TEXT, phone TEXT, password_hash TEXT NOT NULL, name TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'active', created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE UNIQUE INDEX users_login_unique_idx ON users (lower(login));
CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email)) WHERE email IS NOT NULL;
-- +goose Down
DROP TABLE users;