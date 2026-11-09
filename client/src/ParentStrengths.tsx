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

function formatNumber(value: number) {
  const [integer, fraction = ''] = value.toString().split('.')
  return `${integer},${fraction.padEnd(2, '0').slice(0, 2)}`
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
    .sort((a, b) => a.name.localeCompare(b.name, 'ru'))
  const overallAverage = strengths.length ? strengths.reduce((sum, subject) => sum + subject.average, 0) / strengths.length : null
  return <section className="parent-strengths" aria-label={`Средний балл — ${student.name}`}>
    <header><h3>Средний балл{!loading && !error && overallAverage != null ? `: ${formatNumber(overallAverage)}` : ''}</h3><p>По оценкам за уроки за текущий учебный год</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем средний балл…</p>
      : error ? <p className="content-error">{error}</p>
        : strengths.length ? <div className="strength-chart-layout">
          <div className="strength-bar-chart" role="img" aria-label={`Средние баллы по оценкам за уроки за текущий учебный год. ${strengths.map(subject => `${subject.name}: ${formatNumber(subject.average)} из 5, ${gradeCount(subject.count)}`).join('; ')}`}>
            <div className="strength-bar-rows" aria-hidden="true">
              {strengths.map(subject => <div className="strength-bar-row" key={subject.id} title={`${subject.name}: средний балл ${formatNumber(subject.average)} из 5, ${gradeCount(subject.count)}`}>
                <span className="strength-bar-name">{subject.name}</span>
                <div className="strength-bar-track">
                  <div className="strength-bar" style={{ width: `${subject.average / 5 * 100}%`, backgroundColor: subject.color }} />
                </div>
                <strong className="strength-bar-value">{formatNumber(subject.average)}</strong>
              </div>)}
            </div>
          </div>
          <p className="strength-chart-note">Среднее арифметическое оценок за работы на уроках по пятибалльной шкале.</p>
        </div> : <p className="parent-weekly-message">Пока нет оценок за уроки для расчёта.</p>}
  </section>
}
