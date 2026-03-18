import { NavLink, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import {
  LayoutDashboard, CheckSquare, Users, Zap, Sparkles,
  LogOut, ChevronLeft, ChevronRight, Moon, Sun
} from 'lucide-react'
import { cn, getInitials, getAvatarGradient } from '@/lib/utils'
import { useAuth } from '@/contexts/AuthContext'
import { useTheme } from '@/components/theme-provider'
import { useState } from 'react'
import { toast } from '@/hooks/use-toast'

interface NavItem {
  to: string
  icon: React.ElementType
  label: string
  managerOnly?: boolean
}

const navItems: NavItem[] = [
  { to: '/dashboard', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/tasks', icon: CheckSquare, label: 'Tasks' },
  { to: '/employees', icon: Users, label: 'Employees', managerOnly: true },
  { to: '/skills', icon: Zap, label: 'Skills', managerOnly: true },
]

export default function AppSidebar() {
  const { user, logout } = useAuth()
  const { theme, setTheme } = useTheme()
  const navigate = useNavigate()
  const [collapsed, setCollapsed] = useState(false)

  const handleLogout = async () => {
    await logout()
    toast({ title: 'Logged out', variant: 'default' })
    navigate('/login')
  }

  const items = navItems.filter(i => !i.managerOnly || user?.role === 'manager')

  return (
    <motion.aside
      animate={{ width: collapsed ? 68 : 220 }}
      transition={{ duration: 0.25, ease: [0.25, 0.46, 0.45, 0.94] }}
      className="relative flex h-screen flex-shrink-0 flex-col border-r border-sidebar-border bg-sidebar overflow-hidden"
    >
      {/* Logo */}
      <div className="flex h-16 items-center gap-3 px-4 border-b border-sidebar-border">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-gradient-violet shadow-md glow-sm">
          <Sparkles className="h-4 w-4 text-white" />
        </div>
        <AnimatePresence>
          {!collapsed && (
            <motion.span
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -10 }}
              transition={{ duration: 0.2 }}
              className="text-sm font-bold tracking-tight gradient-text whitespace-nowrap"
            >
              SkillTracker
            </motion.span>
          )}
        </AnimatePresence>
      </div>

      {/* Nav */}
      <nav className="flex-1 space-y-1 p-2 pt-4">
        {items.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) => cn(
              'group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-150',
              isActive
                ? 'bg-gradient-violet text-white shadow-md glow-sm'
                : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground'
            )}
            title={collapsed ? item.label : undefined}
          >
            {({ isActive }) => (
              <>
                <item.icon className={cn('h-4 w-4 shrink-0', isActive ? 'text-white' : '')} />
                <AnimatePresence>
                  {!collapsed && (
                    <motion.span
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      exit={{ opacity: 0 }}
                      transition={{ duration: 0.15 }}
                      className="whitespace-nowrap"
                    >
                      {item.label}
                    </motion.span>
                  )}
                </AnimatePresence>
              </>
            )}
          </NavLink>
        ))}
      </nav>

      {/* Bottom */}
      <div className="border-t border-sidebar-border p-2 space-y-1">
        {/* Theme toggle */}
        <button
          onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          className="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-all duration-150"
          title={collapsed ? 'Toggle theme' : undefined}
        >
          {theme === 'dark' ? <Sun className="h-4 w-4 shrink-0" /> : <Moon className="h-4 w-4 shrink-0" />}
          <AnimatePresence>
            {!collapsed && (
              <motion.span initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
                {theme === 'dark' ? 'Light mode' : 'Dark mode'}
              </motion.span>
            )}
          </AnimatePresence>
        </button>

        {/* User */}
        <div className={cn('flex items-center gap-3 rounded-xl px-3 py-2.5', collapsed ? 'justify-center' : '')}>
          <div className={cn('flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gradient-to-br text-white text-xs font-bold shadow', getAvatarGradient(user?.name || 'U'))}>
            {getInitials(user?.name || 'User')}
          </div>
          <AnimatePresence>
            {!collapsed && (
              <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }} className="min-w-0 flex-1">
                <p className="truncate text-xs font-medium text-sidebar-foreground">{user?.name}</p>
                <p className="truncate text-[10px] text-muted-foreground capitalize">{user?.role}</p>
              </motion.div>
            )}
          </AnimatePresence>
          <AnimatePresence>
            {!collapsed && (
              <motion.button
                initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                onClick={handleLogout}
                className="shrink-0 rounded-lg p-1.5 text-muted-foreground hover:text-rose-400 hover:bg-rose-400/10 transition-colors"
                title="Logout"
              >
                <LogOut className="h-3.5 w-3.5" />
              </motion.button>
            )}
          </AnimatePresence>
        </div>
      </div>

      {/* Collapse toggle */}
      <button
        onClick={() => setCollapsed(!collapsed)}
        className="absolute -right-3 top-20 flex h-6 w-6 items-center justify-center rounded-full border border-border bg-background shadow-md text-muted-foreground hover:text-foreground transition-colors z-10"
      >
        {collapsed ? <ChevronRight className="h-3 w-3" /> : <ChevronLeft className="h-3 w-3" />}
      </button>
    </motion.aside>
  )
}
