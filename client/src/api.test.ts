import { afterEach, describe, expect, it, vi } from 'vitest'
import { request } from './api'

afterEach(() => vi.restoreAllMocks())

describe('request', () => {
  it('returns JSON and includes credentials', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ ok: true }), { status: 200 }))
    await expect(request<{ ok: boolean }>('/test')).resolves.toEqual({ ok: true })
    expect(fetchMock).toHaveBeenCalledWith('/test', expect.objectContaining({ credentials: 'include' }))
  })

  it('preserves the API error details', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ code: 'invalid', message: 'Ошибка', fields: { name: 'Обязательно' } }), { status: 400 }))
    await expect(request('/test')).rejects.toMatchObject({ status: 400, code: 'invalid', message: 'Ошибка' })
  })

  it('maps a network failure', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('offline'))
    await expect(request('/test')).rejects.toThrow('Не удалось связаться с сервером')
  })
})
