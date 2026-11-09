import { afterEach, expect, it } from 'vitest'
import { cleanup, render, screen, within } from '@testing-library/react'
import ParentReportCard from './ParentReportCard'
import { gradedLesson, journalFixture } from './test/parentJournalFixture'
import type { QuarterGrade } from './quarterGrades'

afterEach(cleanup)

function grade(id: string, quarterId = 'quarter-1', value = 5): QuarterGrade {
  return { id, quarterId, classSubjectId: 'link-1', studentId: 'student-1', gradingScale: 'five_point', numericValue: value }
}

it('groups final grades by quarter, sorts subjects and hides empty quarters', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [gradedLesson('ordinary-grade', '2026-10-19')]
  journal.subjects[0].quarterGrades = [grade('math', 'quarter-1', 4)]
  journal.subjects.push({ ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id: 'english-link', subject: { id: 'english', code: 'english', name: 'Английский язык' } }, quarterGrades: [{ ...grade('english'), teacherComment: 'Отличная работа' }] })
  render(<ParentReportCard student={journal.student} journals={[journal]} today="2026-11-09" loading={false} error="" />)
  const table = screen.getByRole('table', { name: 'Итоги — Первая четверть' })
  expect(within(table).getAllByRole('rowheader').map(cell => cell.textContent)).toEqual(['Английский язык', 'Математика'])
  expect(screen.queryByText('Отличная работа')).not.toBeInTheDocument()
  expect(within(table).getAllByRole('cell').map(cell => cell.textContent)).toEqual(['5', '4'])
  expect(screen.queryByText('Вторая четверть')).not.toBeInTheDocument()
  expect(screen.queryByText('Работа ordinary-grade')).not.toBeInTheDocument()
  expect(screen.queryByRole('button')).not.toBeInTheDocument()
})

it('orders quarters chronologically and shows points, zero and pass/fail finals', () => {
  const journal = journalFixture()
  journal.academicYear.quarters.reverse()
  journal.subjects[0].quarterGrades = [
    { ...grade('points', 'quarter-1', 0), gradingScale: 'points', maxScore: 10 },
    { ...grade('fail', 'quarter-2'), gradingScale: 'pass_fail', numericValue: null, textValue: 'fail' },
  ]
  journal.subjects.push({ ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id: 'pass-link' }, quarterGrades: [{ ...grade('pass', 'quarter-2'), gradingScale: 'pass_fail', numericValue: null, textValue: 'pass' }] })
  render(<ParentReportCard student={journal.student} journals={[journal]} today="2026-12-01" loading={false} error="" />)
  expect(screen.getAllByRole('table').map(table => table.getAttribute('aria-label'))).toEqual(['Итоги — Первая четверть', 'Итоги — Вторая четверть'])
  expect(screen.getByText('0 / 10')).toBeInTheDocument()
  expect(screen.getByText('Незачёт')).toBeInTheDocument()
  expect(screen.getByText('Зачёт')).toBeInTheDocument()
})

it('does not show historical or future academic years or a report with no finals', () => {
  const journal = journalFixture()
  const past = { ...journal, academicYear: { ...journal.academicYear, id: 'past', startsOn: '2025-09-01', endsOn: '2026-05-31' }, subjects: [{ ...journal.subjects[0], quarterGrades: [grade('past')] }] }
  const future = { ...past, academicYear: { ...past.academicYear, id: 'future', startsOn: '2027-09-01', endsOn: '2028-05-31' } }
  const { container } = render(<ParentReportCard student={journal.student} journals={[past, journal, future]} today="2026-11-09" loading={false} error="" />)
  expect(container).toBeEmptyDOMElement()
})

it('combines classes under one quarter and identifies their subjects', () => {
  const journal = journalFixture()
  journal.subjects[0].quarterGrades = [grade('first')]
  const second = { ...journal, class: { ...journal.class, id: 'class-2', name: '7Б' }, subjects: [{ ...journal.subjects[0], quarterGrades: [grade('second', 'quarter-1', 4)] }] }
  render(<ParentReportCard student={journal.student} journals={[journal, second]} today="2026-11-09" loading={false} error="" />)
  expect(screen.getAllByRole('table')).toHaveLength(1)
  expect(screen.getByText('7А')).toBeInTheDocument()
  expect(screen.getByText('7Б')).toBeInTheDocument()
})

it('distinguishes loading and access errors from no finals', () => {
  const journal = journalFixture()
  const { rerender } = render(<ParentReportCard student={journal.student} journals={[]} today="2026-11-09" loading error="" />)
  expect(screen.getByRole('status')).toHaveTextContent('Загружаем табель…')
  rerender(<ParentReportCard student={journal.student} journals={[]} today="2026-11-09" loading={false} error="Доступ к журналу отозван" />)
  expect(screen.getByText('Доступ к журналу отозван')).toBeInTheDocument()
  expect(screen.queryByRole('table')).not.toBeInTheDocument()
})
