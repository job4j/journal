import { FormEvent, useEffect, useState } from 'react'
import { AcademicYear, AcademicYearDraft, createAcademicYear, listAcademicYears } from './academicYears'

const emptyDraft = ():AcademicYearDraft => ({name:'',startsOn:'',endsOn:'',status:'planned',quarters:[1,2,3,4].map(number=>({number,startsOn:'',endsOn:''}))})
const statusLabels = {planned:'Запланирован',active:'Активный',completed:'Завершён'}

export default function AcademicYearsView(){
 const[items,setItems]=useState<AcademicYear[]>([]);const[loading,setLoading]=useState(true);const[error,setError]=useState('');const[draft,setDraft]=useState<AcademicYearDraft|null>(null);const[saving,setSaving]=useState(false)
 useEffect(()=>{let active=true;listAcademicYears().then(value=>{if(active)setItems(value)}).catch(cause=>{if(active)setError(cause instanceof Error?cause.message:'Не удалось загрузить учебные годы')}).finally(()=>{if(active)setLoading(false)});return()=>{active=false}},[])
 function setQuarter(index:number,field:'startsOn'|'endsOn',value:string){setDraft(current=>current?{...current,quarters:current.quarters.map((quarter,i)=>i===index?{...quarter,[field]:value}:quarter)}:current)}
 async function submit(event:FormEvent){event.preventDefault();if(!draft)return;setSaving(true);setError('');try{const saved=await createAcademicYear(draft);setItems(current=>[saved,...current]);setDraft(null)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось создать учебный год')}finally{setSaving(false)}}
 return <div className="content-stack">
  <div className="content-heading"><div><p className="content-kicker">Учебный календарь</p><h2>Учебные годы</h2></div><button className="primary-action" type="button" onClick={()=>setDraft(emptyDraft())}>+ Создать учебный год</button></div>
  {error&&<p className="content-error" role="alert">{error}</p>}
  {loading?<div className="empty-card">Загружаем учебные годы…</div>:<div className="table-card table-scroll"><table className="roles-table"><thead><tr><th>Название</th><th>Период</th><th>Статус</th><th>Четверти</th></tr></thead><tbody>
   {items.map(item=><tr key={item.id}><td><strong>{item.name}</strong></td><td>{item.startsOn} — {item.endsOn}</td><td>{statusLabels[item.status]}</td><td>{item.quarters.length}</td></tr>)}
   {!items.length&&<tr><td className="empty-cell" colSpan={4}>Учебные годы ещё не созданы</td></tr>}
  </tbody></table></div>}
  {draft&&<div className="dialog-backdrop" role="presentation"><section className="role-dialog" role="dialog" aria-modal="true" aria-labelledby="year-dialog-title"><div className="dialog-heading"><h3 id="year-dialog-title">Новый учебный год</h3><button type="button" aria-label="Закрыть" onClick={()=>setDraft(null)}>×</button></div>
   <form className="role-form" onSubmit={submit}><label>Название<input value={draft.name} onChange={event=>setDraft({...draft,name:event.target.value})} required maxLength={100}/></label><div className="form-columns"><label>Начало<input type="date" value={draft.startsOn} onChange={event=>setDraft({...draft,startsOn:event.target.value})} required/></label><label>Окончание<input type="date" value={draft.endsOn} onChange={event=>setDraft({...draft,endsOn:event.target.value})} required/></label></div><label>Статус<select value={draft.status} onChange={event=>setDraft({...draft,status:event.target.value as AcademicYearDraft['status']})}><option value="planned">Запланирован</option><option value="active">Активный</option><option value="completed">Завершён</option></select></label>
   <fieldset><legend>Четверти</legend>{draft.quarters.map((quarter,index)=><div className="form-columns" key={quarter.number}><label>{quarter.number} четверть: начало<input aria-label={`${quarter.number} четверть: начало`} type="date" value={quarter.startsOn} onChange={event=>setQuarter(index,'startsOn',event.target.value)} required/></label><label>{quarter.number} четверть: окончание<input aria-label={`${quarter.number} четверть: окончание`} type="date" value={quarter.endsOn} onChange={event=>setQuarter(index,'endsOn',event.target.value)} required/></label></div>)}</fieldset>
   <div className="dialog-actions"><button className="secondary-action" type="button" onClick={()=>setDraft(null)}>Отмена</button><button className="primary-action" type="submit" disabled={saving}>{saving?'Сохраняем…':'Сохранить'}</button></div></form>
  </section></div>}
 </div>
}
