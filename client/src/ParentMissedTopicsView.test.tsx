import { afterEach, expect, it, vi } from 'vitest'
import { act, cleanup, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ParentMissedTopicsView from './ParentMissedTopicsView'
import * as journalAPI from './parentJournal'
import * as topicsAPI from './parentMissedTopics'
import * as filesAPI from './lessonAttachments'
import { journalFixture } from './test/parentJournalFixture'

const journal=journalFixture(),student=journal.student
const history={...journal.academicYear,id:'old-year',name:'2025/2026',startsOn:'2025-09-01',endsOn:'2026-05-31',status:'completed' as const}
const topic:topicsAPI.ParentMissedTopic={lessonId:'lesson-math-1',lessonDate:'2026-11-24',position:1,topic:'Дроби',subjectId:'math',subjectName:'Математика',classId:'class-1',className:'7А'}
const topics=[topic,{...topic,lessonId:'lesson-math-2',position:2,topic:'Сложение дробей'},{...topic,lessonId:'lesson-russian',lessonDate:'2026-11-23',subjectId:'russian',subjectName:'Русский язык',topic:'Правописание\nСловарные слова'},{...topic,lessonId:'lesson-empty',lessonDate:'2026-11-22',subjectId:'russian',subjectName:'Русский язык',topic:'  '},{...topic,lessonId:'future',lessonDate:'2026-11-26',topic:'Будущая тема'}]
function setup(){
 vi.useFakeTimers({toFake:['Date']});vi.setSystemTime(new Date(2026,10,25,12))
 vi.spyOn(filesAPI,'listLessonAttachments').mockResolvedValue([])
 vi.spyOn(journalAPI,'listParentStudentPeriods').mockResolvedValue([{academicYear:history,class:journal.class,subjects:[]},{academicYear:journal.academicYear,class:journal.class,subjects:[]},{academicYear:journal.academicYear,class:{...journal.class,id:'class-2'},subjects:[]}])
 return vi.spyOn(topicsAPI,'listParentMissedTopics').mockResolvedValue(topics)
}
afterEach(()=>{cleanup();vi.restoreAllMocks();vi.useRealTimers()})

it('opens the current year, sorts missed lessons, filters subjects and exposes lesson files',async()=>{
 const list=setup()
 const file={id:'file-1',lessonId:topic.lessonId,fileName:'Дроби.pdf',contentType:'application/pdf',sizeBytes:1024,createdAt:'2026-11-24T00:00:00Z'}
 vi.mocked(filesAPI.listLessonAttachments).mockImplementation(async lessonID=>lessonID===topic.lessonId?[file]:[])
 const download=vi.spyOn(filesAPI,'downloadLessonAttachment').mockResolvedValue()
 render(<ParentMissedTopicsView student={student} onBack={vi.fn()}/>);
 expect(await screen.findByText('Сложение дробей')).toBeInTheDocument()
 expect(screen.getByLabelText('Учебный год пропущенных тем')).toHaveValue(journal.academicYear.id)
 expect(screen.getByRole('combobox',{name:'Учебный год пропущенных тем'}).querySelectorAll('option')).toHaveLength(2)
 expect(list).toHaveBeenCalledWith(student.id,journal.academicYear.id)
 expect(screen.getAllByRole('article').map(card=>card.querySelector('.parent-missed-topic')?.textContent)).toEqual(['Сложение дробей','Дроби','Правописание\nСловарные слова','Тема не заполнена'])
 expect(screen.getByText('Пропущено уроков: 4')).toBeInTheDocument()
 expect(screen.queryByText('Будущая тема')).not.toBeInTheDocument()
 expect(filesAPI.listLessonAttachments).not.toHaveBeenCalled()
 const card=screen.getByText('Дроби').closest('article') as HTMLElement
 await userEvent.click(within(card).getByRole('button',{name:'Файлы урока 2026-11-24'}))
 await userEvent.click(await screen.findByRole('button',{name:'Скачать Дроби.pdf'}))
 expect(download).toHaveBeenCalledWith(file,student.id)
 expect(screen.queryByRole('button',{name:/Удалить/})).not.toBeInTheDocument()
 await userEvent.click(screen.getByRole('button',{name:'Закрыть'}))
 await userEvent.selectOptions(screen.getByLabelText('Предмет пропущенных тем'),'math')
 expect(screen.getAllByRole('article')).toHaveLength(2)
 expect(screen.getByText('Пропущено уроков: 2')).toBeInTheDocument()
 expect(screen.queryByText('Правописание Словарные слова')).not.toBeInTheDocument()
})
it('switches year and resets the subject filter',async()=>{
 const list=setup()
 list.mockImplementation(async(_student,year)=>year===history.id?[{...topic,lessonId:'old-lesson',lessonDate:'2025-11-20',topic:'Тема прошлого года',subjectId:'history',subjectName:'История'}]:topics)
 render(<ParentMissedTopicsView student={student} onBack={vi.fn()}/>);
 await screen.findByText('Дроби')
 await userEvent.selectOptions(screen.getByLabelText('Предмет пропущенных тем'),'math')
 await userEvent.selectOptions(screen.getByLabelText('Учебный год пропущенных тем'),history.id)
 expect(await screen.findByText('Тема прошлого года')).toBeInTheDocument()
 expect(screen.getByLabelText('Предмет пропущенных тем')).toHaveValue('')
 expect(screen.queryByText('Дроби')).not.toBeInTheDocument()
 expect(list).toHaveBeenLastCalledWith(student.id,history.id)
})
it('shows an empty state and supports returning to the child list',async()=>{
 setup().mockResolvedValue([]);const back=vi.fn()
 render(<ParentMissedTopicsView student={student} onBack={back}/>);
 expect(await screen.findByText('За выбранный учебный год отмеченных пропусков нет')).toBeInTheDocument()
 await userEvent.click(screen.getByRole('button',{name:'← К ученикам'}));expect(back).toHaveBeenCalledOnce()
})
it('reports revoked journal access and allows a retry',async()=>{
 const list=setup().mockRejectedValueOnce(new Error('Доступ к журналу отозван')).mockResolvedValueOnce([topic])
 render(<ParentMissedTopicsView student={student} onBack={vi.fn()}/>);
 expect(await screen.findByRole('alert')).toHaveTextContent('Доступ к журналу отозван')
 expect(screen.queryByRole('article')).not.toBeInTheDocument()
 await userEvent.click(screen.getByRole('button',{name:'Повторить загрузку'}))
 expect(await screen.findByText('Дроби')).toBeInTheDocument();expect(list).toHaveBeenCalledTimes(2)
})
it('does not request missed topics when there is no accessible academic year',async()=>{
 const list=setup();vi.mocked(journalAPI.listParentStudentPeriods).mockResolvedValue([])
 render(<ParentMissedTopicsView student={student} onBack={vi.fn()}/>);
 expect(await screen.findByText('Нет доступных учебных годов')).toBeInTheDocument();expect(list).not.toHaveBeenCalled()
})
it('ignores a late response after choosing another year',async()=>{
 let finish:(value:topicsAPI.ParentMissedTopic[])=>void=()=>{}
 setup().mockImplementationOnce(()=>new Promise(resolve=>{finish=resolve})).mockResolvedValueOnce([{...topic,lessonId:'old',lessonDate:'2025-11-20',topic:'Правильная история'}])
 render(<ParentMissedTopicsView student={student} onBack={vi.fn()}/>);
 await screen.findByRole('option',{name:history.name})
 await userEvent.selectOptions(screen.getByLabelText('Учебный год пропущенных тем'),history.id)
 expect(await screen.findByText('Правильная история')).toBeInTheDocument()
 await act(async()=>{finish(topics)})
 expect(screen.queryByText('Дроби')).not.toBeInTheDocument()
 expect(screen.getByText('Правильная история')).toBeInTheDocument()
})
it('retries an academic-year loading failure',async()=>{
 setup();vi.mocked(journalAPI.listParentStudentPeriods).mockRejectedValueOnce(new Error('Сеть недоступна')).mockResolvedValueOnce([{academicYear:history,class:journal.class,subjects:[]}])
 vi.mocked(topicsAPI.listParentMissedTopics).mockResolvedValue([])
 render(<ParentMissedTopicsView student={student} onBack={vi.fn()}/>);
 expect(await screen.findByRole('alert')).toHaveTextContent('Сеть недоступна')
 await userEvent.click(screen.getByRole('button',{name:'Повторить загрузку'}))
 await waitFor(()=>expect(screen.getByLabelText('Учебный год пропущенных тем')).toHaveValue(history.id))
 expect(await screen.findByText('За выбранный учебный год отмеченных пропусков нет')).toBeInTheDocument()
})
