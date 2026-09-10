import{FormEvent,useCallback,useEffect,useMemo,useState}from'react'
import{ClassStudent,listClassStudents}from'./classes'
import{ClassSubjectRecord}from'./classSubjects'
import{Absence,GradeItem,GradeItemDraft,LessonDraft,LessonRecord,Score,createGradeItem,createLesson,updateLesson,deleteStudentAbsence,listLessons,putStudentAbsence,putStudentScore}from'./lessons'
import{GradingScale,QuarterGrade,listQuarterGrades,putQuarterGrade}from'./quarterGrades'

const emptyLesson=():LessonDraft=>({lessonDate:new Date().toISOString().slice(0,10),position:1,topic:'',homework:'',materials:[]})
export default function TeacherLessonsView({assignment,quarters=[],className,onBack}:{assignment:ClassSubjectRecord;quarters?:{id:string;number:number;startsOn?:string;endsOn?:string}[];className?:string;onBack:()=>void}){
 const[items,setItems]=useState<LessonRecord[]>([]),[students,setStudents]=useState<ClassStudent[]>([]),[quarterID,setQuarterID]=useState(quarters[0]?.id??''),[draft,setDraft]=useState<LessonDraft|null>(null),[editingLessonID,setEditingLessonID]=useState(''),[gradeLessonID,setGradeLessonID]=useState(''),[loading,setLoading]=useState(true),[saving,setSaving]=useState(false),[error,setError]=useState('')
 const selectedQuarter=quarters.find(item=>item.id===quarterID),from=selectedQuarter?.startsOn??'',to=selectedQuarter?.endsOn??''
 const load=useCallback(async()=>{setLoading(true);setError('');try{const[lessons,roster]=await Promise.all([listLessons(assignment.id,from,to),listClassStudents(assignment.classId)]);setItems(lessons);setStudents(roster)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось загрузить журнал')}finally{setLoading(false)}},[assignment.id,assignment.classId,from,to])
 useEffect(()=>{void load()},[load])
 async function submit(event:FormEvent){event.preventDefault();if(!draft)return;setSaving(true);try{const payload={...draft,homework:draft.homework||undefined};const saved=editingLessonID?await updateLesson(editingLessonID,payload):await createLesson(assignment.id,payload);setItems(current=>editingLessonID?current.map(item=>item.id===editingLessonID?{...saved,gradeItems:item.gradeItems,absences:item.absences}:item):[...current,saved]);setEditingLessonID('');setDraft(null)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось создать урок')}finally{setSaving(false)}}
 function scoreSaved(saved:Score){setItems(current=>current.map(lesson=>({...lesson,gradeItems:lesson.gradeItems.map(item=>item.id!==saved.gradeItemId?item:{...item,scores:[...item.scores.filter(score=>score.studentId!==saved.studentId),saved]})})))}
 function absenceChanged(lessonID:string,studentID:string,absence?:Absence){setItems(current=>current.map(lesson=>lesson.id!==lessonID?lesson:{...lesson,absences:absence?[...(lesson.absences??[]).filter(item=>item.studentId!==studentID),absence]:(lesson.absences??[]).filter(item=>item.studentId!==studentID)}))}
 if(gradeLessonID)return <GradeEditor lessonID={gradeLessonID} onBack={()=>setGradeLessonID('')} onSaved={saved=>setItems(current=>current.map(item=>item.id===gradeLessonID?{...item,gradeItems:[...item.gradeItems,saved]}:item))}/>
 return <div className="content-stack"><button className="back-button" onClick={onBack}>← Предметы класса</button><div className="content-heading"><div><p className="content-kicker">{className?`${className} · ${assignment.subject.name}`:assignment.subject.name}</p><h2>Журнал класса</h2></div></div><div className="lesson-filters">{quarters.length>0?<label>Период<select aria-label="Период" value={quarterID} onChange={e=>setQuarterID(e.target.value)}>{quarters.map(item=><option value={item.id} key={item.id}>{item.number} период</option>)}</select></label>:<span className="summary-empty">Нет периодов</span>}</div>{error&&<p className="content-error" role="alert">{error}</p>}{loading?<div className="empty-card">Загружаем журнал…</div>:<><div className="grade-legend" aria-label="Обозначения оценок">{gradeKinds.map(kind=><span key={kind}><GradeIcon kind={kind}/>{gradeLabels[kind]}</span>)}</div><GradeTable items={items} students={students} from={from} to={to} onHomework={(date,lesson)=>{setEditingLessonID(lesson?.id??'');setDraft(lesson?{lessonDate:lesson.lessonDate,position:lesson.position,topic:lesson.topic,homework:lesson.homework??'',materials:lesson.materials.map(m=>({title:m.title,url:m.url,position:m.position}))}:{...emptyLesson(),lessonDate:date})}} onSaved={scoreSaved} onAbsence={absenceChanged} onAddGrade={setGradeLessonID}/>{quarterID&&<QuarterGradeTable assignmentID={assignment.id} quarterID={quarterID} students={students}/>}</>}{draft&&<LessonDialog draft={draft} setDraft={setDraft} saving={saving} submit={submit}/>}</div>
}

function QuarterGradeTable({assignmentID,quarterID,students}:{assignmentID:string;quarterID:string;students:ClassStudent[]}){const[items,setItems]=useState<QuarterGrade[]>([]),[error,setError]=useState('');useEffect(()=>{setError('');listQuarterGrades(assignmentID,quarterID).then(setItems).catch(cause=>setError(cause instanceof Error?cause.message:'Не удалось загрузить итоги'))},[assignmentID,quarterID]);return <div className="content-stack"><h3>Итоги четверти</h3>{error&&<p className="content-error" role="alert">{error}</p>}<div className="table-card"><table className="classes-table"><thead><tr><th>Ученик</th><th>Итог</th></tr></thead><tbody>{students.map(member=><tr key={member.student.id}><td>{member.student.name}</td><td><QuarterGradeEditor current={(items??[]).find(item=>item.studentId===member.student.id)} onSave={async value=>{const saved=await putQuarterGrade(assignmentID,quarterID,member.student.id,value);setItems(current=>[...current.filter(item=>item.studentId!==saved.studentId),saved])}}/></td></tr>)}</tbody></table></div></div>}
function QuarterGradeEditor({current,onSave}:{current?:QuarterGrade;onSave:(value:{gradingScale:GradingScale;maxScore?:number;numericValue?:number;textValue?:string})=>Promise<void>}){const[scale,setScale]=useState<GradingScale>(current?.gradingScale??'five_point'),[value,setValue]=useState(current?.numericValue?.toString()??current?.textValue??''),[max,setMax]=useState(current?.maxScore?.toString()??''),[saving,setSaving]=useState(false);async function save(){setSaving(true);try{await onSave({gradingScale:scale,maxScore:scale==='points'?Number(max):undefined,...(scale==='pass_fail'?{textValue:value}:{numericValue:Number(value)})})}finally{setSaving(false)}}return <div className="quarter-editor"><select aria-label="Шкала итога" value={scale} onChange={e=>{setScale(e.target.value as GradingScale);setValue('')}}><option value="five_point">2–5</option><option value="points">Баллы</option><option value="pass_fail">Зачёт</option></select>{scale==='points'&&<input aria-label="Максимум итога" type="number" min="0.01" value={max} onChange={e=>setMax(e.target.value)}/>} {scale==='pass_fail'?<select aria-label="Итог" value={value} onChange={e=>setValue(e.target.value)}><option value="">—</option><option value="pass">Зачёт</option><option value="fail">Незачёт</option></select>:<input aria-label="Итог" type="number" min={scale==='five_point'?2:0} max={scale==='five_point'?5:Number(max)||undefined} value={value} onChange={e=>setValue(e.target.value)}/>}<button className="secondary-action" disabled={!value||saving} onClick={save}>{saving?'…':'Сохранить итог'}</button></div>}

type JournalKind = 'homework' | 'classwork' | 'knowledge_check' | 'absence'
const gradeKinds: JournalKind[] = ['homework', 'classwork', 'knowledge_check', 'absence']
const gradeLabels: Record<JournalKind, string> = {
  homework: 'Домашняя работа',
  classwork: 'Работа на уроке',
  knowledge_check: 'Проверочная работа',
  absence: 'Пропуск',
}

function GradeIcon({ kind }: { kind: JournalKind }) {
  if (kind === 'absence') {
    return <svg viewBox="0 0 20 20" aria-hidden="true"><circle cx="10" cy="10" r="7"/><path d="m5 5 10 10"/></svg>
  }
  if (kind === 'homework') {
    return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m3 9 7-6 7 6v8h-5v-5H8v5H3V9Z"/></svg>
  }
  if (kind === 'knowledge_check') {
    return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="M6 4h8v13H6V4Zm2-2h4v4H8V2Zm0 8h4m-4 3h4"/></svg>
  }
  return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m4 14-1 3 3-1L16 6l-2-2L4 14Zm8-8 2 2"/></svg>
}

function formatLessonDate(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', timeZone: 'UTC' })
    .format(new Date(`${value}T00:00:00Z`))
}

function formatWeekday(value: string) {
  const day = new Intl.DateTimeFormat('ru-RU', { weekday: 'long', timeZone: 'UTC' })
    .format(new Date(`${value}T00:00:00Z`))
  return day.charAt(0).toUpperCase() + day.slice(1)
}

function GradeTable({
  items,
  students,
  from,
  to,
  onHomework,
  onSaved,
  onAbsence,
  onAddGrade,
}: {
  items: LessonRecord[]
  students: ClassStudent[]
  from: string
  to: string
  onHomework: (date: string, lesson?: LessonRecord) => void
  onSaved: (score: Score) => void
  onAbsence: (lessonID: string, studentID: string, absence?: Absence) => void
  onAddGrade: (lessonID: string) => void
}) {
  const dates = useMemo(() => {
    if (!from || !to) return items.map((item) => item.lessonDate)
    const result: string[] = []
    for (const day = new Date(from); day <= new Date(to); day.setUTCDate(day.getUTCDate() + 1)) {
      result.push(day.toISOString().slice(0, 10))
    }
    return result
  }, [from, to, items])
  const lessons = useMemo(() => new Map(items.map((item) => [item.lessonDate, item])), [items])
  if (!dates.length) return <div className="empty-card">Добавьте учебный период.</div>
  return (
    <div className="table-card journal-scroll">
      <table className="journal-table">
        <thead>
          <tr>
            <th className="student-heading">Ученик</th>
            {dates.map((date) => {
              const lesson = lessons.get(date)
              return (
                <th key={date}>
                  <span className="lesson-date">{formatLessonDate(date)}</span>
                  <span className="lesson-day">{formatWeekday(date)}</span>
                  {lesson?.topic && <span className="lesson-topic">{lesson.topic}</span>}
                  <button
                    className="homework-button"
                    aria-label={`Домашнее задание ${date}`}
                    title={lesson?.homework ?? 'Добавить задание'}
                    onClick={() => onHomework(date, lesson)}
                  >
                    Домашнее задание
                  </button>
                  <span className="homework-text">{lesson?.homework ?? '—'}</span>
                </th>
              )
            })}
          </tr>
        </thead>
        <tbody>
          {students.map((member) => (
            <tr key={member.student.id}>
              <th><span className="student-name">{member.student.name}</span></th>
              {dates.map((date) => {
                const lesson = lessons.get(date)
                const disabled = date < member.enrolledOn
                  || Boolean(member.leftOn && date > member.leftOn)
                return (
                  <td key={date}>
                    <div className="grade-cell">
                      {gradeKinds.map((kind) => {
                        if (kind === 'absence') {
                          return lesson ? (
                            <AbsenceToggle
                              key={kind}
                              lesson={lesson}
                              studentID={member.student.id}
                              disabled={disabled}
                              onChange={onAbsence}
                            />
                          ) : <EmptyGradeControl kind={kind} key={kind}/>
                        }
                        const gradeItem = lesson?.gradeItems.find((item) => item.kind === kind)
                        return gradeItem ? (
                          <ScoreEditor
                            item={gradeItem}
                            studentID={member.student.id}
                            disabled={disabled}
                            onSaved={onSaved}
                            key={kind}
                          />
                        ) : lesson ? (
                          <button
                            className="grade-control grade-control--empty"
                            aria-label={`Добавить ${gradeLabels[kind]}`}
                            title={gradeLabels[kind]}
                            onClick={() => onAddGrade(lesson.id)}
                            key={kind}
                          >
                            <GradeIcon kind={kind}/><span>—</span>
                          </button>
                        ) : <EmptyGradeControl kind={kind} key={kind}/>
                      })}
                    </div>
                  </td>
                )
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function EmptyGradeControl({ kind }: { kind: JournalKind }) {
  return <span className="grade-control grade-control--empty"><GradeIcon kind={kind}/><span>—</span></span>
}

function AbsenceToggle({ lesson, studentID, disabled, onChange }: {
  lesson: LessonRecord
  studentID: string
  disabled: boolean
  onChange: (lessonID: string, studentID: string, absence?: Absence) => void
}) {
  const initial = (lesson.absences ?? []).some((item) => item.studentId === studentID)
  const [absent, setAbsent] = useState(initial)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState(false)
  async function toggle() {
    setSaving(true)
    setError(false)
    try {
      if (absent) {
        await deleteStudentAbsence(lesson.id, studentID)
        setAbsent(false)
        onChange(lesson.id, studentID)
      } else {
        const saved = await putStudentAbsence(lesson.id, studentID)
        setAbsent(true)
        onChange(lesson.id, studentID, saved)
      }
    } catch {
      setError(true)
    } finally {
      setSaving(false)
    }
  }
  return (
    <button
      className={`grade-control${absent ? ' grade-control--absence' : ' grade-control--empty'}`}
      aria-pressed={absent}
      aria-label={absent ? 'Н' : 'Был'}
      title={error ? 'Ошибка сохранения' : gradeLabels.absence}
      disabled={disabled || saving}
      onClick={toggle}
    >
      <GradeIcon kind="absence"/><span>{absent ? 'Н' : '—'}</span>
    </button>
  )
}

function ScoreEditor({ item, studentID, disabled, onSaved }: {
  item: GradeItem
  studentID: string
  disabled: boolean
  onSaved: (score: Score) => void
}) {
  const existing = item.scores.find((score) => score.studentId === studentID)
  const initial = existing?.numericValue?.toString() ?? existing?.textValue ?? ''
  const [value, setValue] = useState(initial)
  const [comment, setComment] = useState(existing?.teacherComment ?? '')
  const [state, setState] = useState<'idle' | 'saving' | 'saved' | 'error'>('idle')
  async function save() {
    if (!value || disabled) return
    setState('saving')
    try {
      const payload = item.gradingScale === 'pass_fail'
        ? { textValue: value, teacherComment: comment || undefined }
        : { numericValue: Number(value), teacherComment: comment || undefined }
      onSaved(await putStudentScore(item.id, studentID, payload))
      setState('saved')
    } catch {
      setState('error')
    }
  }
  const kind = item.kind === 'knowledge_check' ? 'knowledge_check'
    : item.kind === 'homework' ? 'homework' : 'classwork'
  return (
    <div className="grade-control-wrap">
      <label className={`grade-control${value ? ` grade-control--${kind}` : ''}`}>
        <GradeIcon kind={kind}/>
        {item.gradingScale === 'pass_fail' ? (
          <select
            aria-label={`Оценка ${item.title}`}
            disabled={disabled}
            value={value}
            onChange={(event) => { setValue(event.target.value); setState('idle') }}
          >
            <option value="">—</option><option value="pass">З</option><option value="fail">Н</option>
          </select>
        ) : (
          <input
            aria-label={`Оценка ${item.title}`}
            disabled={disabled}
            type="number"
            min={item.gradingScale === 'five_point' ? 2 : 0}
            max={item.gradingScale === 'five_point' ? 5 : item.maxScore ?? undefined}
            step={item.gradingScale === 'five_point' ? 1 : .01}
            value={value}
            onChange={(event) => { setValue(event.target.value); setState('idle') }}
          />
        )}
      </label>
      <div className="grade-popover">
        <strong>{item.title}</strong>
        <input
          aria-label={`Комментарий ${item.title}`}
          disabled={disabled}
          placeholder="Комментарий"
          value={comment}
          onChange={(event) => { setComment(event.target.value); setState('idle') }}
        />
        <button className="secondary-action" disabled={disabled || !value || state === 'saving'} onClick={save}>
          {state === 'saving' ? '…' : state === 'saved' ? '✓' : 'Сохранить'}
        </button>
        {state === 'error' && <small role="alert">Ошибка сохранения</small>}
      </div>
    </div>
  )
}
function LessonDialog({draft,setDraft,saving,submit}:{draft:LessonDraft;setDraft:(value:LessonDraft|null)=>void;saving:boolean;submit:(event:FormEvent)=>void}){return <div className="dialog-backdrop"><form className="role-dialog role-form" onSubmit={submit}><h3>Новый урок</h3><label>Дата<input aria-label="Дата урока" type="date" value={draft.lessonDate} onChange={e=>setDraft({...draft,lessonDate:e.target.value})}/></label><label>Позиция<input aria-label="Позиция урока" type="number" min="1" value={draft.position} onChange={e=>setDraft({...draft,position:Number(e.target.value)})}/></label><label>Тема<input value={draft.topic} onChange={e=>setDraft({...draft,topic:e.target.value})} required/></label><label>Домашнее задание<input value={draft.homework} onChange={e=>setDraft({...draft,homework:e.target.value})}/></label><button type="button" className="secondary-action" onClick={()=>setDraft({...draft,materials:[...draft.materials,{title:'',url:'',position:draft.materials.length+1}]})}>+ Материал</button>{draft.materials.map((material,i)=><div className="form-columns" key={i}><label>Название материала<input aria-label={`Название материала ${i+1}`} value={material.title} onChange={e=>setDraft({...draft,materials:draft.materials.map((m,j)=>j===i?{...m,title:e.target.value}:m)})}/></label><label>Ссылка<input aria-label={`Ссылка материала ${i+1}`} type="url" value={material.url} onChange={e=>setDraft({...draft,materials:draft.materials.map((m,j)=>j===i?{...m,url:e.target.value}:m)})}/></label></div>)}<button className="primary-action" disabled={saving}>Сохранить</button></form></div>}
function GradeEditor({lessonID,onBack,onSaved}:{lessonID:string;onBack:()=>void;onSaved:(item:GradeItem)=>void}){const[value,setValue]=useState<GradeItemDraft>({title:'',kind:'homework',gradingScale:'five_point'}),[error,setError]=useState('');async function submit(e:FormEvent){e.preventDefault();try{const saved=await createGradeItem(lessonID,value);onSaved(saved);onBack()}catch(cause){setError(cause instanceof Error?cause.message:'Ошибка')}}return <div className="content-stack"><button className="back-button" onClick={onBack}>← К урокам</button><h2>Новая работа</h2>{error&&<p role="alert">{error}</p>}<form className="role-form" onSubmit={submit}><label>Название<input value={value.title} onChange={e=>setValue({...value,title:e.target.value})}/></label><label>Тип<select aria-label="Тип работы" value={value.kind} onChange={e=>setValue({...value,kind:e.target.value as GradeItem['kind']})}><option value="homework">Домашняя</option><option value="classwork">Классная</option><option value="knowledge_check">Проверочная</option><option value="other">Другая</option></select></label><label>Шкала<select aria-label="Шкала" value={value.gradingScale} onChange={e=>{const gradingScale=e.target.value as GradeItem['gradingScale'];setValue({...value,gradingScale,maxScore:undefined})}}><option value="five_point">Пятибалльная</option><option value="points">Баллы</option><option value="pass_fail">Зачёт</option></select></label>{value.gradingScale==='points'&&<label>Максимальный балл<input aria-label="Максимальный балл" type="number" min="0.01" step="0.01" onChange={e=>setValue({...value,maxScore:Number(e.target.value)})}/></label>}<button className="primary-action">Создать работу</button></form></div>}
