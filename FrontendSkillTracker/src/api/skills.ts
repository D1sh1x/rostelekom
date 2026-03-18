import { api } from './client'
import type { Skill, SkillRequest } from '@/types'

export const skillsApi = {
  list: () =>
    api.get<Skill[]>('/skills').then((r) => r.data),

  create: (data: SkillRequest) =>
    api.post<Skill>('/skills', data).then((r) => r.data),

  delete: (id: number) =>
    api.delete(`/skills/${id}`),
}
