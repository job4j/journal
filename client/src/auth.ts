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

interface ErrorResponse {
  code?: string
  message?: string
}

export class LoginError extends Error {}

export async function login(email: string, password: string): Promise<User> {
  let response: Response
  try {
    response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })
  } catch {
    throw new LoginError('Не удалось связаться с сервером')
  }

  if (!response.ok) {
    const error = await response.json().catch(() => null) as ErrorResponse | null
    throw new LoginError(error?.message ?? 'Не удалось войти')
  }

  const payload = await response.json() as LoginResponse
  return payload.user
}
