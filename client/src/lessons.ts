import{request}from'./api'
export interface LessonMaterial{id:string;title:string;url:string;position:number}
export interface Score{id:string;gradeItemId:string;studentId:string;numericValue?:number|null;textValue?:string|null;teacherComment?:string|null}
export interface Absence{id:string;lessonId:string;studentId:string}
export interface GradeItem{id:string;lessonId:string;title:string;kind:'homework'|'classwork'|'knowledge_check'|'other';gradingScale:'five_point'|'points'|'pass_fail';maxScore?:number|null;scores:Score[]}
export type GradeItemDraft=Omit<GradeItem,'id'|'lessonId'|'scores'>
export interface LessonRecord{id:string;classSubjectId:string;lessonDate:string;position:number;topic:string;homework?:string|null;materials:LessonMaterial[];gradeItems:GradeItem[];absences:Absence[]}
export interface LessonDraft{lessonDate:string;position:number;topic:string;homework?:string;materials:{title:string;url:string;position:number}[]}
export async function listLessons(id:string,from:string,to:string):Promise<LessonRecord[]>{const query=new URLSearchParams();if(from)query.set('dateFrom',from);if(to)query.set('dateTo',to);return(await request<{items:LessonRecord[]}>(`/api/v1/class-subjects/${id}/lessons?${query}`)).items}
export async function updateLesson(id:string,value:LessonDraft):Promise<LessonRecord>{return(await request<{lesson:LessonRecord}>('/api/v1/lessons/'+id,{method:'PATCH',body:JSON.stringify(value)})).lesson}
export async function createLesson(id:string,value:LessonDraft):Promise<LessonRecord>{return(await request<{lesson:LessonRecord}>(`/api/v1/class-subjects/${id}/lessons`,{method:'POST',body:JSON.stringify(value)})).lesson}
export async function createGradeItem(lessonID:string,value:GradeItemDraft):Promise<GradeItem>{return(await request<{gradeItem:GradeItem}>(`/api/v1/lessons/${lessonID}/grade-items`,{method:'POST',body:JSON.stringify(value)})).gradeItem}
export async function putStudentScore(gradeItemID:string,studentID:string,value:{numericValue?:number;textValue?:string;teacherComment?:string}):Promise<Score>{return(await request<{score:Score}>(`/api/v1/grade-items/${gradeItemID}/scores/${studentID}`,{method:'PUT',body:JSON.stringify(value)})).score}
export async function putStudentAbsence(lessonID:string,studentID:string):Promise<Absence>{return(await request<{absence:Absence}>(`/api/v1/lessons/${lessonID}/absences/${studentID}`,{method:'PUT'})).absence}
export async function deleteStudentAbsence(lessonID:string,studentID:string):Promise<void>{await request(`/api/v1/lessons/${lessonID}/absences/${studentID}`,{method:'DELETE'})}
