-- +goose Up
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    value TEXT,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX permissions_global_code_unique_idx
    ON permissions (code)
    WHERE value IS NULL;

CREATE UNIQUE INDEX permissions_object_code_value_unique_idx
    ON permissions (code, value)
    WHERE value IS NOT NULL;
-- +goose Down
DROP TABLE permissions;
