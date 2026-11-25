import { useEffect, useState } from 'react'
import type { AcademicYear } from './academicYears'
import { localToday } from './classSchedule'
import { listParentStudentPeriods, type ParentStudent } from './parentJournal'
import { listParentMissedTopics, type ParentMissedTopic } from './parentMissedTopics'
import LessonAttachmentsButton from './LessonAttachmentsButton'

function formatDate(value:string) {
  return new Intl.DateTimeFormat('ru-RU',{day:'numeric',month:'long',year:'numeric',timeZone:'UTC'}).format(new Date(`${value}T00:00:00Z`))
}

export default function ParentMissedTopicsView({student,onBack}:{student:ParentStudent;onBack:()=>void}) {
  const [today]=useState(localToday)
  const [years,setYears]=useState<AcademicYear[]>([])
  const [yearID,setYearID]=useState('')
  const [subjectID,setSubjectID]=useState('')
  const [items,setItems]=useState<ParentMissedTopic[]>([])
  const [loadingYears,setLoadingYears]=useState(true)
  const [loadingTopics,setLoadingTopics]=useState(false)
  const [error,setError]=useState('')
  const [yearsAttempt,setYearsAttempt]=useState(0)
  const [topicsAttempt,setTopicsAttempt]=useState(0)
  useEffect(()=>{
    let active=true
    setLoadingYears(true);setError('');setYears([]);setYearID('');setItems([]);setSubjectID('')
    listParentStudentPeriods(student.id).then(periods=>{
      if(!active)return
      const available=[...new Map(periods.map(period=>[period.academicYear.id,period.academicYear])).values()]
        .sort((a,b)=>b.startsOn.localeCompare(a.startsOn)||a.name.localeCompare(b.name,'ru'))
      const selected=available.find(year=>year.startsOn<=today&&year.endsOn>=today)??available.find(year=>year.startsOn<=today)??available[0]
      setYears(available);setLoadingTopics(Boolean(selected));setYearID(selected?.id??'')
    }).catch(cause=>{if(active)setError(cause instanceof Error?cause.message:'Не удалось загрузить учебные годы')})
      .finally(()=>{if(active)setLoadingYears(false)})
    return()=>{active=false}
  },[student.id,today,yearsAttempt])
  useEffect(()=>{
    if(!yearID)return
    let active=true
    setLoadingTopics(true);setError('');setItems([]);setSubjectID('')
    listParentMissedTopics(student.id,yearID).then(value=>{if(active)setItems(value)})
      .catch(cause=>{if(active)setError(cause instanceof Error?cause.message:'Не удалось загрузить пропущенные темы')})
      .finally(()=>{if(active)setLoadingTopics(false)})
    return()=>{active=false}
  },[student.id,yearID,topicsAttempt])

  const topics=items.filter(item=>item.lessonDate<=today)
    .sort((a,b)=>b.lessonDate.localeCompare(a.lessonDate)||b.position-a.position||a.subjectName.localeCompare(b.subjectName,'ru')||a.lessonId.localeCompare(b.lessonId))
  const subjects=[...new Map(topics.map(item=>[item.subjectId,item.subjectName])).entries()].sort((a,b)=>a[1].localeCompare(b[1],'ru'))
  const visible=topics.filter(item=>!subjectID||item.subjectId===subjectID)
  const loading=loadingYears||loadingTopics
  return <div className="content-stack parent-missed-topics-view">
    <button type="button" className="back-button" onClick={onBack}>← К ученикам</button>
    <header><p className="content-kicker">{student.name}</p><h2>Пропущенные темы</h2></header>
    <p className="schedule-hint">Уроки с отметкой о пропуске, которую поставил учитель. Темы помогут восполнить пробелы.</p>
    {years.length>0&&<div className="parent-missed-filters">
      <label>Учебный год<select aria-label="Учебный год пропущенных тем" value={yearID} disabled={loadingYears} onChange={event=>{setItems([]);setSubjectID('');setError('');setLoadingTopics(true);setYearID(event.target.value)}}>{years.map(year=><option value={year.id} key={year.id}>{year.name}</option>)}</select></label>
      <label>Предмет<select aria-label="Предмет пропущенных тем" value={subjectID} disabled={loading||!subjects.length} onChange={event=>setSubjectID(event.target.value)}><option value="">Все предметы</option>{subjects.map(([id,name])=><option value={id} key={id}>{name}</option>)}</select></label>
    </div>}
    {loading?<p role="status">Загружаем пропущенные темы…</p>:error?<div className="content-error" role="alert">{error}<button type="button" className="secondary-action" onClick={()=>{if(years.length)setTopicsAttempt(value=>value+1);else setYearsAttempt(value=>value+1)}}>Повторить загрузку</button></div>:!years.length?<div className="empty-card">Нет доступных учебных годов</div>:<>
      <p className="schedule-hint" role="status">Пропущено уроков: {visible.length}</p>
      {visible.length?<div className="parent-missed-topic-list">{visible.map(item=><article className="parent-missed-topic-card" aria-label={`${item.subjectName} — ${formatDate(item.lessonDate)} — ${item.topic.trim()||'Тема не заполнена'}`} key={item.lessonId}>
        <header><strong>{item.subjectName}</strong><time dateTime={item.lessonDate}>{formatDate(item.lessonDate)}</time></header>
        <p className="parent-missed-topic">{item.topic.trim()||'Тема не заполнена'}</p>
        <small>Класс {item.className}</small>
        <div className="parent-missed-materials"><span>Файлы урока</span><LessonAttachmentsButton lessonID={item.lessonId} date={item.lessonDate} studentID={student.id}/></div>
      </article>)}</div>:<div className="empty-card">{subjectID?'По этому предмету отмеченных пропусков нет':'За выбранный учебный год отмеченных пропусков нет'}</div>}
    </>}
  </div>
}
