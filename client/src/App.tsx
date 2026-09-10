import { FormEvent, useEffect, useState } from 'react'
import { currentUser, login, logout, LoginError, User } from './auth'
import TeacherClassesView from './TeacherClassesView'
import RolesView from './RolesView'
import UsersView from './UsersView'
import AcademicYearsView from './AcademicYearsView'
import SubjectsView from './SubjectsView'
import AdminClassesView from './AdminClassesView'
import ParentJournalView from './ParentJournalView'
import './styles.css'

type Role = User['roles'][number]
type NavigationItem = {
  id: string
  label: string
  icon: 'users' | 'shield' | 'calendar' | 'classes' | 'student' | 'teacher' | 'book'
  roles: Role[]
}

const navigation: NavigationItem[] = [
  { id: 'users', label: 'Пользователи', icon: 'users', roles: ['admin'] },
  { id: 'roles', label: 'Роли', icon: 'shield', roles: ['admin'] },
  { id: 'academic-years', label: 'Учебные годы', icon: 'calendar', roles: ['admin'] },
  { id: 'classes', label: 'Классы', icon: 'classes', roles: ['admin', 'teacher'] },
  { id: 'students', label: 'Мои дети', icon: 'student', roles: ['parent'] },
  { id: 'subjects', label: 'Предметы', icon: 'book', roles: ['admin'] },
  { id: 'journal', label: 'Мой журнал', icon: 'book', roles: ['student'] },
]

const roleLabels: Record<string, string> = {
  admin: 'Администратор',
  teacher: 'Учитель',
  parent: 'Родитель',
  student: 'Ученик',
}

function NavigationIcon({ name }: { name: NavigationItem['icon'] }) {
  const paths = {
    users: <><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/></>,
    shield: <><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z"/><path d="m9 12 2 2 4-4"/></>,
    calendar: <><rect x="3" y="5" width="18" height="16" rx="2"/><path d="M16 3v4M8 3v4M3 11h18"/></>,
    classes: <><path d="m3 7 9-4 9 4-9 4-9-4Z"/><path d="m7 9.5v5L12 17l5-2.5v-5M21 7v6"/></>,
    student: <><circle cx="12" cy="8" r="4"/><path d="M4 21a8 8 0 0 1 16 0"/></>,
    teacher: <><circle cx="8" cy="8" r="3"/><path d="M2 20a6 6 0 0 1 12 0M15 5h7v11h-5M18.5 8v5"/></>,
    book: <><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2Z"/></>,
  }
  return <svg viewBox="0 0 24 24" aria-hidden="true">{paths[name]}</svg>
}

