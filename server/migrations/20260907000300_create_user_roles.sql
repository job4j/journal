-- +goose Up
CREATE TABLE user_roles (user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(user_id,role_id));
-- +goose Down
DROP TABLE user_roles;