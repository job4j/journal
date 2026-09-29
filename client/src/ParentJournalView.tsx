import { useEffect, useState } from 'react'
import { ClassSubjectRecord } from './classSubjects'
import { ParentJournalSubject, ParentStudent, ParentStudentPeriod, getParentStudentJournal, listCurrentParentStudents, listParentStudentPeriods } from './parentJournal'

function formatLessonDate(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', weekday: 'short' }).format(new Date(`${value}T00:00:00`))
}

function scoreValue(score: { numericValue?: number | null; textValue?: string | null } | undefined) {
  if (!score) return '—'
  return score.numericValue ?? (score.textValue === 'pass' ? 'Зачёт' : 'Незачёт')
}

export default function ParentJournalView() {
  const [students, setStudents] = useState<ParentStudent[]>([])
  const [student, setStudent] = useState<ParentStudent | null>(null)
  const [periods, setPeriods] = useState<ParentStudentPeriod[]>([])
  const [period, setPeriod] = useState<ParentStudentPeriod | null>(null)
  const [subject, setSubject] = useState<ClassSubjectRecord | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    listCurrentParentStudents().then(setStudents).catch(cause => setError(cause instanceof Error ? cause.message : 'Не удалось загрузить учеников')).finally(() => setLoading(false))
  }, [])

  async function openStudent(value: ParentStudent) {
    setStudent(value)
    setPeriod(null)
    setSubject(null)
    setLoading(true)
    setError('')
    try {
      setPeriods(await listParentStudentPeriods(value.id))
    } catch (cause) {
      setPeriods([])
      setError(cause instanceof Error ? cause.message : 'Доступ к журналу отозван')
    } finally {
      setLoading(false)
    }
  }

  if (subject && period && student) return <ParentSubjectJournal student={student} period={period} subject={subject} onBack={() => setSubject(null)} />

  if (period && student) return <div className="content-stack">
    <button className="back-button" onClick={() => setPeriod(null)}>← К учебным периодам</button>
    <div className="content-heading"><div><p className="content-kicker">{student.name}</p><h2>{period.class.name} · {period.academicYear.name}</h2></div></div>
    <div className="parent-option-list">{period.subjects.map(item => <button className="parent-option-card" aria-label={item.subject.name} key={item.id} onClick={() => setSubject(item)}><span><strong>{item.subject.name}</strong><small>{item.responsibleTeacher.name}</small></span><span aria-hidden="true">›</span></button>)}{!period.subjects.length && <div className="empty-card">В классе пока нет предметов</div>}</div>
  </div>

  if (student) return <div className="content-stack">
    <button className="back-button" onClick={() => setStudent(null)}>← К ученикам</button>
    <div className="content-heading"><div><p className="content-kicker">История обучения</p><h2>{student.name}</h2></div></div>
    {error && <p className="content-error" role="alert">{error}</p>}
    {loading ? <div className="empty-card">Загружаем историю…</div> : periods.length ? <div className="parent-option-list">{periods.map(item => <button className="parent-option-card" aria-label={item.class.name} key={`${item.academicYear.id}-${item.class.id}`} onClick={() => setPeriod(item)}><span><strong>{item.class.name}</strong><small>{item.academicYear.name} · {item.subjects.length} предметов</small></span><span aria-hidden="true">›</span></button>)}</div> : <div className="empty-card">История обучения пока пуста</div>}
  </div>

  return <div className="content-stack">
    <div className="content-heading"><div><p className="content-kicker">Родительский кабинет</p><h2>Ученики</h2></div><span className="count-badge">{students.length}</span></div>
    {error && <p className="content-error" role="alert">{error}</p>}
    {loading ? <div className="empty-card">Загружаем учеников…</div> : students.length ? <div className="parent-option-list">{students.map(item => <button className="parent-option-card" aria-label={item.name} key={item.id} onClick={() => void openStudent(item)}><span><strong>{item.name}</strong><small>Открыть журнал</small></span><span aria-hidden="true">›</span></button>)}</div> : <div className="empty-card">Нет доступных учеников</div>}
  </div>
}

function ParentSubjectJournal({ student, period, subject, onBack }: { student: ParentStudent; period: ParentStudentPeriod; subject: ClassSubjectRecord; onBack: () => void }) {
  const [data, setData] = useState<ParentJournalSubject | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getParentStudentJournal(student.id, period.academicYear.id).then(journal => setData(journal.subjects.find(item => item.classSubject.id === subject.id) ?? null)).catch(cause => setError(cause instanceof Error ? cause.message : 'Не удалось загрузить табель')).finally(() => setLoading(false))
  }, [student.id, period.academicYear.id, subject.id])

  return <div className="content-stack">
    <button className="back-button" onClick={onBack}>← К предметам</button>
    <div><p className="content-kicker">{student.name} · {period.class.name}</p><h2>{subject.subject.name}</h2></div>
    {error && <p className="content-error" role="alert">{error}</p>}
    {loading ? <div className="empty-card">Загружаем табель…</div> : <>
      {data?.quarterGrades?.length ? <div className="quarter-summary">{data.quarterGrades.map(item => <span key={item.id}>{period.academicYear.quarters.find(q => q.id === item.quarterId)?.name ?? 'Период'}: <strong>{scoreValue(item)}</strong></span>)}</div> : null}
      {data?.lessons.length ? <div className="parent-lesson-list">{data.lessons.map(({ lesson }) => <article className="parent-lesson-card" key={lesson.id}>
        <header><time dateTime={lesson.lessonDate}>{formatLessonDate(lesson.lessonDate)}</time>{lesson.absences.length ? <span className="absence-mark">Пропуск</span> : <span className="attendance-mark">Был</span>}</header>
        {lesson.topic && <div className="parent-lesson-section"><span>Тема</span><strong>{lesson.topic}</strong></div>}
        <div className="parent-lesson-section"><span>Домашнее задание</span><strong>{lesson.homework ?? 'Не задано'}</strong></div>
        {lesson.gradeItems.length ? <div className="parent-lesson-section"><span>Работы и оценки</span><div className="parent-score-list">{lesson.gradeItems.map(item => { const score = item.scores[0]; return <div className="parent-score-card" key={item.id}><span>{item.title}</span><strong>{scoreValue(score)}</strong>{score?.teacherComment && <small>{score.teacherComment}</small>}</div> })}</div></div> : null}
        {lesson.materials.length ? <div className="parent-lesson-section"><span>Материалы</span><div className="parent-materials">{lesson.materials.map(item => <a href={item.url} key={item.id}>{item.title}</a>)}</div></div> : null}
      </article>)}</div> : <div className="empty-card">По предмету пока нет уроков</div>}
    </>}
  </div>
}
