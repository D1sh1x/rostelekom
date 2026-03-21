import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { Plus, Search, Trash2, Users } from 'lucide-react'
import { usersApi } from '@/api/users'
import { skillsApi } from '@/api/skills'
import AnimatedPage from '@/components/common/AnimatedPage'
import EmptyState from '@/components/common/EmptyState'
import EmployeeCard from '@/components/employees/EmployeeCard'
import EmployeeForm, { type EmployeeFormData } from '@/components/employees/EmployeeForm'
import SkillSelector from '@/components/skills/SkillSelector'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { toast } from '@/hooks/use-toast'
import type { User, Skill } from '@/types'

export default function EmployeesPage() {
  const qc = useQueryClient()
  const [search, setSearch] = useState('')
  const [createOpen, setCreateOpen] = useState(false)
  const [skillsDialogUser, setSkillsDialogUser] = useState<User | null>(null)
  const [userSkillsMap, setUserSkillsMap] = useState<Record<number, Skill[]>>({})

  const { data: employees = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => usersApi.list(),
    select: (data) => data ?? [],
  })
  const { data: allSkills = [] } = useQuery({
    queryKey: ['skills'],
    queryFn: () => skillsApi.list(),
    select: (data) => data ?? [],
  })

  const createMutation = useMutation({
    mutationFn: (data: EmployeeFormData) => usersApi.create({
      name: data.name,
      username: data.username,
      password: data.password || '',
      role: data.role,
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['users'] })
      toast({ title: 'Employee created' })
      setCreateOpen(false)
    },
    onError: () => toast({ title: 'Failed to create employee', variant: 'destructive' }),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => usersApi.delete(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['users'] })
      toast({ title: 'Employee removed' })
    },
    onError: () => toast({ title: 'Failed to delete employee', variant: 'destructive' }),
  })

  const assignSkillMutation = useMutation({
    mutationFn: ({ userId, skillId }: { userId: number; skillId: number }) =>
      usersApi.assignSkill(userId, skillId),
    onSuccess: (_, { userId }) => loadUserSkills(userId),
  })

  const removeSkillMutation = useMutation({
    mutationFn: ({ userId, skillId }: { userId: number; skillId: number }) =>
      usersApi.removeSkill(userId, skillId),
    onSuccess: (_, { userId }) => loadUserSkills(userId),
  })

  const loadUserSkills = async (userId: number) => {
    const skills = await usersApi.getSkills(userId)
    setUserSkillsMap(prev => ({ ...prev, [userId]: skills }))
  }

  const openSkillsDialog = async (u: User) => {
    setSkillsDialogUser(u)
    if (!userSkillsMap[u.id]) {
      await loadUserSkills(u.id)
    }
  }

  const filtered = employees.filter(e =>
    e.name.toLowerCase().includes(search.toLowerCase()) ||
    e.username.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <AnimatedPage>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold">Team</h1>
            <p className="text-sm text-muted-foreground mt-0.5">{employees.length} member{employees.length !== 1 ? 's' : ''}</p>
          </div>
          <Button onClick={() => setCreateOpen(true)} className="gap-1.5" size="sm">
            <Plus className="h-4 w-4" />
            Add Employee
          </Button>
        </div>

        {/* Search */}
        <div className="relative max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
          <Input
            placeholder="Search by name or username..."
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="pl-9 h-9"
          />
        </div>

        {/* Grid */}
        {filtered.length === 0 ? (
          <EmptyState
            icon={Users}
            title="No employees found"
            description="Add team members to get started."
            action={<Button onClick={() => setCreateOpen(true)}>Add employee</Button>}
          />
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <AnimatePresence>
              {filtered.map((emp, i) => (
                <div key={emp.id} className="relative group/card">
                  <EmployeeCard
                    user={emp}
                    skills={userSkillsMap[emp.id]}
                    index={i}
                    onManageSkills={() => openSkillsDialog(emp)}
                  />
                  <button
                    onClick={() => {
                      if (confirm(`Remove ${emp.name}?`)) deleteMutation.mutate(emp.id)
                    }}
                    className="absolute top-3 right-3 opacity-0 group-hover/card:opacity-100 transition-opacity rounded-lg p-1.5 text-muted-foreground hover:text-rose-400 hover:bg-rose-400/10"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </button>
                </div>
              ))}
            </AnimatePresence>
          </div>
        )}
      </div>

      {/* Create employee dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Add Employee</DialogTitle>
          </DialogHeader>
          <EmployeeForm
            onSubmit={async data => { await createMutation.mutateAsync(data) }}
            isLoading={createMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Manage skills dialog */}
      <Dialog open={!!skillsDialogUser} onOpenChange={open => { if (!open) setSkillsDialogUser(null) }}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>
              Skills — {skillsDialogUser?.name}
            </DialogTitle>
          </DialogHeader>
          {skillsDialogUser && (
            <SkillSelector
              allSkills={allSkills}
              selectedSkills={userSkillsMap[skillsDialogUser.id] ?? []}
              onAdd={skillId => assignSkillMutation.mutate({ userId: skillsDialogUser.id, skillId })}
              onRemove={skillId => removeSkillMutation.mutate({ userId: skillsDialogUser.id, skillId })}
            />
          )}
        </DialogContent>
      </Dialog>
    </AnimatedPage>
  )
}
