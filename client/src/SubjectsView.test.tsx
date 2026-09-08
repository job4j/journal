import{afterEach,describe,expect,it,vi}from'vitest'
import{cleanup,render,screen}from'@testing-library/react'
import userEvent from'@testing-library/user-event'
import SubjectsView from'./SubjectsView'
afterEach(()=>{cleanup();vi.restoreAllMocks()})
describe('subjects',()=>{
 it('loads subjects',async()=>{vi.spyOn(globalThis,'fetch').mockResolvedValue(new Response(JSON.stringify({items:[{id:'1',code:'math',name:'Математика'}]}),{status:200}));render(<SubjectsView/>);expect(await screen.findByText('Математика')).toBeInTheDocument();expect(screen.getByText('math')).toBeInTheDocument()})
 it('creates a subject',async()=>{vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({subject:{id:'1',code:'history',name:'История'}}),{status:201}));render(<SubjectsView/>);await userEvent.click(await screen.findByRole('button',{name:/Создать предмет/}));await userEvent.type(screen.getByLabelText('Название'),'История');await userEvent.type(screen.getByLabelText('Код'),'history');await userEvent.click(screen.getByRole('button',{name:'Сохранить'}));expect(await screen.findByText('История')).toBeInTheDocument()})
 it('shows duplicate error',async()=>{vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200})).mockResolvedValueOnce(new Response(JSON.stringify({message:'Предмет уже существует'}),{status:409}));render(<SubjectsView/>);await userEvent.click(await screen.findByRole('button',{name:/Создать предмет/}));await userEvent.type(screen.getByLabelText('Название'),'История');await userEvent.type(screen.getByLabelText('Код'),'history');await userEvent.click(screen.getByRole('button',{name:'Сохранить'}));expect(await screen.findByRole('alert')).toHaveTextContent('Предмет уже существует')})
})
