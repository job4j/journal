import { useEffect, useId, useState, type FormEvent } from 'react'
import { deleteLessonAttachment, downloadLessonAttachment, listLessonAttachments, maxAttachmentBytes, uploadLessonAttachment, type LessonAttachment } from './lessonAttachments'

export default function LessonAttachmentsButton({lessonID,date,studentID,onChanged}:{lessonID?:string;date:string;studentID?:string;onChanged?:()=>void}) {
  const [open,setOpen]=useState(false)
  const [items,setItems]=useState<LessonAttachment[]>([])
  const [file,setFile]=useState<File|null>(null)
  const [loading,setLoading]=useState(false)
  const [busy,setBusy]=useState(false)
  const [error,setError]=useState('')
  const [refresh,setRefresh]=useState(0)
  const [inputVersion,setInputVersion]=useState(0)
  const titleID=useId()
  useEffect(()=>{
    if(!open||!lessonID)return
    let active=true
    setLoading(true);setError('');setItems([])
    listLessonAttachments(lessonID,studentID).then(value=>{if(active)setItems(value)})
      .catch(cause=>{if(active)setError(cause instanceof Error?cause.message:'Не удалось загрузить файлы')})
      .finally(()=>{if(active)setLoading(false)})
    return()=>{active=false}
  },[open,lessonID,studentID,refresh])
  async function upload(event:FormEvent){
    event.preventDefault()
    if(!file||!lessonID||busy)return
    if(!file.size||file.size>maxAttachmentBytes){setError('Выберите непустой файл размером не более 20 МБ');return}
    setBusy(true);setError('')
    try{const saved=await uploadLessonAttachment(lessonID,file);setItems(current=>[...current,saved]);setFile(null);setInputVersion(value=>value+1);onChanged?.()}
    catch(cause){setError(cause instanceof Error?cause.message:'Не удалось загрузить файл')}
    finally{setBusy(false)}
  }
  async function remove(item:LessonAttachment){
    setBusy(true);setError('')
    try{await deleteLessonAttachment(item.id);setItems(current=>current.filter(value=>value.id!==item.id));onChanged?.()}
    catch(cause){setError(cause instanceof Error?cause.message:'Не удалось удалить файл')}
    finally{setBusy(false)}
  }
  async function download(item:LessonAttachment){
    setBusy(true);setError('')
    try{await downloadLessonAttachment(item,studentID)}
    catch(cause){setError(cause instanceof Error?cause.message:'Не удалось скачать файл')}
    finally{setBusy(false)}
  }
  return <>
    <button type="button" className="lesson-tasks-button" aria-label={`Файлы урока ${date}`} title={lessonID?'Файлы урока':'Сначала сохраните тему урока'} disabled={!lessonID} onClick={()=>{setFile(null);setInputVersion(value=>value+1);setOpen(true)}}>
      <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m7 11 5-5a2 2 0 0 1 3 3l-6 6a4 4 0 0 1-6-6l7-7a3 3 0 0 1 4 4l-7 7a1 1 0 0 1-2-2l6-6"/></svg>
    </button>
    {open&&<div className="dialog-backdrop"><section className="role-dialog attachment-dialog" role="dialog" aria-modal="true" aria-labelledby={titleID}>
      <h3 id={titleID}>Файлы урока · {date.split('-').reverse().join('.')}</h3>
      {error&&<p className="content-error" role="alert">{error}</p>}
      {loading?<p role="status">Загружаем файлы…</p>:<>
        {items.length?<ul className="attachment-list">{items.map(item=><li key={item.id}><div><button type="button" className="subject-link attachment-name" disabled={busy} onClick={()=>void download(item)} aria-label={`Скачать ${item.fileName}`}>{item.fileName}</button><small>{new Intl.NumberFormat('ru-RU',{maximumFractionDigits:1}).format(item.sizeBytes/1024)} КБ</small></div>{!studentID&&<button type="button" className="secondary-action danger-action" disabled={busy} onClick={()=>void remove(item)} aria-label={`Удалить ${item.fileName}`}>Удалить</button>}</li>)}</ul>:<p className="summary-empty">Файлы пока не прикреплены</p>}
        {error&&<button type="button" className="secondary-action" disabled={busy} onClick={()=>setRefresh(value=>value+1)}>Обновить список</button>}
        {!studentID&&<form className="role-form" onSubmit={upload}><label>Файл<input key={inputVersion} type="file" disabled={busy} onChange={event=>{setFile(event.target.files?.[0]??null);setError('')}}/></label><p className="schedule-hint">Один файл до 20 МБ</p><button type="submit" className="primary-action" disabled={busy||!file}>{busy?'Подождите…':'Загрузить'}</button></form>}
      </>}
      <div className="dialog-actions"><button type="button" className="secondary-action" disabled={busy} onClick={()=>setOpen(false)}>Закрыть</button></div>
    </section></div>}
  </>
}
