import { useState } from 'react'
import { ParentStudent, getParentStudentJournal, listParentStudentPeriods } from './parentJournal'

export default function ParentStudentExport({ student }: { student: ParentStudent }) {
  const [exporting, setExporting] = useState(false)
  const [error, setError] = useState('')

  async function exportPdf() {
    setExporting(true)
    setError('')
    try {
      const now = new Date()
      const today = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
      const periods = await listParentStudentPeriods(student.id)
      const period = periods.filter(({ academicYear }) => academicYear.startsOn <= today && academicYear.endsOn >= today)
        .sort((a, b) => b.academicYear.startsOn.localeCompare(a.academicYear.startsOn))[0]
      if (!period) {
        setError('У ученика нет журнала за текущий учебный год')
        return
      }
      const journal = await getParentStudentJournal(student.id, period.academicYear.id)
      const { downloadParentJournalPdf } = await import('./parentJournalPdf')
      await downloadParentJournalPdf(journal)
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : 'Не удалось экспортировать оценки')
    } finally {
      setExporting(false)
    }
  }

  return <div className="parent-student-export">
    <button className="secondary-action" aria-label={`Экспорт PDF — ${student.name}`} disabled={exporting} onClick={() => void exportPdf()}>{exporting ? 'Готовим PDF…' : 'Экспорт PDF'}</button>
    <small>Все оценки за текущий учебный год</small>
    {error && <p className="content-error" role="alert">{error}</p>}
  </div>
}
