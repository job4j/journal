import{afterEach,expect,it,vi}from'vitest'
import{cleanup,render,screen}from'@testing-library/react'
import userEvent from'@testing-library/user-event'
import TeacherLessonsView from'./TeacherLessonsView'

const assignment={id:'link-1',classId:'class-1',subject:{id:'s-1',code:'MATH',name:'Математика'},responsibleTeacher:{id:'t-1',name:'Иван'+' '+'Петров',roles:['teacher']}}
const student={student:{id:'student-1',name:'Анна'+' '+'Иванова',roles:['student']},enrolledOn:'2026-09-01',leftOn:null}
afterEach(()=>{cleanup();vi.restoreAllMocks();document.cookie='journal_period_view=; Max-Age=0; Path=/'})
it('creates lesson with a material',async()=>{vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({lesson:{id:'l-1',classSubjectId:'link-1',lessonDate:'2026-09-09',position:1,topic:'Дроби',homework:null,materials:[{id:'m-1',title:'Видео',url:'https://example.test',position:1}],gradeItems:[]}}),{status:201}));render(<TeacherLessonsView assignment={assignment} quarters={[{id:'q-1',number:1,startsOn:'2026-09-09',endsOn:'2026-09-09'}]}/>);await userEvent.click(await screen.findByRole('button',{name:'Домашнее задание 2026-09-09'}));await userEvent.clear(screen.getByLabelText('Дата урока'));await userEvent.type(screen.getByLabelText('Дата урока'),'2026-09-09');await userEvent.type(screen.getByLabelText('Тема'),'Дроби');await userEvent.click(screen.getByRole('button',{name:/Материал/}));await userEvent.type(screen.getByLabelText('Название материала 1'),'Видео');await userEvent.type(screen.getByLabelText('Ссылка материала 1'),'https://example.test');await userEvent.click(screen.getByRole('button',{name:'Сохранить'}));expect(await screen.findByText('Дроби')).toBeInTheDocument()})
it('shows one week by default and saves the period view in a cookie',async()=>{
 vi.spyOn(globalThis,'fetch')
  .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}))
  .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}))
  .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}))
 render(<TeacherLessonsView assignment={assignment} quarters={[{id:'q-1',number:1,startsOn:'2000-01-01',endsOn:'2000-01-03'}]}/>)
 const week=await screen.findByRole('button',{name:'Неделя'})
 expect(week).toHaveAttribute('aria-pressed','true')
 expect(screen.queryByRole('button',{name:'Домашнее задание 2000-01-03'})).not.toBeInTheDocument()
 await userEvent.click(screen.getByRole('button',{name:'Весь период'}))
 expect(screen.getByRole('button',{name:'Домашнее задание 2000-01-03'})).toBeInTheDocument()
 expect(document.cookie).toContain('journal_period_view=period')
})
it('loads roster and saves a score with comment',async()=>{const lesson={id:'l-1',classSubjectId:'link-1',lessonDate:'2026-09-09',position:1,topic:'Дроби',materials:[],absences:[],gradeItems:[{id:'g-1',lessonId:'l-1',title:'Самостоятельная',kind:'classwork',gradingScale:'five_point',scores:[]}]};vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[lesson]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[student]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({score:{id:'score-1',gradeItemId:'g-1',studentId:'student-1',numericValue:5,teacherComment:'Отлично'}}),{status:200}));render(<TeacherLessonsView assignment={assignment} quarters={[{id:'q-1',number:1,startsOn:'2026-09-09',endsOn:'2026-09-09'}]}/>);await screen.findAllByText('Анна Иванова');await userEvent.click(screen.getByRole('button',{name:'Изменить оценку Самостоятельная'}));await userEvent.type(screen.getByLabelText('Оценка Самостоятельная'),'5');await userEvent.type(screen.getByLabelText('Комментарий Самостоятельная'),'Отлично');await userEvent.click(screen.getByRole('button',{name:'Сохранить оценку'}));expect(await screen.findByText('5*')).toBeInTheDocument();expect(screen.getByRole('button',{name:'Изменить оценку Самостоятельная'})).toHaveAttribute('title','Отлично')})
it('hides grade icons when a lesson topic is not set',async()=>{
 vi.spyOn(globalThis,'fetch')
  .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}))
  .mockResolvedValueOnce(new Response(JSON.stringify({items:[student]}),{status:200}))
  .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}))
 render(<TeacherLessonsView assignment={assignment} quarters={[{id:'q-1',number:1,startsOn:'2026-09-09',endsOn:'2026-09-09'}]}/>)
 await screen.findByText('Анна Иванова')
 expect(screen.getByLabelText('Нет темы 2026-09-09')).toBeEmptyDOMElement()
 expect(screen.queryByRole('button',{name:'Добавить Домашняя работа'})).not.toBeInTheDocument()
})
it('sets an absence independently from a score',async()=>{const lesson={id:'l-1',classSubjectId:'link-1',lessonDate:'2026-09-09',position:1,topic:'Дроби',materials:[],gradeItems:[],absences:[]};vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[lesson]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[student]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({absence:{id:'a-1',lessonId:'l-1',studentId:'student-1'}}),{status:200}));render(<TeacherLessonsView assignment={assignment} quarters={[{id:'q-1',number:1,startsOn:'2026-09-09',endsOn:'2026-09-09'}]}/>);await userEvent.click(await screen.findByRole('button',{name:'Был'}));expect(await screen.findByRole('button',{name:'Н'})).toHaveAttribute('aria-pressed','true');expect(globalThis.fetch).toHaveBeenLastCalledWith('/api/v1/lessons/l-1/absences/student-1',expect.objectContaining({method:'PUT'}))})
it('shows a quarter grade without manual editing',async()=>{const total={id:'qg-1',quarterId:'q-1',classSubjectId:'link-1',studentId:'student-1',gradingScale:'five_point',numericValue:5,teacherComment:'Уверенный результат'};vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[student]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({items:[total]}),{status:200}));render(<TeacherLessonsView assignment={assignment} quarters={[{id:'q-1',number:1}]}/>);await screen.findByText('5*');expect(screen.getByTitle('Уверенный результат')).toHaveTextContent('5*');expect(screen.queryByRole('button',{name:/Итог/})).not.toBeInTheDocument()})
