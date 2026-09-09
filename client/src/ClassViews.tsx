import { useState } from 'react'

type GradeKind = 'homework' | 'classwork' | 'test'
type Grade = { value: string; kind: GradeKind }

const classes = [
  { id: '7a', name: '7А', year: '2026 / 2027', teacher: 'Анна Сергеевна Иванова', students: 14 },
  { id: '6b', name: '6Б', year: '2026 / 2027', teacher: 'Михаил Олегович Волков', students: 12 },
  { id: '5a', name: '5А', year: '2026 / 2027', teacher: 'Елена Викторовна Смирнова', students: 16 },
  { id: '4a', name: '4А', year: '2026 / 2027', teacher: 'Ольга Андреевна Лебедева', students: 11 },
]

const subjects = [
  { id: 'math', name: 'Математика', teacher: 'Анна Сергеевна Иванова', lessons: 8 },
  { id: 'russian', name: 'Русский язык', teacher: 'Елена Викторовна Смирнова', lessons: 7 },
  { id: 'literature', name: 'Литература', teacher: 'Елена Викторовна Смирнова', lessons: 5 },
  { id: 'history', name: 'История', teacher: 'Михаил Олегович Волков', lessons: 4 },
  { id: 'english', name: 'Английский язык', teacher: 'Ольга Андреевна Лебедева', lessons: 6 },
]

const students = [
  { id: 'anna', name: 'Анна Белова' },
  { id: 'ivan', name: 'Иван Громов' },
  { id: 'maria', name: 'Мария Кузнецова' },
  { id: 'pavel', name: 'Павел Орлов' },
  { id: 'sofia', name: 'София Романова' },
]

const lessons = [
  { id: '09-07', date: '7 сентября', day: 'Понедельник', homework: '№ 14–18, повторить дроби' },
  { id: '09-08', date: '8 сентября', day: 'Вторник', homework: '№ 19–20, задача 4' },
  { id: '09-09', date: '9 сентября', day: 'Среда', homework: '№ 21, 23 и 25' },
  { id: '09-10', date: '10 сентября', day: 'Четверг', homework: 'Повторить правила, № 27–29' },
  { id: '09-11', date: '11 сентября', day: 'Пятница', homework: 'Подготовиться к проверочной' },
]

const initialGrades: Record<string, Grade[]> = {
  '09-07:anna': [
    { value: '5', kind: 'homework' },
    { value: '4', kind: 'classwork' },
    { value: '5', kind: 'test' },
  ],
  '09-07:ivan': [{ value: '4', kind: 'homework' }],
  '09-07:maria': [{ value: '5', kind: 'classwork' }],
  '09-07:sofia': [{ value: '4', kind: 'homework' }],
  '09-08:anna': [{ value: '4', kind: 'homework' }],
  '09-08:maria': [{ value: '5', kind: 'homework' }],
  '09-08:pavel': [{ value: '4', kind: 'classwork' }],
  '09-08:sofia': [{ value: '5', kind: 'classwork' }],
  '09-09:anna': [{ value: '5', kind: 'classwork' }],
  '09-09:ivan': [{ value: '3', kind: 'homework' }],
  '09-09:pavel': [{ value: '4', kind: 'classwork' }],
  '09-10:ivan': [{ value: '4', kind: 'classwork' }],
  '09-10:maria': [{ value: '5', kind: 'homework' }],
  '09-10:pavel': [{ value: '4', kind: 'homework' }],
  '09-10:sofia': [{ value: '5', kind: 'homework' }],
  '09-11:anna': [{ value: '5', kind: 'test' }],
  '09-11:ivan': [{ value: '4', kind: 'test' }],
  '09-11:maria': [{ value: '5', kind: 'test' }],
  '09-11:pavel': [{ value: '3', kind: 'test' }],
  '09-11:sofia': [{ value: '5', kind: 'test' }],
}

function GradeIcon({ kind }: { kind: GradeKind }) {
  if (kind === 'homework') {
    return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m3 9 7-6 7 6v8h-5v-5H8v5H3V9Z" /></svg>
  }
  if (kind === 'test') {
    return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="M6 4h8v13H6V4Zm2-2h4v4H8V2Zm0 8h4m-4 3h4" /></svg>
  }
  return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m4 14-1 3 3-1L16 6l-2-2L4 14Zm8-8 2 2" /></svg>
}

const gradeLabels: Record<GradeKind, string> = {
  homework: 'Домашняя работа',
  classwork: 'Работа на уроке',
  test: 'Проверочная работа',
}

const gradeKinds = Object.keys(gradeLabels) as GradeKind[]

