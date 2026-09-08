export interface Role { id:string; code:string; name:string; permissions:string[] }
export interface Permission { code:string; description:string }
interface ErrorResponse { message?:string }
export class RolesApiError extends Error {}
async function request<T>(url:string,init?:RequestInit):Promise<T>{
 let response:Response
 try{response=await fetch(url,{credentials:'include',...init,headers:{'Content-Type':'application/json',...(init?.headers??{})}})}catch{throw new RolesApiError('Не удалось связаться с сервером')}
 if(!response.ok){const payload=await response.json().catch(()=>null) as ErrorResponse|null;throw new RolesApiError(payload?.message??'Не удалось выполнить запрос')}
 if(response.status===204)return undefined as T
 return response.json() as Promise<T>
}
export async function listRoles():Promise<Role[]>{return (await request<{roles:Role[]}>('/api/v1/roles')).roles}
export async function listPermissions():Promise<Permission[]>{return (await request<{permissions:Permission[]}>('/api/v1/permissions')).permissions}
export async function createRole(value:Omit<Role,'id'>):Promise<Role>{return (await request<{role:Role}>('/api/v1/roles',{method:'POST',body:JSON.stringify(value)})).role}
export async function updateRole(value:Role):Promise<Role>{return (await request<{role:Role}>(`/api/v1/roles/${value.id}`,{method:'PUT',body:JSON.stringify({code:value.code,name:value.name,permissions:value.permissions})})).role}
export async function deleteRole(id:string):Promise<void>{await request<void>(`/api/v1/roles/${id}`,{method:'DELETE'})}
