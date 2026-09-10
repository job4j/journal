-- +goose Up
INSERT INTO roles(code,name) VALUES ('admin','Администратор'),('teacher','Учитель'),('parent','Родитель'),('student','Ученик');
-- +goose Down
DELETE FROM roles WHERE code IN ('admin','teacher','parent','student');