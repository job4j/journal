import{request}from'./api'
import{ClassRecord}from'./classes'
import{ClassSubjectRecord}from'./classSubjects'
export async function listTeacherClasses():Promise<ClassRecord[]>{return(await request<{items:ClassRecord[]}>('/api/v1/teacher/classes')).items}
export async function listTeacherClassSubjects(classID:string):Promise<ClassSubjectRecord[]>{return(await request<{items:ClassSubjectRecord[]}>(`/api/v1/teacher/classes/${classID}/subjects`)).items}
