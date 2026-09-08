import{request}from'./api'
export interface ClassRecord{id:string;academicYearId:string;name:string;gradeLevel:number;studentCount:number}
export interface ClassDraft{academicYearId:string;name:string;gradeLevel:number}
export async function listClasses(yearID:string):Promise<ClassRecord[]>{return(await request<{items:ClassRecord[]}>(`/api/v1/classes?academicYearId=${encodeURIComponent(yearID)}`)).items}
export async function createClass(value:ClassDraft):Promise<ClassRecord>{return(await request<{class:ClassRecord}>('/api/v1/classes',{method:'POST',body:JSON.stringify(value)})).class}
