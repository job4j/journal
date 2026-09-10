-- +goose Up
CREATE TABLE user_permissions (user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(user_id,permission_id));
-- +goose Down
DROP TABLE user_permissions;