import { request } from './api'

export interface Subject { id:string; code:string; name:string }
export interface SubjectDraft { code:string; name:string }
export async function listSubjects():Promise<Subject[]>{return(await request<{items:Subject[]}>('/api/v1/subjects')).items}
export async function createSubject(value:SubjectDraft):Promise<Subject>{return(await request<{subject:Subject}>('/api/v1/subjects',{method:'POST',body:JSON.stringify(value)})).subject}
