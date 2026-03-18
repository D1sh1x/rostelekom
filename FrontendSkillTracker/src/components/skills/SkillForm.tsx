import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const schema = z.object({
  name: z.string().min(2, 'Min 2 characters'),
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
        <Label>Skill name</Label>
        <Input placeholder="e.g. TypeScript" {...register('name')} />
        {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
      </div>
      <div className="space-y-1.5">
        <Label>Description <span className="text-muted-foreground text-xs">(optional)</span></Label>
        <Input placeholder="Brief description..." {...register('description')} />
      </div>
      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? <><Loader2 className="h-4 w-4 animate-spin" />Creating...</> : 'Create Skill'}
      </Button>
    </form>
  )
}
