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

it('shows horizontal subject bars alphabetically with widths proportional to average grades', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [lesson('math-1',5),lesson('math-2',5)]
  const english = { ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id:'english-link', subject:{id:'english',name:'Английский язык',code:'english'} }, lessons:[lesson('english-1',3),lesson('english-2',3)] }
  journal.subjects.push(english)
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  const chart = screen.getByRole('img')
  const bars = chart.querySelectorAll<HTMLElement>('.strength-bar')
  expect(bars).toHaveLength(2)
  expect([...chart.querySelectorAll('.strength-bar-name')].map(label => label.textContent)).toEqual(['Английский язык', 'Математика'])
  expect(bars[0]).toHaveStyle({width: '60%'})
  expect(bars[1]).toHaveStyle({width: '100%'})
  expect(bars[0].style.backgroundColor).not.toBe(bars[1].style.backgroundColor)
  expect(chart.querySelectorAll('.strength-bar-value')[0]).toHaveTextContent('3')
  expect(chart.querySelectorAll('.strength-bar-value')[1]).toHaveTextContent('5')
  expect(chart).toHaveAccessibleName(expect.stringContaining('Математика: 5 из 5'))
  expect(chart).toHaveAccessibleName(expect.stringContaining('Английский язык: 3 из 5'))
  expect(chart.querySelector('.strength-bar-row')).toHaveAttribute('title', 'Английский язык: средний балл 3 из 5, 2 оценки')
})
it('uses all work grades, including remediated low grades, but excludes other scales and period grades', () => {
  const journal = journalFixture()
  const low = lesson('low',2)
  low.lesson.gradeItems[0].scores[0].isRemediated = true
  journal.subjects[0].lessons = [lesson('first',5),lesson('second',5),low,lesson('points',1,'points'),lesson('pass',0,'pass_fail')]
  journal.subjects[0].quarterGrades = [{id:'quarter',quarterId:'quarter-1',classSubjectId:'link-1',studentId:journal.student.id,gradingScale:'five_point',numericValue:2}]
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByRole('img')).toHaveAccessibleName(expect.stringContaining('Математика: 4 из 5, 3 оценки'))
  expect(screen.getByRole('img').querySelector('.strength-bar')).toHaveStyle({width: '80%'})
})

it('excludes historical and future years and includes subjects with an average below four', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [...Array.from({length:24},(_,index)=>lesson('current-'+index,4)),lesson('low',3)]
  const past = {...journal,academicYear:{...journal.academicYear,id:'past',startsOn:'2025-09-01',endsOn:'2026-05-31'},subjects:[{...journal.subjects[0],lessons:[lesson('past',5)]}]}
  const future = {...past,academicYear:{...past.academicYear,id:'future',startsOn:'2027-09-01',endsOn:'2028-05-31'}}
  render(<ParentStrengths student={journal.student} journals={[past,journal,future]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByRole('img')).toBeInTheDocument()
  expect(screen.getByRole('img')).toHaveAccessibleName(expect.stringContaining('25 оценок'))
  expect(screen.getByRole('img').querySelector('.strength-bar-value')).toHaveTextContent('4')
})

it('combines assessments for the same subject using grade counts rather than averaging class averages', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [lesson('first',5)]
  journal.subjects.push({...journal.subjects[0],classSubject:{...journal.subjects[0].classSubject,id:'second-link'},lessons:[lesson('second',4),lesson('third',4)]})
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getAllByRole('img')).toHaveLength(1)
  expect(screen.getByRole('img')).toHaveAccessibleName(expect.stringContaining('Математика: 4,3 из 5'))
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
