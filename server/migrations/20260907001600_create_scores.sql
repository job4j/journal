-- +goose Up
CREATE TABLE scores (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), grade_item_id UUID NOT NULL REFERENCES grade_items(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, numeric_value SMALLINT, text_value TEXT, teacher_comment TEXT, created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(grade_item_id,user_id));
CREATE INDEX scores_user_id_idx ON scores(user_id);
-- +goose Down
DROP TABLE scores;