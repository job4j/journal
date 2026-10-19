import { afterEach, expect, it, vi } from 'vitest'
import { cleanup, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ParentJournalView from './ParentJournalView'
import * as api from './parentJournal'
import { GradeItem, LessonRecord } from './lessons'

const student = { id: 'student-1', name: 'Анна Иванова', roles: ['student'] }
const subject = { id: 'link-1', classId: 'class-1', subject: { id: 's-1', code: 'math', name: 'Математика' }, responsibleTeacher: { id: 't-1', name: 'Иван Петров', roles: ['teacher'] } }
const period: api.ParentStudentPeriod = {
  academicYear: { id: 'year-1', name: '2026/2027', startsOn: '2026-09-01', endsOn: '2027-05-31', status: 'active', quarters: [] },
  class: { id: 'class-1', academicYearId: 'year-1', name: '7А', gradeLevel: 7, studentCount: 1 },
  subjects: [subject],
}

function lesson(id: string, date: string, numericValue = 5): LessonRecord {
  return {
    id, classSubjectId: subject.id, lessonDate: date, position: 1, topic: 'Дроби', homework: '№ 5', materials: [],
    absences: [{ id: `absence-${id}`, lessonId: id, studentId: student.id }],
    gradeItems: [{ id: `item-${id}`, lessonId: id, title: `Работа ${id}`, kind: 'classwork', gradingScale: 'five_point', scores: [{ id: `score-${id}`, gradeItemId: `item-${id}`, studentId: student.id, numericValue, teacherComment: `Комментарий ${id}` }] }],
  }
}

function mockJournal(lessons: LessonRecord[] = []) {
  vi.useFakeTimers({ toFake: ['Date'] })
  vi.setSystemTime(new Date(2026, 9, 21, 12))
  vi.spyOn(api, 'listCurrentParentStudents').mockResolvedValue([student])
  vi.spyOn(api, 'listParentStudentPeriods').mockResolvedValue([period])
  return vi.spyOn(api, 'getParentStudentJournal').mockResolvedValue({
    student, academicYear: period.academicYear, class: period.class,
    subjects: [{ classSubject: subject, lessons: lessons.map(lesson => ({ lesson })) }],
  })
}

afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.useRealTimers() })

it('preserves navigation from the student name to the read-only subject journal', async () => {
  mockJournal([lesson('detail', '2026-09-09')])
  render(<ParentJournalView />)
  await userEvent.click(await screen.findByRole('button', { name: student.name }))
  await userEvent.click(await screen.findByRole('button', { name: '7А' }))
  await userEvent.click(await screen.findByRole('button', { name: 'Математика' }))
  expect(await screen.findByText('Комментарий detail')).toBeInTheDocument()
  expect(screen.getByRole('article')).toHaveTextContent('Дроби')
  expect(screen.getByText('Пропуск')).toBeInTheDocument()
  expect(screen.queryByRole('table')).not.toBeInTheDocument()
  expect(screen.queryByRole('button', { name: 'Сохранить' })).not.toBeInTheDocument()
})

it('shows subjects, lesson dates, multiple grades and comments only for the current Monday–Sunday week', async () => {
  const monday = lesson('понедельник', '2026-10-19', 0)
  const passItem: GradeItem = { id: 'pass-item', lessonId: monday.id, title: 'Зачётная работа', kind: 'other', gradingScale: 'pass_fail', scores: [{ id: 'pass-score', gradeItemId: 'pass-item', studentId: student.id, textValue: 'pass' }] }
  monday.gradeItems.push(passItem, { ...passItem, id: 'ungraded', title: 'Без оценки', scores: [] })
  mockJournal([lesson('будущая', '2026-10-26'), lesson('воскресенье', '2026-10-25', 4), monday, lesson('прошлая', '2026-10-18')])
  render(<ParentJournalView />)
  const weekly = await screen.findByRole('region', { name: `Оценки за неделю — ${student.name}` })
  expect(await within(weekly).findByText('Математика')).toBeInTheDocument()
  expect(within(weekly).getByText('19 октября — 25 октября')).toBeInTheDocument()
  expect(within(weekly).getByText('0')).toBeInTheDocument()
  expect(within(weekly).getByText('Зачёт')).toBeInTheDocument()
  expect(within(weekly).getByText('4')).toBeInTheDocument()
  expect(within(weekly).getByText('Комментарий понедельник')).toBeInTheDocument()
  expect(within(weekly).getAllByRole('listitem')[0]).toHaveTextContent('Работа понедельник')
  expect(within(weekly).queryByText('Работа прошлая')).not.toBeInTheDocument()
  expect(within(weekly).queryByText('Работа будущая')).not.toBeInTheDocument()
  expect(within(weekly).queryByText('Без оценки')).not.toBeInTheDocument()
  expect(within(weekly).queryByRole('button')).not.toBeInTheDocument()
  expect(within(weekly).queryByRole('link')).not.toBeInTheDocument()
})

