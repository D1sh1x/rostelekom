import { motion } from 'framer-motion'
import { useNavigate } from 'react-router-dom'
import { Calendar, Paperclip } from 'lucide-react'
import { cn, formatDate, getAvatarGradient, getInitials, isDeadlineOverdue, isDeadlineSoon } from '@/lib/utils'
import StatusBadge from '@/components/common/StatusBadge'
import SkillBadge from '@/components/skills/SkillBadge'
import { Progress } from '@/components/ui/progress'
import type { Task } from '@/types'

interface TaskCardProps {
  task: Task
  employeeName?: string
  index?: number
}

export default function TaskCard({ task, employeeName, index = 0 }: TaskCardProps) {
  const navigate = useNavigate()
  const overdue = task.status !== 'completed' && isDeadlineOverdue(task.deadline)
  const soon = !overdue && isDeadlineSoon(task.deadline)

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05, ease: [0.25, 0.46, 0.45, 0.94] }}
      whileHover={{ y: -3, transition: { duration: 0.15 } }}
      onClick={() => navigate(`/tasks/${task.id}`)}
      className="group cursor-pointer rounded-2xl border border-border bg-card p-4 shadow-sm hover:shadow-lg hover:border-violet-500/30 transition-all duration-200"
    >
      {/* Header */}
      <div className="flex items-start justify-between gap-2 mb-3">
        <h3 className="text-sm font-semibold leading-snug line-clamp-2 group-hover:text-primary transition-colors">
          {task.title}
        </h3>
        <StatusBadge status={task.status} />
      </div>

      {/* Description */}
      {task.description && (
        <p className="text-xs text-muted-foreground line-clamp-2 mb-3">{task.description}</p>
      )}

      {/* Progress */}
      {task.progress > 0 && (
        <div className="mb-3">
          <div className="flex justify-between text-[10px] text-muted-foreground mb-1">
            <span>Progress</span>
            <span>{task.progress}%</span>
          </div>
          <Progress value={task.progress} className="h-1.5" />
        </div>
      )}

      {/* Skills */}
      {task.required_skills?.length > 0 && (
        <div className="flex flex-wrap gap-1 mb-3">
          {task.required_skills.slice(0, 3).map(s => (
            <SkillBadge key={s.id} skill={s} size="sm" />
          ))}
          {task.required_skills.length > 3 && (
            <span className="text-[10px] text-muted-foreground">+{task.required_skills.length - 3}</span>
          )}
        </div>
      )}

      {/* Footer */}
      <div className="flex items-center justify-between mt-auto">
        <div className={cn('flex items-center gap-1 text-[11px]', overdue ? 'text-rose-400' : soon ? 'text-amber-400' : 'text-muted-foreground')}>
          <Calendar className="h-3 w-3" />
          {overdue ? '⚠ ' : soon ? '⏰ ' : ''}
          {formatDate(task.deadline, 'MMM d')}
        </div>
        {employeeName && (
          <div className="flex items-center gap-1.5">
            <div className={cn('flex h-5 w-5 items-center justify-center rounded-full text-[9px] font-bold text-white bg-gradient-to-br', getAvatarGradient(employeeName))}>
              {getInitials(employeeName)}
            </div>
            <span className="text-[10px] text-muted-foreground truncate max-w-[80px]">{employeeName}</span>
          </div>
        )}
      </div>
    </motion.div>
  )
}
