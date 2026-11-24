import { afterEach, expect, it, vi } from 'vitest'
import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import LessonAttachmentsButton from './LessonAttachmentsButton'
import * as api from './lessonAttachments'

const item:api.LessonAttachment={id:'a-1',lessonId:'l-1',fileName:'Задание.pdf',contentType:'application/pdf',sizeBytes:123,createdAt:'2026-11-24T00:00:00Z'}
afterEach(()=>{cleanup();vi.restoreAllMocks()})
it('disables attachments until a lesson topic is saved',()=>{
 render(<LessonAttachmentsButton date="2026-11-24"/>);
 expect(screen.getByRole('button',{name:'Файлы урока 2026-11-24'})).toBeDisabled()
 expect(screen.getByRole('button')).toHaveAttribute('title','Сначала сохраните тему урока')
})
it('uploads downloads and removes a lesson file',async()=>{
 vi.spyOn(api,'listLessonAttachments').mockResolvedValue([])
 const upload=vi.spyOn(api,'uploadLessonAttachment').mockResolvedValue(item)
 const remove=vi.spyOn(api,'deleteLessonAttachment').mockResolvedValue()
 const download=vi.spyOn(api,'downloadLessonAttachment').mockResolvedValue()
 render(<LessonAttachmentsButton lessonID="l-1" date="2026-11-24"/>);
 await userEvent.click(screen.getByRole('button',{name:'Файлы урока 2026-11-24'}))
 await screen.findByText('Файлы пока не прикреплены')
 const file=new File(['pdf'],'Задание.pdf',{type:'application/pdf'})
 await userEvent.upload(screen.getByLabelText('Файл'),file)
 await userEvent.click(screen.getByRole('button',{name:'Загрузить'}))
 expect(await screen.findByRole('button',{name:'Скачать Задание.pdf'})).toBeInTheDocument()
 expect(upload).toHaveBeenCalledWith('l-1',file)
 await userEvent.click(screen.getByRole('button',{name:'Скачать Задание.pdf'}))
 expect(download).toHaveBeenCalledWith(item,undefined)
 await userEvent.click(screen.getByRole('button',{name:'Удалить Задание.pdf'}))
 await waitFor(()=>expect(screen.queryByText('Задание.pdf')).not.toBeInTheDocument())
 expect(remove).toHaveBeenCalledWith('a-1')
 await userEvent.click(screen.getByRole('button',{name:'Закрыть'}))
 expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
})
it('gives a parent read-only access through their child endpoint',async()=>{
 const list=vi.spyOn(api,'listLessonAttachments').mockResolvedValue([item])
 const download=vi.spyOn(api,'downloadLessonAttachment').mockResolvedValue()
 render(<LessonAttachmentsButton lessonID="l-1" date="2026-11-24" studentID="child-1"/>);
 await userEvent.click(screen.getByRole('button',{name:'Файлы урока 2026-11-24'}))
 await userEvent.click(await screen.findByRole('button',{name:'Скачать Задание.pdf'}))
 expect(list).toHaveBeenCalledWith('l-1','child-1')
 expect(download).toHaveBeenCalledWith(item,'child-1')
 expect(screen.queryByLabelText('Файл')).not.toBeInTheDocument()
 expect(screen.queryByRole('button',{name:/Удалить/})).not.toBeInTheDocument()
})
it('keeps a failed deletion available for retry',async()=>{
 vi.spyOn(api,'listLessonAttachments').mockResolvedValue([item])
 const remove=vi.spyOn(api,'deleteLessonAttachment').mockRejectedValueOnce(new Error('Хранилище недоступно')).mockResolvedValueOnce()
 render(<LessonAttachmentsButton lessonID="l-1" date="2026-11-24"/>);
 await userEvent.click(screen.getByRole('button',{name:'Файлы урока 2026-11-24'}))
 await userEvent.click(await screen.findByRole('button',{name:'Удалить Задание.pdf'}))
 expect(await screen.findByRole('alert')).toHaveTextContent('Хранилище недоступно')
 await userEvent.click(screen.getByRole('button',{name:'Удалить Задание.pdf'}))
 await waitFor(()=>expect(screen.queryByText('Задание.pdf')).not.toBeInTheDocument())
 expect(remove).toHaveBeenCalledTimes(2)
})
it('rejects an empty file before uploading',async()=>{
 vi.spyOn(api,'listLessonAttachments').mockResolvedValue([])
 const upload=vi.spyOn(api,'uploadLessonAttachment')
 render(<LessonAttachmentsButton lessonID="l-1" date="2026-11-24"/>);
 await userEvent.click(screen.getByRole('button',{name:'Файлы урока 2026-11-24'}))
 await userEvent.upload(await screen.findByLabelText('Файл'),new File([],'empty.txt'))
 await userEvent.click(screen.getByRole('button',{name:'Загрузить'}))
 expect(await screen.findByRole('alert')).toHaveTextContent('непустой файл')
 expect(upload).not.toHaveBeenCalled()
})
