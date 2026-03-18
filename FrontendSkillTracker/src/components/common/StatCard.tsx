import { motion } from 'framer-motion'
import { cn } from '@/lib/utils'
import type { LucideIcon } from 'lucide-react'

interface StatCardProps {
  title: string
  value: string | number
  icon: LucideIcon
  trend?: { value: number; label: string }
  color?: 'violet' | 'sky' | 'emerald' | 'amber' | 'rose'
  delay?: number
}

const colorMap = {
  violet: { icon: 'text-violet-400', bg: 'bg-violet-400/10', glow: 'shadow-violet-500/10' },
  sky: { icon: 'text-sky-400', bg: 'bg-sky-400/10', glow: 'shadow-sky-500/10' },
  emerald: { icon: 'text-emerald-400', bg: 'bg-emerald-400/10', glow: 'shadow-emerald-500/10' },
  amber: { icon: 'text-amber-400', bg: 'bg-amber-400/10', glow: 'shadow-amber-500/10' },
  rose: { icon: 'text-rose-400', bg: 'bg-rose-400/10', glow: 'shadow-rose-500/10' },
}

export default function StatCard({ title, value, icon: Icon, trend, color = 'violet', delay = 0 }: StatCardProps) {
  const colors = colorMap[color]
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay, ease: [0.25, 0.46, 0.45, 0.94] }}
      whileHover={{ y: -3, transition: { duration: 0.2 } }}
      className={cn('rounded-2xl border border-border bg-card p-6 shadow-lg', colors.glow)}
    >
      <div className="flex items-start justify-between">
        <div>
          <p className="text-sm text-muted-foreground">{title}</p>
          <p className="mt-2 text-3xl font-bold tracking-tight">{value}</p>
          {trend && (
            <p className={cn('mt-1 text-xs', trend.value >= 0 ? 'text-emerald-400' : 'text-rose-400')}>
              {trend.value >= 0 ? '↑' : '↓'} {Math.abs(trend.value)}% {trend.label}
            </p>
          )}
        </div>
        <div className={cn('rounded-xl p-3', colors.bg)}>
          <Icon className={cn('h-5 w-5', colors.icon)} />
        </div>
      </div>
    </motion.div>
  )
}
