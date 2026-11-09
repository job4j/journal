import { afterEach, expect, it, vi } from 'vitest'
import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ClassScheduleView from './ClassScheduleView'
import * as api from './classSchedule'
import { ApiError } from './api'
import { journalFixture } from './test/parentJournalFixture'

const journal = journalFixture()
const subjects = journal.subjects.map(s => s.classSubject)
function setup() {
  vi.useFakeTimers({toFake:['Date']})
  vi.setSystemTime(new Date(2026,10,9,12))
  const entry: api.ScheduleEntry = {id:'entry-1',day:'2026-11-09',position:1,classSubjectId:subjects[0].id,title:'Математика',notes:'',startsAt:'08:45',endsAt:'09:30'}
  const weeks = new Map<string, api.ScheduleWeek>([['2026-11-09',{weekStart:'2026-11-09',revision:1,entries:[entry]}]])
  const get = vi.spyOn(api,'getClassSchedule').mockImplementation(async (_id,start) => structuredClone(weeks.get(start) ?? {weekStart:start,revision:0,entries:[]}))
  const save = vi.spyOn(api,'replaceClassSchedule').mockImplementation(async (_id,week) => {
    const value = {...structuredClone(week),revision:week.revision+1,entries:week.entries.map((e,i) => ({...e,id:e.id || `new-${i}`}))}
    weeks.set(week.weekStart,value)
    return value
  })
  return {weeks,get,save}
}
function view() {return render(<ClassScheduleView item={journal.class} year={journal.academicYear} subjects={subjects} onBack={vi.fn()} />)}
afterEach(() => {cleanup();vi.restoreAllMocks();vi.useRealTimers()})

it('creates an extracurricular entry and saves only the selected week', async () => {
  const {save} = setup();view()
  await screen.findByRole('button',{name:'Изменить Математика, 2026-11-09, урок 1'})
  await userEvent.click(screen.getByRole('button',{name:'Добавить занятие 2026-11-11, урок 5'}))
  await userEvent.selectOptions(screen.getByLabelText('Предмет занятия'),'')
  await userEvent.type(screen.getByLabelText('Название занятия'),'Классный час')
  await userEvent.type(screen.getByLabelText('Примечание занятия'),'Раз в месяц')
  await userEvent.click(screen.getByRole('button',{name:'Применить'}))
  expect(save).not.toHaveBeenCalled()
  expect(screen.getByRole('button',{name:'Следующая неделя расписания'})).toBeDisabled()
  await userEvent.click(screen.getByRole('button',{name:'Сохранить неделю'}))
  await waitFor(() => expect(save).toHaveBeenCalledWith(journal.class.id,expect.objectContaining({weekStart:'2026-11-09',revision:1,entries:expect.arrayContaining([expect.objectContaining({day:'2026-11-11',position:5,title:'Классный час',classSubjectId:null,notes:'Раз в месяц'})])})))
  expect(await screen.findByText('Расписание этой недели сохранено')).toBeInTheDocument()
  await userEvent.click(screen.getByRole('button',{name:'Следующая неделя расписания'}))
  await waitFor(() => expect(screen.queryByText('Классный час')).not.toBeInTheDocument())
  await userEvent.click(screen.getByRole('button',{name:'Предыдущая неделя расписания'}))
  expect(await screen.findByRole('button',{name:'Изменить Классный час, 2026-11-11, урок 5'})).toBeInTheDocument()
})

it('edits and removes an entry while keeping another group in the same slot', async () => {
  const {weeks,save} = setup()
  weeks.get('2026-11-09')!.entries.push({id:'group',day:'2026-11-09',position:1,classSubjectId:null,title:'Труд',notes:'Мальчики',startsAt:'08:45',endsAt:'09:30'})
  view()
  await userEvent.click(await screen.findByRole('button',{name:'Изменить Математика, 2026-11-09, урок 1'}))
  await userEvent.type(screen.getByLabelText('Примечание занятия'),'Другая группа')
  await userEvent.click(screen.getByRole('button',{name:'Применить'}))
  expect(screen.getByText('Другая группа')).toBeInTheDocument()
  await userEvent.click(screen.getByRole('button',{name:'Изменить Математика, 2026-11-09, урок 1'}))
  await userEvent.click(screen.getByRole('button',{name:'Удалить занятие'}))
  expect(screen.getByText('Труд')).toBeInTheDocument()
  await userEvent.click(screen.getByRole('button',{name:'Сохранить неделю'}))
  await waitFor(() => expect(save).toHaveBeenCalledWith(journal.class.id,expect.objectContaining({entries:[expect.objectContaining({id:'group',title:'Труд'})]})))
})

it('keeps unsaved edits after a conflict and allows an explicit reload', async () => {
  const {save,get,weeks} = setup();save.mockRejectedValue(new ApiError('Расписание изменилось. Загрузите неделю заново',409))
  view()
  await userEvent.click(await screen.findByRole('button',{name:'Изменить Математика, 2026-11-09, урок 1'}))
  await userEvent.type(screen.getByLabelText('Примечание занятия'),'Моя правка')
  await userEvent.click(screen.getByRole('button',{name:'Применить'}))
  await userEvent.click(screen.getByRole('button',{name:'Сохранить неделю'}))
  expect(await screen.findByRole('alert')).toHaveTextContent('Расписание изменилось')
  expect(screen.getByText('Моя правка')).toBeInTheDocument()
  weeks.get('2026-11-09')!.revision=3
  await userEvent.click(screen.getByRole('button',{name:'Загрузить неделю заново'}))
  await waitFor(() => expect(get).toHaveBeenCalledTimes(2))
  await waitFor(() => expect(screen.queryByText('Моя правка')).not.toBeInTheDocument())
})

it.each([{today:new Date(2026,7,31,12),inside:'2026-09-01',outside:'2026-08-31'},{today:new Date(2027,5,1,12),inside:'2027-05-31',outside:'2027-06-01'}])('does not offer lessons outside the academic year: $outside', async ({today,inside,outside}) => {
  setup();vi.setSystemTime(today);view()
  expect(await screen.findByRole('button',{name:`Добавить занятие ${inside}, урок 1`})).toBeInTheDocument()
  expect(screen.queryByRole('button',{name:`Добавить занятие ${outside}, урок 1`})).not.toBeInTheDocument()
})

it('shows access errors without displaying a timetable', async () => {
  setup().get.mockRejectedValue(new ApiError('Расписание доступно только администратору',403));view()
  expect(await screen.findByRole('alert')).toHaveTextContent('только администратору')
  expect(screen.queryByRole('table')).not.toBeInTheDocument()
})
