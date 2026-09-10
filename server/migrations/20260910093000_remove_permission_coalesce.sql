-- +goose Up
DROP INDEX permissions_code_value_unique_idx;
CREATE UNIQUE INDEX permissions_code_value_unique_idx
    ON permissions (code, value) NULLS NOT DISTINCT;

-- +goose Down
DROP INDEX permissions_code_value_unique_idx;
CREATE UNIQUE INDEX permissions_code_value_unique_idx
    ON permissions (code, (CASE WHEN value IS NULL THEN '' ELSE value END));