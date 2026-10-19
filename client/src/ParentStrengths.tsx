import type { ParentJournal, ParentStudent } from './parentJournal'

interface Props {
  student: ParentStudent
  journals: ParentJournal[]
  today: string
  loading: boolean
  error: string
}

function gradeCount(value: number) {
  const remainder = value % 100
  const word = remainder >= 11 && remainder <= 14 ? 'оценок'
    : value % 10 === 1 ? 'оценка' : value % 10 >= 2 && value % 10 <= 4 ? 'оценки' : 'оценок'
  return `${value} ${word}`
}

export default function ParentStrengths({ student, journals, today, loading, error }: Props) {
  const subjects = new Map<string, { name: string; sum: number; count: number }>()
  for (const journal of journals) {
    if (journal.academicYear.startsOn > today || journal.academicYear.endsOn < today) continue
    for (const subject of journal.subjects) {
      const id = subject.classSubject.subject.id
      for (const { lesson } of subject.lessons) {
        for (const item of lesson.gradeItems) {
          if (item.gradingScale !== 'five_point') continue
          for (const score of item.scores) {
            if (score.numericValue == null || !Number.isFinite(score.numericValue)) continue
            const summary = subjects.get(id) ?? { name: subject.classSubject.subject.name, sum: 0, count: 0 }
            summary.sum += score.numericValue
            summary.count++
            subjects.set(id, summary)
          }
        }
      }
    }
  }
  const colors = new Map([...subjects.keys()].sort().map((id, index) => [id, `hsl(${(155 + index * 137.508) % 360} 58% 43%)`]))
  const strengths = [...subjects].map(([id, subject]) => ({ id, ...subject, average: subject.sum / subject.count, color: colors.get(id)! }))
    .sort((a, b) => b.average - a.average || b.count - a.count || a.name.localeCompare(b.name, 'ru'))
  const totalAverage = strengths.reduce((sum, subject) => sum + subject.average, 0)
  const formatNumber = (value: number) => new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 1 }).format(value)
  let offset = 0
  const segments = strengths.map(subject => {
    const share = subject.average / totalAverage * 100
    const segment = { ...subject, share, offset }
    offset += share
    return segment
  })

  function wrapName(name: string) {
    const lines: string[] = []
    let line = ''
    for (const word of name.split(/\s+/).flatMap(word => word.match(/.{1,12}/g) ?? [])) {
      if (line && `${line} ${word}`.length > 12) { lines.push(line); line = word }
      else line = line ? `${line} ${word}` : word
    }
    if (line) lines.push(line)
    return lines
  }
  const labels = segments.map(subject => {
    const angle = (subject.offset + subject.share / 2) / 100 * Math.PI * 2 - Math.PI / 2
    const lines = wrapName(subject.name)
    return { ...subject, angle, right: Math.cos(angle) >= 0, lines, height: (lines.length + 1) * 13 + 14, y: 0 }
  })
  const sides = [labels.filter(label => !label.right), labels.filter(label => label.right)]
  const chartHeight = Math.max(220, ...sides.map(side => side.reduce((sum, label) => sum + label.height, 0) + 24))
  const centerY = chartHeight / 2
  for (const side of sides) {
    side.sort((a, b) => Math.sin(a.angle) - Math.sin(b.angle))
    let top = (chartHeight - side.reduce((sum, label) => sum + label.height, 0)) / 2
    for (const label of side) { label.y = top; top += label.height }
  }
  return <section className="parent-strengths" aria-label={`Сильные стороны — ${student.name}`}>
    <header><h3>Сильные стороны</h3><p>Все предметы · за текущий учебный год</p></header>
    {loading ? <p className="parent-weekly-message" role="status">Загружаем успехи…</p>
      : error ? <p className="content-error">{error}</p>
        : segments.length ? <div className="strength-chart-layout">
          <div className="strength-ring" role="img" aria-label={`Средние баллы по предметам. ${segments.map(subject => `${subject.name}: ${formatNumber(subject.average)} из 5, ${gradeCount(subject.count)}, доля ${formatNumber(subject.share)}%`).join('; ')}`}>
            <svg viewBox={`0 0 320 ${chartHeight}`} aria-hidden="true">
              <circle className="strength-ring-track" cx="160" cy={centerY} r="62" />
              {segments.map(subject => <circle className="strength-ring-value" key={subject.id} cx="160" cy={centerY} r="62" pathLength="100" stroke={subject.color} strokeDasharray={`${subject.share} ${100 - subject.share}`} strokeDashoffset={-subject.offset} transform={`rotate(-90 160 ${centerY})`}><title>{subject.name}: средний балл {formatNumber(subject.average)} из 5, {gradeCount(subject.count)}</title></circle>)}
              {labels.map(label => {
                const x = label.right ? 244 : 4
                const labelY = label.y + (label.height - 14) / 2
                return <g className="strength-chart-label" key={label.id}>
                  <polyline points={`${160 + Math.cos(label.angle) * 74},${centerY + Math.sin(label.angle) * 74} ${label.right ? 236 : 84},${labelY} ${label.right ? 240 : 80},${labelY}`} stroke={label.color} />
                  <text x={x} y={label.y + 11}>{label.lines.map((line, index) => <tspan x={x} dy={index ? 13 : 0} key={index}>{line}</tspan>)}<tspan className="strength-chart-percent" x={x} dy="15" fill={label.color}>{formatNumber(label.share)}%</tspan></text>
                </g>
              })}
            </svg>
          </div>
          <p className="strength-chart-note">Длина сектора пропорциональна среднему баллу предмета.</p>
        </div> : <p className="parent-weekly-message">Пока недостаточно оценок для расчёта.</p>}
  </section>
}