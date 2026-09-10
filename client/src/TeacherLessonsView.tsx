import{FormEvent,useCallback,useEffect,useMemo,useState}from'react'
import{ClassStudent,listClassStudents}from'./classes'
import{ClassSubjectRecord}from'./classSubjects'
import{Absence,GradeItem,LessonDraft,LessonRecord,Score,createGradeItem,createLesson,updateLesson,deleteStudentAbsence,listLessons,putStudentAbsence,putStudentScore}from'./lessons'
import{GradingScale,QuarterGrade,listQuarterGrades}from'./quarterGrades'

const emptyLesson=():LessonDraft=>({lessonDate:new Date().toISOString().slice(0,10),position:1,topic:'',homework:'',materials:[]})

type PeriodView = 'week' | 'period'
const periodViewCookie = 'journal_period_view'
function readPeriodView(): PeriodView {
 const value=document.cookie.split('; ').find(item=>item.startsWith(periodViewCookie+'='))?.split('=')[1]
 return value==='period'?'period':'week'
}
export default function TeacherLessonsView({assignment,quarters=[],className,onBack}:{assignment:ClassSubjectRecord;quarters?:{id:string;number:number;startsOn?:string;endsOn?:string}[];className?:string;onBack:()=>void}){
 const[items,setItems]=useState<LessonRecord[]>([]),[students,setStudents]=useState<ClassStudent[]>([]),[quarterGrades,setQuarterGrades]=useState<QuarterGrade[]>([]),[quarterID,setQuarterID]=useState(quarters[0]?.id??''),[periodView,setPeriodView]=useState<PeriodView>(readPeriodView),[draft,setDraft]=useState<LessonDraft|null>(null),[editingLessonID,setEditingLessonID]=useState(''),[newScore,setNewScore]=useState<{date:string;lesson?:LessonRecord;studentID:string;kind:GradeItem['kind']}|null>(null),[loading,setLoading]=useState(true),[saving,setSaving]=useState(false),[error,setError]=useState('')
 const selectedQuarter=quarters.find(item=>item.id===quarterID),from=selectedQuarter?.startsOn??'',to=selectedQuarter?.endsOn??''
 const load=useCallback(async()=>{setLoading(true);setError('');try{const[lessons,roster,totals]=await Promise.all([listLessons(assignment.id,from,to),listClassStudents(assignment.classId),quarterID?listQuarterGrades(assignment.id,quarterID):Promise.resolve([])]);setItems(lessons);setStudents(roster);setQuarterGrades(totals)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось загрузить журнал')}finally{setLoading(false)}},[assignment.id,assignment.classId,from,to,quarterID])
 useEffect(()=>{void load()},[load])
 async function submit(event:FormEvent){event.preventDefault();if(!draft)return;setSaving(true);try{const payload={...draft,homework:draft.homework||undefined};const saved=editingLessonID?await updateLesson(editingLessonID,payload):await createLesson(assignment.id,payload);setItems(current=>editingLessonID?current.map(item=>item.id===editingLessonID?{...saved,gradeItems:item.gradeItems,absences:item.absences}:item):[...current,saved]);setEditingLessonID('');setDraft(null)}catch(cause){setError(cause instanceof Error?cause.message:'Не удалось создать урок')}finally{setSaving(false)}}
 function scoreSaved(saved:Score){setItems(current=>current.map(lesson=>({...lesson,gradeItems:lesson.gradeItems.map(item=>item.id!==saved.gradeItemId?item:{...item,scores:[...item.scores.filter(score=>score.studentId!==saved.studentId),saved]})})))}
 function absenceChanged(lessonID:string,studentID:string,absence?:Absence){setItems(current=>current.map(lesson=>lesson.id!==lessonID?lesson:{...lesson,absences:absence?[...(lesson.absences??[]).filter(item=>item.studentId!==studentID),absence]:(lesson.absences??[]).filter(item=>item.studentId!==studentID)}))}
 function changePeriodView(value:PeriodView){setPeriodView(value);document.cookie=periodViewCookie+'='+value+'; Max-Age=31536000; Path=/; SameSite=Lax'}
 return <div className="content-stack"><button className="back-button" onClick={onBack}>← Предметы класса</button><div className="content-heading"><div><p className="content-kicker">{className?`${className} · ${assignment.subject.name}`:assignment.subject.name}</p><h2>Журнал класса</h2></div></div>{error&&<p className="content-error" role="alert">{error}</p>}{loading?<div className="empty-card">Загружаем журнал…</div>:<><div className="journal-toolbar"><div className="grade-legend" aria-label="Обозначения оценок">{gradeKinds.map(kind=><span key={kind}><GradeIcon kind={kind}/>{gradeLabels[kind]}</span>)}</div><div className="lesson-filters"><div className="period-view-toggle" aria-label="Отображение периода"><button type="button" aria-pressed={periodView==='week'} onClick={()=>changePeriodView('week')}>Неделя</button><button type="button" aria-pressed={periodView==='period'} onClick={()=>changePeriodView('period')}>Весь период</button></div>{quarters.length>0?<label><span className="visually-hidden">Период</span><select aria-label="Период" value={quarterID} onChange={e=>setQuarterID(e.target.value)}>{quarters.map(item=><option value={item.id} key={item.id}>{item.number} период</option>)}</select></label>:<span className="summary-empty">Нет периодов</span>}</div></div><GradeTable periodView={periodView} items={items} students={students} from={from} to={to} quarterID={quarterID} quarterGrades={quarterGrades} onHomework={(date,lesson)=>{setEditingLessonID(lesson?.id??'');setDraft(lesson?{lessonDate:lesson.lessonDate,position:lesson.position,topic:lesson.topic,homework:lesson.homework??'',materials:lesson.materials.map(m=>({title:m.title,url:m.url,position:m.position}))}:{...emptyLesson(),lessonDate:date})}} onSaved={scoreSaved} onAbsence={absenceChanged} onAddGrade={(date,lesson,studentID,kind)=>setNewScore({date,lesson,studentID,kind})}/></>}{draft&&<LessonDialog draft={draft} setDraft={setDraft} saving={saving} submit={submit}/>}{newScore&&<NewScoreDialog assignmentID={assignment.id} date={newScore.date} lesson={newScore.lesson} studentID={newScore.studentID} kind={newScore.kind} onBack={()=>setNewScore(null)} onSaved={(lesson,saved)=>{setItems(current=>current.some(item=>item.id===lesson.id)?current.map(item=>item.id===lesson.id?{...item,gradeItems:[...item.gradeItems,saved]}:item):[...current,{...lesson,gradeItems:[saved],absences:lesson.absences??[]}]);setNewScore(null)}}/>}</div>
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
  periodView,
  items,
  students,
  from,
  to,
  quarterID,
  quarterGrades,
  onHomework,
  onSaved,
  onAbsence,
  onAddGrade,
}: {
  periodView: PeriodView
  items: LessonRecord[]
  students: ClassStudent[]
  from: string
  to: string
  quarterID: string
  quarterGrades: QuarterGrade[]
  onHomework: (date: string, lesson?: LessonRecord) => void
  onSaved: (score: Score) => void
  onAbsence: (lessonID: string, studentID: string, absence?: Absence) => void
  onAddGrade: (
    date: string,
    lesson: LessonRecord | undefined,
    studentID: string,
    kind: GradeItem['kind'],
  ) => void
}) {
  const dates = useMemo(() => {
    if (!from || !to) return items.map((item) => item.lessonDate)
    let rangeStart = new Date(from)
    let rangeEnd = new Date(to)
    if (periodView === 'week') {
      const today = new Date().toISOString().slice(0, 10)
      const anchor = new Date(today >= from && today <= to ? today : from)
      const weekdayFromMonday = (anchor.getUTCDay() + 6) % 7
      const weekStart = new Date(anchor)
      weekStart.setUTCDate(weekStart.getUTCDate() - weekdayFromMonday)
      const weekEnd = new Date(weekStart)
      weekEnd.setUTCDate(weekEnd.getUTCDate() + 6)
      if (weekStart > rangeStart) rangeStart = weekStart
      if (weekEnd < rangeEnd) rangeEnd = weekEnd
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
  }, [from, to, items, periodView])
  const lessons = useMemo(() => new Map(items.map((item) => [item.lessonDate, item])), [items])
  if (!dates.length && !quarterID) return <div className="empty-card">Добавьте учебный период.</div>
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
                  return <td key={date} aria-label={`Нет темы ${date}`}/>
                }
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
                <QuarterGradeValue
                  current={quarterGrades.find((item) => (
                    item.studentId === member.student.id
                  ))}
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
      <GradeIcon kind="absence"/><span>{absent ? 'Н' : '—'}</span>
    </button>
  )
}

function displayGrade(value?: { numericValue?: number | null; textValue?: string | null }) {
  if (value?.numericValue !== undefined) return String(value.numericValue)
  if (value?.textValue === 'pass') return 'З'
  if (value?.textValue === 'fail') return 'Н'
  return ''
}

function ScoreEditor({ item, studentID, disabled, onSaved }: {
  item: GradeItem
  studentID: string
  disabled: boolean
  onSaved: (score: Score) => void
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

  return (
    <>
      <button
        className={`grade-control${existing ? ` grade-control--${kind}` : ''}`}
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
            {state === 'error' && <p className="content-error" role="alert">Ошибка сохранения</p>}
            <div className="dialog-actions">
              <button type="button" className="secondary-action" onClick={() => setOpen(false)}>
                Отмена
              </button>
              <button className="primary-action" disabled={!value || state === 'saving'}>
                {state === 'saving' ? 'Сохраняем…' : 'Сохранить оценку'}
              </button>
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

function QuarterGradeValue({ current }: { current?: QuarterGrade }) {
  return (
    <span
      className={`grade-control${current ? ' grade-control--total' : ' grade-control--empty'}`}
      title={current?.teacherComment || 'Итог будет рассчитан автоматически'}
    >
      <span>{displayGrade(current)}{current?.teacherComment ? '*' : ''}</span>
    </span>
  )
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