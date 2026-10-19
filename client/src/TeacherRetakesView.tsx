import type { ReactNode } from 'react'
import type { ClassStudent } from './classes'
import type { GradeItem, LessonRecord } from './lessons'

export default function TeacherRetakesView({ lessons, students, renderScore }: { lessons: LessonRecord[]; students: ClassStudent[]; renderScore?: (item: GradeItem, member: ClassStudent, date: string) => ReactNode }) {
  const members = new Map(students.map(member => [member.student.id, member]))
  const retakes = lessons.flatMap(lesson => lesson.gradeItems
    .filter(item => item.gradingScale === 'five_point')
    .flatMap(item => item.scores
      .filter(score => !score.isRemediated && (score.numericValue === 2 || score.numericValue === 3) && members.has(score.studentId))
      .map(score => ({ id: score.id, gradeItem: item, member: members.get(score.studentId)!, student: members.get(score.studentId)!.student.name, value: score.numericValue, date: lesson.lessonDate, position: lesson.position, topic: lesson.topic, work: item.title }))))
    .sort((a, b) => b.date.localeCompare(a.date) || b.position - a.position || a.student.localeCompare(b.student, 'ru'))

  if (!retakes.length) return null

  return <section className="teacher-retakes" aria-label="Список учеников с отработками">
    <header><div><h3>Отработки</h3><span className="teacher-retakes-period">за учебный год</span></div><span className="count-badge" aria-label={`Оценок к отработке: ${retakes.length}`}>{retakes.length}</span></header>
    <ul>{retakes.map(item => <li key={item.id}>
      <div className="teacher-retake-student"><strong>{item.student}</strong></div>
      <div className="teacher-retake-grade">{renderScore ? renderScore(item.gradeItem, item.member, item.date) : <strong className={`teacher-retake-score teacher-retake-score-${item.value}`}>{item.value}</strong>}</div>
      <div className="teacher-retake-date"><time dateTime={item.date}>{new Intl.DateTimeFormat('ru-RU').format(new Date(`${item.date}T00:00:00`))}</time></div>
      <div className="teacher-retake-topic"><strong>{item.topic.trim() || 'Без темы'}</strong><small>{item.work}</small></div>
    </li>)}</ul>
  </section>
}
