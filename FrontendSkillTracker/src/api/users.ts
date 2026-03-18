import { api } from './client'
import type { User, UserRequest, UpdateUserRequest, Skill } from '@/types'

export const usersApi = {
  list: () =>
    api.get<User[]>('/users').then((r) => r.data),

  get: (id: number) =>
    api.get<User>(`/users/${id}`).then((r) => r.data),

  create: (data: UserRequest) =>
    api.post<User>('/users', data).then((r) => r.data),

  update: (id: number, data: UpdateUserRequest) =>
    api.put(`/users/${id}`, data),

  delete: (id: number) =>
    api.delete(`/users/${id}`),

  getSkills: (id: number) =>
    api.get<Skill[]>(`/users/${id}/skills`).then((r) => r.data),

  assignSkill: (userId: number, skillId: number) =>
    api.post(`/users/${userId}/skills/${skillId}`),

  removeSkill: (userId: number, skillId: number) =>
    api.delete(`/users/${userId}/skills/${skillId}`),
}
