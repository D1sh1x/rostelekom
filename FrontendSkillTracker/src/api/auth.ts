import { api } from './client'
import type { LoginRequest, LoginResponse } from '@/types'

export const authApi = {
  login: (data: LoginRequest) =>
    api.post<LoginResponse>('/login', data).then((r) => r.data),

  logout: () => api.post('/logout'),

  refresh: (refreshToken: string) =>
    api.post<LoginResponse>('/refresh', { refresh_token: refreshToken }).then((r) => r.data),
}
