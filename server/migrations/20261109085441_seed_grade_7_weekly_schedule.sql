-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM classes c JOIN academic_years y ON y.id=c.academic_year_id
        WHERE c.id='60000000-0000-4000-8000-000000000001'
          AND y.id='50000000-0000-4000-8000-000000000001'
          AND y.starts_on=DATE '2026-09-01' AND y.ends_on=DATE '2027-05-31'
    ) THEN
        RAISE EXCEPTION 'Expected grade 7 class and 2026/2027 academic year are missing';
    END IF;
    IF (SELECT count(DISTINCT s.code) FROM class_subjects cs JOIN subjects s ON s.id=cs.subject_id
        WHERE cs.class_id='60000000-0000-4000-8000-000000000001'
          AND s.code IN ('physics','algebra','geometry','english_language','informatics',
                         'russian_language','history','geography','literature','botany')) <> 10 THEN
        RAISE EXCEPTION 'Assign all ten core subjects to grade 7 before seeding its schedule';
    END IF;
END;
$$;
-- +goose StatementEnd

-- Only newly inserted weeks are populated; existing weeks are never overwritten.
WITH weeks AS (
    INSERT INTO class_schedule_weeks(id,class_id,week_start)
    SELECT md5('grade7-schedule-' || d::DATE::TEXT)::UUID,
        '60000000-0000-4000-8000-000000000001', d::DATE
    FROM generate_series(DATE '2026-08-31', DATE '2027-05-31', INTERVAL '7 days') d
    ON CONFLICT(class_id,week_start) DO NOTHING
    RETURNING id,week_start
), template(dow,position,code,title,notes) AS (
    VALUES
    (1,1,'physics','',''),(1,2,'algebra','',''),(1,3,'russian_language','',''),
    (1,4,'literature','',''),(1,5,'russian_language','',''),(1,6,NULL,'Рукоделие',''),
    (2,1,'algebra','',''),(2,2,'english_language','',''),(2,3,'history','',''),
    (2,4,'russian_language','',''),(2,5,'literature','',''),(2,6,NULL,'Физкультура','Мальчики'),
    (3,1,'geometry','',''),(3,2,'physics','',''),(3,3,'russian_language','',''),
    (3,4,'literature','',''),(3,5,NULL,'Шахматы / классный час','Классный час — 1 раз в месяц; конкретную неделю уточняет администратор. Шахматы до 14:00.'),
    (4,1,'algebra','',''),(4,2,'informatics','',''),(4,3,'geography','',''),
    (4,4,'geography','',''),(4,5,'botany','',''),(4,6,NULL,'Физкультура','Мальчики'),(4,7,NULL,'Тир',''),
    (5,1,'geometry','',''),(5,2,'english_language','',''),(5,3,'history','',''),
    (5,4,'russian_language','',''),(5,5,NULL,'Труд','Мальчики'),
    (5,5,NULL,'Физкультура','Девочки; начало 14:30 по исправлению на фото, окончание уточнить')
)
INSERT INTO class_schedule_entries(id,week_id,lesson_date,position,class_subject_id,title,notes,starts_at,ends_at)
SELECT md5(w.id::TEXT || '-' || t.dow || '-' || t.position || '-' || coalesce(t.code,t.title))::UUID,
    w.id,w.week_start+t.dow-1,t.position,cs.id,t.title,t.notes,
    CASE WHEN t.dow=5 AND t.title='Физкультура' THEN '14:30'
         ELSE (ARRAY['08:45','09:45','10:45','12:00','13:00','14:00','15:00'])[t.position] END,
    CASE WHEN t.dow=5 AND t.title='Физкультура' THEN NULL
         WHEN t.dow=3 AND t.position=5 THEN '14:00'
         ELSE (ARRAY['09:30','10:30','11:30','12:45','13:45','14:45','15:45'])[t.position] END
FROM weeks w CROSS JOIN template t
LEFT JOIN subjects s ON s.code=t.code
LEFT JOIN class_subjects cs ON cs.subject_id=s.id AND cs.class_id='60000000-0000-4000-8000-000000000001'
WHERE w.week_start+t.dow-1 BETWEEN DATE '2026-09-01' AND DATE '2027-05-31';

-- +goose Down
-- Seeded weeks may have been edited independently. Do not erase those edits.
-- +goose StatementBegin
DO $$
BEGIN
    RAISE EXCEPTION 'Seeded schedule cannot be safely reversed automatically; use a reviewed backup';
END;
$$;
-- +goose StatementEnd
