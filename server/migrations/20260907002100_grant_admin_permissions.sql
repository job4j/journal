-- +goose Up
INSERT INTO role_permissions(role_id,permission_id) SELECT r.id,p.id FROM roles r CROSS JOIN permissions p WHERE r.code='admin';
-- +goose Down
DELETE FROM role_permissions WHERE role_id=(SELECT id FROM roles WHERE code='admin');