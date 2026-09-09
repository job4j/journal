import{request}from'./api'
import{AcademicYear}from'./academicYears'
import{ClassRecord}from'./classes'
import{ClassSubjectRecord}from'./classSubjects'

export interface ParentStudent{id:string;firstName:string;lastName:string;roles:string[]}
export interface ParentStudentPeriod{academicYear:AcademicYear;class:ClassRecord;subjects:ClassSubjectRecord[]}
export async function listCurrentParentStudents():Promise<ParentStudent[]>{return(await request<{items:ParentStudent[]}>('/api/v1/parent/students')).items}
export async function listParentStudentPeriods(studentID:string):Promise<ParentStudentPeriod[]>{return(await request<{items:ParentStudentPeriod[]}>(`/api/v1/parent/students/${studentID}/periods`)).items}
