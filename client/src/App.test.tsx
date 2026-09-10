import { afterEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
  location.hash=''
})

describe('role routes',()=>{it('shows 403 for a known inaccessible route',async()=>{location.hash='roles';vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({user:{id:'s-1',email:'s@test',name:'А'+' '+'Б',status:'active',roles:['student']}}),{status:200}));render(<App/>);expect(await screen.findByText('403')).toBeInTheDocument()});it('shows 404 for an unknown route',async()=>{location.hash='missing';vi.spyOn(globalThis,'fetch').mockResolvedValueOnce(new Response(JSON.stringify({user:{id:'a-1',email:'a@test',name:'А'+' '+'Б',status:'active',roles:['admin']}}),{status:200}));render(<App/>);expect(await screen.findByText('404')).toBeInTheDocument()})})

describe('login', () => {
  it('opens the empty workspace after successful login', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(new Response(null, { status: 401 })).mockResolvedValue(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'admin@example.ru', name: 'Admin'+' '+ 'User', status: 'active', roles: ['admin'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    render(<App />)
    await userEvent.type(await screen.findByLabelText('Логин'), 'admin@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByLabelText('Рабочая область')).toBeInTheDocument()
    expect(screen.getByRole('navigation', { name: 'Основная навигация' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Пользователи' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Роли' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Классы' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Предметы' })).toBeInTheDocument()
  })

  it('shows the server error and keeps the form', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(new Response(null, { status: 401 })).mockResolvedValue(new Response(JSON.stringify({ message: 'Неверный email или пароль' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' },
    }))

    render(<App />)
    await userEvent.type(await screen.findByLabelText('Логин'), 'admin@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'wrong-password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Неверный email или пароль')
    expect(screen.getByRole('button', { name: 'Войти' })).toBeEnabled()
  })

  it('shows only classes to a teacher', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(new Response(null, { status: 401 })).mockResolvedValueOnce(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'teacher@example.ru', name: 'Анна'+' '+ 'Иванова', status: 'active', roles: ['teacher'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } })).mockResolvedValue(new Response(JSON.stringify({ items: [] }), { status: 200 }))

    render(<App />)
    await userEvent.type(await screen.findByLabelText('Логин'), 'teacher@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByRole('button', { name: 'Классы' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Пользователи' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Ученики' })).not.toBeInTheDocument()
  })

  it('shows only students to a parent', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(new Response(null, { status: 401 })).mockResolvedValue(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'parent@example.ru', name: 'Пётр'+' '+ 'Иванов', status: 'active', roles: ['parent'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    render(<App />)
    await userEvent.type(await screen.findByLabelText('Логин'), 'parent@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))

    expect(await screen.findByRole('button', { name: 'Мои дети' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Классы' })).not.toBeInTheDocument()
  })

  it('opens the real teacher class subjects', async () => {
    const klass = { id: 'class-1', academicYearId: 'year-1', name: '7А', gradeLevel: 7, studentCount: 12 }
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce(new Response(null, { status: 401 })).mockResolvedValueOnce(new Response(JSON.stringify({
      user: { id: crypto.randomUUID(), email: 'teacher@example.ru', name: 'Анна'+' '+ 'Иванова', status: 'active', roles: ['teacher'] },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } })).mockResolvedValueOnce(new Response(JSON.stringify({ items: [klass] }), { status: 200 })).mockResolvedValueOnce(new Response(JSON.stringify({ items: [{ id: 'link-1', classId: 'class-1', subject: { id: 'subject-1', code: 'MATH', name: 'Математика' }, responsibleTeacher: { id: 'teacher-1', name: 'Анна'+' '+ 'Иванова', roles: ['teacher'] } }] }), { status: 200 }))

    render(<App />)
    await userEvent.type(await screen.findByLabelText('Логин'), 'teacher@example.ru')
    await userEvent.type(screen.getByLabelText('Пароль'), 'password')
    await userEvent.click(screen.getByRole('button', { name: 'Войти' }))
    await userEvent.click(await screen.findByRole('button', { name: /7А/ }))

    expect(await screen.findByLabelText('Хлебные крошки')).toHaveTextContent('Классы7А')
    expect(screen.getByText('Математика')).toBeInTheDocument()
  })
})
