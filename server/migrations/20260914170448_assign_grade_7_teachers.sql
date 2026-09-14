-- +goose Up
WITH assignments (id, subject_code, teacher_login) AS (
    VALUES
        ('65000000-0000-4000-8000-000000000001'::UUID, 'history', 'margatskaya.zhanna'),
        ('65000000-0000-4000-8000-000000000002'::UUID, 'botany', 'mural.natalia'),
        ('65000000-0000-4000-8000-000000000003'::UUID, 'geography', 'mural.natalia'),
        ('65000000-0000-4000-8000-000000000004'::UUID, 'algebra', 'burakova.maria'),
        ('65000000-0000-4000-8000-000000000005'::UUID, 'geometry', 'burakova.maria'),
        ('65000000-0000-4000-8000-000000000006'::UUID, 'physics', 'burakova.maria'),
        ('65000000-0000-4000-8000-000000000007'::UUID, 'informatics', 'burakova.maria'),
        ('65000000-0000-4000-8000-000000000008'::UUID, 'russian_language', 'gerasina.yulia'),
        ('65000000-0000-4000-8000-000000000009'::UUID, 'literature', 'gerasina.yulia'),
        ('65000000-0000-4000-8000-000000000010'::UUID, 'english_language', 'bogolyubskaya.darya')
)
INSERT INTO class_subjects (
    id,
    class_id,
    subject_id,
    responsible_teacher_id
)
SELECT
    assignments.id,
    '60000000-0000-4000-8000-000000000001',
    subjects.id,
    users.id
FROM assignments
JOIN subjects ON subjects.code = assignments.subject_code
JOIN users ON users.login = assignments.teacher_login;

-- +goose Down
DELETE FROM class_subjects
WHERE id::TEXT LIKE '65000000-0000-4000-8000-%';
