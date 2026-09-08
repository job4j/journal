import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import RolesView from './RolesView'

afterEach(()=>{cleanup();vi.restoreAllMocks()})
describe('roles',()=>{
 it('loads roles and exposes CRUD actions',async()=>{
  vi.spyOn(globalThis,'fetch').mockImplementation(async input=>{
   const url=String(input)
   if(url.endsWith('/permissions'))return new Response(JSON.stringify({permissions:[{code:'can_view_user',description:'Просмотр пользователей'}]}),{status:200,headers:{'Content-Type':'application/json'}})
   return new Response(JSON.stringify({roles:[{id:crypto.randomUUID(),code:'teacher',name:'Учитель',permissions:['can_view_user']}]}),{status:200,headers:{'Content-Type':'application/json'}})
  })
  render(<RolesView/>)
  expect(await screen.findByText('Учитель')).toBeInTheDocument()
  expect(screen.getByRole('button',{name:'Создать роль'})).toBeInTheDocument()
  expect(screen.getByRole('button',{name:'Изменить роль Учитель'})).toBeInTheDocument()
  expect(screen.getByRole('button',{name:'Удалить роль Учитель'})).toBeInTheDocument()
 })
 it('creates a role and adds it to the list',async()=>{
  vi.spyOn(globalThis,'fetch').mockImplementation(async(input,init)=>{
   const url=String(input)
   if(init?.method==='POST')return new Response(JSON.stringify({role:{id:crypto.randomUUID(),code:'editor',name:'Редактор',permissions:['can_view_user']}}),{status:201,headers:{'Content-Type':'application/json'}})
   if(url.endsWith('/permissions'))return new Response(JSON.stringify({permissions:[{code:'can_view_user',description:'Просмотр пользователей'}]}),{status:200,headers:{'Content-Type':'application/json'}})
   return new Response(JSON.stringify({roles:[]}),{status:200,headers:{'Content-Type':'application/json'}})
  })
  render(<RolesView/>);await screen.findByText('Роли ещё не созданы')
  await userEvent.click(screen.getByRole('button',{name:'Создать роль'}))
  await userEvent.type(screen.getByLabelText('Название'),'Редактор')
  await userEvent.type(screen.getByLabelText('Код'),'editor')
  await userEvent.click(screen.getByText('Просмотр пользователей'))
  await userEvent.click(screen.getByRole('button',{name:'Сохранить'}))
  expect(await screen.findByText('Редактор')).toBeInTheDocument()
 })
})
