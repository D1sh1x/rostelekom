import { cn, getStatusConfig } from '@/lib/utils'
import type { TaskStatus } from '@/types'

export default function StatusBadge({ status }: { status: TaskStatus }) {
  const cfg = getStatusConfig(status)
  return (
    <span className={cn('inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium', cfg.color, cfg.bg, cfg.border)}>
      <span className={cn('h-1.5 w-1.5 rounded-full', cfg.dot)} />
      {cfg.label}
    </span>
  )
}
