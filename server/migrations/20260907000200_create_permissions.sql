-- +goose Up
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    value TEXT,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE NULLS NOT DISTINCT INDEX permissions_code_value_unique_idx
    ON permissions (code, value);
-- +goose Down
DROP TABLE permissions;
