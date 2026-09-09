import{request}from'./api'
export type GradingScale='five_point'|'points'|'pass_fail'
export interface QuarterGrade{id:string;quarterId:string;classSubjectId:string;studentId:string;gradingScale:GradingScale;maxScore?:number|null;numericValue?:number|null;textValue?:string|null;teacherComment?:string|null}
export interface QuarterGradeDraft{gradingScale:GradingScale;maxScore?:number;numericValue?:number;textValue?:string;teacherComment?:string}
export async function listQuarterGrades(assignmentID:string,quarterID:string):Promise<QuarterGrade[]>{return(await request<{items:QuarterGrade[]}>(`/api/v1/class-subjects/${assignmentID}/quarters/${quarterID}/grades`)).items}
export async function putQuarterGrade(assignmentID:string,quarterID:string,studentID:string,value:QuarterGradeDraft):Promise<QuarterGrade>{return(await request<{quarterGrade:QuarterGrade}>(`/api/v1/class-subjects/${assignmentID}/quarters/${quarterID}/grades/${studentID}`,{method:'PUT',body:JSON.stringify(value)})).quarterGrade}
