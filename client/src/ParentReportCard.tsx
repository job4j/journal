import type { ParentJournal, ParentStudent } from './parentJournal'
import type { QuarterGrade } from './quarterGrades'

interface Props {
  student: ParentStudent
  journals: ParentJournal[]
  today: string
  loading: boolean
  error: string
}

function gradeValue(grade: QuarterGrade) {
  if (grade.gradingScale === 'pass_fail') return grade.textValue === 'pass' ? 'Зачёт' : 'Незачёт'
  if (grade.numericValue == null) return '—'
  const value = new Intl.NumberFormat('ru-RU').format(grade.numericValue)
  return grade.gradingScale === 'points' && grade.maxScore != null ? `${value} / ${grade.maxScore}` : value
}

export default function ParentReportCard({ student, journals, today, loading, error }: Props) {
  const quarters = new Map<string, { name: string; number: number; startsOn: string; grades: { subject: string; className: string; grade: QuarterGrade }[] }>()
  for (const journal of journals) {
    if (journal.academicYear.startsOn > today || journal.academicYear.endsOn < today) continue
    for (const quarter of journal.academicYear.quarters) {
      const grades = journal.subjects.flatMap(subject => (subject.quarterGrades ?? [])
        .filter(grade => grade.quarterId === quarter.id)
        .map(grade => ({ subject: subject.classSubject.subject.name, className: journal.class.name, grade })))
      if (!grades.length) continue
      const group = quarters.get(quarter.id) ?? { name: quarter.name, number: quarter.number, startsOn: quarter.startsOn, grades: [] }
      group.grades.push(...grades)
      quarters.set(quarter.id, group)
    }
  }
  const groups = [...quarters].sort(([, a], [, b]) => a.startsOn.localeCompare(b.startsOn) || a.number - b.number)
  const multipleClasses = new Set(journals.map(journal => journal.class.id)).size > 1

  if (!loading && !error && !groups.length) return null

  return <section className="parent-report-card" aria-label={`Табель — ${student.name}`}>
    <header><h3>Табель</h3><p>Итоговые оценки за текущий учебный год</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем табель…</p>
      : error ? <p className="content-error">{error}</p>
        : groups.map(([id, quarter]) => <div className="parent-report-quarter" key={id}>
          <h4>{quarter.name}</h4>
          <table aria-label={`Итоги — ${quarter.name}`}>
            <thead><tr><th scope="col">Предмет</th><th scope="col">Оценка</th></tr></thead>
            <tbody>{[...quarter.grades].sort((a, b) => a.subject.localeCompare(b.subject, 'ru') || a.className.localeCompare(b.className, 'ru')).map(({ subject, className, grade }) => <tr key={grade.id}>
              <th scope="row">{subject}{multipleClasses && <small>{className}</small>}</th>
              <td><strong>{gradeValue(grade)}</strong></td>
            </tr>)}</tbody>
          </table>
        </div>)}
  </section>
}
