import{request}from'./api'
import{UserRecord}from'./users'
export interface SubjectRecord{id:string;code:string;name:string}
export interface ClassSubjectRecord{id:string;classId:string;subject:SubjectRecord;responsibleTeacher:Pick<UserRecord,'id'|'name'|'roles'>}
export async function listClassSubjects(classID:string):Promise<ClassSubjectRecord[]>{return(await request<{items:ClassSubjectRecord[]}>(`/api/v1/classes/${classID}/subjects`)).items}
export async function assignSubjectToClass(classID:string,subjectId:string,responsibleTeacherId:string):Promise<ClassSubjectRecord>{return(await request<{classSubject:ClassSubjectRecord}>(`/api/v1/classes/${classID}/subjects`,{method:'POST',body:JSON.stringify({subjectId,responsibleTeacherId})})).classSubject}
export async function updateClassSubject(id:string,responsibleTeacherId:string):Promise<ClassSubjectRecord>{return(await request<{classSubject:ClassSubjectRecord}>(`/api/v1/class-subjects/${id}`,{method:'PATCH',body:JSON.stringify({responsibleTeacherId})})).classSubject}
