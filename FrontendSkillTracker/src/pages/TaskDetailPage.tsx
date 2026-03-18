import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Edit2, Trash2, Upload, Users, Calendar, ChevronRight, Zap } from 'lucide-react'
import { tasksApi } from '@/api/tasks'
import { usersApi } from '@/api/users'
import { skillsApi } from '@/api/skills'
import { useAuth } from '@/contexts/AuthContext'
import AnimatedPage from '@/components/common/AnimatedPage'
import StatusBadge from '@/components/common/StatusBadge'
import ProgressRing from '@/components/common/ProgressRing'
import SkillSelector from '@/components/skills/SkillSelector'
import TaskComments from '@/components/tasks/TaskComments'
import TaskHistory from '@/components/tasks/TaskHistory'
import RecommendedList from '@/components/employees/RecommendedList'
import TaskForm, { type TaskFormData } from '@/components/tasks/TaskForm'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { formatDate, formatDateTime, formatFileSize, getAvatarGradient, getInitials, isDeadlineOverdue, cn } from '@/lib/utils'
import { toast } from '@/hooks/use-toast'

export default function TaskDetailPage() {
  const { id } = useParams<{ id: string }>()
  const taskId = Number(id)
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { user } = useAuth()
  const isManager = user?.role === 'manager'

  const [editOpen, setEditOpen] = useState(false)

  const { data: task, isLoading } = useQuery({
    queryKey: ['task', taskId],
    queryFn: () => tasksApi.get(taskId),
  })
  const { data: employees = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => usersApi.list(),
    enabled: isManager,
  })
  const { data: allSkills = [] } = useQuery({
    queryKey: ['skills'],
    queryFn: () => skillsApi.list(),
    enabled: isManager,
  })
  const { data: recommended = [], isFetching: recLoading } = useQuery({
    queryKey: ['recommended', taskId],
    queryFn: () => tasksApi.getRecommendedEmployees(taskId),
    enabled: isManager,
  })

  const employeeMap = Object.fromEntries(employees.map(e => [e.id, e.name]))
  const assigneeName = task ? (employeeMap[task.employee_id] ?? `User #${task.employee_id}`) : ''

  const updateMutation = useMutation({
    mutationFn: (data: TaskFormData) => tasksApi.update(taskId, {
      title: data.title,
      description: data.description,
      employee_id: data.employee_id,
      deadline: data.deadline,
      status: data.status,
      progress: data.progress,
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['task', taskId] })
      qc.invalidateQueries({ queryKey: ['tasks'] })
      toast({ title: 'Task updated' })
      setEditOpen(false)
    },
    onError: () => toast({ title: 'Failed to update task', variant: 'destructive' }),
  })

  const deleteMutation = useMutation({
    mutationFn: () => tasksApi.delete(taskId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['tasks'] })
      navigate('/tasks')
      toast({ title: 'Task deleted' })
    },
    onError: () => toast({ title: 'Failed to delete task', variant: 'destructive' }),
  })

  const addSkillMutation = useMutation({
    mutationFn: (skillId: number) => tasksApi.addSkill(taskId, skillId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['task', taskId] })
      qc.invalidateQueries({ queryKey: ['recommended', taskId] })
    },
  })

  const removeSkillMutation = useMutation({
    mutationFn: (skillId: number) => tasksApi.removeSkill(taskId, skillId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['task', taskId] })
      qc.invalidateQueries({ queryKey: ['recommended', taskId] })
    },
  })

  const uploadMutation = useMutation({
    mutationFn: (file: File) => tasksApi.uploadAttachment(taskId, file),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['task', taskId] })
      toast({ title: 'File uploaded' })
    },
    onError: () => toast({ title: 'Upload failed', variant: 'destructive' }),
  })

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) uploadMutation.mutate(file)
    e.target.value = ''
  }

  const handleDelete = () => {
    if (confirm('Delete this task? This cannot be undone.')) {
      deleteMutation.mutate()
    }
  }

  if (isLoading) {
    return (
      <AnimatedPage>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-24 rounded-2xl bg-muted/40 animate-pulse" />
          ))}
        </div>
      </AnimatedPage>
    )
  }

  if (!task) {
    return (
      <AnimatedPage>
        <div className="flex flex-col items-center gap-4 py-20">
          <p className="text-muted-foreground">Task not found</p>
          <Button variant="outline" onClick={() => navigate('/tasks')}>Back to Tasks</Button>
        </div>
      </AnimatedPage>
    )
  }

  const overdue = task.status !== 'completed' && isDeadlineOverdue(task.deadline)
  const taskSkills = task.required_skills ?? []

  return (
    <AnimatedPage>
      <div className="space-y-6 max-w-5xl">
        {/* Breadcrumb */}
        <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
          <button onClick={() => navigate('/tasks')} className="hover:text-foreground transition-colors">Tasks</button>
          <ChevronRight className="h-4 w-4" />
          <span className="text-foreground font-medium truncate max-w-[300px]">{task.title}</span>
        </div>

        {/* Hero */}
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          className="rounded-2xl border border-border bg-card p-6"
        >
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div className="flex items-start gap-4 flex-1 min-w-0">
              <ProgressRing value={task.progress} size={72} strokeWidth={6} />
              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2 mb-2">
                  <h1 className="text-xl font-bold leading-tight">{task.title}</h1>
                  <StatusBadge status={task.status} />
                </div>
                {task.description && (
                  <p className="text-sm text-muted-foreground mb-3 leading-relaxed">{task.description}</p>
                )}
                <div className="flex flex-wrap gap-4 text-xs text-muted-foreground">
                  <div className="flex items-center gap-1.5">
                    <Calendar className={cn('h-3.5 w-3.5', overdue ? 'text-rose-400' : '')} />
                    <span className={overdue ? 'text-rose-400 font-medium' : ''}>
                      {overdue ? '⚠ Overdue · ' : ''}{formatDateTime(task.deadline)}
                    </span>
                  </div>
                  {assigneeName && (
                    <div className="flex items-center gap-1.5">
                      <div className={cn('flex h-4 w-4 items-center justify-center rounded-full text-[8px] font-bold text-white bg-gradient-to-br', getAvatarGradient(assigneeName))}>
                        {getInitials(assigneeName)}
                      </div>
                      <span>{assigneeName}</span>
                    </div>
                  )}
                  <span>Created {formatDate(task.created_at)}</span>
                </div>
              </div>
            </div>

            {isManager && (
              <div className="flex items-center gap-2">
                <Button variant="outline" size="sm" className="gap-1.5" onClick={() => setEditOpen(true)}>
                  <Edit2 className="h-3.5 w-3.5" />
                  Edit
                </Button>
                <Button variant="destructive" size="sm" className="gap-1.5" onClick={handleDelete} disabled={deleteMutation.isPending}>
                  <Trash2 className="h-3.5 w-3.5" />
                  Delete
                </Button>
              </div>
            )}
          </div>

          {/* Progress bar */}
          <div className="mt-5">
            <div className="flex justify-between text-xs text-muted-foreground mb-1.5">
              <span>Progress</span>
              <span className="font-medium">{task.progress}%</span>
            </div>
            <Progress value={task.progress} className="h-2" />
          </div>
        </motion.div>

        <div className="grid gap-6 lg:grid-cols-3">
          {/* Main content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Required Skills */}
            <motion.div
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="rounded-2xl border border-border bg-card p-5"
            >
              <div className="flex items-center gap-2 mb-4">
                <Zap className="h-4 w-4 text-amber-400" />
                <h2 className="font-semibold">Required Skills</h2>
                {taskSkills.length > 0 && (
                  <span className="rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground">{taskSkills.length}</span>
                )}
              </div>
              <SkillSelector
                allSkills={allSkills}
                selectedSkills={taskSkills}
                onAdd={skillId => addSkillMutation.mutate(skillId)}
                onRemove={skillId => removeSkillMutation.mutate(skillId)}
                readonly={!isManager}
              />
            </motion.div>

            {/* Tabs: Comments / History / Attachments */}
            <motion.div
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.15 }}
              className="rounded-2xl border border-border bg-card p-5"
            >
              <Tabs defaultValue="comments">
                <TabsList className="mb-4">
                  <TabsTrigger value="comments">Comments</TabsTrigger>
                  <TabsTrigger value="history">History</TabsTrigger>
                  <TabsTrigger value="attachments">Attachments</TabsTrigger>
                </TabsList>

                <TabsContent value="comments">
                  <TaskComments taskId={taskId} />
                </TabsContent>

                <TabsContent value="history">
                  <TaskHistory taskId={taskId} />
                </TabsContent>

                <TabsContent value="attachments">
                  <div className="space-y-3">
                    {/* Upload */}
                    <label className="flex cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed border-border/60 p-6 text-center hover:border-violet-500/40 hover:bg-muted/20 transition-colors">
                      <Upload className="h-8 w-8 text-muted-foreground/50" />
                      <div>
                        <p className="text-sm font-medium">Click to upload file</p>
                        <p className="text-xs text-muted-foreground">Any file type supported</p>
                      </div>
                      <input type="file" className="hidden" onChange={handleFileInput} disabled={uploadMutation.isPending} />
                    </label>

                    {/* Attachments list from task (task type may not include it, shown separately) */}
                    <p className="text-xs text-muted-foreground text-center">
                      {uploadMutation.isPending ? 'Uploading...' : 'Uploaded files appear after refresh'}
                    </p>
                  </div>
                </TabsContent>
              </Tabs>
            </motion.div>
          </div>

          {/* Sidebar: Recommended Employees */}
          {isManager && (
            <motion.div
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 }}
              className="rounded-2xl border border-border bg-card p-5 h-fit"
            >
              <div className="flex items-center gap-2 mb-4">
                <Users className="h-4 w-4 text-violet-400" />
                <h2 className="font-semibold">Recommended</h2>
              </div>
              <RecommendedList employees={recommended} isLoading={recLoading} />
            </motion.div>
          )}
        </div>
      </div>

      {/* Edit dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Task</DialogTitle>
          </DialogHeader>
          <TaskForm
            employees={employees}
            defaultValues={{
              title: task.title,
              description: task.description,
              employee_id: task.employee_id,
              deadline: task.deadline.slice(0, 16),
              status: task.status,
              progress: task.progress,
            }}
            onSubmit={async data => { await updateMutation.mutateAsync(data) }}
            isLoading={updateMutation.isPending}
          />
        </DialogContent>
      </Dialog>
    </AnimatedPage>
  )
}
