import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import type { User } from '@/types'

const schema = z.object({
  title: z.string().min(3, 'Min 3 characters'),
  description: z.string().optional(),
  employee_id: z.coerce.number().min(1, 'Select an employee'),
  deadline: z.string().min(1, 'Deadline is required'),
  status: z.enum(['pending', 'in_progress', 'completed']).optional(),
  progress: z.coerce.number().min(0).max(100).optional(),
})
export type TaskFormData = z.infer<typeof schema>

interface TaskFormProps {
  employees: User[]
  onSubmit: (data: TaskFormData) => Promise<void>
  isLoading?: boolean
  defaultValues?: Partial<TaskFormData>
}

export default function TaskForm({ employees, onSubmit, isLoading, defaultValues }: TaskFormProps) {
  const { register, handleSubmit, setValue, watch, formState: { errors } } = useForm<TaskFormData>({
    resolver: zodResolver(schema),
    defaultValues: { progress: 0, status: 'pending', ...defaultValues },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-1.5">
        <Label>Title</Label>
        <Input placeholder="Task title" {...register('title')} />
        {errors.title && <p className="text-xs text-destructive">{errors.title.message}</p>}
      </div>

      <div className="space-y-1.5">
        <Label>Description</Label>
        <textarea
          placeholder="Task description..."
          {...register('description')}
          rows={3}
          className="flex w-full rounded-lg border border-border bg-background/50 px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none"
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <Label>Assign to</Label>
          <Select onValueChange={(v) => setValue('employee_id', Number(v))} defaultValue={defaultValues?.employee_id?.toString()}>
            <SelectTrigger>
              <SelectValue placeholder="Select employee" />
            </SelectTrigger>
            <SelectContent>
              {employees.map(e => (
                <SelectItem key={e.id} value={e.id.toString()}>{e.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.employee_id && <p className="text-xs text-destructive">{errors.employee_id.message}</p>}
        </div>

        <div className="space-y-1.5">
          <Label>Deadline</Label>
          <Input type="datetime-local" {...register('deadline')} />
          {errors.deadline && <p className="text-xs text-destructive">{errors.deadline.message}</p>}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <Label>Status</Label>
          <Select onValueChange={(v) => setValue('status', v as TaskFormData['status'])} defaultValue={defaultValues?.status || 'pending'}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="pending">Pending</SelectItem>
              <SelectItem value="in_progress">In Progress</SelectItem>
              <SelectItem value="completed">Completed</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-1.5">
          <Label>Progress ({watch('progress') ?? 0}%)</Label>
          <Input type="range" min={0} max={100} step={5} {...register('progress')} className="h-9 cursor-pointer" />
        </div>
      </div>

      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? <><Loader2 className="h-4 w-4 animate-spin" />Saving...</> : 'Save Task'}
      </Button>
    </form>
  )
}
