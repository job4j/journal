import{FormEvent,useCallback,useEffect,useMemo,useState}from'react'
import{ClassStudent,listClassStudents}from'./classes'
import{ClassSubjectRecord}from'./classSubjects'
import{Absence,GradeItem,LessonDraft,LessonRecord,Score,createGradeItem,createLesson,updateLesson,deleteStudentAbsence,deleteStudentScore,listLessons,putStudentAbsence,putStudentScore}from'./lessons'
import{GradingScale,QuarterGrade,deleteQuarterGrade,listQuarterGrades,putQuarterGrade}from'./quarterGrades'

const emptyLesson=():LessonDraft=>({lessonDate:new Date().toISOString().slice(0,10),position:1,topic:'',homework:'',materials:[]})

type PeriodView = 'week' | 'period'
const periodViewCookie = 'journal_period_view'
function readPeriodView(): PeriodView {
 const value=document.cookie.split('; ').find(item=>item.startsWith(periodViewCookie+'='))?.split('=')[1]
 return value==='period'?'period':'week'
}
function getWeekWindow(from:string,to:string,offset:number){
 if(!from||!to)return null
 const rangeStart=new Date(`${from}T00:00:00Z`),rangeEnd=new Date(`${to}T00:00:00Z`)
 const today=new Date().toISOString().slice(0,10),anchor=new Date(`${today>=from&&today<=to?today:from}T00:00:00Z`)
 const weekStart=new Date(anchor)
 weekStart.setUTCDate(weekStart.getUTCDate()-(weekStart.getUTCDay()+6)%7+offset*7)
 const weekEnd=new Date(weekStart)
 weekEnd.setUTCDate(weekEnd.getUTCDate()+6)
 return{start:weekStart>rangeStart?weekStart:rangeStart,end:weekEnd<rangeEnd?weekEnd:rangeEnd,hasPrevious:weekStart>rangeStart,hasNext:weekEnd<rangeEnd}
}
function formatPeriodDate(value: string) {
 const [year,month,day]=value.split('-')
 return `${day}.${month}.${year}`
}
function formatPeriodLabel(item:{number:number;name?:string;startsOn?:string;endsOn?:string}){
 const name=item.name||item.number+' период'
 return item.startsOn&&item.endsOn?`${name} · ${formatPeriodDate(item.startsOn)} — ${formatPeriodDate(item.endsOn)}`:name
}
export default function TeacherLessonsView({assignment,quarters=[]}:{assignment:ClassSubjectRecord;quarters?:{id:string;number:number;name?:string;startsOn?:string;endsOn?:string}[]}){
 const[items,setItems]=useState<LessonRecord[]>([]),[students,setStudents]=useState<ClassStudent[]>([]),[quarterGrades,setQuarterGrades]=useState<QuarterGrade[]>([]),[quarterID,setQuarterID]=useState(quarters[0]?.id??''),[periodView,setPeriodView]=useState<PeriodView>(readPeriodView),[weekOffset,setWeekOffset]=useState(0),[draft,setDraft]=useState<LessonDraft|null>(null),[editingLessonID,setEditingLessonID]=useState(''),[newScore,setNewScore]=useState<{date:string;lesson?:LessonRecord;studentID:string;kind:GradeItem['kind']}|null>(null),[loading,setLoading]=useState(true),[saving,setSaving]=useState(false),[error,setError]=useState('')
 const selectedQuarter=quarters.find(item=>item.id===quarterID),from=selectedQuarter?.startsOn??'',to=selectedQuarter?.endsOn??''
 const weekWindow=useMemo(()=>getWeekWindow(from,to,weekOffset),[from,to,weekOffset])
 const load=useCallback(async()=>{setLoading(true);setError('');try{const[lessons,roster,totals]=await Promise.all([listLessons(assignment.id,from,to),listClassStudents(assignment.classId),quarterID?listQuarterGrades(assignment.id,quarterID):Promise.resolve([])]);setItems(lessons);setStudents(roster);setQuarterGrades(totals)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось загрузить журнал')}finally{setLoading(false)}},[assignment.id,assignment.classId,from,to,quarterID])
 useEffect(()=>{void load()},[load])
 async function submit(event:FormEvent){event.preventDefault();if(!draft)return;setSaving(true);try{const payload={...draft,homework:draft.homework||undefined};const saved=editingLessonID?await updateLesson(editingLessonID,payload):await createLesson(assignment.id,payload);setItems(current=>editingLessonID?current.map(item=>item.id===editingLessonID?{...saved,gradeItems:item.gradeItems,absences:item.absences}:item):[...current,saved]);setEditingLessonID('');setDraft(null)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось создать урок')}finally{setSaving(false)}}
 function scoreSaved(saved:Score){setItems(current=>current.map(lesson=>({...lesson,gradeItems:lesson.gradeItems.map(item=>item.id!==saved.gradeItemId?item:{...item,scores:[...item.scores.filter(score=>score.studentId!==saved.studentId),saved]})})))}
 function scoreDeleted(gradeItemID:string,studentID:string){setItems(current=>current.map(lesson=>({...lesson,gradeItems:lesson.gradeItems.map(item=>item.id!==gradeItemID?item:{...item,scores:item.scores.filter(score=>score.studentId!==studentID)})})))}
 function absenceChanged(lessonID:string,studentID:string,absence?:Absence){setItems(current=>current.map(lesson=>lesson.id!==lessonID?lesson:{...lesson,absences:absence?[...(lesson.absences??[]).filter(item=>item.studentId!==studentID),absence]:(lesson.absences??[]).filter(item=>item.studentId!==studentID)}))}
 function quarterGradeSaved(saved:QuarterGrade){setQuarterGrades(current=>[...current.filter(item=>item.studentId!==saved.studentId),saved])}
 function quarterGradeDeleted(studentID:string){setQuarterGrades(current=>current.filter(item=>item.studentId!==studentID))}
 function changePeriodView(value:PeriodView){setPeriodView(value);document.cookie=periodViewCookie+'='+value+'; Max-Age=31536000; Path=/; SameSite=Lax'}
 return <div className="content-stack">{error&&<p className="content-error" role="alert">{error}</p>}{loading?<div className="empty-card">Загружаем журнал…</div>:<><div className="journal-toolbar"><div className="grade-legend" aria-label="Обозначения оценок">{gradeKinds.map(kind=><span key={kind}><GradeIcon kind={kind}/>{gradeLabels[kind]}</span>)}</div><div className="lesson-filters"><div className="period-view-toggle" aria-label="Отображение периода"><button type="button" aria-pressed={periodView==='week'} onClick={()=>changePeriodView('week')}>Неделя</button><button type="button" aria-pressed={periodView==='period'} onClick={()=>changePeriodView('period')}>Весь период</button></div>{periodView==='week'&&<div className="week-pagination" aria-label="Навигация по неделям"><button type="button" aria-label="Предыдущая неделя" disabled={!weekWindow?.hasPrevious} onClick={()=>setWeekOffset(value=>value-1)}>‹</button><button type="button" aria-label="Следующая неделя" disabled={!weekWindow?.hasNext} onClick={()=>setWeekOffset(value=>value+1)}>›</button></div>}{quarters.length>0?<label><span className="visually-hidden">Период</span><select aria-label="Период" value={quarterID} onChange={e=>{setQuarterID(e.target.value);setWeekOffset(0)}}>{quarters.map(item=><option value={item.id} key={item.id}>{formatPeriodLabel(item)}</option>)}</select></label>:<span className="summary-empty">Нет периодов</span>}</div></div><GradeTable assignmentID={assignment.id} periodView={periodView} weekOffset={weekOffset} items={items} students={students} from={from} to={to} quarterID={quarterID} quarterGrades={quarterGrades} onHomework={(date,lesson)=>{setEditingLessonID(lesson?.id??'');setDraft(lesson?{lessonDate:lesson.lessonDate,position:lesson.position,topic:lesson.topic,homework:lesson.homework??'',materials:lesson.materials.map(m=>({title:m.title,url:m.url,position:m.position}))}:{...emptyLesson(),lessonDate:date})}} onSaved={scoreSaved} onDeleted={scoreDeleted} onAbsence={absenceChanged} onQuarterGradeSaved={quarterGradeSaved} onQuarterGradeDeleted={quarterGradeDeleted} onAddGrade={(date,lesson,studentID,kind)=>setNewScore({date,lesson,studentID,kind})}/></>}{draft&&<LessonDialog draft={draft} setDraft={setDraft} saving={saving} submit={submit}/>}{newScore&&<NewScoreDialog assignmentID={assignment.id} date={newScore.date} lesson={newScore.lesson} studentID={newScore.studentID} kind={newScore.kind} onBack={()=>setNewScore(null)} onSaved={(lesson,saved)=>{setItems(current=>current.some(item=>item.id===lesson.id)?current.map(item=>item.id===lesson.id?{...item,gradeItems:[...item.gradeItems,saved]}:item):[...current,{...lesson,gradeItems:[saved],absences:lesson.absences??[]}]);setNewScore(null)}}/>}</div>
}

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

function TasksIcon() {
  return <svg viewBox="0 0 20 20" aria-hidden="true"><path d="M6 3h9a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2Z"/><path d="m7 8 1.2 1.2L10.5 7M12 8h2M7 13l1.2 1.2 2.3-2.2M12 13h2"/></svg>
}

function formatLessonDate(value: string) {
  return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', timeZone: 'UTC' })
    .format(new Date(`${value}T00:00:00Z`))
}

function formatWeekday(value: string) {
  const day = new Intl.DateTimeFormat('ru-RU', { weekday: 'short', timeZone: 'UTC' })
    .format(new Date(`${value}T00:00:00Z`))
    .replace(/\.$/, '')
  return `${day.charAt(0).toUpperCase()}${day.slice(1)}.`
}

function GradeTable({
  assignmentID,
  periodView,
  weekOffset,
  items,
  students,
  from,
  to,
  quarterID,
  quarterGrades,
  onHomework,
  onSaved,
  onDeleted,
  onAbsence,
  onQuarterGradeSaved,
  onQuarterGradeDeleted,
  onAddGrade,
}: {
  assignmentID: string
  periodView: PeriodView
  weekOffset: number
  items: LessonRecord[]
  students: ClassStudent[]
  from: string
  to: string
  quarterID: string
  quarterGrades: QuarterGrade[]
  onHomework: (date: string, lesson?: LessonRecord) => void
  onSaved: (score: Score) => void
  onDeleted: (gradeItemID: string, studentID: string) => void
  onAbsence: (lessonID: string, studentID: string, absence?: Absence) => void
  onQuarterGradeSaved: (grade: QuarterGrade) => void
  onQuarterGradeDeleted: (studentID: string) => void
  onAddGrade: (
    date: string,
    lesson: LessonRecord | undefined,
    studentID: string,
    kind: GradeItem['kind'],
  ) => void
}) {
  const today = new Date().toISOString().slice(0, 10)
  const dates = useMemo(() => {
    if (!from || !to) return items.map((item) => item.lessonDate)
    let rangeStart = new Date(from)
    let rangeEnd = new Date(to)
    if (periodView === 'week') {
      const window = getWeekWindow(from, to, weekOffset)
      if (window) {
        rangeStart = window.start
        rangeEnd = window.end
      }
    }
    const result: string[] = []
    for (
      const day = new Date(rangeStart);
      day <= rangeEnd;
      day.setUTCDate(day.getUTCDate() + 1)
    ) {
      result.push(day.toISOString().slice(0, 10))
    }
    return result
  }, [from, to, items, periodView, weekOffset])
  const lessons = useMemo(() => new Map(items.map((item) => [item.lessonDate, item])), [items])
  if (!dates.length && !quarterID) return <div className="empty-card">Добавьте учебный период.</div>
  return (
    <div className="table-card journal-scroll">
      <table className="journal-table">
        <thead>
          <tr>
            <th className="student-heading" aria-label="Ученик"/>
            {dates.map((date) => {
              const lesson = lessons.get(date)
              return (
                <th className={date === today ? 'journal-today' : undefined} key={date}>
                  <div className="lesson-date-row">
                    <span className="lesson-date">{formatLessonDate(date)}, {formatWeekday(date)}</span>
                    <button
                      type="button"
                      className="lesson-tasks-button"
                      aria-label={`Домашнее задание ${date}`}
                      title="Заполнить тему и домашнее задание"
                      onClick={() => onHomework(date, lesson)}
                    >
                      <TasksIcon/>
                    </button>
                  </div>
                  {lesson?.topic && <span className="lesson-topic">{lesson.topic}</span>}
                  {lesson?.homework && <span className="homework-text"><strong>Д.З.:</strong> {lesson.homework}</span>}
                </th>
              )
            })}
            <th className="journal-total-heading">Итоги</th>
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
                if (!lesson?.topic.trim()) {
                  return (
                    <td
                      className={date === today ? 'journal-today' : undefined}
                      key={date}
                      aria-label={`Нет темы ${date}`}
                    />
                  )
                }
                return (
                  <td className={date === today ? 'journal-today' : undefined} key={date}>
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
                            onDeleted={onDeleted}
                            key={kind}
                          />
                        ) : lesson ? (
                          <button
                            className="grade-control grade-control--empty"
                            aria-label={`Добавить ${gradeLabels[kind]}`}
                            title={gradeLabels[kind]}
                            onClick={() => onAddGrade(date, lesson, member.student.id, kind)}
                            key={kind}
                          >
                            <GradeIcon kind={kind}/><span className="grade-value"/>
                          </button>
                        ) : (
                          <button
                            className="grade-control grade-control--empty"
                            aria-label={`Добавить ${gradeLabels[kind]}`}
                            title={gradeLabels[kind]}
                            disabled={disabled}
                            onClick={() => onAddGrade(
                              date,
                              undefined,
                              member.student.id,
                              kind,
                            )}
                            key={kind}
                          >
                            <GradeIcon kind={kind}/><span className="grade-value"/>
                          </button>
                        )
                      })}
                    </div>
                  </td>
                )
              })}
              <td className="journal-total-cell">
                <QuarterGradeEditor
                  assignmentID={assignmentID}
                  quarterID={quarterID}
                  studentID={member.student.id}
                  studentName={member.student.name}
                  current={quarterGrades.find((item) => item.studentId === member.student.id)}
                  onSaved={onQuarterGradeSaved}
                  onDeleted={onQuarterGradeDeleted}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function EmptyGradeControl({ kind }: { kind: JournalKind }) {
  return <span className="grade-control grade-control--empty"><GradeIcon kind={kind}/><span className="grade-value"/></span>
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
      <GradeIcon kind="absence"/>
    </button>
  )
}

