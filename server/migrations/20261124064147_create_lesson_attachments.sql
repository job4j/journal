-- +goose Up
CREATE TABLE lesson_attachments (
    id UUID PRIMARY KEY,
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE RESTRICT,
    object_key TEXT NOT NULL UNIQUE,
    file_name TEXT NOT NULL CHECK (length(file_name) BETWEEN 1 AND 255),
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes BETWEEN 1 AND 20971520),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX lesson_attachments_lesson_idx ON lesson_attachments(lesson_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE lesson_attachments;
