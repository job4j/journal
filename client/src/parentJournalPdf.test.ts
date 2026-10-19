import { expect, it } from 'vitest'
import { createParentJournalPdfDefinition } from './parentJournalPdf'
import { gradedLesson, journalFixture } from './test/parentJournalFixture'

it('groups by subject and sorts grades ASC by lesson date and position without mutating the journal', () => {
  const journal = journalFixture()
  journal.subjects[0].lessons = [gradedLesson('new-second', '2026-10-21', 2), gradedLesson('old', '2026-09-01'), gradedLesson('new-first', '2026-10-21')]
  journal.subjects.unshift({ ...journal.subjects[0], classSubject: { ...journal.subjects[0].classSubject, id: 'russian', subject: { id: 'russian', code: 'russian', name: 'Русский язык' } }, lessons: [gradedLesson('dictation', '2026-09-02')] })
  const original = JSON.stringify(journal)
  const report = JSON.stringify(createParentJournalPdfDefinition(journal).content)
  expect(report.indexOf('Математика')).toBeLessThan(report.indexOf('Русский язык'))
  expect(report.indexOf('Работа old')).toBeLessThan(report.indexOf('Работа new-first'))
  expect(report.indexOf('Работа new-first')).toBeLessThan(report.indexOf('Работа new-second'))
  expect(report).toContain('Анна Иванова')
  expect(report).toContain('2026/2027')
  expect(report).toContain('Комментарий old')
  expect(JSON.stringify(journal)).toBe(original)
})

it('includes all graded works, zero, points, pass/fail and comments, but skips ungraded work', () => {
  const journal = journalFixture()
  const { lesson } = gradedLesson('zero', '2026-10-21')
  lesson.gradeItems[0].scores[0].numericValue = 0
  const item = lesson.gradeItems[0]
  lesson.gradeItems.push(
    { ...item, id: 'points', title: 'Баллы', gradingScale: 'points', maxScore: 10, scores: [{ ...item.scores[0], id: 'points-score', numericValue: 8 }] },
    { ...item, id: 'pass', title: 'Зачётная работа', gradingScale: 'pass_fail', scores: [{ ...item.scores[0], id: 'pass-score', numericValue: null, textValue: 'pass' }] },
    { ...item, id: 'fail', title: 'Пересдача', gradingScale: 'pass_fail', scores: [{ ...item.scores[0], id: 'fail-score', numericValue: null, textValue: 'fail' }] },
    { ...item, id: 'ungraded', title: 'Без оценки', scores: [] },
  )
  journal.subjects[0].lessons = [{ lesson }]
  const report = JSON.stringify(createParentJournalPdfDefinition(journal).content)
  for (const value of ['"text":"0"', '8 / 10', 'Зачёт', 'Незачёт', 'Комментарий zero']) expect(report).toContain(value)
  expect(report).not.toContain('Без оценки')
})

it('includes period grades in chronological order and repeats table headers on subsequent pages', () => {
  const journal = journalFixture()
  journal.subjects[0].quarterGrades = ['quarter-2', 'quarter-1'].map((quarterId, index) => ({
    id: quarterId, quarterId, classSubjectId: 'link-1', studentId: journal.student.id, gradingScale: 'five_point', numericValue: 5 - index,
  }))
  const report = JSON.stringify(createParentJournalPdfDefinition(journal).content)
  expect(report).toContain('Итоги периодов')
  expect(report.indexOf('Первая четверть')).toBeLessThan(report.indexOf('Вторая четверть'))
  expect(report).toContain('"headerRows":2')
})

it('creates an explicit empty report when the current year has no grades', () => {
  const report = JSON.stringify(createParentJournalPdfDefinition(journalFixture()).content)
  expect(report).toContain('Оценок за этот учебный год пока нет.')
})
