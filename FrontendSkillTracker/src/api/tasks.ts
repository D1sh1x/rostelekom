import api from './client';
import type { Task, TaskFilters, Comment, TaskStatusHistory } from '../types';

export const taskApi = {
  getTasks: async (filters?: TaskFilters): Promise<Task[]> => {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== '') {
          params.append(key, value.toString());
        }
      });
    }
    const response = await api.get<Task[]>('/tasks', { params });
    return response.data;
  },

  getTaskById: async (id: number): Promise<Task> => {
    const response = await api.get<Task>(`/tasks/${id}`);
    return response.data;
  },

  updateTask: async (id: number, data: Partial<Task>): Promise<Task> => {
    const response = await api.put<Task>(`/tasks/${id}`, data);
    return response.data;
  },

  createTask: async (data: Omit<Task, 'id' | 'createdAt' | 'updatedAt' | 'progress'>): Promise<Task> => {
    const response = await api.post<Task>('/tasks', data);
    return response.data;
  },

  deleteTask: async (id: number): Promise<void> => {
    await api.delete(`/tasks/${id}`);
  },

  getComments: async (taskId: number): Promise<Comment[]> => {
    const response = await api.get<Comment[]>(`/tasks/${taskId}/comments`);
    return response.data;
  },

  addComment: async (taskId: number, text: string): Promise<Comment> => {
    const response = await api.post<Comment>('/comments', { taskId, text });
    return response.data;
  },

  getHistory: async (taskId: number): Promise<TaskStatusHistory[]> => {
    const response = await api.get<TaskStatusHistory[]>(`/tasks/${taskId}/history`);
    return response.data;
  }
};
