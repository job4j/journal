-- +goose Up
INSERT INTO classes (id, academic_year_id, name, grade_level)
VALUES (
    '60000000-0000-4000-8000-000000000001',
    '50000000-0000-4000-8000-000000000001',
    '7 класс',
    7
);

-- +goose Down
DELETE FROM classes
WHERE id = '60000000-0000-4000-8000-000000000001';
