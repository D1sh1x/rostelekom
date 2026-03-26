import { motion } from 'framer-motion'
import { Shield, User as UserIcon } from 'lucide-react'
import { getAvatarGradient, getInitials, cn } from '@/lib/utils'
import SkillBadge from '@/components/skills/SkillBadge'
import type { User, Skill } from '@/types'

interface EmployeeCardProps {
  user: User
  skills?: Skill[]
  index?: number
  onManageSkills?: () => void
}

export default function EmployeeCard({ user, skills = [], index = 0, onManageSkills }: EmployeeCardProps) {
  const isManager = user.role === 'manager'

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05 }}
      className="group rounded-2xl border border-border bg-card p-5 hover:border-violet-500/30 hover:shadow-lg transition-all duration-200"
    >
      {/* Header */}
      <div className="flex items-start gap-3 mb-4">
        <div className={cn(
          'flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl text-base font-bold text-white bg-gradient-to-br shadow-md',
          getAvatarGradient(user.name)
        )}>
          {getInitials(user.name)}
        </div>
        <div className="min-w-0 flex-1">
          <h3 className="font-semibold leading-tight truncate">{user.name}</h3>
          <p className="text-sm text-muted-foreground truncate">@{user.username}</p>
          <div className={cn(
            'mt-1 inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-medium',
            isManager ? 'bg-violet-500/10 text-violet-400' : 'bg-sky-500/10 text-sky-400'
          )}>
            {isManager ? <Shield className="h-2.5 w-2.5" /> : <UserIcon className="h-2.5 w-2.5" />}
            {isManager ? 'Менеджер' : 'Сотрудник'}
          </div>
        </div>
      </div>

      {/* Skills */}
      <div className="min-h-[32px]">
        {skills.length > 0 ? (
          <div className="flex flex-wrap gap-1">
            {skills.slice(0, 5).map(s => <SkillBadge key={s.id} skill={s} size="sm" />)}
            {skills.length > 5 && (
              <span className="text-[10px] text-muted-foreground self-center">+{skills.length - 5}</span>
            )}
          </div>
        ) : (
          <p className="text-xs text-muted-foreground italic">Навыки не назначены</p>
        )}
      </div>

      {onManageSkills && (
        <button
          onClick={onManageSkills}
          className="mt-3 w-full rounded-lg border border-border/50 py-1.5 text-xs text-muted-foreground hover:border-violet-500/40 hover:text-violet-400 transition-colors"
        >
          Управление навыками
        </button>
      )}
    </motion.div>
  )
}
