-- +goose Up
CREATE TABLE sessions (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, token_hash TEXT NOT NULL UNIQUE, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), last_used_at TIMESTAMPTZ NOT NULL DEFAULT now(), revoked_at TIMESTAMPTZ);
CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_active_expiry_idx ON sessions(expires_at) WHERE revoked_at IS NULL;
-- +goose Down
DROP TABLE sessions;