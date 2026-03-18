import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AnimatePresence } from 'framer-motion'
import { AuthProvider, useAuth } from '@/contexts/AuthContext'
import { ThemeProvider } from '@/components/theme-provider'
import { Toaster } from '@/components/ui/toaster'
import PageLayout from '@/components/layout/PageLayout'
import LoginPage from '@/pages/LoginPage'
import DashboardPage from '@/pages/DashboardPage'
import TasksPage from '@/pages/TasksPage'
import TaskDetailPage from '@/pages/TaskDetailPage'
import EmployeesPage from '@/pages/EmployeesPage'
import SkillsPage from '@/pages/SkillsPage'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth()
  if (isLoading) return (
    <div className="flex h-screen items-center justify-center bg-background">
      <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
    </div>
  )
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function ManagerRoute({ children }: { children: React.ReactNode }) {
  const { user } = useAuth()
  if (user?.role !== 'manager') return <Navigate to="/dashboard" replace />
  return <>{children}</>
}

function AppRoutes() {
  const { isAuthenticated } = useAuth()

  return (
    <AnimatePresence mode="wait">
      <Routes>
        <Route
          path="/login"
          element={isAuthenticated ? <Navigate to="/dashboard" replace /> : <LoginPage />}
        />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <PageLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="tasks" element={<TasksPage />} />
          <Route path="tasks/:id" element={<TaskDetailPage />} />
          <Route
            path="employees"
            element={<ManagerRoute><EmployeesPage /></ManagerRoute>}
          />
          <Route
            path="skills"
            element={<ManagerRoute><SkillsPage /></ManagerRoute>}
          />
        </Route>
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </AnimatePresence>
  )
}

export default function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="skilltracker-theme">
      <AuthProvider>
        <BrowserRouter>
          <AppRoutes />
          <Toaster />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  )
}
