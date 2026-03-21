import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Clock, ArrowRight, History } from 'lucide-react'
import { tasksApi } from '@/api/tasks'
import { getStatusConfig, formatDateTime, cn } from '@/lib/utils'
import type { TaskStatus } from '@/types'

interface TaskHistoryProps {
  taskId: number
}

function StatusDot({ status }: { status: TaskStatus }) {
  const cfg = getStatusConfig(status)
  return <span className={cn('h-2.5 w-2.5 rounded-full shrink-0', cfg.dot)} />
}

export default function TaskHistory({ taskId }: TaskHistoryProps) {
  const { data: history = [], isLoading } = useQuery({
    queryKey: ['task-history', taskId],
    queryFn: () => tasksApi.getHistory(taskId),
    select: (data) => data ?? [],
  })

  if (isLoading) return <div className="py-4 text-sm text-muted-foreground text-center">Loading...</div>
  if (history.length === 0) return (
    <div className="flex flex-col items-center gap-2 py-8 text-center text-sm text-muted-foreground">
      <History className="h-8 w-8 opacity-20" />
      <p>No status changes recorded yet.</p>
    </div>
  )

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center gap-2 mb-3">
        <History className="h-4 w-4 text-violet-400" />
        <h3 className="font-semibold">Status History</h3>
      </div>
      <div className="relative pl-4">
        {/* Vertical line */}
        <div className="absolute left-[7px] top-2 bottom-2 w-px bg-border" />

        {history.map((entry, i) => {
          const from = getStatusConfig(entry.old_status)
          const to = getStatusConfig(entry.new_status)
          return (
            <motion.div
              key={entry.id}
              initial={{ opacity: 0, x: -8 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.06 }}
              className="relative flex items-start gap-3 pb-5 last:pb-0"
            >
              {/* Timeline dot */}
              <div className="absolute left-[-9px] top-0.5 flex h-[14px] w-[14px] items-center justify-center rounded-full border-2 border-border bg-background">
                <span className={cn('h-1.5 w-1.5 rounded-full', to.dot)} />
              </div>

              <div className="min-w-0 flex-1 pt-0.5 pl-2">
                <div className="flex flex-wrap items-center gap-1.5 text-sm">
                  <span className={cn('font-medium', from.color)}>{from.label}</span>
                  <ArrowRight className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                  <span className={cn('font-semibold', to.color)}>{to.label}</span>
                </div>
                <div className="mt-0.5 flex items-center gap-1 text-[11px] text-muted-foreground">
                  <Clock className="h-3 w-3" />
                  {formatDateTime(entry.created_at)}
                </div>
              </div>
            </motion.div>
          )
        })}
      </div>
    </div>
  )
}
