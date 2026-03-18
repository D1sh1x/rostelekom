import { motion } from 'framer-motion'

interface ProgressRingProps {
  value: number
  size?: number
  strokeWidth?: number
  color?: string
  label?: string
}

export default function ProgressRing({ value, size = 80, strokeWidth = 6, color = '#7c3aed', label }: ProgressRingProps) {
  const r = (size - strokeWidth) / 2
  const circumference = 2 * Math.PI * r
  const offset = circumference - (value / 100) * circumference

  return (
    <div className="relative inline-flex items-center justify-center">
      <svg width={size} height={size} className="-rotate-90">
        <circle cx={size / 2} cy={size / 2} r={r} fill="none" stroke="currentColor" strokeWidth={strokeWidth} className="text-border" />
        <motion.circle
          cx={size / 2} cy={size / 2} r={r}
          fill="none" stroke={color} strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          initial={{ strokeDashoffset: circumference }}
          animate={{ strokeDashoffset: offset }}
          transition={{ duration: 1, ease: 'easeOut', delay: 0.2 }}
        />
      </svg>
      <div className="absolute flex flex-col items-center">
        <span className="text-lg font-bold">{value}%</span>
        {label && <span className="text-[10px] text-muted-foreground">{label}</span>}
      </div>
    </div>
  )
}
