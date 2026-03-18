import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { Plus, LayoutGrid, List } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { tasksApi } from '@/api/tasks'
import { usersApi } from '@/api/users'
import { useAuth } from '@/contexts/AuthContext'
import AnimatedPage from '@/components/common/AnimatedPage'
import EmptyState from '@/components/common/EmptyState'
import StatusBadge from '@/components/common/StatusBadge'
import TaskKanban from '@/components/tasks/TaskKanban'
import TaskFilters from '@/components/tasks/TaskFilters'
import TaskForm, { type TaskFormData } from '@/components/tasks/TaskForm'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { formatDate, isDeadlineOverdue, isDeadlineSoon, getAvatarGradient, getInitials, cn } from '@/lib/utils'
import { toast } from '@/hooks/use-toast'
import type { TaskFilter, TaskStatus } from '@/types'
import { Calendar } from 'lucide-react'

export default function TasksPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const isManager = user?.role === 'manager'

  const [view, setView] = useState<'kanban' | 'list'>('kanban')
  const [filters, setFilters] = useState<TaskFilter>({})
  const [dialogOpen, setDialogOpen] = useState(false)
  const [defaultStatus, setDefaultStatus] = useState<TaskStatus>('pending')

  const { data: allTasks = [] } = useQuery({
    queryKey: ['tasks'],
    queryFn: () => tasksApi.list(),
    enabled: isManager,
  })
  const { data: myTasks = [] } = useQuery({
    queryKey: ['my-tasks'],
    queryFn: () => tasksApi.myTasks(),
    enabled: !isManager,
  })
  const { data: employees = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => usersApi.list(),
    enabled: isManager,
  })

  const rawTasks = isManager ? allTasks : myTasks

  // Client-side filter
  const tasks = rawTasks.filter(t => {
    if (filters.status && t.status !== filters.status) return false
    if (filters.employee_id && t.employee_id !== filters.employee_id) return false
    if (filters.search) {
      const q = filters.search.toLowerCase()
      if (!t.title.toLowerCase().includes(q) && !t.description?.toLowerCase().includes(q)) return false
    }
    return true
  })

  const employeeMap = Object.fromEntries(employees.map(e => [e.id, e.name]))

  const createMutation = useMutation({
    mutationFn: (data: TaskFormData) => tasksApi.create({
      title: data.title,
      description: data.description,
      employee_id: data.employee_id,
      deadline: data.deadline,
      status: data.status,
      progress: data.progress,
    }),
    onSuccess: (task) => {
      qc.invalidateQueries({ queryKey: ['tasks'] })
      toast({ title: 'Task created' })
      setDialogOpen(false)
      navigate(`/tasks/${task.id}`)
    },
    onError: () => toast({ title: 'Failed to create task', variant: 'destructive' }),
  })

  const openDialog = (status: TaskStatus = 'pending') => {
    setDefaultStatus(status)
    setDialogOpen(true)
  }

  return (
    <AnimatedPage>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold">Tasks</h1>
            <p className="text-sm text-muted-foreground mt-0.5">
              {tasks.length} task{tasks.length !== 1 ? 's' : ''}
            </p>
          </div>
          <div className="flex items-center gap-2">
            {/* View toggle */}
            <div className="flex rounded-lg border border-border p-0.5 bg-muted/40">
              <button
                onClick={() => setView('kanban')}
                className={cn('flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors', view === 'kanban' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground')}
              >
                <LayoutGrid className="h-3.5 w-3.5" />
                Kanban
              </button>
              <button
                onClick={() => setView('list')}
                className={cn('flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors', view === 'list' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground')}
              >
                <List className="h-3.5 w-3.5" />
                List
              </button>
            </div>

            {isManager && (
              <Button onClick={() => openDialog()} className="gap-1.5" size="sm">
                <Plus className="h-4 w-4" />
                New Task
              </Button>
            )}
          </div>
        </div>

        {/* Filters */}
        <TaskFilters
          filters={filters}
          onChange={setFilters}
          employees={employees}
          showEmployeeFilter={isManager}
        />

        {/* Content */}
        <AnimatePresence mode="wait">
          {tasks.length === 0 ? (
            <EmptyState
              icon={Plus}
              title="No tasks found"
              description={isManager ? 'Create your first task to get started.' : 'No tasks assigned to you yet.'}
              action={isManager ? <Button onClick={() => openDialog()}>Create task</Button> : undefined}
            />
          ) : view === 'kanban' ? (
            <motion.div key="kanban" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
              <TaskKanban
                tasks={tasks}
                employeeMap={employeeMap}
                isManager={isManager}
                onAddTask={openDialog}
              />
            </motion.div>
          ) : (
            <motion.div key="list" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
              <div className="rounded-2xl border border-border overflow-hidden">
                {tasks.map((task, i) => {
                  const overdue = task.status !== 'completed' && isDeadlineOverdue(task.deadline)
                  const soon = !overdue && isDeadlineSoon(task.deadline)
                  const empName = employeeMap[task.employee_id]
                  return (
                    <motion.div
                      key={task.id}
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      transition={{ delay: i * 0.03 }}
                      onClick={() => navigate(`/tasks/${task.id}`)}
                      className="flex items-center gap-4 px-4 py-3.5 hover:bg-muted/40 cursor-pointer transition-colors border-b border-border/50 last:border-0 group"
                    >
                      <div className="min-w-0 flex-1">
                        <p className="font-medium text-sm truncate group-hover:text-primary transition-colors">{task.title}</p>
                        {task.description && <p className="text-xs text-muted-foreground truncate mt-0.5">{task.description}</p>}
                      </div>
                      {task.progress > 0 && (
                        <div className="hidden sm:flex items-center gap-2 w-24 shrink-0">
                          <Progress value={task.progress} className="h-1.5 flex-1" />
                          <span className="text-[10px] text-muted-foreground w-6 text-right">{task.progress}%</span>
                        </div>
                      )}
                      <div className={cn('hidden md:flex items-center gap-1 text-xs shrink-0', overdue ? 'text-rose-400' : soon ? 'text-amber-400' : 'text-muted-foreground')}>
                        <Calendar className="h-3.5 w-3.5" />
                        {formatDate(task.deadline, 'MMM d')}
                      </div>
                      {empName && (
                        <div className="hidden md:flex items-center gap-1.5 shrink-0">
                          <div className={cn('flex h-5 w-5 items-center justify-center rounded-full text-[9px] font-bold text-white bg-gradient-to-br', getAvatarGradient(empName))}>
                            {getInitials(empName)}
                          </div>
                          <span className="text-xs text-muted-foreground">{empName}</span>
                        </div>
                      )}
                      <StatusBadge status={task.status} />
                    </motion.div>
                  )
                })}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Create task dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Create Task</DialogTitle>
          </DialogHeader>
          <TaskForm
            employees={employees}
            defaultValues={{ status: defaultStatus }}
            onSubmit={async (data) => { await createMutation.mutateAsync(data) }}
            isLoading={createMutation.isPending}
          />
        </DialogContent>
      </Dialog>
    </AnimatedPage>
  )
}
