import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { authApi } from '@/api/auth'
import type { User, LoginRequest } from '@/types'

interface AuthState {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (data: LoginRequest) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthState | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const storedUser = localStorage.getItem('user')
    const token = localStorage.getItem('access_token')
    if (storedUser && token) {
      try {
        setUser(JSON.parse(storedUser))
      } catch {
        localStorage.clear()
      }
    }
    setIsLoading(false)
  }, [])

  const login = useCallback(async (data: LoginRequest) => {
    const res = await authApi.login(data)
    localStorage.setItem('access_token', res.access_token)
    localStorage.setItem('refresh_token', res.refresh_token)
    localStorage.setItem('user', JSON.stringify(res.user))
    setUser(res.user)
  }, [])

  const logout = useCallback(async () => {
    try {
      await authApi.logout()
    } catch {
      // ignore
    } finally {
      localStorage.clear()
      setUser(null)
    }
  }, [])

  return (
    <AuthContext.Provider value={{ user, isLoading, isAuthenticated: !!user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
