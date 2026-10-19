import { afterEach, expect, it } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import ParentRiskGrades from './ParentRiskGrades'
import { gradedLesson, journalFixture } from './test/parentJournalFixture'
import type { GradeItem } from './lessons'

afterEach(cleanup)

function riskLesson(id: string, date: string, value: number, position = 1, scale: GradeItem['gradingScale'] = 'five_point') {
  const entry = gradedLesson(id, date, position)
  entry.lesson.gradeItems[0].gradingScale = scale
  entry.lesson.gradeItems[0].scores[0].numericValue = value
  return entry
}

it('shows only twos and threes on the five-point scale, with subject, date, work and comment', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [
    riskLesson('двойка', '2026-10-01', 2), riskLesson('тройка', '2026-10-19', 3),
    riskLesson('четвёрка', '2026-10-20', 4), riskLesson('пятёрка', '2026-10-21', 5),
    riskLesson('баллы', '2026-10-21', 2, 1, 'points'),
  ]
  const pass = riskLesson('зачёт', '2026-10-21', 3, 1, 'pass_fail')
  pass.lesson.gradeItems[0].scores[0].numericValue = null
  pass.lesson.gradeItems[0].scores[0].textValue = 'fail'
  journal.subjects[0].lessons.push(pass)
  render(<ParentRiskGrades student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  const items = screen.getAllByRole('listitem')
  expect(items).toHaveLength(2)
  expect(items[0]).toHaveTextContent('Работа тройка')
  expect(items[0]).toHaveTextContent('Математика')
  expect(items[0]).toHaveTextContent('19 октября 2026')
  expect(items[0]).toHaveTextContent('Комментарий тройка')
  expect(items[1]).toHaveTextContent('Работа двойка')
  for (const work of ['четвёрка', 'пятёрка', 'баллы', 'зачёт']) expect(screen.queryByText(`Работа ${work}`)).not.toBeInTheDocument()
  expect(screen.queryByRole('button')).not.toBeInTheDocument()
  expect(screen.queryByRole('link')).not.toBeInTheDocument()
})

it('includes period grades and orders dates and lesson positions from newest to oldest', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [riskLesson('первый урок', '2026-10-01', 2), riskLesson('второй урок', '2026-10-01', 3, 2)]
  journal.subjects[0].quarterGrades = [
    { id: 'quarter-grade', quarterId: 'quarter-1', classSubjectId: 'link-1', studentId: journal.student.id, gradingScale: 'five_point', numericValue: 3, teacherComment: 'Нужно повторить материал' },
    { id: 'good-quarter', quarterId: 'quarter-2', classSubjectId: 'link-1', studentId: journal.student.id, gradingScale: 'five_point', numericValue: 4 },
  ]
  render(<ParentRiskGrades student={journal.student} journals={[journal]} today="2026-11-01" loading={false} error="" />)
  const items = screen.getAllByRole('listitem')
  expect(items).toHaveLength(3)
  expect(items[0]).toHaveTextContent('Итоговая оценка · Первая четверть')
  expect(items[0]).toHaveTextContent('Нужно повторить материал')
  expect(items[1]).toHaveTextContent('Работа второй урок')
  expect(items[2]).toHaveTextContent('Работа первый урок')
})

it('excludes previous and future academic years', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [riskLesson('текущий год', '2026-09-01', 2)]
  const past = { ...journal, academicYear: { ...journal.academicYear, id: 'past', startsOn: '2025-09-01', endsOn: '2026-05-31' }, subjects: [{ ...journal.subjects[0], lessons: [riskLesson('прошлый год', '2025-10-01', 2)] }] }
  const future = { ...past, academicYear: { ...past.academicYear, id: 'future', startsOn: '2027-09-01', endsOn: '2028-05-31' } }
  render(<ParentRiskGrades student={journal.student} journals={[past, journal, future]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getAllByRole('listitem')).toHaveLength(1)
  expect(screen.getByText('Работа текущий год')).toBeInTheDocument()
  expect(screen.queryByText('Работа прошлый год')).not.toBeInTheDocument()
})

it.each([{ values: [] }, { values: [4, 5] }])('hides the entire risk block when grades are $values', ({ values }) => {
  const journal = journalFixture()
  journal.subjects[0].lessons = values.map(value => riskLesson(`оценка ${value}`, '2026-10-21', value))
  const { container } = render(<ParentRiskGrades student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.queryByRole('region', { name: `Зона риска — ${journal.student.name}` })).not.toBeInTheDocument()
  expect(container).toBeEmptyDOMElement()
})
it('shows loading and then the access error instead of an empty state', () => {
  const journal = journalFixture()
  const { rerender } = render(<ParentRiskGrades student={journal.student} journals={[]} today="2026-10-21" loading error="" />)
  expect(screen.getByRole('status')).toHaveTextContent('Загружаем оценки…')
  rerender(<ParentRiskGrades student={journal.student} journals={[]} today="2026-10-21" loading={false} error="Доступ к журналу отозван" />)
  expect(screen.getByText('Доступ к журналу отозван')).toBeInTheDocument()
  expect(screen.queryByText('Оценок 2 и 3 за текущий учебный год нет')).not.toBeInTheDocument()
})

it('excludes remediated work grades from the risk zone while keeping period grades', () => {
  const journal = journalFixture()
  const remediated = riskLesson('отработана', '2026-10-21', 2)
  remediated.lesson.gradeItems[0].scores[0].isRemediated = true
  journal.subjects[0].lessons = [remediated]
  const { rerender } = render(<ParentRiskGrades student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.queryByRole('region')).not.toBeInTheDocument()
  remediated.lesson.gradeItems[0].scores[0].isRemediated = false
  rerender(<ParentRiskGrades student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByText('Работа отработана')).toBeInTheDocument()
})
