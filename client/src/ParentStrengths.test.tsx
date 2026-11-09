import { afterEach, expect, it } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import ParentStrengths from './ParentStrengths'
import { gradedLesson, journalFixture } from './test/parentJournalFixture'
import type { QuarterGrade } from './quarterGrades'

afterEach(cleanup)

function finalGrade(id: string, value: number, scale: QuarterGrade['gradingScale'] = 'five_point'): QuarterGrade {
  return { id, quarterId: 'quarter-1', classSubjectId: 'link-1', studentId: 'student-1', numericValue: value, gradingScale: scale }
}

it('shows horizontal subject bars alphabetically with widths proportional to average grades', () => {
  const journal = journalFixture()
  journal.subjects[0].quarterGrades = [finalGrade('math-1',5),{...finalGrade('math-2',5),quarterId:'quarter-2'}]
  const english = { ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id:'english-link', subject:{id:'english',name:'Английский язык',code:'english'} }, quarterGrades:[finalGrade('english-1',3),{...finalGrade('english-2',3),quarterId:'quarter-2'}] }
  journal.subjects.push(english)
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  const chart = screen.getByRole('img')
  expect(screen.getByRole('region', {name: `Средний балл — ${journal.student.name}`})).toBeInTheDocument()
  expect(screen.getByRole('heading', {name: 'Средний балл'})).toBeInTheDocument()
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
it('averages only numeric five-point finals and ignores work grades and subjects without finals', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [gradedLesson('work', '2026-10-19')]
  journal.subjects[0].quarterGrades = [finalGrade('first',4),{...finalGrade('second',2),quarterId:'quarter-2'},finalGrade('points',1,'points'),finalGrade('pass',0,'pass_fail'),finalGrade('invalid',NaN),{...finalGrade('missing',5),numericValue:null}]
  journal.subjects.push({...journal.subjects[0],classSubject:{...journal.subjects[0].classSubject,id:'physics-link',subject:{id:'physics',name:'Физика',code:'physics'}},quarterGrades:undefined})
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByRole('img')).toHaveAccessibleName(expect.stringContaining('Математика: 3 из 5, 2 оценки'))
  expect(screen.getByRole('img').querySelector('.strength-bar')).toHaveStyle({width: '60%'})
  expect(screen.queryByText('Физика')).not.toBeInTheDocument()
})

it('excludes historical and future years and includes subjects with an average below four', () => {
  const journal = journalFixture()
  journal.subjects[0].quarterGrades = [finalGrade('low',3)]
  const past = {...journal,academicYear:{...journal.academicYear,id:'past',startsOn:'2025-09-01',endsOn:'2026-05-31'},subjects:[{...journal.subjects[0],quarterGrades:[finalGrade('past',5)]}]}
  const future = {...past,academicYear:{...past.academicYear,id:'future',startsOn:'2027-09-01',endsOn:'2028-05-31'}}
  render(<ParentStrengths student={journal.student} journals={[past,journal,future]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getByRole('img')).toBeInTheDocument()
  expect(screen.getByRole('img')).toHaveAccessibleName(expect.stringContaining('Математика: 3 из 5, 1 оценка'))
  expect(screen.getByRole('img').querySelector('.strength-bar-value')).toHaveTextContent('3')
})

it('combines finals for the same subject using grade counts rather than averaging class averages', () => {
  const journal = journalFixture()
  journal.subjects[0].quarterGrades = [finalGrade('first',5)]
  journal.subjects.push({...journal.subjects[0],classSubject:{...journal.subjects[0].classSubject,id:'second-link'},quarterGrades:[finalGrade('second',4),{...finalGrade('third',4),quarterId:'quarter-2'}]})
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.getAllByRole('img')).toHaveLength(1)
  expect(screen.getByRole('img')).toHaveAccessibleName(expect.stringContaining('Математика: 4,3 из 5'))
})

it('does not calculate an average from work grades when there are no finals', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [gradedLesson('work', '2026-10-19')]
  render(<ParentStrengths student={journal.student} journals={[journal]} today="2026-10-21" loading={false} error="" />)
  expect(screen.queryByRole('img')).not.toBeInTheDocument()
  expect(screen.getByText('Пока нет итоговых оценок для расчёта.')).toBeInTheDocument()
})

it('distinguishes loading and errors from a computed result', () => {
  const journal = journalFixture()
  const {rerender} = render(<ParentStrengths student={journal.student} journals={[]} today="2026-10-21" loading error="" />)
  expect(screen.getByRole('status')).toHaveTextContent('Загружаем средний балл…')
  rerender(<ParentStrengths student={journal.student} journals={[]} today="2026-10-21" loading={false} error="Доступ к журналу отозван" />)
  expect(screen.getByText('Доступ к журналу отозван')).toBeInTheDocument()
  expect(screen.queryByText('Пока нет итоговых оценок для расчёта.')).not.toBeInTheDocument()
})
