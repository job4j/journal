import { afterEach, expect, it, vi } from 'vitest'
import { act, cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ParentLessonAttachments from './ParentLessonAttachments'
import * as api from './lessonAttachments'

const file:api.LessonAttachment={id:'file-1',lessonId:'l-1',fileName:'Задание.pdf',contentType:'application/pdf',sizeBytes:2048,createdAt:'2026-11-24T00:00:00Z'}
afterEach(()=>{cleanup();vi.restoreAllMocks()})
it('hides an empty file list',async()=>{
 vi.spyOn(api,'listLessonAttachments').mockResolvedValue([])
 const {container}=render(<ParentLessonAttachments lessonID="l-1" studentID="child-1"/>);
 await waitFor(()=>expect(screen.queryByRole('status')).not.toBeInTheDocument())
 expect(container).toBeEmptyDOMElement()
})
it('reports loading errors and allows retry without exposing management controls',async()=>{
 const list=vi.spyOn(api,'listLessonAttachments').mockRejectedValueOnce(new Error('Недостаточно прав')).mockResolvedValueOnce([file])
 const download=vi.spyOn(api,'downloadLessonAttachment').mockRejectedValue(new Error('Доступ к файлу отозван'))
 render(<ParentLessonAttachments lessonID="l-1" studentID="child-1"/>);
 expect(await screen.findByRole('alert')).toHaveTextContent('Недостаточно прав')
 await userEvent.click(screen.getByRole('button',{name:'Повторить загрузку файлов'}))
 await userEvent.click(await screen.findByRole('button',{name:'Скачать Задание.pdf'}))
 expect(list).toHaveBeenCalledWith('l-1','child-1')
 expect(download).toHaveBeenCalledWith(file,'child-1')
 expect(await screen.findByRole('alert')).toHaveTextContent('Доступ к файлу отозван')
 expect(screen.queryByRole('button',{name:/Удалить|Загрузить/})).not.toBeInTheDocument()
})
it('ignores a late response after changing child and lesson',async()=>{
 let finish:(items:api.LessonAttachment[])=>void=()=>{}
 vi.spyOn(api,'listLessonAttachments').mockImplementationOnce(()=>new Promise(resolve=>{finish=resolve})).mockResolvedValueOnce([])
 const {rerender}=render(<ParentLessonAttachments key="old" lessonID="l-1" studentID="child-1"/>);
 rerender(<ParentLessonAttachments key="new" lessonID="l-2" studentID="child-2"/>);
 await act(async()=>{finish([file])})
 await waitFor(()=>expect(screen.queryByRole('status')).not.toBeInTheDocument())
 expect(screen.queryByText('Задание.pdf')).not.toBeInTheDocument()
})