it('uses the same week on Sunday and does not duplicate requests for classes in one year', async () => {
  const getJournal = mockJournal([lesson('воскресенье', '2026-10-25')])
  vi.setSystemTime(new Date(2026, 9, 25, 23, 59))
  vi.mocked(api.listParentStudentPeriods).mockResolvedValue([period, { ...period, class: { ...period.class, id: 'class-2' } }])
  render(<ParentJournalView />)
  expect(await screen.findByText('Работа воскресенье')).toBeInTheDocument()
  expect(screen.getByText('19 октября — 25 октября')).toBeInTheDocument()
  expect(getJournal).toHaveBeenCalledTimes(1)
})

it('shows an empty state and does not request journals outside the week’s academic years', async () => {
  const getJournal = mockJournal()
  vi.setSystemTime(new Date(2026, 7, 10, 12))
  render(<ParentJournalView />)
  expect(await screen.findByText('За текущую неделю оценок пока нет')).toBeInTheDocument()
  expect(getJournal).not.toHaveBeenCalled()
})

it('shows an empty state when this week has no graded work', async () => {
  mockJournal([lesson('old', '2026-10-18')])
  render(<ParentJournalView />)
  expect(await screen.findByText('За текущую неделю оценок пока нет')).toBeInTheDocument()
})

it('shows access errors without hiding the student navigation', async () => {
  const getJournal = mockJournal()
  getJournal.mockRejectedValue(new Error('Доступ к журналу отозван'))
  render(<ParentJournalView />)
  expect(await screen.findByRole('alert')).toHaveTextContent('Доступ к журналу отозван')
  expect(screen.getByRole('button', { name: student.name })).toBeInTheDocument()
  expect(screen.queryByText('За текущую неделю оценок пока нет')).not.toBeInTheDocument()
})

it('loads each child’s weekly grades independently', async () => {
  mockJournal([lesson('Анны', '2026-10-20')])
  const secondStudent = { ...student, id: 'student-2', name: 'Пётр Иванов' }
  vi.mocked(api.listCurrentParentStudents).mockResolvedValue([student, secondStudent])
  vi.mocked(api.listParentStudentPeriods).mockImplementation(async id => {
    if (id === secondStudent.id) throw new Error('Нет доступа к журналу Петра')
    return [period]
  })
  render(<ParentJournalView />)
  const anna = await screen.findByRole('region', { name: `Оценки за неделю — ${student.name}` })
  expect(await within(anna).findByText('Работа Анны')).toBeInTheDocument()
  const petr = screen.getByRole('region', { name: `Оценки за неделю — ${secondStudent.name}` })
  expect(await within(petr).findByRole('alert')).toHaveTextContent('Нет доступа к журналу Петра')
  expect(within(petr).queryByText('Работа Анны')).not.toBeInTheDocument()
})

it('shows empty students state', async () => {
  vi.spyOn(api, 'listCurrentParentStudents').mockResolvedValue([])
  render(<ParentJournalView />)
  expect(await screen.findByText('Нет доступных учеников')).toBeInTheDocument()
})

it('shows subject lessons and their grades from newest to oldest, including lesson order within a day', async () => {
  mockJournal([
    lesson('oldest', '2026-09-09', 3),
    lesson('newest-first', '2026-10-21', 4),
    lesson('middle', '2026-10-01', 2),
    { ...lesson('newest-second', '2026-10-21', 5), position: 2 },
  ])
  render(<ParentJournalView />)
  await userEvent.click(await screen.findByRole('button', { name: student.name }))
  await userEvent.click(await screen.findByRole('button', { name: '7А' }))
  await userEvent.click(await screen.findByRole('button', { name: 'Математика' }))
  const articles = await screen.findAllByRole('article')
  expect(articles).toHaveLength(4)
  for (const [index, id] of ['newest-second', 'newest-first', 'middle', 'oldest'].entries()) {
    expect(articles[index]).toHaveTextContent(`Работа ${id}`)
    expect(articles[index]).toHaveTextContent(`Комментарий ${id}`)
  }
})

it('shows risk grades outside this week using the same journal request as the weekly summary', async () => {
  const getJournal = mockJournal([lesson('ранняя двойка', '2026-09-09', 2), lesson('свежая тройка', '2026-10-20', 3)])
  render(<ParentJournalView />)
  const risk = await screen.findByRole('region', { name: `Зона риска — ${student.name}` })
  expect(await within(risk).findByText('Работа ранняя двойка')).toBeInTheDocument()
  expect(within(risk).getAllByRole('listitem')[0]).toHaveTextContent('Работа свежая тройка')
  const weekly = screen.getByRole('region', { name: `Оценки за неделю — ${student.name}` })
  expect(await within(weekly).findByText('Работа свежая тройка')).toBeInTheDocument()
  expect(within(weekly).queryByText('Работа ранняя двойка')).not.toBeInTheDocument()
  expect(getJournal).toHaveBeenCalledTimes(1)
})
