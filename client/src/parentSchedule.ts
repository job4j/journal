import { request } from './api'
import type { ScheduleEntry } from './classSchedule'

export interface ParentScheduleEntry extends ScheduleEntry { homework: string | null; homeworkFrom: string | null }
export interface ParentScheduleWeek { weekStart: string; entries: ParentScheduleEntry[] }
export async function getParentSchedule(studentID: string, classID: string, weekStart: string): Promise<ParentScheduleWeek> {
  return (await request<{week: ParentScheduleWeek}>(`/api/v1/parent/students/${studentID}/classes/${classID}/schedule/${weekStart}`)).week
}
