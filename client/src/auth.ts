export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  status: 'active' | 'blocked'
  roles: string[]
}

interface LoginResponse {
  user: User
}

export { ApiError as LoginError }


export async function login(email: string, password: string): Promise<User> {
  return (await request<LoginResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })).user
}

export async function currentUser(): Promise<User | null> {
  try { return (await request<LoginResponse>('/api/v1/me')).user }
  catch (error) { if (error instanceof ApiError && error.status === 401) return null; throw error }
}

export async function logout(): Promise<void> {
  await request('/api/v1/auth/logout', { method: 'POST' })
}
import { ApiError, request } from './api'
