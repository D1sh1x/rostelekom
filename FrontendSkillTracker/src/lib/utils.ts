import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { format, formatDistanceToNow, isPast, isToday } from 'date-fns'
import type { TaskStatus } from '@/types'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string | Date, fmt = 'MMM d, yyyy') {
  return format(new Date(date), fmt)
}

export function formatRelative(date: string | Date) {
  return formatDistanceToNow(new Date(date), { addSuffix: true })
}

export function formatDateTime(date: string | Date) {
  return format(new Date(date), 'MMM d, yyyy · HH:mm')
}

export function isDeadlineSoon(deadline: string) {
  const d = new Date(deadline)
  const diff = d.getTime() - Date.now()
  return diff > 0 && diff < 1000 * 60 * 60 * 24 * 2 // within 2 days
}

export function isDeadlineOverdue(deadline: string) {
  return isPast(new Date(deadline))
}

export function isDeadlineToday(deadline: string) {
  return isToday(new Date(deadline))
}

export function getStatusConfig(status: TaskStatus) {
  const map = {
    pending: {
      label: 'Pending',
      color: 'text-amber-400',
      bg: 'bg-amber-400/10',
      border: 'border-amber-400/20',
      dot: 'bg-amber-400',
    },
    in_progress: {
      label: 'In Progress',
      color: 'text-sky-400',
      bg: 'bg-sky-400/10',
      border: 'border-sky-400/20',
      dot: 'bg-sky-400',
    },
    completed: {
      label: 'Completed',
      color: 'text-emerald-400',
      bg: 'bg-emerald-400/10',
      border: 'border-emerald-400/20',
      dot: 'bg-emerald-400',
    },
  }
  return map[status] ?? map.pending
}

export function getMatchScoreConfig(score: number) {
  if (score >= 80) return { color: 'text-emerald-400', bg: 'bg-emerald-400/10', label: 'Great match' }
  if (score >= 50) return { color: 'text-amber-400', bg: 'bg-amber-400/10', label: 'Partial match' }
  return { color: 'text-rose-400', bg: 'bg-rose-400/10', label: 'Low match' }
}

// Generate consistent color for a skill based on its name
export function getSkillColor(name: string): string {
  const colors = [
    'bg-violet-500/15 text-violet-300 border-violet-500/20',
    'bg-sky-500/15 text-sky-300 border-sky-500/20',
    'bg-emerald-500/15 text-emerald-300 border-emerald-500/20',
    'bg-rose-500/15 text-rose-300 border-rose-500/20',
    'bg-amber-500/15 text-amber-300 border-amber-500/20',
    'bg-pink-500/15 text-pink-300 border-pink-500/20',
    'bg-cyan-500/15 text-cyan-300 border-cyan-500/20',
    'bg-indigo-500/15 text-indigo-300 border-indigo-500/20',
    'bg-orange-500/15 text-orange-300 border-orange-500/20',
    'bg-teal-500/15 text-teal-300 border-teal-500/20',
  ]
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

export function getInitials(name: string): string {
  return name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

export function getAvatarGradient(name: string): string {
  const gradients = [
    'from-violet-500 to-purple-600',
    'from-sky-500 to-blue-600',
    'from-emerald-500 to-teal-600',
    'from-rose-500 to-pink-600',
    'from-amber-500 to-orange-600',
    'from-indigo-500 to-blue-700',
    'from-cyan-500 to-sky-600',
  ]
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return gradients[Math.abs(hash) % gradients.length]
}

export function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
