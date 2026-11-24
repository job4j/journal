import { useEffect, useState } from 'react'
import { downloadLessonAttachment, listLessonAttachments, type LessonAttachment } from './lessonAttachments'

export default function LessonAttachmentsList({lessonID,studentID,revision=0}:{lessonID:string;studentID?:string;revision?:number}) {
  const [items,setItems]=useState<LessonAttachment[]>([])
  const [loading,setLoading]=useState(true)
  const [busy,setBusy]=useState(false)
  const [error,setError]=useState('')
  const [refresh,setRefresh]=useState(0)
  useEffect(()=>{
    let active=true
    setLoading(true);setItems([]);setError('')
    listLessonAttachments(lessonID,studentID).then(value=>{if(active)setItems(value)})
      .catch(cause=>{if(active)setError(cause instanceof Error?cause.message:'Не удалось загрузить файлы урока')})
      .finally(()=>{if(active)setLoading(false)})
    return()=>{active=false}
  },[lessonID,studentID,refresh,revision])
  async function download(item:LessonAttachment){
    setBusy(true);setError('')
    try{await downloadLessonAttachment(item,studentID)}
    catch(cause){setError(cause instanceof Error?cause.message:'Не удалось скачать файл')}
    finally{setBusy(false)}
  }
  if(loading)return <span className="visually-hidden" role="status">Загружаем файлы урока…</span>
  if(!items.length&&!error)return null
  return <div className="lesson-attachments-list">
    <strong>Файлы урока</strong>
    {error&&<div className="content-error" role="alert">{error} <button type="button" className="secondary-action" disabled={busy} onClick={()=>setRefresh(value=>value+1)}>Повторить загрузку файлов</button></div>}
    {items.length>0&&<ul aria-label="Файлы урока">{items.map(item=><li key={item.id}><button type="button" className="subject-link attachment-name" disabled={busy} aria-label={`Скачать ${item.fileName}`} onClick={()=>void download(item)}><svg viewBox="0 0 20 20" aria-hidden="true"><path d="m7 11 5-5a2 2 0 0 1 3 3l-6 6a4 4 0 0 1-6-6l7-7a3 3 0 0 1 4 4l-7 7a1 1 0 0 1-2-2l6-6"/></svg><span>{item.fileName}</span></button><small>{new Intl.NumberFormat('ru-RU',{maximumFractionDigits:1}).format(item.sizeBytes/1024)} КБ</small></li>)}</ul>}
  </div>
}
