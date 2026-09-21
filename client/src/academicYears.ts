import { request } from './api'
export interface AcademicYearQuarter{id:string;number:number;name:string;startsOn:string;endsOn:string}
export interface AcademicYear{id:string;name:string;startsOn:string;endsOn:string;status:'planned'|'active'|'completed';quarters:AcademicYearQuarter[]}
export type AcademicYearDraft=Omit<AcademicYear,'id'|'quarters'>
export async function listAcademicYears():Promise<AcademicYear[]>{return(await request<{items:AcademicYear[]}>('/api/v1/academic-years')).items}
export async function createAcademicYear(value:AcademicYearDraft):Promise<AcademicYear>{return(await request<{academicYear:AcademicYear}>('/api/v1/academic-years',{method:'POST',body:JSON.stringify(value)})).academicYear}
export async function createAcademicYearQuarter(yearID:string,value:{name:string;startsOn:string;endsOn:string}):Promise<AcademicYearQuarter>{return(await request<{quarter:AcademicYearQuarter}>(`/api/v1/academic-years/${yearID}/quarters`,{method:'POST',body:JSON.stringify(value)})).quarter}

export async function updateAcademicYearQuarter(id:string,value:{name:string;startsOn:string;endsOn:string}):Promise<AcademicYearQuarter>{return(await request<{quarter:AcademicYearQuarter}>(`/api/v1/academic-year-quarters/${id}`,{method:'PUT',body:JSON.stringify(value)})).quarter}
export async function deleteAcademicYearQuarter(id:string):Promise<void>{await request(`/api/v1/academic-year-quarters/${id}`,{method:'DELETE'})}