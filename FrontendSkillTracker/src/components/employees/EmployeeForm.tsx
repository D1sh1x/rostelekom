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
  name: z.string().min(2, 'Min 2 characters'),
  username: z.string().min(3, 'Min 3 characters'),
  password: z.string().min(6, 'Min 6 characters').optional().or(z.literal('')),
  role: z.enum(['manager', 'employee']),
})
export type EmployeeFormData = z.infer<typeof schema>

interface EmployeeFormProps {
  defaultValues?: Partial<EmployeeFormData>
  onSubmit: (data: EmployeeFormData) => Promise<void>
  isLoading?: boolean
  isEdit?: boolean
}

export default function EmployeeForm({ defaultValues, onSubmit, isLoading, isEdit }: EmployeeFormProps) {
  const { register, handleSubmit, setValue, formState: { errors } } = useForm<EmployeeFormData>({
    resolver: zodResolver(schema),
    defaultValues: { role: 'employee', ...defaultValues },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <Label>Full name</Label>
          <Input placeholder="John Doe" {...register('name')} />
          {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label>Username</Label>
          <Input placeholder="johndoe" {...register('username')} />
          {errors.username && <p className="text-xs text-destructive">{errors.username.message}</p>}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <Label>Password {isEdit && <span className="text-muted-foreground text-xs">(leave blank to keep)</span>}</Label>
          <Input type="password" placeholder={isEdit ? '••••••' : 'Min 6 chars'} {...register('password')} />
          {errors.password && <p className="text-xs text-destructive">{errors.password.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label>Role</Label>
          <Select onValueChange={v => setValue('role', v as 'manager' | 'employee')} defaultValue={defaultValues?.role ?? 'employee'}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="employee">Employee</SelectItem>
              <SelectItem value="manager">Manager</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? <><Loader2 className="h-4 w-4 animate-spin" />{isEdit ? 'Saving...' : 'Creating...'}</> : isEdit ? 'Save Changes' : 'Create Employee'}
      </Button>
    </form>
  )
}
