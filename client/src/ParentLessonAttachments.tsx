import LessonAttachmentsList from './LessonAttachmentsList'

export default function ParentLessonAttachments({lessonID,studentID}:{lessonID:string;studentID:string}) {
  return <LessonAttachmentsList lessonID={lessonID} studentID={studentID}/>
}
