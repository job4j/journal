-- +goose Up
INSERT INTO roles(code,name) VALUES ('admin','Администратор'),('teacher','Учитель'),('parent','Родитель'),('student','Ученик');
INSERT INTO permissions(code,value,description) VALUES
('can_create_user',NULL,'Создание пользователей'),('can_view_user',NULL,'Просмотр всех пользователей'),('can_update_user',NULL,'Изменение пользователей'),
('can_manage_academic_year',NULL,'Управление учебными годами'),('can_create_class',NULL,'Создание классов'),('can_view_class',NULL,'Просмотр всех классов'),
('can_update_class',NULL,'Изменение классов'),('can_manage_class_students',NULL,'Управление составом класса'),('can_create_subject',NULL,'Создание предметов'),
('can_view_subject',NULL,'Просмотр всех предметов'),('can_update_subject',NULL,'Изменение предметов'),('can_manage_class_subjects',NULL,'Назначение предметов и учителей'),
('can_view_class_subject',NULL,'Просмотр всех назначений предметов классам'),('can_create_lesson',NULL,'Создание уроков'),('can_view_lesson',NULL,'Просмотр всех уроков'),
('can_update_lesson',NULL,'Изменение уроков'),('can_create_score',NULL,'Создание оценок'),('can_view_score',NULL,'Просмотр всех оценок'),
('can_update_score',NULL,'Изменение оценок'),('can_view_journal',NULL,'Просмотр всего журнала'),('can_manage_roles',NULL,'Управление ролями и их правами');
INSERT INTO role_permissions(role_id,permission_id) SELECT r.id,p.id FROM roles r CROSS JOIN permissions p WHERE r.code='admin';
WITH admin_user AS (
 INSERT INTO users(login,email,phone,password_hash,name,status) VALUES
 ('parsentev','parsentev@yandex.ru',NULL,'$argon2id$v=19$m=65536,t=3,p=2$gCnDHRz1hL5HdvQZHk6aGw$w7RWGkpSszSaYxkAwYjb0YalihF6mrYxHrP93ynF5g0','Администратор','active') RETURNING id
)
INSERT INTO user_roles(user_id,role_id) SELECT admin_user.id,roles.id FROM admin_user CROSS JOIN roles WHERE roles.code='admin';

-- +goose Down
DELETE FROM users WHERE login='parsentev';
DELETE FROM roles WHERE code IN ('admin','teacher','parent','student');
DELETE FROM permissions;