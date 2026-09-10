-- +goose Up
CREATE TABLE roles (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT NOT NULL UNIQUE, name TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
-- +goose Down
DROP TABLE roles;