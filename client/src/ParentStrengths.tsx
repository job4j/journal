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
  const colors = new Map([...subjects.keys()].sort().map((id, index) => [id, `hsl(${(155 + index * 137.508) % 360} 58% 43%)`]))
  const strengths = [...subjects].map(([id, subject]) => ({ id, ...subject, average: subject.sum / subject.count, color: colors.get(id)! }))
    .sort((a, b) => b.average - a.average || b.count - a.count || a.name.localeCompare(b.name, 'ru'))
  const totalAverage = strengths.reduce((sum, subject) => sum + subject.average, 0)
  const totalCount = strengths.reduce((sum, subject) => sum + subject.count, 0)
  const overall = totalCount ? strengths.reduce((sum, subject) => sum + subject.sum, 0) / totalCount : 0
  const formatNumber = (value: number) => new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 1 }).format(value)
  let offset = 0
  const segments = strengths.map(subject => {
    const share = subject.average / totalAverage * 100
    const segment = { ...subject, share, offset }
    offset += share
    return segment
  })

  return <section className="parent-strengths" aria-label={`Сильные стороны — ${student.name}`}>
    <header><h3>Сильные стороны</h3><p>Все предметы · за текущий учебный год</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем успехи…</p>
      : error ? <p className="content-error">{error}</p>
        : segments.length ? <div className="strength-chart-layout">
          <div className="strength-ring" role="img" aria-label={`Средние баллы по предметам. ${segments.map(subject => `${subject.name}: ${formatNumber(subject.average)} из 5, доля ${formatNumber(subject.share)}%`).join('; ')}`}>
            <svg viewBox="0 0 160 160" aria-hidden="true">
              <circle className="strength-ring-track" cx="80" cy="80" r="62" />
              {segments.map(subject => <circle className="strength-ring-value" key={subject.id} cx="80" cy="80" r="62" pathLength="100" stroke={subject.color} strokeDasharray={`${subject.share} ${100 - subject.share}`} strokeDashoffset={-subject.offset} transform="rotate(-90 80 80)"><title>{subject.name}: {formatNumber(subject.average)} из 5</title></circle>)}
            </svg>
            <span aria-hidden="true"><strong>{formatNumber(overall)}</strong><small>средний балл</small></span>
          </div>
          <ul>{segments.map(subject => <li key={subject.id}>
            <span className="strength-color" style={{ backgroundColor: subject.color }} aria-hidden="true" />
            <div className="strength-subject"><strong>{subject.name}</strong><span>{gradeCount(subject.count)} · {formatNumber(subject.share)}%</span></div>
            <strong className="strength-average" aria-label={`${subject.name}: средний балл ${formatNumber(subject.average)} из 5`}>{formatNumber(subject.average)}</strong>
          </li>)}</ul>
          <p className="strength-chart-note">Длина сектора пропорциональна среднему баллу предмета.</p>
        </div> : <p className="parent-weekly-message">Пока недостаточно оценок для расчёта.</p>}
  </section>
}