export interface UserRecord{id:string;email:string;firstName:string;lastName:string;status:'active'|'blocked';roles:string[]}
export interface UserDraft{email:string;password:string;firstName:string;lastName:string;status:'active'|'blocked';roles:string[]}
interface ErrorResponse{message?:string}
export class UsersApiError extends Error{}
async function request<T>(url:string,init?:RequestInit):Promise<T>{let response:Response;try{response=await fetch(url,{credentials:'include',...init,headers:{'Content-Type':'application/json',...(init?.headers??{})}})}catch{throw new UsersApiError('Не удалось связаться с сервером')}if(!response.ok){const payload=await response.json().catch(()=>null) as ErrorResponse|null;throw new UsersApiError(payload?.message??'Не удалось выполнить запрос')}if(response.status===204)return undefined as T;return response.json() as Promise<T>}
export async function listUsers():Promise<UserRecord[]>{return(await request<{items:UserRecord[]}>('/api/v1/users?limit=100&offset=0')).items}
export async function createUser(draft:UserDraft):Promise<UserRecord>{const body={email:draft.email,password:draft.password,firstName:draft.firstName,lastName:draft.lastName,roles:draft.roles};return(await request<{user:UserRecord}>('/api/v1/users',{method:'POST',body:JSON.stringify(body)})).user}
export async function updateUser(id:string,draft:UserDraft):Promise<UserRecord>{const body={...draft,password:draft.password||undefined};return(await request<{user:UserRecord}>(`/api/v1/users/${id}`,{method:'PUT',body:JSON.stringify(body)})).user}
export async function deleteUser(id:string):Promise<void>{await request<void>(`/api/v1/users/${id}`,{method:'DELETE'})}

