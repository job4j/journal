import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})

describe('login', () => {
  it('opens the empty workspace after successful login', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'admin@example.ru', firstName: 'Admin', lastName: 'User', status: 'active', roles: ['admin'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    render(<App />)
    await userEvent.type(screen.getByLabelText('Электронная почта'), 'admin@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByLabelText('Рабочая область')).toBeInTheDocument()
    expect(screen.getByRole('navigation', { name: 'Основная навигация' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Пользователи' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Роли' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Классы' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Ученики' })).toBeInTheDocument()
  })

  it('shows the server error and keeps the form', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({ message: 'Неверный email или пароль' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    }))

    render(<App />)
    await userEvent.type(screen.getByLabelText('Электронная почта'), 'admin@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'wrong-password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Неверный email или пароль')
    expect(screen.getByRole('button', { name: 'Войти' })).toBeEnabled()
  })

  it('shows only classes to a teacher', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'teacher@example.ru', firstName: 'Анна', lastName: 'Иванова', status: 'active', roles: ['teacher'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    render(<App />)
    await userEvent.type(screen.getByLabelText('Электронная почта'), 'teacher@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByRole('button', { name: 'Классы' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Пользователи' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Ученики' })).not.toBeInTheDocument()
  })

  it('shows only students to a parent', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'parent@example.ru', firstName: 'Пётр', lastName: 'Иванов', status: 'active', roles: ['parent'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    render(<App />)
    await userEvent.type(screen.getByLabelText('Электронная почта'), 'parent@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByRole('button', { name: 'Ученики' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Классы' })).not.toBeInTheDocument()
  })

  it('opens a class journal and adds a mock grade', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'teacher@example.ru', firstName: 'Анна', lastName: 'Иванова', status: 'active', roles: ['teacher'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    render(<App />)
    await userEvent.type(screen.getByLabelText('Электронная почта'), 'teacher@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))
    await userEvent.click(await screen.findByRole('button', { name: /7А/ }))

    expect(screen.getByRole('heading', { name: 'Предметы класса' })).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: /Математика/ }))

    expect(screen.getByRole('heading', { name: 'Журнал класса' })).toBeInTheDocument()
    expect(screen.getByRole('rowheader', { name: 'Анна Белова' })).toBeInTheDocument()
    expect(screen.getByRole('columnheader', { name: /7 сентября/ })).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: 'Пропуск: Иван Громов, 8 сентября' })).toHaveValue('Н')
    const gradeSelect = screen.getByRole('combobox', { name: 'Домашняя работа: Павел Орлов, 7 сентября' })
    expect(gradeSelect).toHaveValue('')
    await userEvent.selectOptions(gradeSelect, '5')
    expect(gradeSelect).toHaveValue('5')
  })
})
