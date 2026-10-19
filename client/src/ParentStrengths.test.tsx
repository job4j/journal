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

it('shows one ring for all subjects with sectors proportional to their average grades', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [lesson('math-1',5),lesson('math-2',5)]
  const english = { ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id:'english-link', subject:{id:'english',name:'Английский язык',code:'english'} }, lessons:[lesson('english-1',3),lesson('english-2',3)] }
  journal.subjects.unshift(english)
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  const ring = screen.getByRole('img')
  const sectors = ring.querySelectorAll('.strength-ring-value')
  expect(sectors).toHaveLength(2)
  expect(sectors[0]).toHaveAttribute('stroke-dasharray','62.5 37.5')
  expect(sectors[1]).toHaveAttribute('stroke-dasharray','37.5 62.5')
  expect(sectors[1]).toHaveAttribute('stroke-dashoffset','-62.5')
  expect(sectors[0].getAttribute('stroke')).not.toBe(sectors[1].getAttribute('stroke'))
  expect(screen.getByLabelText('Математика: средний балл 5 из 5')).toBeInTheDocument()
  expect(screen.getByLabelText('Английский язык: средний балл 3 из 5')).toBeInTheDocument()
  expect(ring).toHaveAccessibleName(expect.stringContaining('доля 62,5%'))
})
it('uses all work grades, including remediated low grades, but excludes other scales and period grades', () => {
  const journal = journalFixture()
  const low = lesson('low',2)
  low.lesson.gradeItems[0].scores[0].isRemediated = true
  journal.subjects[0].lessons = [lesson('first',5),lesson('second',5),low,lesson('points',1,'points'),lesson('pass',0,'pass_fail')]
  journal.subjects[0].quarterGrades = [{id:'quarter',quarterId:'quarter-1',classSubjectId:'link-1',studentId:journal.student.id,gradingScale:'five_point',numericValue:2}]
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByLabelText('Математика: средний балл 4 из 5')).toBeInTheDocument()
  expect(screen.getByText('3 оценки · 100%')).toBeInTheDocument()
})

it('excludes historical years and includes subjects with an average below four', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [...Array.from({length:24},(_,index)=>lesson('current-'+index,4)),lesson('low',3)]
  const past = {...journal,academicYear:{...journal.academicYear,id:'past',startsOn:'2025-09-01',endsOn:'2026-05-31'},subjects:[{...journal.subjects[0],lessons:[lesson('past',5)]}]}
  render(<ParentStrengths student={journal.student} journals={[past,journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByRole('img')).toBeInTheDocument()
  expect(screen.getByText('25 оценок · 100%')).toBeInTheDocument()
})

it('combines assessments for the same subject using grade counts rather than averaging class averages', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [lesson('first',5)]
  journal.subjects.push({...journal.subjects[0],classSubject:{...journal.subjects[0].classSubject,id:'second-link'},lessons:[lesson('second',4),lesson('third',4)]})
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getAllByRole('img')).toHaveLength(1)
  expect(screen.getByLabelText('Математика: средний балл 4,3 из 5')).toBeInTheDocument()
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
