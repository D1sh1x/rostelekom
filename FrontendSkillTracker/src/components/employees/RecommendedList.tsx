import { motion } from 'framer-motion'
import { CheckCircle2, XCircle, Users } from 'lucide-react'
import { getAvatarGradient, getInitials, getMatchScoreConfig, cn } from '@/lib/utils'
import type { RecommendedEmployee } from '@/types'

interface RecommendedListProps {
  employees: RecommendedEmployee[]
  isLoading?: boolean
}

function MatchScore({ score }: { score: number }) {
  const cfg = getMatchScoreConfig(score)
  return (
    <div className={cn('flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold', cfg.bg, cfg.color)}>
      <div className="relative h-3 w-3">
        <svg viewBox="0 0 12 12" className="absolute inset-0 -rotate-90">
          <circle cx="6" cy="6" r="4.5" fill="none" stroke="currentColor" strokeOpacity={0.2} strokeWidth="2" />
          <circle
            cx="6" cy="6" r="4.5"
            fill="none" stroke="currentColor" strokeWidth="2"
            strokeDasharray={`${score * 0.283} 28.3`}
            strokeLinecap="round"
          />
        </svg>
      </div>
      {score}% {cfg.label}
    </div>
  )
}

export default function RecommendedList({ employees, isLoading }: RecommendedListProps) {
  if (isLoading) {
    return <div className="py-6 text-center text-sm text-muted-foreground">Loading recommendations...</div>
  }

  if (employees.length === 0) {
    return (
      <div className="flex flex-col items-center gap-2 py-8 text-center text-sm text-muted-foreground">
        <Users className="h-8 w-8 opacity-20" />
        <p>No required skills set. Add skills to get recommendations.</p>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-3">
      {employees.map((emp, i) => (
        <motion.div
          key={emp.id}
          initial={{ opacity: 0, x: -12 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: i * 0.06 }}
          className="rounded-xl border border-border bg-card/50 p-3.5"
        >
          <div className="flex items-start justify-between gap-3">
            <div className="flex items-center gap-2.5">
              <div className={cn(
                'flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-xs font-bold text-white bg-gradient-to-br',
                getAvatarGradient(emp.name)
              )}>
                {getInitials(emp.name)}
              </div>
              <div>
                <p className="font-medium text-sm leading-tight">{emp.name}</p>
                <p className="text-xs text-muted-foreground">@{emp.username}</p>
              </div>
            </div>
            <MatchScore score={emp.match_score} />
          </div>

          {(emp.matched_skills.length > 0 || emp.missing_skills.length > 0) && (
            <div className="mt-3 flex flex-col gap-1.5">
              {emp.matched_skills.length > 0 && (
                <div className="flex items-start gap-1.5">
                  <CheckCircle2 className="h-3.5 w-3.5 text-emerald-400 mt-0.5 shrink-0" />
                  <p className="text-xs text-emerald-400">{emp.matched_skills.join(', ')}</p>
                </div>
              )}
              {emp.missing_skills.length > 0 && (
                <div className="flex items-start gap-1.5">
                  <XCircle className="h-3.5 w-3.5 text-rose-400 mt-0.5 shrink-0" />
                  <p className="text-xs text-rose-400">{emp.missing_skills.join(', ')}</p>
                </div>
              )}
            </div>
          )}
        </motion.div>
      ))}
    </div>
  )
}
