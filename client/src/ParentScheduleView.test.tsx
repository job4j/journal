import { afterEach, expect, it, vi } from 'vitest'
import { act, cleanup, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ParentScheduleView from './ParentScheduleView'
import * as journalAPI from './parentJournal'
import * as scheduleAPI from './parentSchedule'
import { journalFixture } from './test/parentJournalFixture'
import type { ParentScheduleWeek } from './parentSchedule'

const journal=journalFixture()
function entry(id: string,day='2026-11-09',position=3): scheduleAPI.ParentScheduleEntry {
  return {id,day,position,classSubjectId:'link-1',title:'География',notes:'',startsAt:'10:45',endsAt:'11:30',homework:'§ 2, упражнение 7\nПовторить тему',homeworkFrom:'2026-11-06'}
}
function setup() {
  vi.useFakeTimers({toFake:['Date']});vi.setSystemTime(new Date(2026,10,9,12))
  vi.spyOn(journalAPI,'listParentStudentPeriods').mockResolvedValue([{academicYear:journal.academicYear,class:journal.class,subjects:journal.subjects.map(s=>s.classSubject)}])
  return vi.spyOn(scheduleAPI,'getParentSchedule').mockImplementation(async(_student,_class,start)=>({weekStart:start,entries:start==='2026-11-09' ? [entry('first'),{...entry('second','2026-11-09',4),startsAt:'12:00',endsAt:'12:45',homework:'Не должно появиться у второго урока'}, {...entry('extra','2026-11-11',5),classSubjectId:null,title:'Классный час',notes:'Раз в месяц',homework:null,homeworkFrom:null}] : []}))
}
afterEach(()=>{cleanup();vi.restoreAllMocks();vi.useRealTimers()})

it('shows day cards with times and homework only on the first daily subject lesson',async()=>{
  const get=setup()
  render(<ParentScheduleView student={journal.student} onBack={vi.fn()} />)
  const monday=await screen.findByRole('region',{name:'понедельник, 9 ноября'})
  expect(within(monday).getAllByText('География')).toHaveLength(2)
  expect(within(monday).getByText('10:45–11:30')).toBeInTheDocument()
  expect(within(monday).getByText('12:00–12:45')).toBeInTheDocument()
  expect(within(monday).getAllByText('Домашнее задание')).toHaveLength(1)
  expect(within(monday).getByText('§ 2, упражнение 7 Повторить тему')).toBeInTheDocument()
  expect(within(monday).getByText('Задано 6 ноября')).toBeInTheDocument()
  expect(screen.queryByText('Не должно появиться у второго урока')).not.toBeInTheDocument()
  const wednesday=screen.getByRole('region',{name:'среда, 11 ноября'})
  expect(within(wednesday).getByText('Классный час')).toBeInTheDocument()
  expect(within(wednesday).getByText('Раз в месяц')).toBeInTheDocument()
  expect(within(wednesday).queryByText('Домашнее задание')).not.toBeInTheDocument()
  expect(screen.queryByRole('table')).not.toBeInTheDocument()
  expect(screen.queryByRole('button',{name:'Сохранить неделю'})).not.toBeInTheDocument()
  expect(get).toHaveBeenCalledWith(journal.student.id,journal.class.id,'2026-11-09')
})

it('switches weeks without retaining previous homework',async()=>{
  const get=setup()
  render(<ParentScheduleView student={journal.student} onBack={vi.fn()} />)
  await screen.findByText('Задано 6 ноября')
  await userEvent.click(screen.getByRole('button',{name:'Следующая неделя'}))
  expect(await screen.findByRole('region',{name:'понедельник, 16 ноября'})).toHaveTextContent('Занятий нет')
  expect(screen.queryByText('Задано 6 ноября')).not.toBeInTheDocument()
  expect(get).toHaveBeenLastCalledWith(journal.student.id,journal.class.id,'2026-11-16')
  await userEvent.click(screen.getByRole('button',{name:'Текущая неделя'}))
  expect(await screen.findByText('Задано 6 ноября')).toBeInTheDocument()
})

it('shows unassigned homework without inventing an older assignment',async()=>{
  setup().mockResolvedValue({weekStart:'2026-11-09',entries:[{...entry('empty'),homework:null,homeworkFrom:null}]})
  render(<ParentScheduleView student={journal.student} onBack={vi.fn()} />)
  expect(await screen.findByText('Не задано')).toBeInTheDocument()
  expect(screen.queryByText(/Задано /)).not.toBeInTheDocument()
})

it('hides dates outside the academic year and disables the last-week arrow',async()=>{
  setup();vi.setSystemTime(new Date(2027,5,1,12))
  render(<ParentScheduleView student={journal.student} onBack={vi.fn()} />)
  expect(await screen.findByRole('region',{name:'понедельник, 31 мая'})).toBeInTheDocument()
  expect(screen.queryByRole('region',{name:'вторник, 1 июня'})).not.toBeInTheDocument()
  expect(screen.getByRole('button',{name:'Следующая неделя'})).toBeDisabled()
})

it('shows access errors and no lessons for an unavailable child',async()=>{
  setup().mockRejectedValue(new Error('Нет доступа к расписанию'))
  render(<ParentScheduleView student={journal.student} onBack={vi.fn()} />)
  expect(await screen.findByRole('alert')).toHaveTextContent('Нет доступа к расписанию')
  expect(screen.queryByText('География')).not.toBeInTheDocument()
})

it('ignores a late response after switching to another accessible class',async()=>{
  const get=setup()
  let resolveFirst!: (week: ParentScheduleWeek)=>void
  get.mockImplementationOnce(()=>new Promise(resolve=>{resolveFirst=resolve})).mockResolvedValue({weekStart:'2026-11-09',entries:[{...entry('new-class'),title:'Литература',homework:'Новое задание'}]})
  vi.mocked(journalAPI.listParentStudentPeriods).mockResolvedValue([
    {academicYear:journal.academicYear,class:journal.class,subjects:[]},
    {academicYear:journal.academicYear,class:{...journal.class,id:'other-class',name:'7Б'},subjects:[]},
  ])
  render(<ParentScheduleView student={journal.student} onBack={vi.fn()} />)
  await userEvent.selectOptions(await screen.findByLabelText('Класс расписания'),'other-class')
  expect(await screen.findByText('Новое задание')).toBeInTheDocument()
  await act(async()=>{resolveFirst({weekStart:'2026-11-09',entries:[{...entry('late'),homework:'Задание из старого ответа'}]})})
  expect(screen.queryByText('Задание из старого ответа')).not.toBeInTheDocument()
})
