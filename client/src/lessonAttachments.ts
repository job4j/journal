import { ApiError, request, type ErrorResponse } from './api'

export const maxAttachmentBytes = 20 * 1024 * 1024
export interface LessonAttachment { id: string; lessonId: string; fileName: string; contentType: string; sizeBytes: number; createdAt: string }
function parentPrefix(studentID?: string) { return studentID ? `/parent/students/${encodeURIComponent(studentID)}` : '' }
export async function listLessonAttachments(lessonID: string, studentID?: string): Promise<LessonAttachment[]> {
  return (await request<{items: LessonAttachment[]}>(`/api/v1${parentPrefix(studentID)}/lessons/${encodeURIComponent(lessonID)}/attachments`)).items
}
export async function uploadLessonAttachment(lessonID: string, file: File): Promise<LessonAttachment> {
  const body = new FormData()
  body.append('file', file)
  return (await request<{attachment: LessonAttachment}>(`/api/v1/lessons/${encodeURIComponent(lessonID)}/attachments`, { method: 'POST', body })).attachment
}
export async function deleteLessonAttachment(id: string): Promise<void> {
  await request(`/api/v1/lesson-attachments/${encodeURIComponent(id)}`, { method: 'DELETE' })
}
export async function downloadLessonAttachment(item: LessonAttachment, studentID?: string): Promise<void> {
  const basePath = import.meta.env.BASE_URL.replace(/\/$/, '')
  let response: Response
  try {
    response = await fetch(`${basePath}/api/v1${parentPrefix(studentID)}/lesson-attachments/${encodeURIComponent(item.id)}/download`, { credentials: 'include' })
  } catch { throw new ApiError('Не удалось связаться с сервером') }
  if (!response.ok) {
    if (response.status === 401) window.dispatchEvent(new CustomEvent('journal:session-expired'))
    const error = await response.json().catch(() => null) as ErrorResponse | null
    throw new ApiError(error?.message ?? 'Не удалось скачать файл', response.status, error?.code)
  }
  const url = URL.createObjectURL(await response.blob())
  const link = document.createElement('a')
  link.href = url; link.download = item.fileName
  document.body.append(link); link.click(); link.remove()
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}
