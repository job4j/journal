import { useEffect, useMemo, useState } from 'react'
import { AcademicYear } from './academicYears'
import { ClassRecord, ClassStudent, listClassStudents } from './classes'
import { ClassSubjectRecord, listClassSubjects } from './classSubjects'
import { QuarterGrade, listQuarterGrades } from './quarterGrades'
import ClassRosterView from './ClassRosterView'
import ClassSubjectsView from './ClassSubjectsView'
import TeacherLessonsView from './TeacherLessonsView'

interface Props {
  item: ClassRecord
  year?: AcademicYear
  onBack: () => void
  onCountChange: (count: number) => void
}

type Section = 'summary' | 'roster' | 'subjects'

function gradeLabel(grade: QuarterGrade) {
  if (grade.numericValue != null) return String(grade.numericValue)
  if (grade.textValue === 'pass') return 'Зачёт'
  if (grade.textValue === 'fail') return 'Незачёт'
  return '—'
}

export default function ClassOverviewView({ item, year, onBack, onCountChange }: Props) {
  const quarters = useMemo(() => item.quarters ?? [], [item.quarters])
  const [section, setSection] = useState<Section>('summary')
  const [quarterID, setQuarterID] = useState(quarters[0]?.id ?? '')
  const [journal, setJournal] = useState<ClassSubjectRecord | null>(null)
  const [students, setStudents] = useState<ClassStudent[]>([])
  const [subjects, setSubjects] = useState<ClassSubjectRecord[]>([])
  const [grades, setGrades] = useState<QuarterGrade[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setQuarterID((current) =>
      quarters.some((quarter) => quarter.id === current) ? current : quarters[0]?.id ?? '',
    )
  }, [item.id, quarters])

  useEffect(() => {
    if (section !== 'summary' || journal) return undefined
    let active = true
    setLoading(true)
    setError('')
    Promise.all([listClassStudents(item.id), listClassSubjects(item.id)])
      .then(async ([members, assignments]) => {
        const gradeLists = quarterID
          ? await Promise.all(
            assignments.map((assignment) => listQuarterGrades(assignment.id, quarterID)),
          )
          : []
        if (active) {
          setStudents(members)
          setSubjects(assignments)
          setGrades(gradeLists.flat())
          setJournal((current) => current ?? assignments[0] ?? null)
        }
      })
      .catch((cause) => {
        if (active) {
          setError(cause instanceof Error ? cause.message : 'Не удалось загрузить сводку класса')
        }
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [item.id, journal, quarterID, section])

  if (section === 'roster') {
    return (
      <ClassRosterView
        item={item}
        year={year}
        onBack={() => setSection('summary')}
        onCountChange={onCountChange}
      />
    )
  }
  if (section === 'subjects') {
    return <ClassSubjectsView item={item} onBack={() => setSection('summary')} />
  }
  if (journal) {
    return (
      <TeacherLessonsView
        assignment={journal}
        className={item.name}
        quarters={item.quarters ?? []}
        onBack={() => setSection('subjects')}
      />
    )
  }

  return (
    <div className="content-stack">
      <button className="back-button" type="button" onClick={onBack}>
        ← К списку классов
      </button>
      <div className="content-heading">
        <div>
          <p className="content-kicker">{year?.name ?? 'Учебный год'}</p>
          <h2>Сводка класса {item.name}</h2>
        </div>
        <div className="heading-actions">
          <button className="secondary-action" type="button" onClick={() => setSection('roster')}>
            Состав класса
          </button>
          <button className="secondary-action" type="button" onClick={() => setSection('subjects')}>
            Предметы класса
          </button>
        </div>
      </div>
      <div className="lesson-filters">
        {quarters.length ? (
          <label>
            Период
            <select
              aria-label="Период"
              value={quarterID}
              onChange={(event) => setQuarterID(event.target.value)}
            >
              {quarters.map((quarter) => (
                <option value={quarter.id} key={quarter.id}>
                  {quarter.number} период
                </option>
              ))}
            </select>
          </label>
        ) : <span className="summary-empty">Нет периодов</span>}
      </div>
      {error && <p className="content-error" role="alert">{error}</p>}
      {loading ? (
        <div className="empty-card">Загружаем сводку…</div>
      ) : !students.length ? (
        <div className="empty-card">В классе пока нет учеников</div>
      ) : !subjects.length ? (
        <div className="empty-card">Классу пока не назначены предметы</div>
      ) : (
        <div className="table-card table-scroll">
          <table className="class-summary-table">
            <thead>
              <tr>
                <th>Ученик</th>
                {subjects.map((assignment) => (
                  <th key={assignment.id}>
                    <button
                      className="subject-link"
                      type="button"
                      onClick={() => setJournal(assignment)}
                    >
                      {assignment.subject.name}
                    </button>
                    <span>{assignment.responsibleTeacher.name}</span>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {students.map((member) => (
                <tr key={member.student.id}>
                  <th>{member.student.name}</th>
                  {subjects.map((assignment) => {
                    const grade = grades.find((grade) =>
                      grade.classSubjectId === assignment.id
                      && grade.studentId === member.student.id
                      && grade.quarterId === quarterID)
                    return (
                      <td key={assignment.id}>
                        {quarterID ? (
                          <span className="summary-grade">
                            <strong>{grade ? gradeLabel(grade) : '—'}</strong>
                          </span>
                        ) : <span className="summary-empty">Нет периодов</span>}
                      </td>
                    )
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}