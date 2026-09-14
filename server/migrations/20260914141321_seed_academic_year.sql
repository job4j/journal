-- +goose Up
INSERT INTO academic_years (id, name, starts_on, ends_on, status)
VALUES (
    '50000000-0000-4000-8000-000000000001',
    '2026/2027',
    DATE '2026-09-01',
    DATE '2027-05-31',
    'active'
);

-- +goose Down
DELETE FROM academic_years
WHERE id = '50000000-0000-4000-8000-000000000001';
