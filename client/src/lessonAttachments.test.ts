import { afterEach, expect, it, vi } from 'vitest'
import { deleteLessonAttachment, downloadLessonAttachment, listLessonAttachments, uploadLessonAttachment, type LessonAttachment } from './lessonAttachments'

const item:LessonAttachment={id:'a-1',lessonId:'l-1',fileName:'work.pdf',contentType:'application/pdf',sizeBytes:3,createdAt:'2026-11-24T00:00:00Z'}
afterEach(()=>vi.restoreAllMocks())
it('sends a multipart file without overriding its boundary content type',async()=>{
 const fetch=vi.spyOn(globalThis,'fetch').mockResolvedValue(new Response(JSON.stringify({attachment:item}),{status:201}))
 const file=new File(['abc'],'work.pdf',{type:'application/pdf'})
 expect(await uploadLessonAttachment('l-1',file)).toEqual(item)
 const [url,init]=fetch.mock.calls[0]
 expect(url).toBe('/api/v1/lessons/l-1/attachments')
 expect(init?.credentials).toBe('include')
 expect(init?.method).toBe('POST')
 expect(init?.body).toBeInstanceOf(FormData)
 expect((init?.body as FormData).get('file')).toBe(file)
 expect(new Headers(init?.headers).has('Content-Type')).toBe(false)
})
it('lists through the child endpoint and deletes through the managed endpoint',async()=>{
 const fetch=vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({items:[item]}))).mockResolvedValueOnce(new Response(null,{status:204}))
 expect(await listLessonAttachments('l-1','child-1')).toEqual([item])
 expect(fetch).toHaveBeenLastCalledWith('/api/v1/parent/students/child-1/lessons/l-1/attachments',expect.anything())
 await deleteLessonAttachment('a-1')
 expect(fetch).toHaveBeenLastCalledWith('/api/v1/lesson-attachments/a-1',expect.objectContaining({method:'DELETE',credentials:'include'}))
})
it('reports a revoked download without creating a local file',async()=>{
 const fetch=vi.spyOn(globalThis,'fetch').mockResolvedValue(new Response(JSON.stringify({message:'Недостаточно прав'}),{status:403}))
 await expect(downloadLessonAttachment(item,'child-1')).rejects.toMatchObject({status:403,message:'Недостаточно прав'})
 expect(fetch).toHaveBeenCalledWith('/api/v1/parent/students/child-1/lesson-attachments/a-1/download',expect.objectContaining({credentials:'include'}))
})
