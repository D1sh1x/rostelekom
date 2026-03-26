import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const schema = z.object({
  name: z.string().min(2, 'Минимум 2 символа'),
  description: z.string().optional(),
})
export type SkillFormData = z.infer<typeof schema>

interface SkillFormProps {
  onSubmit: (data: SkillFormData) => Promise<void>
  isLoading?: boolean
}

export default function SkillForm({ onSubmit, isLoading }: SkillFormProps) {
  const { register, handleSubmit, formState: { errors } } = useForm<SkillFormData>({
    resolver: zodResolver(schema),
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-1.5">
        <Label>Название навыка</Label>
        <Input placeholder="например, TypeScript" {...register('name')} />
        {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
      </div>
      <div className="space-y-1.5">
        <Label>Описание <span className="text-muted-foreground text-xs">(необязательно)</span></Label>
        <Input placeholder="Краткое описание..." {...register('description')} />
      </div>
      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? <><Loader2 className="h-4 w-4 animate-spin" />Создание...</> : 'Создать навык'}
      </Button>
    </form>
  )
}
