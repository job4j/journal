import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import AcademicYearsView from './AcademicYearsView'

afterEach(()=>{cleanup();vi.restoreAllMocks()})

describe('academic years',()=>{
 it('loads the academic year list',async()=>{
  vi.spyOn(globalThis,'fetch').mockResolvedValue(new Response(JSON.stringify({items:[{id:'year-1',name:'2026/2027',startsOn:'2026-09-01',endsOn:'2027-05-31',status:'active',quarters:[1,2,3,4].map(number=>({id:`q-${number}`,number,startsOn:'2026-09-01',endsOn:'2026-10-01'}))}]}),{status:200}))
  render(<AcademicYearsView/>)
  expect(await screen.findByText('2026/2027')).toBeInTheDocument()
  expect(screen.getByText('Активный')).toBeInTheDocument()
 })

 it('creates an academic year without periods',async()=>{
  const starts=['2027-09-01','2027-11-01','2028-01-10','2028-03-20'];const ends=['2027-10-20','2027-12-20','2028-03-10','2028-05-31']
  const saved={id:'year-1',name:'2027/2028',startsOn:'2027-09-01',endsOn:'2028-05-31',status:'planned',quarters:[1,2,3,4].map((number,index)=>({id:`q-${number}`,number,startsOn:starts[index],endsOn:ends[index]}))}
  vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({academicYear:saved}),{status:201}))
  render(<AcademicYearsView/>);await userEvent.click(await screen.findByRole('button',{name:/Создать учебный год/}))
  await userEvent.type(screen.getByLabelText('Название'),'2027/2028');await userEvent.type(screen.getByLabelText('Начало'),'2027-09-01');await userEvent.type(screen.getByLabelText('Окончание'),'2028-05-31');await userEvent.click(screen.getByRole('button',{name:'Сохранить'}))
  expect(await screen.findByText('2027/2028')).toBeInTheDocument()
 })

 it('creates and displays a named period',async()=>{
  const year={id:'year-1',name:'2026/2027',startsOn:'2026-09-01',endsOn:'2027-05-31',status:'active',quarters:[]}
  const period={id:'q-1',number:1,name:'Осень',startsOn:'2026-09-01',endsOn:'2026-10-31'}
  vi.spyOn(globalThis,'fetch')
   .mockResolvedValueOnce(new Response(JSON.stringify({items:[year]}),{status:200}))
   .mockResolvedValueOnce(new Response(JSON.stringify({quarter:period}),{status:201}))
  render(<AcademicYearsView/>)
  await userEvent.click(await screen.findByRole('button',{name:'+ Период'}))
  await userEvent.type(screen.getByLabelText('Название периода'),'Осень')
  await userEvent.clear(screen.getByLabelText('Окончание'))
  await userEvent.type(screen.getByLabelText('Окончание'),'2026-10-31')
  await userEvent.click(screen.getByRole('button',{name:'Сохранить'}))
  expect(await screen.findByText(/Осень:/)).toBeInTheDocument()
  expect(globalThis.fetch).toHaveBeenLastCalledWith(
   '/api/v1/academic-years/year-1/quarters',
   expect.objectContaining({body:expect.stringContaining('"name":"Осень"')}),
  )
 })})
