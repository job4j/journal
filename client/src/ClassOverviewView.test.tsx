import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ClassOverviewView from './ClassOverviewView'

const item = {
  id: 'class-1',
  academicYearId: 'year-1',
  name: '7А',
  gradeLevel: 7,
  studentCount: 1,
  quarters: [
    { id: 'quarter-1', number: 1, startsOn: '2026-09-01', endsOn: '2026-10-31' },
    { id: 'quarter-2', number: 2, startsOn: '2026-11-01', endsOn: '2026-12-31' },
  ],
}

const student = {
  student: { id: 'student-1', name: 'Анна Иванова', roles: ['student'] },
  enrolledOn: '2026-09-01',
  leftOn: null,
}

const assignment = {
  id: 'subject-link-1',
  classId: 'class-1',
  subject: { id: 'subject-1', code: 'math', name: 'Математика' },
  responsibleTeacher: { id: 'teacher-1', name: 'Иван Петров', roles: ['teacher'] },
}

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})

describe('class overview', () => {
  it('shows students, subjects and period grades', async () => {
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)
      if (url.endsWith('/classes/class-1/students')) {
        return new Response(JSON.stringify({ items: [student] }), { status: 200 })
      }
      if (url.endsWith('/classes/class-1/subjects')) {
        return new Response(JSON.stringify({ items: [assignment] }), { status: 200 })
      }
      if (url.includes('/class-subjects/subject-link-1/lessons?')) {
        return new Response(JSON.stringify({ items: [] }), { status: 200 })
      }      if (url.includes('/quarters/quarter-2/grades')) {
        return new Response(JSON.stringify({
          items: [{
            id: 'grade-2',
            quarterId: 'quarter-2',
            classSubjectId: 'subject-link-1',
            studentId: 'student-1',
            gradingScale: 'five_point',
            numericValue: 4,
          }],
        }), { status: 200 })
      }      if (url.includes('/quarters/quarter-1/grades')) {
        return new Response(JSON.stringify({
          items: [{
            id: 'grade-1',
            quarterId: 'quarter-1',
            classSubjectId: 'subject-link-1',
            studentId: 'student-1',
            gradingScale: 'five_point',
            numericValue: 5,
          }],
        }), { status: 200 })
      }
      return new Response(null, { status: 404 })
    })

    render(<ClassOverviewView item={item} onBack={() => {}} onCountChange={() => {}} />)

    expect(await screen.findByText('Журнал класса')).toBeInTheDocument()
    expect(screen.getByText('7А · Математика')).toBeInTheDocument()
    expect(await screen.findAllByText('Анна Иванова')).toHaveLength(2)
    const period = screen.getByLabelText('Период')
    expect(period).toHaveValue('quarter-1')

    await userEvent.selectOptions(period, 'quarter-2')
    expect(period).toHaveValue('quarter-2')
  })
})