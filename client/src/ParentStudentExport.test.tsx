import { afterEach, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ParentStudentExport from './ParentStudentExport'
import ParentJournalView from './ParentJournalView'
import * as api from './parentJournal'
import { downloadParentJournalPdf } from './parentJournalPdf'
import { journalFixture } from './test/parentJournalFixture'

vi.mock('./parentJournalPdf', () => ({ downloadParentJournalPdf: vi.fn() }))
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.clearAllMocks(); vi.useRealTimers() })

function setup() {
  vi.useFakeTimers({ toFake: ['Date'] })
  vi.setSystemTime(new Date(2026, 9, 21, 12))
  const journal = journalFixture()
  const period = { academicYear: journal.academicYear, class: journal.class, subjects: journal.subjects.map(s => s.classSubject) }
  vi.spyOn(api, 'listCurrentParentStudents').mockResolvedValue([journal.student])
  const periods = vi.spyOn(api, 'listParentStudentPeriods').mockResolvedValue([period])
  const getJournal = vi.spyOn(api, 'getParentStudentJournal').mockResolvedValue(journal)
  vi.mocked(downloadParentJournalPdf).mockResolvedValue()
  return { journal, period, periods, getJournal }
}

it('provides a separate PDF export on the home page without opening the student journal', async () => {
  const { journal } = setup()
  render(<ParentJournalView />)
  await userEvent.click(await screen.findByRole('button', { name: `Экспорт PDF — ${journal.student.name}` }))
  await vi.waitFor(() => expect(downloadParentJournalPdf).toHaveBeenCalledWith(journal))
  expect(screen.getByRole('heading', { name: 'Ученики' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: journal.student.name })).toBeInTheDocument()
})

it('exports the selected student and current academic year even when history is unsorted', async () => {
  const { journal, period, periods, getJournal } = setup()
  periods.mockResolvedValue([
    { ...period, academicYear: { ...period.academicYear, id: 'future', startsOn: '2027-09-01', endsOn: '2028-05-31' } },
    { ...period, academicYear: { ...period.academicYear, id: 'past', startsOn: '2025-09-01', endsOn: '2026-05-31' } },
    period,
  ])
  render(<ParentStudentExport student={journal.student} />)
  expect(periods).not.toHaveBeenCalled()
  await userEvent.click(screen.getByRole('button'))
  await vi.waitFor(() => expect(downloadParentJournalPdf).toHaveBeenCalledWith(journal))
  expect(periods).toHaveBeenCalledWith(journal.student.id)
  expect(getJournal).toHaveBeenCalledExactlyOnceWith(journal.student.id, journal.academicYear.id)
})

it('disables duplicate exports while the PDF is being prepared', async () => {
  const { journal, getJournal } = setup()
  let finish!: (value: api.ParentJournal) => void
  getJournal.mockReturnValue(new Promise(resolve => { finish = resolve }))
  render(<ParentStudentExport student={journal.student} />)
  const button = screen.getByRole('button')
  await userEvent.click(button)
  expect(button).toBeDisabled()
  expect(button).toHaveTextContent('Готовим PDF…')
  finish(journal)
  await vi.waitFor(() => expect(button).toBeEnabled())
  expect(downloadParentJournalPdf).toHaveBeenCalledTimes(1)
})

it('reports a missing current academic year without downloading an empty report', async () => {
  const { journal, periods, getJournal } = setup()
  periods.mockResolvedValue([])
  render(<ParentStudentExport student={journal.student} />)
  await userEvent.click(screen.getByRole('button'))
  expect(await screen.findByRole('alert')).toHaveTextContent('У ученика нет журнала за текущий учебный год')
  expect(getJournal).not.toHaveBeenCalled()
  expect(downloadParentJournalPdf).not.toHaveBeenCalled()
})

it('does not generate a PDF when the API denies access', async () => {
  const { journal, getJournal } = setup()
  getJournal.mockRejectedValue(new Error('Доступ к журналу отозван'))
  render(<ParentStudentExport student={journal.student} />)
  await userEvent.click(screen.getByRole('button'))
  expect(await screen.findByRole('alert')).toHaveTextContent('Доступ к журналу отозван')
  expect(downloadParentJournalPdf).not.toHaveBeenCalled()
  expect(screen.getByRole('button')).toBeEnabled()
})

it('allows retrying after PDF generation fails', async () => {
  const { journal } = setup()
  vi.mocked(downloadParentJournalPdf).mockRejectedValueOnce(new Error('Не удалось подготовить PDF'))
  render(<ParentStudentExport student={journal.student} />)
  await userEvent.click(screen.getByRole('button'))
  expect(await screen.findByRole('alert')).toHaveTextContent('Не удалось подготовить PDF')
  await userEvent.click(screen.getByRole('button'))
  await vi.waitFor(() => expect(downloadParentJournalPdf).toHaveBeenCalledTimes(2))
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
})
