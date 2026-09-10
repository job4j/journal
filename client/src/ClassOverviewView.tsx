import { useEffect, useState } from 'react'
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
  const [section, setSection] = useState<Section>('summary')
  const [journal, setJournal] = useState<ClassSubjectRecord | null>(null)
  const [students, setStudents] = useState<ClassStudent[]>([])
  const [subjects, setSubjects] = useState<ClassSubjectRecord[]>([])
  const [grades, setGrades] = useState<QuarterGrade[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    setLoading(true)
    setError('')
    Promise.all([listClassStudents(item.id), listClassSubjects(item.id)])
      .then(async ([members, assignments]) => {
        const gradeLists = await Promise.all(
          assignments.flatMap((assignment) =>
            (item.quarters ?? []).map((quarter) =>
              listQuarterGrades(assignment.id, quarter.id),
            ),
          ),
        )
        if (active) {
          setStudents(members)
          setSubjects(assignments)
          setGrades(gradeLists.flat())
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
  }, [item.id, item.quarters, journal, section])

  if (journal) {
    return (
      <TeacherLessonsView
        assignment={journal}
        quarters={item.quarters ?? []}
        onBack={() => setJournal(null)}
      />
    )
  }
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
                    const values = (item.quarters ?? []).map((quarter) => ({
                      quarter,
                      grade: grades.find((grade) =>
                        grade.classSubjectId === assignment.id
                        && grade.studentId === member.student.id
                        && grade.quarterId === quarter.id),
                    }))
                    return (
                      <td key={assignment.id}>
                        {values.length ? values.map(({ quarter, grade }) => (
                          <span className="summary-grade" key={quarter.id}>
                            <small>{quarter.number} период</small>
                            <strong>{grade ? gradeLabel(grade) : '—'}</strong>
                          </span>
                        )) : <span className="summary-empty">Нет периодов</span>}
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