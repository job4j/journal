-- +goose Up
WITH password_hash (value) AS (
    VALUES (
        '$argon2id$v=19$m=65536,t=3,p=2$D/i9w/OKhPXHL85ALOUmcw$'
            || 'engTax6n6lAMRsYLrVbzIKfEUg/A1R3pyCD/v/NVUvE'
    )
), seed (id, login, name) AS (
    VALUES
        ('10000000-0000-4000-8000-000000000001'::UUID, 'tatyana.bulgakova', 'Татьяна Булгакова'),
        ('20000000-0000-4000-8000-000000000001'::UUID, 'arsentev.matvey', 'Арсентьев Матвей'),
        ('20000000-0000-4000-8000-000000000002'::UUID, 'bulgakova.darina', 'Булгакова Дарина'),
        ('20000000-0000-4000-8000-000000000003'::UUID, 'bokova.zlata', 'Бокова Злата'),
        ('20000000-0000-4000-8000-000000000004'::UUID, 'volkova.zoya', 'Волкова Зоя'),
        ('20000000-0000-4000-8000-000000000005'::UUID, 'gertsik.andrey', 'Герцик Андрей'),
        (
            '20000000-0000-4000-8000-000000000006'::UUID,
            'direktorenko.arseniy',
            'Директоренко Арсений'
        ),
        ('20000000-0000-4000-8000-000000000007'::UUID, 'krivopishina.vera', 'Кривопишина Вера'),
        ('20000000-0000-4000-8000-000000000008'::UUID, 'kovtun.maria', 'Ковтун Мария'),
        ('20000000-0000-4000-8000-000000000009'::UUID, 'pomogaeva.inna', 'Помогаева Инна'),
        (
            '20000000-0000-4000-8000-000000000010'::UUID,
            'strumenskoy.vladislav',
            'Струменской Владислав'
        ),
        ('20000000-0000-4000-8000-000000000011'::UUID, 'shkurko.alexander', 'Шкурко Александр'),
        (
            '25000000-0000-4000-8000-000000000001'::UUID,
            'margatskaya.zhanna',
            'Маргацкая Жанна Александровна'
        ),
        (
            '25000000-0000-4000-8000-000000000002'::UUID,
            'mural.natalia',
            'Мураль Наталья Павловна'
        ),
        (
            '25000000-0000-4000-8000-000000000003'::UUID,
            'burakova.maria',
            'Буракова Мария Сергеевна'
        ),
        (
            '25000000-0000-4000-8000-000000000004'::UUID,
            'gerasina.yulia',
            'Герасина Юлия Андреевна'
        ),
        (
            '25000000-0000-4000-8000-000000000005'::UUID,
            'bogolyubskaya.darya',
            'Боголюбская Дарья Дмитриевна'
        ),
        ('30000000-0000-4000-8000-000000000001'::UUID, 'arsenteva.svetlana', 'Арсентьева Светлана'),
        ('30000000-0000-4000-8000-000000000002'::UUID, 'arsentev.petr', 'Арсентьев Петр'),
        ('30000000-0000-4000-8000-000000000003'::UUID, 'bokova.ludmila', 'Бокова Людмила'),
        ('30000000-0000-4000-8000-000000000004'::UUID, 'volkova.elena', 'Волкова Елена'),
        ('30000000-0000-4000-8000-000000000005'::UUID, 'gertsik.galina', 'Герцик Галина'),
        ('30000000-0000-4000-8000-000000000006'::UUID, 'pershina.alexandra', 'Першина Александра'),
        ('30000000-0000-4000-8000-000000000007'::UUID, 'krivopishina.olga', 'Кривопишина Ольга'),
        ('30000000-0000-4000-8000-000000000008'::UUID, 'kovtun.svetlana', 'Ковтун Светлана'),
        ('30000000-0000-4000-8000-000000000009'::UUID, 'pomogaeva.natalia', 'Помогаева Наталья'),
        ('30000000-0000-4000-8000-000000000010'::UUID, 'strumenskaya.irina', 'Струменская Ирина'),
        ('30000000-0000-4000-8000-000000000011'::UUID, 'shkurko.anna', 'Шкурко Анна')
)
INSERT INTO users (id, login, password_hash, name, status)
SELECT seed.id, seed.login, password_hash.value, seed.name, 'active'
FROM seed
CROSS JOIN password_hash;

WITH assignments (login, role_code) AS (
    VALUES
        ('tatyana.bulgakova', 'admin'),
        ('tatyana.bulgakova', 'parent'),
        ('arsentev.matvey', 'student'),
        ('bulgakova.darina', 'student'),
        ('bokova.zlata', 'student'),
        ('volkova.zoya', 'student'),
        ('gertsik.andrey', 'student'),
        ('direktorenko.arseniy', 'student'),
        ('krivopishina.vera', 'student'),
        ('kovtun.maria', 'student'),
        ('pomogaeva.inna', 'student'),
        ('strumenskoy.vladislav', 'student'),
        ('shkurko.alexander', 'student'),
        ('margatskaya.zhanna', 'teacher'),
        ('mural.natalia', 'teacher'),
        ('burakova.maria', 'teacher'),
        ('gerasina.yulia', 'teacher'),
        ('bogolyubskaya.darya', 'teacher'),
        ('arsenteva.svetlana', 'parent'),
        ('arsentev.petr', 'parent'),
        ('bokova.ludmila', 'parent'),
        ('volkova.elena', 'parent'),
        ('gertsik.galina', 'parent'),
        ('pershina.alexandra', 'parent'),
        ('krivopishina.olga', 'parent'),
        ('kovtun.svetlana', 'parent'),
        ('pomogaeva.natalia', 'parent'),
        ('strumenskaya.irina', 'parent'),
        ('shkurko.anna', 'parent')
)
INSERT INTO user_roles (user_id, role_id)
SELECT users.id, roles.id
FROM assignments
JOIN users ON users.login = assignments.login
JOIN roles ON roles.code = assignments.role_code;

-- +goose Down
DELETE FROM users
WHERE id = '10000000-0000-4000-8000-000000000001'
    OR id::TEXT LIKE '20000000-0000-4000-8000-%'
    OR id::TEXT LIKE '25000000-0000-4000-8000-%'
    OR id::TEXT LIKE '30000000-0000-4000-8000-%';