function Workspace({ user, onLogout }: { user: User; onLogout: () => void }) {
  const items = navigation.filter((item) => item.roles.some((role) => user.roles.includes(role)))
  const [activeID, setActiveID] = useState(()=>location.hash.slice(1).split('/')[0]||items[0]?.id||'')
  const knownItem=navigation.find(item=>item.id===activeID),activeItem=items.find(item=>item.id===activeID)
  const initials = `${user.name.charAt(0)}${user.name.charAt(0)}`.toUpperCase()
  const [breadcrumbs, setBreadcrumbs] = useState<{ label: string; action?: () => void }[]>([])
  useEffect(()=>{const sync=()=>setActiveID(location.hash.slice(1).split('/')[0]||items[0]?.id||'');window.addEventListener('hashchange',sync);return()=>window.removeEventListener('hashchange',sync)},[items])
  function navigate(id:string){setBreadcrumbs([]);location.hash=id;setActiveID(id)}

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="sidebar-brand">
          <div className="brand-mark brand-mark--small" aria-hidden="true">Р</div>
          <span>Родное слово</span>
        </div>

        <nav className="navigation" aria-label="Основная навигация">
          <p className="navigation-label">Управление</p>
          {items.map((item) => (
            <button
              className={`navigation-item${item.id === activeItem?.id ? ' navigation-item--active' : ''}`}
              key={item.id}
              type="button"
              aria-current={item.id === activeItem?.id ? 'page' : undefined}
              onClick={() => navigate(item.id)}
            >
              <NavigationIcon name={item.icon} />
              <span>{item.label}</span>
            </button>
          ))}
        </nav>

        <div className="sidebar-user">
          <div className="avatar" aria-hidden="true">{initials}</div>
          <div className="user-copy">
            <strong>{user.name}</strong>
            <span>{user.roles.map((role) => roleLabels[role] ?? role).join(' · ')}</span>
          </div>
          <button type="button" aria-label="Выйти" onClick={onLogout}>↪</button>
        </div>
      </aside>

      <main className="workspace" aria-label="Рабочая область">
        <header className="workspace-header">
          <div>

            <nav className="breadcrumbs" aria-label="Хлебные крошки">
              {(breadcrumbs.length
                ? breadcrumbs
                : [{ label: activeItem?.label ?? 'Рабочая область' }]
              ).map((item, index, values) => (
                <span className="breadcrumb" key={`${item.label}-${index}`}>
                  {item.action && index < values.length - 1 ? (
                    <button type="button" onClick={item.action}>{item.label}</button>
                  ) : (
                    <strong aria-current="page">{item.label}</strong>
                  )}
                </span>
              ))}
            </nav>
          </div>
          <div className="header-avatar" aria-label={`${user.name}`}>{initials}</div>
        </header>
        <section className="workspace-content" aria-label={activeItem?.label}>
          {!knownItem&&<RouteState code="404" title="Страница не найдена" action={()=>navigate(items[0]?.id??'')}/>} 
          {knownItem&&!activeItem&&<RouteState code="403" title="Этот раздел недоступен вашей роли" action={()=>navigate(items[0]?.id??'')}/>} 
          {activeItem?.id === 'classes' && (user.roles.includes('admin') ? (
            <AdminClassesView setBreadcrumbs={setBreadcrumbs} />
          ) : (
            <TeacherClassesView setBreadcrumbs={setBreadcrumbs} />
          ))}
          {activeItem?.id === 'roles' && <RolesView />}
          {activeItem?.id === 'users' && <UsersView />}
          {activeItem?.id === 'academic-years' && <AcademicYearsView />}
          {activeItem?.id === 'subjects' && <SubjectsView />}
          {activeItem?.id === 'students' && user.roles.includes('parent') && <ParentJournalView />}
          {activeItem?.id === 'journal' && <div className="empty-card">Для ученика пока нет доступных действий.</div>}
        </section>
      </main>
    </div>
  )
}

function RouteState({code,title,action}:{code:string;title:string;action:()=>void}){return <div className="route-state"><strong>{code}</strong><h2>{title}</h2><button className="primary-action" onClick={action}>В рабочую область</button></div>}

export default function App() {
  const [user, setUser] = useState<User | null>(null)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isRestoring, setIsRestoring] = useState(true)

  useEffect(() => {
    currentUser().then(setUser).catch(() => setError('Не удалось восстановить сессию')).finally(() => setIsRestoring(false))
  }, [])
  useEffect(()=>{const expired=()=>{setUser(null);setError('Сессия завершена. Войдите снова.')};window.addEventListener('journal:session-expired',expired);return()=>window.removeEventListener('journal:session-expired',expired)},[])

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    setError('')
    setIsSubmitting(true)

    try {
      setUser(await login(String(form.get('login')), String(form.get('password'))))
    } catch (cause) {
      setError(cause instanceof LoginError ? cause.message : 'Не удалось войти')
    } finally {
      setIsSubmitting(false)
    }
  }

  if (isRestoring) return <main className="login-page" aria-label="Загрузка">Загрузка…</main>
  if (user) return <Workspace user={user} onLogout={() => { void logout().then(() => setUser(null)).catch(() => setError('Не удалось выйти')) }} />

  return (
    <main className="login-page">
      <section className="login-card" aria-labelledby="login-title">
        <div className="brand-mark" aria-hidden="true">J</div>
        <div className="heading">
          <p className="eyebrow">Электронный журнал</p>
          <h1 id="login-title">Вход в Journal</h1>
        </div>

        <form className="login-form" onSubmit={handleSubmit}>
          <label htmlFor="login">Логин</label>
          <input id="login" name="login" type="text" autoComplete="username" placeholder="Введите логин" required autoFocus />

          <label htmlFor="password">Пароль</label>
          <input id="password" name="password" type="password" autoComplete="current-password" placeholder="Введите пароль" required />

          {error && <p className="form-error" role="alert">{error}</p>}

          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Входим…' : 'Войти'}
          </button>
        </form>
      </section>
    </main>
  )
}
