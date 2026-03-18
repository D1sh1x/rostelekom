import { useState, useCallback } from 'react'
import type { ToastProps } from '@/components/ui/toast'

type ToastItem = ToastProps & {
  id: string
  title?: React.ReactNode
  description?: React.ReactNode
  action?: React.ReactElement
}

let count = 0
const genId = () => `toast-${++count}`

const toastState = { toasts: [] as ToastItem[], listeners: new Set<() => void>() }
const notify = () => toastState.listeners.forEach((l) => l())

export function toast(props: Omit<ToastItem, 'id'>) {
  const id = genId()
  const item: ToastItem = { ...props, id, open: true, onOpenChange: (open) => { if (!open) dismiss(id) } }
  toastState.toasts = [item, ...toastState.toasts].slice(0, 5)
  notify()
  return id
}

function dismiss(id: string) {
  toastState.toasts = toastState.toasts.map((t) => t.id === id ? { ...t, open: false } : t)
  notify()
  setTimeout(() => {
    toastState.toasts = toastState.toasts.filter((t) => t.id !== id)
    notify()
  }, 300)
}

export function useToast() {
  const [toasts, setToasts] = useState<ToastItem[]>(toastState.toasts)
  const update = useCallback(() => setToasts([...toastState.toasts]), [])

  useState(() => {
    toastState.listeners.add(update)
    return () => { toastState.listeners.delete(update) }
  })

  return { toasts, toast, dismiss }
}
