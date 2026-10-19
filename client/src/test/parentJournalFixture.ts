import type { ParentJournal } from '../parentJournal'
import type { LessonRecord } from '../lessons'

export function journalFixture(): ParentJournal {
  return {
    student: { id: 'student-1', name: 'Анна Иванова', roles: ['student'] },
    academicYear: { id: 'year-1', name: '2026/2027', startsOn: '2026-09-01', endsOn: '2027-05-31', status: 'active', quarters: [
      { id: 'quarter-1', number: 1, name: 'Первая четверть', startsOn: '2026-09-01', endsOn: '2026-10-30' },
      { id: 'quarter-2', number: 2, name: 'Вторая четверть', startsOn: '2026-11-01', endsOn: '2026-12-25' },
    ] },
    class: { id: 'class-1', academicYearId: 'year-1', name: '7А', gradeLevel: 7, studentCount: 1 },
    subjects: [{ classSubject: { id: 'link-1', classId: 'class-1', subject: { id: 'subject-1', code: 'math', name: 'Математика' }, responsibleTeacher: { id: 'teacher-1', name: 'Иван Петров', roles: ['teacher'] } }, lessons: [] }],
  }
}

export function gradedLesson(id: string, date: string, position = 1): { lesson: LessonRecord } {
  return { lesson: { id, classSubjectId: 'link-1', lessonDate: date, position, topic: 'Дроби', materials: [], absences: [], gradeItems: [
    { id: `item-${id}`, lessonId: id, title: `Работа ${id}`, kind: 'classwork', gradingScale: 'five_point', scores: [
      { id: `score-${id}`, gradeItemId: `item-${id}`, studentId: 'student-1', numericValue: 5, teacherComment: `Комментарий ${id}` },
    ] },
  ] } }
}
