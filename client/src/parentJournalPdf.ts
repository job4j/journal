import type { Content, TableCell, TDocumentDefinitions } from 'pdfmake/interfaces'
import type { ParentJournal } from './parentJournal'
import type { GradeItem, Score } from './lessons'
import type { QuarterGrade } from './quarterGrades'

function formatDate(value: string) {
  return new Intl.DateTimeFormat('ru-RU').format(new Date(`${value}T00:00:00`))
}

function gradeValue(score: Score | QuarterGrade, grading: Pick<GradeItem, 'gradingScale' | 'maxScore'>) {
  if (score.numericValue != null) {
    return grading.gradingScale === 'points' && grading.maxScore != null
      ? `${score.numericValue} / ${grading.maxScore}` : String(score.numericValue)
  }
  return score.textValue === 'pass' ? 'Зачёт' : score.textValue === 'fail' ? 'Незачёт' : (score.textValue ?? '-')
}

function heading(text: string): TableCell {
  return { text, bold: true, color: '#ffffff', fillColor: '#234a68', margin: [5, 6, 5, 6] }
}

function gradeTable(name: string, columns: string[], rows: TableCell[][]): Content {
  return {
    margin: [0, 16, 0, 0],
    layout: {
      hLineWidth: () => 0.5,
      vLineWidth: () => 0,
      hLineColor: () => '#dbe5ed',
      paddingTop: () => 6,
      paddingBottom: () => 6,
      fillColor: row => row > 1 && row % 2 === 0 ? '#f3f7fa' : null,
    },
    table: {
      headerRows: 2,
      keepWithHeaderRows: 1,
      widths: [64, '*', 60, '*'],
      body: [
        [{ text: name, colSpan: 4, fontSize: 13, bold: true, color: '#183f5d', margin: [0, 4, 0, 6] }, {}, {}, {}],
        columns.map(heading),
        ...rows,
      ],
    },
  }
}

export function createParentJournalPdfDefinition(journal: ParentJournal): TDocumentDefinitions {
  const content: Content[] = [
    { text: 'JOURNAL', fontSize: 9, bold: true, color: '#60798c', characterSpacing: 2, margin: [0, 0, 0, 8] },
    { text: 'Сводка оценок', fontSize: 23, bold: true, color: '#183f5d', margin: [0, 0, 0, 12] },
    { text: journal.student.name, fontSize: 16, bold: true, margin: [0, 0, 0, 5] },
    { text: `${journal.academicYear.name} учебный год · ${journal.class.name}`, fontSize: 11, margin: [0, 0, 0, 4] },
    { text: `${formatDate(journal.academicYear.startsOn)} - ${formatDate(journal.academicYear.endsOn)}`, color: '#60798c', margin: [0, 0, 0, 6] },
  ]
  let gradeCount = 0
  for (const subject of [...journal.subjects].sort((a, b) => a.classSubject.subject.name.localeCompare(b.classSubject.subject.name, 'ru'))) {
    const rows: TableCell[][] = [...subject.lessons]
      .sort((a, b) => a.lesson.lessonDate.localeCompare(b.lesson.lessonDate) || a.lesson.position - b.lesson.position)
      .flatMap(({ lesson }) => lesson.gradeItems.flatMap(item => item.scores.map(score => [
        formatDate(lesson.lessonDate), item.title,
        { text: gradeValue(score, item), bold: true, alignment: 'center' as const },
        score.teacherComment ?? '-',
      ])))
    gradeCount += rows.length
    if (rows.length) content.push(gradeTable(subject.classSubject.subject.name, ['Дата урока', 'Работа', 'Оценка', 'Комментарий'], rows))

    const quarters = journal.academicYear.quarters
    const quarterRows: TableCell[][] = [...(subject.quarterGrades ?? [])]
      .sort((a, b) => (quarters.find(q => q.id === a.quarterId)?.endsOn ?? '').localeCompare(quarters.find(q => q.id === b.quarterId)?.endsOn ?? ''))
      .map(grade => {
        const quarter = quarters.find(q => q.id === grade.quarterId)
        return [quarter ? formatDate(quarter.endsOn) : '-', quarter?.name ?? 'Учебный период', { text: gradeValue(grade, grade), bold: true, alignment: 'center' }, grade.teacherComment ?? '-']
      })
    gradeCount += quarterRows.length
    if (quarterRows.length) content.push(gradeTable(`${subject.classSubject.subject.name} · Итоги периодов`, ['Окончание', 'Учебный период', 'Оценка', 'Комментарий'], quarterRows))
  }
  if (!gradeCount) content.push({ text: 'Оценок за этот учебный год пока нет.', margin: [0, 18, 0, 0], color: '#60798c' })

  return {
    info: { title: `Сводка оценок - ${journal.student.name} - ${journal.academicYear.name}`, author: 'Journal' },
    pageSize: 'A4',
    pageMargins: [40, 40, 40, 48],
    defaultStyle: { font: 'Roboto', fontSize: 10, color: '#29465e', lineHeight: 1.2 },
    content,
    footer: (page, pages) => ({ text: `${page} / ${pages}`, alignment: 'right', color: '#60798c', fontSize: 9, margin: [40, 14, 40, 0] }),
  }
}

export async function downloadParentJournalPdf(journal: ParentJournal) {
  const [{ default: pdfMake }, { default: fonts }] = await Promise.all([
    import('pdfmake/build/pdfmake'), import('pdfmake/build/vfs_fonts'),
  ])
  pdfMake.addVirtualFileSystem(fonts)
  const name = `Оценки-${journal.student.name}-${journal.academicYear.name}`.replace(/[<>:"/\\|?*]/g, '-').replace(/\p{Cc}/gu, '-').slice(0, 150)
  await pdfMake.createPdf(createParentJournalPdfDefinition(journal)).download(`${name}.pdf`)
}
