import { useEffect, useState } from 'react'
import { Search, X } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import type { TaskFilter, User } from '@/types'

interface TaskFiltersProps {
  filters: TaskFilter
  onChange: (filters: TaskFilter) => void
  employees?: User[]
  showEmployeeFilter?: boolean
}

export default function TaskFilters({ filters, onChange, employees = [], showEmployeeFilter = false }: TaskFiltersProps) {
  const [search, setSearch] = useState(filters.search ?? '')

  useEffect(() => {
    const timer = setTimeout(() => {
      onChange({ ...filters, search: search || undefined })
    }, 300)
    return () => clearTimeout(timer)
  }, [search])

  const hasActiveFilters = filters.status || filters.employee_id || filters.search

  return (
    <div className="flex flex-wrap items-center gap-2">
      <div className="relative flex-1 min-w-[180px] max-w-xs">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
        <Input
          placeholder="Поиск задач..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="pl-9 h-9"
        />
      </div>

      <Select
        value={filters.status ?? 'all'}
        onValueChange={v => onChange({ ...filters, status: v === 'all' ? undefined : (v as TaskFilter['status']) })}
      >
        <SelectTrigger className="h-9 w-36">
          <SelectValue placeholder="Все статусы" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Все статусы</SelectItem>
          <SelectItem value="pending">Ожидает</SelectItem>
          <SelectItem value="in_progress">В процессе</SelectItem>
          <SelectItem value="completed">Завершено</SelectItem>
        </SelectContent>
      </Select>

      {showEmployeeFilter && employees.length > 0 && (
        <Select
          value={filters.employee_id?.toString() ?? 'all'}
          onValueChange={v => onChange({ ...filters, employee_id: v === 'all' ? undefined : Number(v) })}
        >
          <SelectTrigger className="h-9 w-40">
            <SelectValue placeholder="Все сотрудники" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Все сотрудники</SelectItem>
            {employees.map(e => (
              <SelectItem key={e.id} value={e.id.toString()}>{e.name}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      )}

      {hasActiveFilters && (
        <Button
          variant="ghost"
          size="sm"
          className="h-9 gap-1.5 text-muted-foreground"
          onClick={() => { onChange({}); setSearch('') }}
        >
          <X className="h-3.5 w-3.5" />
          Очистить
        </Button>
      )}
    </div>
  )
}
