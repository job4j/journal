import type { ParentJournal, ParentStudent } from './parentJournal'

interface Props {
  student: ParentStudent
  journals: ParentJournal[]
  today: string
  loading: boolean
  error: string
}

export default function ParentRiskGrades({ student, journals, today, loading, error }: Props) {
  const grades = journals.filter(journal => journal.academicYear.startsOn <= today && journal.academicYear.endsOn >= today)
    .flatMap(journal => journal.subjects.flatMap(subject => {
      const subjectName = subject.classSubject.subject.name
      const lessonGrades = subject.lessons.flatMap(({ lesson }) => lesson.gradeItems
        .filter(item => item.gradingScale === 'five_point')
        .flatMap(item => item.scores.filter(score => !score.isRemediated && (score.numericValue === 2 || score.numericValue === 3))
          .map(score => ({ id: score.id, subjectName, work: item.title, date: lesson.lessonDate, position: lesson.position, value: score.numericValue, comment: score.teacherComment }))))
      const periodGrades = (subject.quarterGrades ?? [])
        .filter(grade => grade.gradingScale === 'five_point' && (grade.numericValue === 2 || grade.numericValue === 3))
        .map(grade => {
          const quarter = journal.academicYear.quarters.find(item => item.id === grade.quarterId)
          return { id: grade.id, subjectName, work: `Итоговая оценка · ${quarter?.name ?? 'Учебный период'}`, date: quarter?.endsOn ?? '', position: 0, value: grade.numericValue, comment: grade.teacherComment }
        })
      return [...lessonGrades, ...periodGrades]
    })).sort((a, b) => b.date.localeCompare(a.date) || b.position - a.position || a.subjectName.localeCompare(b.subjectName, 'ru'))

  if (!loading && !error && !grades.length) return null

  return <section className="parent-risk-grades" aria-label={`Зона риска — ${student.name}`}>
    <header><h3>Зона риска</h3><p>Оценки 2 и 3 за текущий учебный год</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем оценки…</p>
      : error ? <p className="content-error">{error}</p>
        : <ul className="parent-risk-list">{grades.map(grade => <li key={grade.id}>
          <div><strong>{grade.subjectName}</strong>{grade.date && <time dateTime={grade.date}>{new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }).format(new Date(`${grade.date}T00:00:00`))}</time>}<span>{grade.work}</span></div>
          <strong className={`parent-risk-score parent-risk-score-${grade.value}`}>{grade.value}</strong>
          {grade.comment && <p>{grade.comment}</p>}
        </li>)}</ul>}
  </section>
}
