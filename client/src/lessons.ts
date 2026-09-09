import{request}from'./api'
export interface LessonMaterial{id:string;title:string;url:string;position:number}
export interface GradeItem{id:string;lessonId:string;title:string;kind:'homework'|'classwork'|'knowledge_check'|'other';gradingScale:'five_point'|'points'|'pass_fail';maxScore?:number|null}
export interface LessonRecord{id:string;classSubjectId:string;lessonDate:string;position:number;topic:string;homework?:string|null;materials:LessonMaterial[];gradeItems:GradeItem[]}
export interface LessonDraft{lessonDate:string;position:number;topic:string;homework?:string;materials:{title:string;url:string;position:number}[]}
export async function listLessons(id:string,from:string,to:string):Promise<LessonRecord[]>{const query=new URLSearchParams();if(from)query.set('dateFrom',from);if(to)query.set('dateTo',to);return(await request<{items:LessonRecord[]}>(`/api/v1/class-subjects/${id}/lessons?${query}`)).items}
export async function createLesson(id:string,value:LessonDraft):Promise<LessonRecord>{return(await request<{lesson:LessonRecord}>(`/api/v1/class-subjects/${id}/lessons`,{method:'POST',body:JSON.stringify(value)})).lesson}
export async function createGradeItem(lessonID:string,value:Omit<GradeItem,'id'|'lessonId'>):Promise<GradeItem>{return(await request<{gradeItem:GradeItem}>(`/api/v1/lessons/${lessonID}/grade-items`,{method:'POST',body:JSON.stringify(value)})).gradeItem}
