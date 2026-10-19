import { useEffect, useState } from 'react'
import ParentRiskGrades from './ParentRiskGrades'
import ParentStrengths from './ParentStrengths'
import { ParentJournal, ParentStudent, getParentStudentJournal, listParentStudentPeriods } from './parentJournal'

function dateString(date: Date) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function currentWeek() {
  const today = new Date()
  const start = new Date(today)
  start.setDate(start.getDate() - (start.getDay() + 6) % 7)
  const end = new Date(start)
  end.setDate(end.getDate() + 6)
  return { from: dateString(start), to: dateString(end), today: dateString(today) }
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long' }).format(new Date(`${value}T00:00:00`))
}

export default function ParentStudentGrades({ student }: { student: ParentStudent }) {
  const [week] = useState(currentWeek)
  const [journals, setJournals] = useState<ParentJournal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    async function load() {
      try {
        const periods = await listParentStudentPeriods(student.id)
        const yearIDs = [...new Set(periods.filter(({ academicYear }) => academicYear.startsOn <= week.to && academicYear.endsOn >= week.from).map(({ academicYear }) => academicYear.id))]
        const data = await Promise.all(yearIDs.map(yearID => getParentStudentJournal(student.id, yearID)))
        if (active) setJournals(data)
      } catch (cause) {
        if (active) setError(cause instanceof Error ? cause.message : 'Не удалось загрузить оценки')
      } finally {
        if (active) setLoading(false)
      }
    }
    void load()
    return () => { active = false }
  }, [student.id, week])

  const subjects = journals.flatMap(journal => journal.subjects.map(subject => ({
    id: subject.classSubject.id,
    name: subject.classSubject.subject.name,
    grades: subject.lessons.flatMap(({ lesson }) => lesson.lessonDate >= week.from && lesson.lessonDate <= week.to
      ? lesson.gradeItems.flatMap(item => item.scores.map(score => ({ score, work: item.title, date: lesson.lessonDate, position: lesson.position })))
      : []).sort((a, b) => a.date.localeCompare(b.date) || a.position - b.position),
  }))).filter(subject => subject.grades.length)

  return <><section className="parent-weekly-grades" aria-label={`Оценки за неделю — ${student.name}`}>
    <header><h3>Оценки за текущую неделю</h3><p>{formatDate(week.from)} — {formatDate(week.to)}</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем оценки…</p>
      : error ? <p className="content-error" role="alert">{error}</p>
        : subjects.length ? subjects.map(subject => <div className="parent-weekly-subject" key={subject.id}>
          <h4>{subject.name}</h4>
          <ul>{subject.grades.map(({ score, work, date }) => <li key={score.id}>
            <div><time dateTime={date}>{formatDate(date)}</time><span>{work}</span></div>
            <strong className="parent-weekly-score">{score.numericValue ?? (score.textValue === 'pass' ? 'Зачёт' : 'Незачёт')}</strong>
            {score.teacherComment && <p>{score.teacherComment}</p>}
          </li>)}</ul>
        </div>) : <p className="parent-weekly-message">За текущую неделю оценок пока нет</p>}
  </section><ParentStrengths student={student} journals={journals} today={week.today} loading={loading} error={error} /><ParentRiskGrades student={student} journals={journals} today={week.today} loading={loading} error={error} /></>
}
