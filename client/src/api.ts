export interface ErrorResponse { code?: string; message?: string; fields?: Record<string, string> }

export class ApiError extends Error {
  constructor(message: string, readonly status = 0, readonly code?: string, readonly fields?: Record<string, string>) { super(message) }
}

export async function request<T>(url: string, init?: RequestInit): Promise<T> {
  let response: Response
  try {
    response = await fetch(url, { credentials: 'include', ...init, headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) } })
  } catch {
    throw new ApiError('Не удалось связаться с сервером')
  }
  if (!response.ok) {
    const payload = await response.json().catch(() => null) as ErrorResponse | null
    if (response.status === 401) window.dispatchEvent(new CustomEvent('journal:session-expired'))
    throw new ApiError(payload?.message ?? 'Не удалось выполнить запрос', response.status, payload?.code, payload?.fields)
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}
