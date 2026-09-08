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

 it('creates an academic year with four quarters',async()=>{
  const starts=['2027-09-01','2027-11-01','2028-01-10','2028-03-20'];const ends=['2027-10-20','2027-12-20','2028-03-10','2028-05-31']
  const saved={id:'year-1',name:'2027/2028',startsOn:'2027-09-01',endsOn:'2028-05-31',status:'planned',quarters:[1,2,3,4].map((number,index)=>({id:`q-${number}`,number,startsOn:starts[index],endsOn:ends[index]}))}
  vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({academicYear:saved}),{status:201}))
  render(<AcademicYearsView/>);await userEvent.click(await screen.findByRole('button',{name:/Создать учебный год/}))
  await userEvent.type(screen.getByLabelText('Название'),'2027/2028');await userEvent.type(screen.getByLabelText('Начало'),'2027-09-01');await userEvent.type(screen.getByLabelText('Окончание'),'2028-05-31')
  for(const number of [1,2,3,4]){await userEvent.type(screen.getByLabelText(`${number} четверть: начало`),starts[number-1]);await userEvent.type(screen.getByLabelText(`${number} четверть: окончание`),ends[number-1])}
  await userEvent.click(screen.getByRole('button',{name:'Сохранить'}))
  expect(await screen.findByText('2027/2028')).toBeInTheDocument()
 })
})
