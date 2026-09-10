import{request}from'./api'
export interface ClassRecord{id:string;academicYearId:string;name:string;gradeLevel:number;studentCount:number;quarters?:{id:string;number:number;name?:string;startsOn:string;endsOn:string}[]}
export interface ClassDraft{academicYearId:string;name:string;gradeLevel:number}
export interface ClassStudent{student:{id:string;name:string;roles:string[]};enrolledOn:string;leftOn?:string|null}
export async function listClasses(yearID:string):Promise<ClassRecord[]>{return(await request<{items:ClassRecord[]}>(`/api/v1/classes?academicYearId=${encodeURIComponent(yearID)}`)).items}
export async function getClass(id:string):Promise<ClassRecord>{return(await request<{class:ClassRecord}>(`/api/v1/classes/${id}`)).class}
export async function createClass(value:ClassDraft):Promise<ClassRecord>{return(await request<{class:ClassRecord}>('/api/v1/classes',{method:'POST',body:JSON.stringify(value)})).class}
export async function listClassStudents(classID:string):Promise<ClassStudent[]>{return(await request<{items:ClassStudent[]}>(`/api/v1/classes/${classID}/students`)).items}
export async function addStudentToClass(classID:string,studentId:string,enrolledOn:string):Promise<ClassStudent>{return(await request<{classStudent:ClassStudent}>(`/api/v1/classes/${classID}/students`,{method:'POST',body:JSON.stringify({studentId,enrolledOn})})).classStudent}
export async function updateClassStudent(classID:string,studentID:string,leftOn:string|null):Promise<ClassStudent>{return(await request<{classStudent:ClassStudent}>(`/api/v1/classes/${classID}/students/${studentID}`,{method:'PATCH',body:JSON.stringify({leftOn})})).classStudent}
