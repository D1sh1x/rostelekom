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
  name: z.string().min(2, 'Минимум 2 символа'),
  username: z.string().min(3, 'Минимум 3 символа'),
  password: z.string().min(6, 'Минимум 6 символов').optional().or(z.literal('')),
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
          <Label>ФИО</Label>
          <Input placeholder="Иван Иванов" {...register('name')} />
          {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label>Логин</Label>
          <Input placeholder="ivanov" {...register('username')} />
          {errors.username && <p className="text-xs text-destructive">{errors.username.message}</p>}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1.5">
          <Label>Пароль {isEdit && <span className="text-muted-foreground text-xs">(оставьте пустым, чтобы не менять)</span>}</Label>
          <Input type="password" placeholder={isEdit ? '••••••' : 'Мин. 6 симв.'} {...register('password')} />
          {errors.password && <p className="text-xs text-destructive">{errors.password.message}</p>}
        </div>
        <div className="space-y-1.5">
          <Label>Роль</Label>
          <Select onValueChange={v => setValue('role', v as 'manager' | 'employee')} defaultValue={defaultValues?.role ?? 'employee'}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="employee">Сотрудник</SelectItem>
              <SelectItem value="manager">Менеджер</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? <><Loader2 className="h-4 w-4 animate-spin" />{isEdit ? 'Сохранение...' : 'Создание...'}</> : isEdit ? 'Сохранить изменения' : 'Создать сотрудника'}
      </Button>
    </form>
  )
}