function displayGrade(value?: { numericValue?: number | null; textValue?: string | null }) {
  if (value?.numericValue !== undefined) return String(value.numericValue)
  if (value?.textValue === 'pass') return 'З'
  if (value?.textValue === 'fail') return 'Н'
  return ''
}

function ScoreEditor({ item, studentID, disabled, onSaved, onDeleted }: {
  item: GradeItem
  studentID: string
  disabled: boolean
  onSaved: (score: Score) => void
  onDeleted: (gradeItemID: string, studentID: string) => void
}) {
  const existing = item.scores.find((score) => score.studentId === studentID)
  const [open, setOpen] = useState(false)
  const [value, setValue] = useState('')
  const [comment, setComment] = useState('')
  const [state, setState] = useState<'idle' | 'saving' | 'error'>('idle')
  const kind = item.kind === 'knowledge_check' ? 'knowledge_check'
    : item.kind === 'homework' ? 'homework' : 'classwork'

  function show() {
    setValue(existing?.numericValue?.toString() ?? existing?.textValue ?? '')
    setComment(existing?.teacherComment ?? '')
    setState('idle')
    setOpen(true)
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    if (!value || disabled) return
    setState('saving')
    try {
      const payload = item.gradingScale === 'pass_fail'
        ? { textValue: value, teacherComment: comment || undefined }
        : { numericValue: Number(value), teacherComment: comment || undefined }
      onSaved(await putStudentScore(item.id, studentID, payload))
      setOpen(false)
    } catch {
      setState('error')
    }
  }

  async function remove() {
    if (!existing || disabled) return
    setState('saving')
    try {
      await deleteStudentScore(item.id, studentID)
      onDeleted(item.id, studentID)
      setOpen(false)
    } catch {
      setState('error')
    }
  }

  return (
    <>
      <button
        className={`grade-control ${existing ? `grade-control--${kind}` : 'grade-control--empty'}`}
        aria-label={`Изменить оценку ${item.title}`}
        title={existing?.teacherComment || item.title}
        disabled={disabled}
        onClick={show}
      >
        <GradeIcon kind={kind}/>
        <span className="grade-value">{displayGrade(existing)}{existing?.teacherComment ? '*' : ''}</span>
      </button>
      {open && (
        <div className="dialog-backdrop">
          <form className="role-dialog role-form" onSubmit={save}>
            <h3>{item.title}</h3>
            <GradeValueField
              label={`Оценка ${item.title}`}
              scale={item.gradingScale}
              maxScore={item.maxScore ?? undefined}
              value={value}
              setValue={setValue}
            />
            <label>
              Комментарий
              <textarea
                aria-label={`Комментарий ${item.title}`}
                value={comment}
                onChange={(event) => setComment(event.target.value)}
              />
            </label>
            {state === 'error' && <p className="content-error" role="alert">Не удалось изменить оценку</p>}
            <div className="dialog-actions">
              <button className="primary-action" disabled={!value || state === 'saving'}>
                {state === 'saving' ? 'Сохраняем…' : 'Сохранить'}
              </button>
              <button type="button" className="secondary-action" onClick={() => setOpen(false)}>
                Отмена
              </button>
              {existing && <button type="button" className="secondary-action danger-action dialog-delete-action" disabled={state === 'saving'} onClick={remove}>Удалить</button>}
            </div>
          </form>
        </div>
      )}
    </>
  )
}

