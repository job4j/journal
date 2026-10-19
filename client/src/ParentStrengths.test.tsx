import { afterEach, expect, it } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import ParentStrengths from './ParentStrengths'
import { gradedLesson, journalFixture } from './test/parentJournalFixture'
import type { GradeItem } from './lessons'

afterEach(cleanup)

function lesson(id: string, value: number, scale: GradeItem['gradingScale'] = 'five_point') {
  const entry = gradedLesson(id, '2026-10-19')
  entry.lesson.gradeItems[0].scores[0].numericValue = value
  entry.lesson.gradeItems[0].gradingScale = scale
  return entry
}

it('shows rings with the arithmetic mean out of five and sample size, sorted by average', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [lesson('math-1',5),lesson('math-2',4)]
  const russian = { ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id:'russian-link', subject:{id:'russian',name:'Русский язык',code:'russian'} }, lessons:[lesson('russian-1',4),lesson('russian-2',4)] }
  journal.subjects.unshift(russian)
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  const rings = screen.getAllByRole('img')
  expect(rings).toHaveLength(2)
  expect(rings[0]).toHaveAccessibleName('Математика: средний балл 4,5 из 5, 2 оценки')
  expect(rings[0].querySelector('.strength-ring-value')).toHaveAttribute('stroke-dasharray','90 100')
  expect(rings[1]).toHaveAccessibleName('Русский язык: средний балл 4 из 5, 2 оценки')
})

it('uses all work grades, including remediated low grades, but excludes other scales and period grades', () => {
  const journal = journalFixture()
  const low = lesson('low',2)
  low.lesson.gradeItems[0].scores[0].isRemediated = true
  journal.subjects[0].lessons = [lesson('first',5),lesson('second',5),low,lesson('points',1,'points'),lesson('pass',0,'pass_fail')]
  journal.subjects[0].quarterGrades = [{id:'quarter',quarterId:'quarter-1',classSubjectId:'link-1',studentId:journal.student.id,gradingScale:'five_point',numericValue:2}]
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByRole('img')).toHaveAccessibleName('Математика: средний балл 4 из 5, 3 оценки')
})

it('does not include historical years or round a below-threshold mean up to four', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [...Array.from({length:24},(_,index)=>lesson('current-'+index,4)),lesson('low',3)]
  const past = {...journal,academicYear:{...journal.academicYear,id:'past',startsOn:'2025-09-01',endsOn:'2026-05-31'},subjects:[{...journal.subjects[0],lessons:[lesson('past',5)]}]}
  render(<ParentStrengths student={journal.student} journals={[past,journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.queryByRole('img')).not.toBeInTheDocument()
  expect(screen.getByText('Здесь появятся предметы со средним баллом от 4.')).toBeInTheDocument()
})

it('combines assessments for the same subject using grade counts rather than averaging class averages', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [lesson('first',5)]
  journal.subjects.push({...journal.subjects[0],classSubject:{...journal.subjects[0].classSubject,id:'second-link'},lessons:[lesson('second',4),lesson('third',4)]})
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getAllByRole('img')).toHaveLength(1)
  expect(screen.getByRole('img')).toHaveAccessibleName('Математика: средний балл 4,3 из 5, 3 оценки')
})

it('does not invent a result when there are no graded works', () => {
  const journal = journalFixture()
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.queryByRole('img')).not.toBeInTheDocument()
  expect(screen.getByText('Пока недостаточно оценок для расчёта.')).toBeInTheDocument()
})

it('distinguishes loading and errors from a computed result', () => {
  const journal = journalFixture()
  const {rerender} = render(<ParentStrengths student={journal.student} journals={[]} today="2026-10-21" loading error="" />)
  expect(screen.getByRole('status')).toHaveTextContent('Загружаем успехи…')
  rerender(<ParentStrengths student={journal.student} journals={[]} today="2026-10-21" loading={false} error="Доступ к журналу отозван" />)
  expect(screen.getByText('Доступ к журналу отозван')).toBeInTheDocument()
  expect(screen.queryByText('Пока недостаточно оценок для расчёта.')).not.toBeInTheDocument()
})
