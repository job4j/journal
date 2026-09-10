import { ApiError, request } from './api'

export interface UserRecord{id:string;login?:string;email?:string|null;phone?:string|null;name:string;status:'active'|'blocked';roles:string[]}
export interface UserDraft{login:string;email:string;phone:string;password:string;name:string;status:'active'|'blocked';roles:string[]}
export { ApiError as UsersApiError }
export async function listUsers():Promise<UserRecord[]>{return(await request<{items:UserRecord[]}>('/api/v1/users?limit=100&offset=0')).items}
export async function createUser(draft:UserDraft):Promise<UserRecord>{const body={login:draft.login,email:draft.email||null,phone:draft.phone||null,password:draft.password,name:draft.name,roles:draft.roles};return(await request<{user:UserRecord}>('/api/v1/users',{method:'POST',body:JSON.stringify(body)})).user}
export async function updateUser(id:string,draft:UserDraft):Promise<UserRecord>{const body={...draft,password:draft.password||undefined};return(await request<{user:UserRecord}>(`/api/v1/users/${id}`,{method:'PUT',body:JSON.stringify(body)})).user}
export async function deleteUser(id:string):Promise<void>{await request<void>(`/api/v1/users/${id}`,{method:'DELETE'})}
export async function listParentStudents(parentID:string):Promise<UserRecord[]>{return(await request<{items:UserRecord[]}>(`/api/v1/parents/${parentID}/students`)).items as UserRecord[]}
export async function grantParentStudent(parentID:string,studentID:string):Promise<void>{await request(`/api/v1/parents/${parentID}/students/${studentID}`,{method:'PUT',body:JSON.stringify({canViewProfile:true,canViewJournal:true})})}
export async function revokeParentStudent(parentID:string,studentID:string):Promise<void>{await request(`/api/v1/parents/${parentID}/students/${studentID}`,{method:'DELETE'})}