function GradeValueField({ label, scale, maxScore, value, setValue }: {
  label: string
  scale: GradingScale
  maxScore?: number
  value: string
  setValue: (value: string) => void
}) {
  if (scale === 'pass_fail') {
    return (
      <label>
        Оценка
        <select aria-label={label} value={value} onChange={(event) => setValue(event.target.value)}>
          <option value="">—</option>
          <option value="pass">Зачёт</option>
          <option value="fail">Незачёт</option>
        </select>
      </label>
    )
  }
  return (
    <label>
      Оценка
      <input
        aria-label={label}
        type="number"
        min={scale === 'five_point' ? 2 : 0}
        max={scale === 'five_point' ? 5 : maxScore}
        step={scale === 'five_point' ? 1 : .01}
        value={value}
        onChange={(event) => setValue(event.target.value)}
      />
    </label>
  )
}

function QuarterGradeEditor({assignmentID,quarterID,studentID,studentName,current,onSaved,onDeleted}:{assignmentID:string;quarterID:string;studentID:string;studentName:string;current?:QuarterGrade;onSaved:(grade:QuarterGrade)=>void;onDeleted:(studentID:string)=>void}){
 const[open,setOpen]=useState(false),[value,setValue]=useState(''),[comment,setComment]=useState(''),[state,setState]=useState<'idle'|'saving'|'error'>('idle')
 const scale=current?.gradingScale??'five_point'
 function show(){setValue(current?.numericValue?.toString()??current?.textValue??'');setComment(current?.teacherComment??'');setState('idle');setOpen(true)}
 async function save(event:FormEvent){event.preventDefault();if(!value)return;setState('saving');try{const payload=scale==='pass_fail'?{gradingScale:scale,textValue:value,teacherComment:comment||undefined}:{gradingScale:scale,maxScore:current?.maxScore??undefined,numericValue:Number(value),teacherComment:comment||undefined};onSaved(await putQuarterGrade(assignmentID,quarterID,studentID,payload));setOpen(false)}catch{setState('error')}}
 async function remove(){if(!current)return;setState('saving');try{await deleteQuarterGrade(assignmentID,quarterID,studentID);onDeleted(studentID);setOpen(false)}catch{setState('error')}}
 return <><button className={`grade-control${current?' grade-control--total':' grade-control--empty'}`} aria-label={`Изменить итог ${studentName}`} title={current?.teacherComment||'Задать итог'} onClick={show}><span>{displayGrade(current)}{current?.teacherComment?'*':''}</span></button>{open&&<div className="dialog-backdrop"><form className="role-dialog role-form" onSubmit={save}><h3>Итог · {studentName}</h3><GradeValueField label={`Итоговая оценка ${studentName}`} scale={scale} maxScore={current?.maxScore??undefined} value={value} setValue={setValue}/><label>Комментарий<textarea aria-label={`Комментарий к итогу ${studentName}`} value={comment} onChange={event=>setComment(event.target.value)}/></label>{state==='error'&&<p className="content-error" role="alert">Не удалось изменить итог</p>}<div className="dialog-actions"><button className="primary-action" disabled={!value||state==='saving'}>{state==='saving'?'Сохраняем…':'Сохранить'}</button><button type="button" className="secondary-action" onClick={()=>setOpen(false)}>Отменить</button><button type="button" className="secondary-action danger-action dialog-delete-action" disabled={!current||state==='saving'} onClick={remove}>Удалить</button></div></form></div>}</>
}
function LessonDialog({draft,setDraft,saving,submit}:{draft:LessonDraft;setDraft:(value:LessonDraft|null)=>void;saving:boolean;submit:(event:FormEvent)=>void}){return <div className="dialog-backdrop"><form className="role-dialog role-form" onSubmit={submit}><h3>Новый урок</h3><label>Дата<input aria-label="Дата урока" type="date" value={draft.lessonDate} onChange={e=>setDraft({...draft,lessonDate:e.target.value})}/></label><label>Позиция<input aria-label="Позиция урока" type="number" min="1" value={draft.position} onChange={e=>setDraft({...draft,position:Number(e.target.value)})}/></label><label>Тема<input value={draft.topic} onChange={e=>setDraft({...draft,topic:e.target.value})} required/></label><label>Домашнее задание<input value={draft.homework} onChange={e=>setDraft({...draft,homework:e.target.value})}/></label><button type="button" className="secondary-action" onClick={()=>setDraft({...draft,materials:[...draft.materials,{title:'',url:'',position:draft.materials.length+1}]})}>+ Материал</button>{draft.materials.map((material,i)=><div className="form-columns" key={i}><label>Название материала<input aria-label={`Название материала ${i+1}`} value={material.title} onChange={e=>setDraft({...draft,materials:draft.materials.map((m,j)=>j===i?{...m,title:e.target.value}:m)})}/></label><label>Ссылка<input aria-label={`Ссылка материала ${i+1}`} type="url" value={material.url} onChange={e=>setDraft({...draft,materials:draft.materials.map((m,j)=>j===i?{...m,url:e.target.value}:m)})}/></label></div>)}<button className="primary-action" disabled={saving}>Сохранить</button></form></div>}
function NewScoreDialog({ assignmentID, date, lesson, studentID, kind, onBack, onSaved }: {
  assignmentID: string
  date: string
  lesson?: LessonRecord
  studentID: string
  kind: GradeItem['kind']
  onBack: () => void
  onSaved: (lesson: LessonRecord, item: GradeItem) => void
}) {
  const title = gradeLabels[kind as JournalKind] ?? 'Работа'
  const [value, setValue] = useState('')
  const [comment, setComment] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  async function submit(event: FormEvent) {
    event.preventDefault()
    if (!value) return
    setSaving(true)
    setError('')
    try {
      const currentLesson = lesson ?? await createLesson(assignmentID, {
        lessonDate: date,
        position: 1,
        topic: 'Без темы',
        materials: [],
      })
      const item = await createGradeItem(currentLesson.id, {
        title,
        kind,
        gradingScale: 'five_point',
      })
      const score = await putStudentScore(item.id, studentID, {
        numericValue: Number(value),
        teacherComment: comment || undefined,
      })
      onSaved(currentLesson, { ...item, scores: [score] })
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : 'Не удалось сохранить оценку')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="dialog-backdrop">
      <form className="role-dialog role-form" onSubmit={submit}>
        <h3>{title}</h3>
        <label>
          Оценка
          <input
            aria-label={`Оценка ${title}`}
            type="number"
            min="2"
            max="5"
            step="1"
            value={value}
            onChange={(event) => setValue(event.target.value)}
            autoFocus
          />
        </label>
        <label>
          Комментарий
          <textarea
            aria-label={`Комментарий ${title}`}
            value={comment}
            onChange={(event) => setComment(event.target.value)}
          />
        </label>
        {error && <p className="content-error" role="alert">{error}</p>}
        <div className="dialog-actions">
          <button type="button" className="secondary-action" onClick={onBack}>
            Отмена
          </button>
          <button className="primary-action" disabled={!value || saving}>
            {saving ? 'Сохраняем…' : 'Сохранить оценку'}
          </button>
        </div>
      </form>
    </div>
  )
}
