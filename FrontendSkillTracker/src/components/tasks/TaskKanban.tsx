import { motion } from 'framer-motion'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import TaskCard from './TaskCard'
import { getStatusConfig } from '@/lib/utils'
import type { Task, TaskStatus, User } from '@/types'

interface Column {
  status: TaskStatus
  label: string
}

const COLUMNS: Column[] = [
  { status: 'pending', label: 'Pending' },
  { status: 'in_progress', label: 'In Progress' },
  { status: 'completed', label: 'Completed' },
]

interface TaskKanbanProps {
  tasks: Task[]
  employeeMap: Record<number, string>
  isManager?: boolean
  onAddTask?: (status: TaskStatus) => void
}

export default function TaskKanban({ tasks, employeeMap, isManager, onAddTask }: TaskKanbanProps) {
  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
      {COLUMNS.map(col => {
        const colTasks = tasks.filter(t => t.status === col.status)
        const cfg = getStatusConfig(col.status)

        return (
          <div key={col.status} className="flex flex-col gap-3">
            {/* Column header */}
            <div className="flex items-center justify-between px-1">
              <div className="flex items-center gap-2">
                <span className={`h-2 w-2 rounded-full ${cfg.dot}`} />
                <span className="text-sm font-semibold">{col.label}</span>
                <span className={`rounded-full px-2 py-0.5 text-[11px] font-medium ${cfg.bg} ${cfg.color}`}>
                  {colTasks.length}
                </span>
              </div>
              {isManager && onAddTask && (
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6 rounded-md text-muted-foreground hover:text-foreground"
                  onClick={() => onAddTask(col.status)}
                >
                  <Plus className="h-4 w-4" />
                </Button>
              )}
            </div>

            {/* Cards */}
            <div className="flex flex-col gap-2 min-h-[120px]">
              {colTasks.length === 0 ? (
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="flex items-center justify-center rounded-2xl border border-dashed border-border/50 p-8 text-xs text-muted-foreground"
                >
                  No tasks
                </motion.div>
              ) : (
                colTasks.map((task, i) => (
                  <TaskCard
                    key={task.id}
                    task={task}
                    employeeName={employeeMap[task.employee_id]}
                    index={i}
                  />
                ))
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}
