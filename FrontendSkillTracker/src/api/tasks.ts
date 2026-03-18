import { api } from './client'
import type { Task, TaskRequest, TaskFilter, TaskHistory, Comment, Attachment, RecommendedEmployee, Skill } from '@/types'

export const tasksApi = {
  list: (filter?: TaskFilter) =>
    api.get<Task[]>('/tasks', { params: filter }).then((r) => r.data),

  myTasks: () =>
    api.get<Task[]>('/tasks/my').then((r) => r.data),

  get: (id: number) =>
    api.get<Task>(`/tasks/${id}`).then((r) => r.data),

  create: (data: TaskRequest) =>
    api.post<Task>('/tasks', data).then((r) => r.data),

  update: (id: number, data: Partial<TaskRequest>) =>
    api.put(`/tasks/${id}`, data),

  delete: (id: number) =>
    api.delete(`/tasks/${id}`),

  getHistory: (id: number) =>
    api.get<TaskHistory[]>(`/tasks/${id}/history`).then((r) => r.data),

  getSkills: (id: number) =>
    api.get<Skill[]>(`/tasks/${id}/skills`).then((r) => r.data),

  addSkill: (taskId: number, skillId: number) =>
    api.post(`/tasks/${taskId}/skills/${skillId}`),

  removeSkill: (taskId: number, skillId: number) =>
    api.delete(`/tasks/${taskId}/skills/${skillId}`),

  getRecommendedEmployees: (id: number) =>
    api.get<RecommendedEmployee[]>(`/tasks/${id}/recommended-employees`).then((r) => r.data),

  uploadAttachment: (id: number, file: File) => {
    const form = new FormData()
    form.append('file', file)
    return api.post<Attachment>(`/tasks/${id}/attachments`, form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }).then((r) => r.data)
  },
}

export const commentsApi = {
  list: (taskId: number) =>
    api.get<Comment[]>(`/tasks/${taskId}/comments`).then((r) => r.data),

  create: (taskId: number, text: string) =>
    api.post<Comment>('/comments', { task_id: taskId, text }).then((r) => r.data),

  update: (id: number, text: string) =>
    api.put(`/comments/${id}`, { text }),

  delete: (id: number) =>
    api.delete(`/comments/${id}`),
}
