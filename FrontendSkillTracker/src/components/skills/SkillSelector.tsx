import { useState } from 'react'
import { Plus, Search, Check } from 'lucide-react'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import SkillBadge from './SkillBadge'
import type { Skill } from '@/types'

interface SkillSelectorProps {
  allSkills: Skill[]
  selectedSkills: Skill[]
  onAdd: (skillId: number) => void
  onRemove: (skillId: number) => void
  readonly?: boolean
}

export default function SkillSelector({ allSkills, selectedSkills, onAdd, onRemove, readonly }: SkillSelectorProps) {
  const [open, setOpen] = useState(false)
  const [search, setSearch] = useState('')

  const selectedIds = new Set(selectedSkills.map(s => s.id))
  const filtered = allSkills.filter(s =>
    s.name.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="flex flex-wrap gap-1.5 items-center">
      {selectedSkills.map(skill => (
        <SkillBadge
          key={skill.id}
          skill={skill}
          onRemove={readonly ? undefined : () => onRemove(skill.id)}
        />
      ))}

      {!readonly && (
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm" className="h-6 gap-1 rounded-full border-dashed text-xs text-muted-foreground">
              <Plus className="h-3 w-3" />
              Add skill
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-64 p-3" align="start">
            <div className="relative flex-1">
              <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
              <Input
                placeholder="Поиск навыков..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                className="pl-8 h-8 text-xs"
              />
            </div>
            <div className="flex flex-col gap-0.5 max-h-48 overflow-y-auto">
              {filtered.length === 0 && (
                <p className="py-3 text-center text-xs text-muted-foreground">No skills found</p>
              )}
              {filtered.map(skill => (
                <button
                  key={skill.id}
                  onClick={() => {
                    if (selectedIds.has(skill.id)) onRemove(skill.id)
                    else onAdd(skill.id)
                  }}
                  className="flex items-center justify-between rounded-lg px-2 py-1.5 text-sm hover:bg-muted transition-colors text-left"
                >
                  <span className="truncate">{skill.name}</span>
                  {selectedIds.has(skill.id) && <Check className="h-3.5 w-3.5 text-emerald-400 shrink-0" />}
                </button>
              ))}
            </div>
          </PopoverContent>
        </Popover>
      )}
    </div>
  )
}
