import { request } from './api'

export interface ParentMissedTopic {
  lessonId:string
  lessonDate:string
  position:number
  topic:string
  subjectId:string
  subjectName:string
  classId:string
  className:string
}
export async function listParentMissedTopics(studentID:string,yearID:string):Promise<ParentMissedTopic[]> {
  const query=new URLSearchParams({academicYearId:yearID})
  return (await request<{items:ParentMissedTopic[]}>(`/api/v1/parent/students/${encodeURIComponent(studentID)}/missed-topics?${query}`)).items
}
