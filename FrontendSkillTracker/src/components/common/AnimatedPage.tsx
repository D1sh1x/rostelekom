import { motion } from 'framer-motion'
import type { ReactNode } from 'react'

const variants = {
  initial: { opacity: 0, y: 16 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -8 },
}

export default function AnimatedPage({ children }: { children: ReactNode }) {
  return (
    <motion.div
      variants={variants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.25, ease: [0.25, 0.46, 0.45, 0.94] }}
      className="h-full"
    >
      {children}
    </motion.div>
  )
}
