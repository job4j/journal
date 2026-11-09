import { useEffect, useState, type FormEvent } from 'react'
import type { AcademicYear } from './academicYears'
import type { ClassRecord } from './classes'
import type { ClassSubjectRecord } from './classSubjects'
import { addDays, getClassSchedule, localToday, monday, replaceClassSchedule, type ScheduleEntry, type ScheduleWeek } from './classSchedule'

function formatDate(day: string) {
  return new Intl.DateTimeFormat('ru-RU', {day:'numeric',month:'long',timeZone:'UTC'}).format(new Date(`${day}T00:00:00Z`))
}

export default function ClassScheduleView({ item, year, subjects, onBack }: {item: ClassRecord; year: AcademicYear; subjects: ClassSubjectRecord[]; onBack: () => void}) {
  const today = localToday()
  const currentWeek = monday(today < year.startsOn ? year.startsOn : today > year.endsOn ? year.endsOn : today)
  const [weekStart, setWeekStart] = useState(currentWeek)
  const [week, setWeek] = useState<ScheduleWeek | null>(null)
  const [original, setOriginal] = useState<ScheduleWeek | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [refresh, setRefresh] = useState(0)
  const [editing, setEditing] = useState<{index: number; entry: ScheduleEntry} | null>(null)
  const dirty = JSON.stringify(week) !== JSON.stringify(original)
  const invalidTime = Boolean(editing?.entry.endsAt && (!editing.entry.startsAt || editing.entry.endsAt <= editing.entry.startsAt))
  const busy = loading || saving
  const days = Array.from({length:7}, (_, index) => addDays(weekStart,index))
  const firstDay = days.find(day => day >= year.startsOn) ?? weekStart
  const lastDay = [...days].reverse().find(day => day <= year.endsOn) ?? days[6]

  useEffect(() => {
    let active = true
    setLoading(true); setError(''); setNotice(''); setEditing(null)
    getClassSchedule(item.id, weekStart).then(value => {
      if (active) {setWeek(value); setOriginal(value)}
    }).catch(cause => {if (active) {setWeek(null);setOriginal(null);setError(cause instanceof Error ? cause.message : 'Не удалось загрузить расписание')}})
      .finally(() => {if (active) setLoading(false)})
    return () => {active = false}
  }, [item.id, weekStart, refresh])

  async function save() {
    if (!week) return
    setSaving(true); setError('');setNotice('')
    try {
      const saved = await replaceClassSchedule(item.id,week)
      setWeek(saved); setOriginal(saved); setNotice('Расписание этой недели сохранено')
    } catch (cause) {setError(cause instanceof Error ? cause.message : 'Не удалось сохранить расписание')}
    finally {setSaving(false)}
  }
  function open(day: string, position: number, index = -1) {
    setEditing({index, entry:index >= 0 ? {...week!.entries[index]} : {id:'',day,position,classSubjectId:subjects[0]?.id ?? null,title:subjects[0]?.subject.name ?? '',notes:'',startsAt:null,endsAt:null}})
    setNotice('')
  }
  function apply(event: FormEvent) {
    event.preventDefault()
    if (!week || !editing) return
    const entries = [...week.entries]
    if (editing.index < 0) entries.push(editing.entry)
    else entries[editing.index] = editing.entry
    setWeek({...week,entries}); setEditing(null)
  }
  const positions = Array.from({length:Math.max(7,...(week?.entries.map(e => e.position) ?? []))}, (_,index) => index + 1)
  return <div className="content-stack class-schedule">
    <button className="back-button" disabled={dirty || busy} onClick={onBack}>← К классу</button>
    <div className="content-heading"><div><p className="content-kicker">{item.name} · {year.name}</p><h2>Расписание</h2></div></div>
    <div className="schedule-toolbar">
      <button className="secondary-action" aria-label="Предыдущая неделя расписания" disabled={busy || dirty || weekStart <= monday(year.startsOn)} onClick={() => setWeekStart(addDays(weekStart,-7))}>←</button>
      <label>Неделя<input aria-label="Неделя расписания" type="date" min={monday(year.startsOn)} max={monday(year.endsOn)} value={weekStart} disabled={busy || dirty} onChange={event => {if (event.target.value) setWeekStart(monday(event.target.value))}} /></label>
      <button className="secondary-action" aria-label="Следующая неделя расписания" disabled={busy || dirty || weekStart >= monday(year.endsOn)} onClick={() => setWeekStart(addDays(weekStart,7))}>→</button>
      <button className="secondary-action" disabled={busy || dirty || weekStart === currentWeek} onClick={() => setWeekStart(currentWeek)}>Текущая неделя</button>
    </div>
    <p className="schedule-hint">{formatDate(weekStart)} — {formatDate(days[6])}. Изменения касаются только этой недели.</p>
    {error && <div className="content-error" role="alert">{error}<button type="button" className="secondary-action" disabled={busy} onClick={() => setRefresh(value => value + 1)}>Загрузить неделю заново</button></div>}
    {notice && <p role="status" className="schedule-hint">{notice}</p>}
    {loading ? <p role="status">Загружаем расписание…</p> : week && <>
      <div className="schedule-actions">
        <button className="primary-action" disabled={!dirty || saving} onClick={() => void save()}>{saving ? 'Сохраняем…' : 'Сохранить неделю'}</button>
        <button className="secondary-action" disabled={!dirty || saving} onClick={() => {setWeek(original);setEditing(null);setError('')}}>Отменить изменения</button>
        <button className="secondary-action" disabled={saving} onClick={() => open(firstDay,1)}>Добавить занятие</button>
        {dirty && <span className="schedule-hint">Сохраните или отмените изменения перед переходом к другой неделе.</span>}
      </div>
      <div className="table-card schedule-scroll"><table className="schedule-table" aria-label={`Расписание ${item.name}`}>
        <thead><tr><th scope="col">№</th>{days.map(day => <th scope="col" key={day}><span>{new Intl.DateTimeFormat('ru-RU',{weekday:'long',timeZone:'UTC'}).format(new Date(`${day}T00:00:00Z`))}</span><small>{formatDate(day)}</small></th>)}</tr></thead>
        <tbody>{positions.map(position => <tr key={position}><th scope="row">{position}</th>{days.map(day => <td key={day} className={day < year.startsOn || day > year.endsOn ? 'schedule-outside' : ''}>
          {day >= year.startsOn && day <= year.endsOn && <>
            {week.entries.map((entry,index) => entry.day === day && entry.position === position ? <button type="button" className={`schedule-entry${entry.classSubjectId ? '' : ' schedule-entry-extra'}`} disabled={saving} key={entry.id || `draft-${index}`} aria-label={`Изменить ${entry.title}, ${day}, урок ${position}`} onClick={() => open(day,position,index)}>
              <strong>{entry.title}</strong><small>{entry.startsAt ?? ''}{entry.endsAt ? `–${entry.endsAt}` : ''}</small>{entry.notes && <small>{entry.notes}</small>}{!entry.classSubjectId && <span>Внеурочное</span>}
            </button> : null)}
            <button type="button" className="schedule-add" disabled={saving} aria-label={`Добавить занятие ${day}, урок ${position}`} onClick={() => open(day,position)}>+</button>
          </>}
        </td>)}</tr>)}</tbody>
      </table></div>
      {item.id === '60000000-0000-4000-8000-000000000001' && <p className="schedule-hint">Обед по исходному расписанию: 11:30–12:00. Готовность к первому уроку: 08:40.</p>}
    </>}
    {editing && <div className="dialog-backdrop"><form className="role-dialog role-form" role="dialog" aria-modal="true" aria-label="Занятие расписания" onSubmit={apply}>
      <h3>{editing.index < 0 ? 'Добавить занятие' : 'Изменить занятие'}</h3>
      <label>Дата<input aria-label="Дата занятия" type="date" min={firstDay} max={lastDay} value={editing.entry.day} required onChange={event => setEditing({...editing,entry:{...editing.entry,day:event.target.value}})} /></label>
      <label>Номер урока<input aria-label="Номер урока" type="number" min={1} max={20} required value={editing.entry.position} onChange={event => setEditing({...editing,entry:{...editing.entry,position:Number(event.target.value)}})} /></label>
      <label>Предмет<select aria-label="Предмет занятия" value={editing.entry.classSubjectId ?? ''} onChange={event => setEditing({...editing,entry:{...editing.entry,classSubjectId:event.target.value || null,title:subjects.find(s => s.id === event.target.value)?.subject.name ?? ''}})}>
        <option value="">Внеурочное занятие</option>{subjects.map(s => <option value={s.id} key={s.id}>{s.subject.name}</option>)}
      </select></label>
      {!editing.entry.classSubjectId && <label>Название<input aria-label="Название занятия" required maxLength={200} value={editing.entry.title} onChange={event => setEditing({...editing,entry:{...editing.entry,title:event.target.value}})} /></label>}
      <div className="form-columns"><label>Начало<input aria-label="Начало занятия" type="time" value={editing.entry.startsAt ?? ''} onChange={event => setEditing({...editing,entry:{...editing.entry,startsAt:event.target.value || null}})} /></label><label>Окончание<input aria-label="Окончание занятия" type="time" min={editing.entry.startsAt ?? undefined} value={editing.entry.endsAt ?? ''} onChange={event => setEditing({...editing,entry:{...editing.entry,endsAt:event.target.value || null}})} /></label></div>
      <label>Примечание<textarea aria-label="Примечание занятия" maxLength={1000} value={editing.entry.notes} onChange={event => setEditing({...editing,entry:{...editing.entry,notes:event.target.value}})} /></label>
      {invalidTime && <p className="content-error">Окончание должно быть позже начала занятия.</p>}
      <div className="dialog-actions"><button className="primary-action" type="submit" disabled={invalidTime || (!editing.entry.classSubjectId && !editing.entry.title.trim())}>Применить</button><button type="button" className="secondary-action" onClick={() => setEditing(null)}>Отмена</button>{editing.index >= 0 && <button type="button" className="secondary-action danger-action" onClick={() => {setWeek({...week!,entries:week!.entries.filter((_,index) => index !== editing.index)});setEditing(null)}}>Удалить занятие</button>}</div>
    </form></div>}
  </div>
}
