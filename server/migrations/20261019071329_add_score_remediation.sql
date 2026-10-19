-- +goose Up
ALTER TABLE scores ADD COLUMN is_remediated boolean NOT NULL DEFAULT false;
ALTER TABLE scores ADD CONSTRAINT scores_remediation_comment_required
    CHECK (NOT is_remediated OR (teacher_comment IS NOT NULL AND teacher_comment ~ '[^[:space:]]'));

-- +goose Down
ALTER TABLE scores DROP CONSTRAINT scores_remediation_comment_required;
ALTER TABLE scores DROP COLUMN is_remediated;
