import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { Plus, Trash2, Zap, Search } from 'lucide-react'
import { skillsApi } from '@/api/skills'
import AnimatedPage from '@/components/common/AnimatedPage'
import EmptyState from '@/components/common/EmptyState'
import SkillBadge from '@/components/skills/SkillBadge'
import SkillForm, { type SkillFormData } from '@/components/skills/SkillForm'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { getSkillColor } from '@/lib/utils'
import { toast } from '@/hooks/use-toast'

export default function SkillsPage() {
  const qc = useQueryClient()
  const [search, setSearch] = useState('')
  const [createOpen, setCreateOpen] = useState(false)

  const { data: skills = [] } = useQuery({
    queryKey: ['skills'],
    queryFn: () => skillsApi.list(),
    select: (data) => data ?? [],
  })

  const createMutation = useMutation({
    mutationFn: (data: SkillFormData) => skillsApi.create(data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['skills'] })
      toast({ title: 'Навык создан' })
      setCreateOpen(false)
    },
    onError: () => toast({ title: 'Ошибка создания навыка', variant: 'destructive' }),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => skillsApi.delete(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['skills'] })
      toast({ title: 'Навык удален' })
    },
    onError: () => toast({ title: 'Ошибка удаления навыка', variant: 'destructive' }),
  })

  const filtered = skills.filter(s =>
    s.name.toLowerCase().includes(search.toLowerCase()) ||
    s.description?.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <AnimatedPage>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold">Навыки</h1>
            <p className="text-sm text-muted-foreground mt-0.5">Всего навыков: {skills.length}</p>
          </div>
          <Button onClick={() => setCreateOpen(true)} className="gap-1.5" size="sm">
            <Plus className="h-4 w-4" />
            Новый навык
          </Button>
        </div>

        {/* Search */}
        <div className="relative max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
          <Input
            placeholder="Поиск навыков..."
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="pl-9 h-9"
          />
        </div>

        {/* Skills grid */}
        {filtered.length === 0 ? (
          <EmptyState
            icon={Zap}
            title="No skills yet"
            description="Create skills to assign to employees and tasks."
            action={<Button onClick={() => setCreateOpen(true)}>Create skill</Button>}
          />
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            <AnimatePresence>
              {filtered.map((skill, i) => (
                <motion.div
                  key={skill.id}
                  initial={{ opacity: 0, scale: 0.95 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.95 }}
                  transition={{ delay: i * 0.04 }}
                  className="group relative rounded-2xl border border-border bg-card p-4 hover:border-violet-500/30 hover:shadow-md transition-all duration-200"
                >
                  {/* Skill icon */}
                  <div className="mb-3 flex items-start justify-between gap-2">
                    <div className={`inline-flex items-center gap-1.5 rounded-xl px-3 py-1.5 text-sm font-medium border ${getSkillColor(skill.name)}`}>
                      <Zap className="h-3.5 w-3.5" />
                      {skill.name}
                    </div>
                    <button
                      onClick={() => {
                        if (confirm(`Удалить навык "${skill.name}"?`)) {
                          deleteMutation.mutate(skill.id)
                        }
                      }}
                      className="opacity-0 group-hover:opacity-100 transition-opacity p-1.5 rounded-lg text-muted-foreground hover:text-rose-400 hover:bg-rose-400/10"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </div>

                  {skill.description ? (
                    <p className="text-xs text-muted-foreground leading-relaxed line-clamp-2">{skill.description}</p>
                  ) : (
                    <p className="text-xs text-muted-foreground/50 italic">Нет описания</p>
                  )}

                  <p className="mt-3 text-[10px] text-muted-foreground/50">
                    ID: {skill.id}
                  </p>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        )}
      </div>

      {/* Create skill dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>Новый навык</DialogTitle>
          </DialogHeader>
          <SkillForm
            onSubmit={async data => { await createMutation.mutateAsync(data) }}
            isLoading={createMutation.isPending}
          />
        </DialogContent>
      </Dialog>
    </AnimatedPage>
  )
}
