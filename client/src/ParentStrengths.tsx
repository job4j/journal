import type { ParentJournal, ParentStudent } from './parentJournal'

interface Props {
  student: ParentStudent
  journals: ParentJournal[]
  today: string
  loading: boolean
  error: string
}

function gradeCount(value: number) {
  const remainder = value % 100
  const word = remainder >= 11 && remainder <= 14 ? 'оценок'
    : value % 10 === 1 ? 'оценка' : value % 10 >= 2 && value % 10 <= 4 ? 'оценки' : 'оценок'
  return `${value} ${word}`
}

export default function ParentStrengths({ student, journals, today, loading, error }: Props) {
  const subjects = new Map<string, { name: string; sum: number; count: number }>()
  for (const journal of journals) {
    if (journal.academicYear.startsOn > today || journal.academicYear.endsOn < today) continue
    for (const subject of journal.subjects) {
      const id = subject.classSubject.subject.id
      for (const { lesson } of subject.lessons) {
        for (const item of lesson.gradeItems) {
          if (item.gradingScale !== 'five_point') continue
          for (const score of item.scores) {
            if (score.numericValue == null || !Number.isFinite(score.numericValue)) continue
            const summary = subjects.get(id) ?? { name: subject.classSubject.subject.name, sum: 0, count: 0 }
            summary.sum += score.numericValue
            summary.count++
            subjects.set(id, summary)
          }
        }
      }
    }
  }
  const strengths = [...subjects].map(([id, subject]) => ({ id, ...subject, average: subject.sum / subject.count }))
    .filter(subject => subject.average >= 4)
    .sort((a, b) => b.average - a.average || b.count - a.count || a.name.localeCompare(b.name, 'ru'))

  return <section className="parent-strengths" aria-label={`Сильные стороны — ${student.name}`}>
    <header><h3>Сильные стороны</h3><p>Средний балл от 4 · за текущий учебный год</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем успехи…</p>
      : error ? <p className="content-error">{error}</p>
        : strengths.length ? <ul>{strengths.map(subject => {
          const average = new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 1 }).format(subject.average)
          return <li key={subject.id}>
            <div className="strength-ring" role="img" aria-label={`${subject.name}: средний балл ${average} из 5, ${gradeCount(subject.count)}`}>
              <svg viewBox="0 0 80 80" aria-hidden="true"><circle className="strength-ring-track" cx="40" cy="40" r="32" /><circle className="strength-ring-value" cx="40" cy="40" r="32" pathLength="100" strokeDasharray={`${subject.average / 5 * 100} 100`} transform="rotate(-90 40 40)" /></svg>
              <span aria-hidden="true"><strong>{average}</strong><small>из 5</small></span>
            </div>
            <div className="strength-subject"><strong>{subject.name}</strong><span>{gradeCount(subject.count)}</span></div>
          </li>
        })}</ul> : <p className="parent-weekly-message">{subjects.size ? 'Здесь появятся предметы со средним баллом от 4.' : 'Пока недостаточно оценок для расчёта.'}</p>}
  </section>
}
