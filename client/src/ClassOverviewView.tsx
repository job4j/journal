import { useEffect, useState } from 'react'
import { AcademicYear } from './academicYears'
import { ClassRecord, ClassStudent, listClassStudents } from './classes'
import { ClassSubjectRecord, listClassSubjects } from './classSubjects'
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

export default function ClassOverviewView({ item, year, onBack, onCountChange }: Props) {
  const [section, setSection] = useState<Section>('summary')
  const [journal, setJournal] = useState<ClassSubjectRecord | null>(null)
  const [students, setStudents] = useState<ClassStudent[]>([])
  const [subjects, setSubjects] = useState<ClassSubjectRecord[]>([])
  const [refresh, setRefresh] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    setLoading(true)
    setError('')
    Promise.all([listClassStudents(item.id), listClassSubjects(item.id)])
      .then(([members, assignments]) => {
        if (!active) return
        setStudents(members)
        setSubjects(assignments)
        const assignmentID = location.hash.match(/\/subjects\/([0-9a-f-]+)$/i)?.[1]
        if (assignmentID) {
          setJournal(assignments.find((assignment) => assignment.id === assignmentID) ?? null)
        }
      })
      .catch((cause) => {
        if (active) {
          setError(cause instanceof Error ? cause.message : 'Не удалось загрузить класс')
        }
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [item.id, refresh])

  function showSummary() {
    location.hash = `classes/${item.id}`
    setJournal(null)
    setSection('summary')
    setRefresh((current) => current + 1)
  }

  function openJournal(assignment: ClassSubjectRecord) {
    location.hash = `classes/${item.id}/subjects/${assignment.id}`
    setJournal(assignment)
  }

  if (section === 'roster') {
    return (
      <ClassRosterView
        item={item}
        year={year}
        onBack={showSummary}
        onCountChange={onCountChange}
      />
    )
  }
  if (section === 'subjects') {
    return <ClassSubjectsView item={item} onBack={showSummary} />
  }
  if (journal) {
    return (
      <TeacherLessonsView
        assignment={journal}
        className={item.name}
        quarters={item.quarters ?? []}
        onBack={showSummary}
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
          <h2>Класс {item.name}</h2>
        </div>
        <div className="heading-actions">
          <button
            className="primary-action"
            type="button"
            onClick={() => setSection('roster')}
          >
            + Добавить ученика
          </button>
          <button
            className="primary-action"
            type="button"
            onClick={() => setSection('subjects')}
          >
            + Добавить предмет
          </button>
        </div>
      </div>
      {error && <p className="content-error" role="alert">{error}</p>}
      {loading ? (
        <div className="empty-card">Загружаем класс…</div>
      ) : (
        <div className="class-overview-columns">
          <section className="class-overview-card" aria-labelledby="class-students-title">
            <div className="class-overview-heading">
              <h3 id="class-students-title">Ученики</h3>
              <span className="count-badge">{students.length}</span>
            </div>
            {students.length ? (
              <ul className="class-overview-list">
                {students.map((member) => (
                  <li key={member.student.id}>
                    <strong>{member.student.name}</strong>
                    {member.leftOn && <span>Выбыл {member.leftOn}</span>}
                  </li>
                ))}
              </ul>
            ) : <p className="class-overview-empty">В классе пока нет учеников</p>}
          </section>

          <section className="class-overview-card" aria-labelledby="class-subjects-title">
            <div className="class-overview-heading">
              <h3 id="class-subjects-title">Предметы</h3>
              <span className="count-badge">{subjects.length}</span>
            </div>
            {subjects.length ? (
              <ul className="class-overview-list class-subject-list">
                {subjects.map((assignment) => (
                  <li key={assignment.id}>
                    <button type="button" onClick={() => openJournal(assignment)}>
                      <strong>{assignment.subject.name}</strong>
                      <span>{assignment.responsibleTeacher.name}</span>
                    </button>
                  </li>
                ))}
              </ul>
            ) : <p className="class-overview-empty">Классу пока не назначены предметы</p>}
          </section>
        </div>
      )}
    </div>
  )
}