import { cn, getSkillColor } from '@/lib/utils'
import type { Skill } from '@/types'

interface SkillBadgeProps {
  skill: Skill
  onRemove?: () => void
  size?: 'sm' | 'md'
}

export default function SkillBadge({ skill, onRemove, size = 'md' }: SkillBadgeProps) {
  const colors = getSkillColor(skill.name)
  return (
    <span className={cn(
      'inline-flex items-center gap-1 rounded-full border font-mono font-medium',
      colors,
      size === 'sm' ? 'px-2 py-0.5 text-[10px]' : 'px-2.5 py-1 text-xs'
    )}>
      {skill.name}
      {onRemove && (
        <button
          onClick={(e) => { e.stopPropagation(); onRemove() }}
          className="ml-0.5 rounded-full hover:opacity-70 transition-opacity"
        >
          ×
        </button>
      )}
    </span>
  )
}
