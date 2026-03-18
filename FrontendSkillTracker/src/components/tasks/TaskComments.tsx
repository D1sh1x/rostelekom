import { useState, useRef, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { Send, MessageSquare } from 'lucide-react'
import { commentsApi } from '@/api/tasks'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/contexts/AuthContext'
import { formatRelative, getAvatarGradient, getInitials, cn } from '@/lib/utils'
import { toast } from '@/hooks/use-toast'

interface TaskCommentsProps {
  taskId: number
}

export default function TaskComments({ taskId }: TaskCommentsProps) {
  const { user } = useAuth()
  const qc = useQueryClient()
  const [text, setText] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)

  const { data: comments = [], isLoading } = useQuery({
    queryKey: ['comments', taskId],
    queryFn: () => commentsApi.list(taskId),
  })

  const addMutation = useMutation({
    mutationFn: (t: string) => commentsApi.create(taskId, t),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['comments', taskId] })
      setText('')
    },
    onError: () => toast({ title: 'Failed to send comment', variant: 'destructive' }),
  })

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [comments.length])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = text.trim()
    if (!trimmed) return
    addMutation.mutate(trimmed)
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <MessageSquare className="h-4 w-4 text-violet-400" />
        <h3 className="font-semibold">Comments</h3>
        <span className="rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground">{comments.length}</span>
      </div>

      {/* Comments list */}
      <div className="flex flex-col gap-3 max-h-80 overflow-y-auto pr-1 scrollbar-thin">
        {isLoading && (
          <div className="flex justify-center py-8 text-sm text-muted-foreground">Loading...</div>
        )}
        {!isLoading && comments.length === 0 && (
          <div className="flex flex-col items-center gap-2 py-8 text-center text-sm text-muted-foreground">
            <MessageSquare className="h-8 w-8 opacity-20" />
            <p>No comments yet. Be the first!</p>
          </div>
        )}
        <AnimatePresence initial={false}>
          {comments.map(comment => {
            const isOwn = comment.user_id === user?.id
            return (
              <motion.div
                key={comment.id}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                className={cn('flex gap-2.5', isOwn && 'flex-row-reverse')}
              >
                <div className={cn(
                  'flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-[10px] font-bold text-white bg-gradient-to-br mt-0.5',
                  getAvatarGradient(user?.name ?? 'U')
                )}>
                  {getInitials(user?.name ?? 'U')}
                </div>
                <div className={cn('flex flex-col gap-1 max-w-[75%]', isOwn && 'items-end')}>
                  <div className={cn(
                    'rounded-2xl px-3.5 py-2 text-sm leading-relaxed',
                    isOwn
                      ? 'bg-violet-500/20 text-foreground rounded-tr-sm'
                      : 'bg-muted text-foreground rounded-tl-sm'
                  )}>
                    {comment.text}
                  </div>
                  <span className="text-[10px] text-muted-foreground px-1">{formatRelative(comment.created_at)}</span>
                </div>
              </motion.div>
            )
          })}
        </AnimatePresence>
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSubmit} className="flex items-end gap-2">
        <textarea
          value={text}
          onChange={e => setText(e.target.value)}
          onKeyDown={e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSubmit(e as unknown as React.FormEvent) } }}
          placeholder="Write a comment... (Enter to send)"
          rows={2}
          className="flex-1 rounded-xl border border-border bg-background/50 px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none"
        />
        <Button
          type="submit"
          size="icon"
          className="h-10 w-10 rounded-xl bg-violet-600 hover:bg-violet-700 shrink-0"
          disabled={!text.trim() || addMutation.isPending}
        >
          <Send className="h-4 w-4" />
        </Button>
      </form>
    </div>
  )
}
