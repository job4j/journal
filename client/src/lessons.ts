import{request}from'./api'
export interface LessonMaterial{id:string;title:string;url:string;position:number}
export interface LessonRecord{id:string;classSubjectId:string;lessonDate:string;position:number;topic:string;homework?:string|null;materials:LessonMaterial[];gradeItems:unknown[]}
export interface LessonDraft{lessonDate:string;position:number;topic:string;homework?:string;materials:{title:string;url:string;position:number}[]}
export async function listLessons(id:string,from:string,to:string):Promise<LessonRecord[]>{const query=new URLSearchParams();if(from)query.set('dateFrom',from);if(to)query.set('dateTo',to);return(await request<{items:LessonRecord[]}>(`/api/v1/class-subjects/${id}/lessons?${query}`)).items}
export async function createLesson(id:string,value:LessonDraft):Promise<LessonRecord>{return(await request<{lesson:LessonRecord}>(`/api/v1/class-subjects/${id}/lessons`,{method:'POST',body:JSON.stringify(value)})).lesson}
