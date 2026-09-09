import{request}from'./api'
import{AcademicYear}from'./academicYears'
import{ClassRecord}from'./classes'
import{ClassSubjectRecord}from'./classSubjects'
import{LessonRecord}from'./lessons'
import{QuarterGrade}from'./quarterGrades'

export interface ParentStudent{id:string;firstName:string;lastName:string;roles:string[]}
export interface ParentStudentPeriod{academicYear:AcademicYear;class:ClassRecord;subjects:ClassSubjectRecord[]}
export interface ParentJournalSubject{classSubject:ClassSubjectRecord;lessons:Array<{lesson:LessonRecord}>;quarterGrades?:QuarterGrade[]}
export interface ParentJournal{student:ParentStudent;academicYear:AcademicYear;class:ClassRecord;subjects:ParentJournalSubject[]}
export async function listCurrentParentStudents():Promise<ParentStudent[]>{return(await request<{items:ParentStudent[]}>('/api/v1/parent/students')).items}
export async function listParentStudentPeriods(studentID:string):Promise<ParentStudentPeriod[]>{return(await request<{items:ParentStudentPeriod[]}>(`/api/v1/parent/students/${studentID}/periods`)).items}
export async function getParentStudentJournal(studentID:string,yearID:string):Promise<ParentJournal>{return request<ParentJournal>(`/api/v1/parent/students/${studentID}/journal?academicYearId=${yearID}`)}
