import { afterEach, expect, it } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import TeacherRetakesView from './TeacherRetakesView'
import { gradedLesson, journalFixture } from './test/parentJournalFixture'
import type { ClassStudent } from './classes'

afterEach(cleanup)
const students: ClassStudent[] = [{ student: journalFixture().student, enrolledOn: '2026-09-01' }]

it('shows every low grade with student, value, date, topic and work, newest first', () => {
  const older = gradedLesson('старшая', '2026-09-09')
  older.lesson.gradeItems[0].scores[0].numericValue = 2
  older.lesson.topic = 'Обыкновенные дроби'
  const newer = gradedLesson('новая', '2026-10-01')
  newer.lesson.gradeItems[0].scores[0].numericValue = 3
  newer.lesson.topic = 'Десятичные дроби'
  render(<TeacherRetakesView lessons={[older.lesson, newer.lesson]} students={students} />)
  const items = screen.getAllByRole('listitem')
  expect(items).toHaveLength(2)
  expect(items[0]).toHaveTextContent('Анна Иванова')
  expect(items[0]).toHaveTextContent('3')
  expect(items[0]).toHaveTextContent('01.10.2026')
  expect(items[0]).toHaveTextContent('Десятичные дроби')
  expect(items[0]).toHaveTextContent('Работа новая')
  expect(items[1]).toHaveTextContent('Обыкновенные дроби')
  expect(items[1]).toHaveTextContent('09.09.2026')
})

it('does not include good grades, points, pass/fail or scores without a student in the roster', () => {
  const lessons = [2, 3, 4, 5].map(value => {
    const { lesson } = gradedLesson(String(value), '2026-09-09')
    lesson.gradeItems[0].scores[0].numericValue = value
    return lesson
  })
  lessons[0].gradeItems[0].gradingScale = 'points'
  lessons[1].gradeItems[0].gradingScale = 'pass_fail'
  const unknown = gradedLesson('unknown', '2026-09-09').lesson
  unknown.gradeItems[0].scores[0].numericValue = 2
  unknown.gradeItems[0].scores[0].studentId = 'unknown'
  const { container } = render(<TeacherRetakesView lessons={[...lessons, unknown]} students={students} />)
  expect(container).toBeEmptyDOMElement()
})

it('hides the banner when there are no retakes', () => {
  const { container } = render(<TeacherRetakesView lessons={[]} students={students} />)
  expect(container).toBeEmptyDOMElement()
})

it('hides remediated grades from the teacher retakes list', () => {
  const { lesson } = gradedLesson('отработана', '2026-09-09')
  lesson.gradeItems[0].scores[0].numericValue = 2
  lesson.gradeItems[0].scores[0].isRemediated = true
  const { container } = render(<TeacherRetakesView lessons={[lesson]} students={students} />)
  expect(container).toBeEmptyDOMElement()
})
