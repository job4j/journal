import { request } from './api'

export interface ScheduleEntry {
  id: string
  day: string
  position: number
  classSubjectId: string | null
  title: string
  notes: string
  startsAt: string | null
  endsAt: string | null
}
export interface ScheduleWeek { weekStart: string; revision: number; entries: ScheduleEntry[] }
export async function getClassSchedule(classID: string, weekStart: string): Promise<ScheduleWeek> {
  return (await request<{week: ScheduleWeek}>(`/api/v1/classes/${classID}/schedule/${weekStart}`)).week
}
export async function replaceClassSchedule(classID: string, week: ScheduleWeek): Promise<ScheduleWeek> {
  return (await request<{week: ScheduleWeek}>(`/api/v1/classes/${classID}/schedule/${week.weekStart}`, { method: 'PUT', body: JSON.stringify({ revision: week.revision, entries: week.entries }) })).week
}
export function addDays(day: string, count: number) {
  const date = new Date(`${day}T00:00:00Z`)
  date.setUTCDate(date.getUTCDate() + count)
  return date.toISOString().slice(0, 10)
}
export function monday(day: string) {
  return addDays(day, -(new Date(`${day}T00:00:00Z`).getUTCDay() + 6) % 7)
}
export function localToday() {
  const date = new Date()
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2,'0')}-${String(date.getDate()).padStart(2,'0')}`
}