function GradeControl({ grade, kind, label, onChange }: { grade?: Grade; kind: GradeKind; label: string; onChange: (value: string) => void }) {
  const options = ['5', '4', '3', '2']
  return (
    <label className={`grade-control${grade ? ` grade-control--${kind}` : ''}`} title={gradeLabels[kind]}>
      <GradeIcon kind={kind} />
      <select aria-label={`${gradeLabels[kind]}: ${label}`} value={grade?.value ?? ''} onChange={(event) => onChange(event.target.value)}>
        <option value="">—</option>
        {options.map((value) => <option value={value} key={value}>{value}</option>)}
      </select>
    </label>
  )
}

function ClassJournal({ className, subjectName, onBack }: { className: string; subjectName: string; onBack: () => void }) {
  const [grades, setGrades] = useState(initialGrades)

  function setGrade(lessonID: string, studentID: string, kind: GradeKind, value: string) {
    const key = `${lessonID}:${studentID}`
    setGrades((current) => {
      const otherGrades = (current[key] ?? []).filter((grade) => grade.kind !== kind)
      return { ...current, [key]: value ? [...otherGrades, { value, kind }] : otherGrades }
    })
  }

  return (
    <div className="content-stack">
      <button className="back-button" type="button" onClick={onBack}>
        <span aria-hidden="true">←</span> Предметы класса
      </button>

      <div className="content-heading">
        <div>
          <p className="content-kicker">{className} · {subjectName}</p>
          <h2>Журнал класса</h2>
        </div>
        <span className="term-badge">I четверть</span>
      </div>

      <div className="grade-legend" aria-label="Обозначения оценок">
        {gradeKinds.map((kind) => (
          <span key={kind}><GradeIcon kind={kind} />{gradeLabels[kind]}</span>
        ))}
      </div>

      <div className="table-card journal-scroll">
        <table className="journal-table">
          <thead>
            <tr>
              <th className="student-heading" scope="col">Ученик</th>
              {lessons.map((lesson) => (
                <th scope="col" key={lesson.id}>
                  <span className="lesson-date">{lesson.date}</span>
                  <span className="lesson-day">{lesson.day}</span>
                  <span className="homework-label">Домашнее задание</span>
                  <span className="homework-text">{lesson.homework}</span>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {students.map((student) => (
              <tr key={student.id}>
                <th scope="row"><span className="student-name">{student.name}</span></th>
                {lessons.map((lesson) => {
                  const cellGrades = grades[`${lesson.id}:${student.id}`] ?? []
                  return (
                    <td key={lesson.id}>
                      <div className="grade-cell">
                        {gradeKinds.map((kind) => (
                          <GradeControl
                            grade={cellGrades.find((grade) => grade.kind === kind)}
                            kind={kind}
                            key={kind}
                            label={`${student.name}, ${lesson.date}`}
                            onChange={(value) => setGrade(lesson.id, student.id, kind, value)}
                          />
                        ))}
                      </div>
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default function ClassesView() {
  const [selectedClass, setSelectedClass] = useState<string | null>(null)
  const [selectedSubject, setSelectedSubject] = useState<string | null>(null)

  if (selectedClass && selectedSubject) {
    return <ClassJournal className={selectedClass} subjectName={selectedSubject} onBack={() => setSelectedSubject(null)} />
  }

  if (selectedClass) {
    return (
      <div className="content-stack">
        <button className="back-button" type="button" onClick={() => setSelectedClass(null)}>
          <span aria-hidden="true">←</span> Все классы
        </button>

        <div className="content-heading">
          <div>
            <p className="content-kicker">{selectedClass} · 2026 / 2027 учебный год</p>
            <h2>Предметы класса</h2>
          </div>
          <span className="count-badge">{subjects.length} предметов</span>
        </div>

        <div className="table-card table-scroll">
          <table className="subjects-table">
            <thead>
              <tr><th>Предмет</th><th>Ответственный учитель</th><th>Уроков</th></tr>
            </thead>
            <tbody>
              {subjects.map((subject) => (
                <tr key={subject.id}>
                  <td><button className="subject-link" type="button" onClick={() => setSelectedSubject(subject.name)}>{subject.name}<span aria-hidden="true">→</span></button></td>
                  <td>{subject.teacher}</td>
                  <td><span className="student-count">{subject.lessons}</span></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    )
  }

  return (
    <div className="content-stack">
      <div className="content-heading">
        <div>
          <p className="content-kicker">2026 / 2027 учебный год</p>
          <h2>Классы</h2>
        </div>
        <span className="count-badge">{classes.length} класса</span>
      </div>

      <div className="table-card table-scroll">
        <table className="classes-table">
          <thead>
            <tr><th>Класс</th><th>Учебный год</th><th>Классный руководитель</th><th>Учеников</th></tr>
          </thead>
          <tbody>
            {classes.map((item) => (
              <tr key={item.id}>
                <td><button className="class-link" type="button" onClick={() => setSelectedClass(item.name)}>{item.name}<span aria-hidden="true">→</span></button></td>
                <td>{item.year}</td>
                <td>{item.teacher}</td>
                <td><span className="student-count">{item.students}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
