import { FormEvent, useEffect, useState } from 'react'
import { createRole, deleteRole, listPermissions, listRoles, Permission, Role, RolesApiError, updateRole } from './roles'

type Draft={id?:string;code:string;name:string;permissions:string[]}
const emptyDraft:Draft={code:'',name:'',permissions:[]}

export default function RolesView(){
 const [roles,setRoles]=useState<Role[]>([]);const [permissions,setPermissions]=useState<Permission[]>([])
 const [draft,setDraft]=useState<Draft|null>(null);const [loading,setLoading]=useState(true);const [saving,setSaving]=useState(false);const [error,setError]=useState('')
 useEffect(()=>{let active=true;Promise.all([listRoles(),listPermissions()]).then(([nextRoles,nextPermissions])=>{if(active){setRoles(nextRoles);setPermissions(nextPermissions)}}).catch(cause=>{if(active)setError(cause instanceof RolesApiError?cause.message:'Не удалось загрузить роли')}).finally(()=>{if(active)setLoading(false)});return()=>{active=false}},[])
 function openEditor(role?:Role){setError('');setDraft(role?{...role,permissions:[...role.permissions]}:{...emptyDraft,permissions:[]})}
 function togglePermission(code:string){setDraft(current=>current?{...current,permissions:current.permissions.includes(code)?current.permissions.filter(item=>item!==code):[...current.permissions,code]}:current)}
 async function submit(event:FormEvent){event.preventDefault();if(!draft)return;setSaving(true);setError('');try{const saved=draft.id?await updateRole(draft as Role):await createRole(draft);setRoles(current=>draft.id?current.map(role=>role.id===saved.id?saved:role):[...current,saved]);setDraft(null)}catch(cause){setError(cause instanceof RolesApiError?cause.message:'Не удалось сохранить роль')}finally{setSaving(false)}}
 async function remove(role:Role){if(!window.confirm(`Удалить роль «${role.name}»?`))return;setError('');try{await deleteRole(role.id);setRoles(current=>current.filter(item=>item.id!==role.id))}catch(cause){setError(cause instanceof RolesApiError?cause.message:'Не удалось удалить роль')}}
 return <div className="content-stack">
  <div className="content-heading"><div><p className="content-kicker">Управление доступом</p><h2>Роли</h2></div><button className="primary-action" type="button" aria-label="Создать роль" onClick={()=>openEditor()}>+ Создать роль</button></div>
  {error&&<p className="content-error" role="alert">{error}</p>}
  {loading?<div className="empty-card">Загружаем роли…</div>:
  <div className="table-card table-scroll"><table className="roles-table"><thead><tr><th>Название</th><th>Код</th><th>Разрешения</th><th><span className="visually-hidden">Действия</span></th></tr></thead><tbody>
   {roles.map(role=><tr key={role.id}><td><strong>{role.name}</strong></td><td><code>{role.code}</code></td><td><div className="permission-tags">{role.permissions.length?role.permissions.map(code=><span key={code}>{code}</span>):<em>Нет разрешений</em>}</div></td><td><div className="row-actions"><button type="button" onClick={()=>openEditor(role)} aria-label={`Изменить роль ${role.name}`}>Изменить</button><button className="danger-action" type="button" disabled={role.code==='admin'} onClick={()=>remove(role)} aria-label={`Удалить роль ${role.name}`}>Удалить</button></div></td></tr>)}
   {!roles.length&&<tr><td colSpan={4} className="empty-cell">Роли ещё не созданы</td></tr>}
  </tbody></table></div>}
  {draft&&<div className="dialog-backdrop" role="presentation"><section className="role-dialog" role="dialog" aria-modal="true" aria-labelledby="role-dialog-title"><div className="dialog-heading"><h3 id="role-dialog-title">{draft.id?'Изменить роль':'Новая роль'}</h3><button type="button" onClick={()=>setDraft(null)} aria-label="Закрыть">×</button></div>
   <form className="role-form" onSubmit={submit}><label>Название<input value={draft.name} onChange={event=>setDraft({...draft,name:event.target.value})} required maxLength={100}/></label><label>Код<input value={draft.code} onChange={event=>setDraft({...draft,code:event.target.value})} required pattern="[a-z][a-z0-9_]*" maxLength={50} disabled={draft.code==='admin'}/></label>
   <fieldset><legend>Разрешения</legend><div className="permissions-grid">{permissions.map(permission=><label className="permission-option" key={permission.code}><input type="checkbox" checked={draft.permissions.includes(permission.code)} disabled={draft.code==='admin'&&permission.code==='can_manage_roles'} onChange={()=>togglePermission(permission.code)}/><span><strong>{permission.description}</strong><code>{permission.code}</code></span></label>)}</div></fieldset>
   <div className="dialog-actions"><button className="secondary-action" type="button" onClick={()=>setDraft(null)}>Отмена</button><button className="primary-action" type="submit" disabled={saving}>{saving?'Сохраняем…':'Сохранить'}</button></div></form>
  </section></div>}
 </div>
}

