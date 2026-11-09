import { useEffect, useState } from 'react'
import { addDays, localToday, monday } from './classSchedule'
import { listParentStudentPeriods, type ParentStudent, type ParentStudentPeriod } from './parentJournal'
import { getParentSchedule, type ParentScheduleWeek } from './parentSchedule'

function formatDate(day: string) {
  return new Intl.DateTimeFormat('ru-RU',{day:'numeric',month:'long',timeZone:'UTC'}).format(new Date(`${day}T00:00:00Z`))
}
function initialWeek(period: ParentStudentPeriod) {
  const today=localToday(), year=period.academicYear
  return monday(today < year.startsOn ? year.startsOn : today > year.endsOn ? year.endsOn : today)
}
export default function ParentScheduleView({student,onBack}:{student:ParentStudent;onBack:()=>void}) {
  const [periods,setPeriods]=useState<ParentStudentPeriod[]>([])
  const [selected,setSelected]=useState<ParentStudentPeriod | null>(null)
  const [weekStart,setWeekStart]=useState('')
  const [week,setWeek]=useState<ParentScheduleWeek | null>(null)
  const [loadingPeriods,setLoadingPeriods]=useState(true)
  const [loadingWeek,setLoadingWeek]=useState(false)
  const [error,setError]=useState('')
  const [refresh,setRefresh]=useState(0)

  useEffect(()=>{
    let active=true
    setLoadingPeriods(true);setError('');setSelected(null);setWeek(null)
    listParentStudentPeriods(student.id).then(values=>{
      if (!active) return
      setPeriods(values)
      const today=localToday()
      const choice=values.find(p=>p.academicYear.startsOn<=today && p.academicYear.endsOn>=today) ?? values[0] ?? null
      setSelected(choice);setWeekStart(choice ? initialWeek(choice) : '')
    }).catch(cause=>{if(active)setError(cause instanceof Error ? cause.message : 'Не удалось загрузить классы')})
      .finally(()=>{if(active)setLoadingPeriods(false)})
    return ()=>{active=false}
  },[student.id])
  useEffect(()=>{
    if (!selected || !weekStart) return
    let active=true
    setLoadingWeek(true);setError('');setWeek(null)
    getParentSchedule(student.id,selected.class.id,weekStart).then(value=>{if(active)setWeek(value)})
      .catch(cause=>{if(active)setError(cause instanceof Error ? cause.message : 'Не удалось загрузить расписание')})
      .finally(()=>{if(active)setLoadingWeek(false)})
    return ()=>{active=false}
  },[student.id,selected,weekStart,refresh])

  const days=weekStart ? Array.from({length:7},(_,index)=>addDays(weekStart,index)).filter(day=>selected && day>=selected.academicYear.startsOn && day<=selected.academicYear.endsOn) : []
  return <div className="content-stack parent-schedule-view">
    <button className="back-button" onClick={onBack}>← К ученикам</button>
    <header><p className="content-kicker">{student.name}</p><h2>Расписание</h2></header>
    {periods.length>1 && <label className="parent-schedule-period">Класс и учебный год<select aria-label="Класс расписания" value={selected?.class.id ?? ''} onChange={event=>{
      const choice=periods.find(p=>p.class.id===event.target.value)
      if(choice){setWeek(null);setSelected(choice);setWeekStart(initialWeek(choice))}
    }}>{periods.map(p=><option value={p.class.id} key={p.class.id}>{p.class.name} · {p.academicYear.name}</option>)}</select></label>}
    {selected && <>
      <p className="schedule-hint">{selected.class.name} · {selected.academicYear.name}</p>
      <nav className="parent-schedule-navigation" aria-label="Недели расписания">
        <button className="secondary-action" aria-label="Предыдущая неделя" disabled={loadingWeek || weekStart<=monday(selected.academicYear.startsOn)} onClick={()=>setWeekStart(addDays(weekStart,-7))}>←</button>
        <strong>{formatDate(weekStart)} — {formatDate(addDays(weekStart,6))}</strong>
        <button className="secondary-action" aria-label="Следующая неделя" disabled={loadingWeek || weekStart>=monday(selected.academicYear.endsOn)} onClick={()=>setWeekStart(addDays(weekStart,7))}>→</button>
      </nav>
      <button className="secondary-action parent-schedule-today" disabled={loadingWeek || weekStart===initialWeek(selected)} onClick={()=>setWeekStart(initialWeek(selected))}>Текущая неделя</button>
    </>}
    {error && <div className="content-error" role="alert">{error}{selected && <button className="secondary-action" onClick={()=>setRefresh(value=>value+1)}>Повторить загрузку</button>}</div>}
    {(loadingPeriods || loadingWeek) ? <p role="status">Загружаем расписание…</p> : !error && !selected ? <div className="empty-card">Нет доступных классов</div> : null}
    {!loadingWeek && week && days.map(day=>{
      const entries=week.entries.filter(e=>e.day===day).sort((a,b)=>a.position-b.position || (a.startsAt ?? '').localeCompare(b.startsAt ?? '') || a.id.localeCompare(b.id))
      const weekday=new Intl.DateTimeFormat('ru-RU',{weekday:'long',timeZone:'UTC'}).format(new Date(`${day}T00:00:00Z`))
      return <section className={`parent-schedule-day${day===localToday() ? ' parent-schedule-day-today' : ''}`} aria-label={`${weekday}, ${formatDate(day)}`} key={day}>
        <header><h3>{weekday}</h3><time dateTime={day}>{formatDate(day)}</time></header>
        {entries.length ? <ol>{entries.map(entry=>{
          const firstOfSubject=entry.classSubjectId && entries.find(e=>e.classSubjectId===entry.classSubjectId)?.id===entry.id
          return <li key={entry.id}>
            <div className="parent-schedule-lesson"><span className="parent-schedule-number">{entry.position}</span><div><h4>{entry.title}</h4><p className="parent-schedule-time">{entry.startsAt ? `${entry.startsAt}${entry.endsAt ? `–${entry.endsAt}` : ''}` : 'Время не указано'}</p></div></div>
            {entry.notes && <p className="parent-schedule-notes">{entry.notes}</p>}
            {firstOfSubject && <div className="parent-schedule-homework"><strong>Домашнее задание</strong><p>{entry.homework || 'Не задано'}</p>{entry.homework && entry.homeworkFrom && <small>Задано {formatDate(entry.homeworkFrom)}</small>}</div>}
          </li>
        })}</ol> : <p className="parent-weekly-message">Занятий нет</p>}
      </section>
    })}
  </div>
}
