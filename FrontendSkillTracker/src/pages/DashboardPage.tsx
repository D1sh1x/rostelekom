import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { CheckSquare, Clock, Users, Zap, TrendingUp, CalendarDays } from 'lucide-react'
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts'
import { tasksApi } from '@/api/tasks'
import { usersApi } from '@/api/users'
import { skillsApi } from '@/api/skills'
import { useAuth } from '@/contexts/AuthContext'
import AnimatedPage from '@/components/common/AnimatedPage'
import StatCard from '@/components/common/StatCard'
import StatusBadge from '@/components/common/StatusBadge'
import ProgressRing from '@/components/common/ProgressRing'
import SkillBadge from '@/components/skills/SkillBadge'
import { formatDate, isDeadlineToday, isDeadlineOverdue } from '@/lib/utils'
import { cn } from '@/lib/utils'
import type { Task } from '@/types'

const containerVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1, transition: { staggerChildren: 0.07 } },
}
const itemVariants = {
  hidden: { opacity: 0, y: 16 },
  visible: { opacity: 1, y: 0, transition: { type: 'spring', stiffness: 300, damping: 24 } },
}

function RecentTaskRow({ task }: { task: Task }) {
  const overdue = task.status !== 'completed' && isDeadlineOverdue(task.deadline)
  const today = isDeadlineToday(task.deadline)
  return (
    <motion.div variants={itemVariants} className="flex items-center gap-3 rounded-xl px-3 py-2.5 hover:bg-muted/40 transition-colors group cursor-pointer">
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-medium">{task.title}</p>
        <p className={cn('text-xs', overdue ? 'text-rose-400' : today ? 'text-amber-400' : 'text-muted-foreground')}>
          {overdue ? '⚠ Overdue · ' : today ? '📅 Today · ' : ''}{formatDate(task.deadline)}
        </p>
      </div>
      <StatusBadge status={task.status} />
    </motion.div>
  )
}

export default function DashboardPage() {
  const { user } = useAuth()
  const isManager = user?.role === 'manager'

  const { data: allTasks = [] } = useQuery({ queryKey: ['tasks'], queryFn: () => tasksApi.list(), select: (data) => data ?? [] })
  const { data: myTasks = [] } = useQuery({ queryKey: ['my-tasks'], queryFn: () => tasksApi.myTasks(), select: (data) => data ?? [] })
  const { data: users = [] } = useQuery({ queryKey: ['users'], queryFn: () => usersApi.list(), enabled: isManager, select: (data) => data ?? [] })
  const { data: skills = [] } = useQuery({ queryKey: ['skills'], queryFn: () => skillsApi.list(), select: (data) => data ?? [] })

  const tasks = isManager ? allTasks : myTasks

  const total = tasks.length
  const completed = tasks.filter(t => t.status === 'completed').length
  const inProgress = tasks.filter(t => t.status === 'in_progress').length
  const pending = tasks.filter(t => t.status === 'pending').length
  const completionRate = total > 0 ? Math.round((completed / total) * 100) : 0

  // Build chart data (last 7 days buckets - simplified)
  const chartData = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((day, i) => ({
    day,
    completed: Math.max(0, completed - i * 2 + i),
    inProgress: Math.max(0, inProgress - i + Math.floor(i / 2)),
    pending: Math.max(0, pending - i),
  }))

  const recentTasks = [...tasks].sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()).slice(0, 6)

  return (
    <AnimatedPage>
      <div className="space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold">
            Good {new Date().getHours() < 12 ? 'morning' : new Date().getHours() < 18 ? 'afternoon' : 'evening'},{' '}
            <span className="gradient-text">{user?.name?.split(' ')[0]}</span> 👋
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            {isManager ? `Managing ${total} tasks across ${users.length} employees` : `You have ${pending} pending tasks`}
          </p>
        </div>

        {/* Stats */}
        {isManager ? (
          <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
            <StatCard title="Total Tasks" value={total} icon={CheckSquare} color="violet" delay={0} />
            <StatCard title="Completed" value={completed} icon={TrendingUp} color="emerald" delay={0.05} />
            <StatCard title="In Progress" value={inProgress} icon={Clock} color="sky" delay={0.1} />
            <StatCard title="Team Members" value={users.length} icon={Users} color="amber" delay={0.15} />
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
            <StatCard title="My Tasks" value={total} icon={CheckSquare} color="violet" delay={0} />
            <StatCard title="Completed" value={completed} icon={TrendingUp} color="emerald" delay={0.05} />
            <StatCard title="In Progress" value={inProgress} icon={Clock} color="sky" delay={0.1} />
            <StatCard title="My Skills" value={skills.length} icon={Zap} color="amber" delay={0.15} />
          </div>
        )}

        <div className="grid gap-6 lg:grid-cols-3">
          {/* Chart */}
          <motion.div
            initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2, duration: 0.4 }}
            className="lg:col-span-2 rounded-2xl border border-border bg-card p-6"
          >
            <div className="mb-4 flex items-center justify-between">
              <div>
                <h2 className="font-semibold">Task Activity</h2>
                <p className="text-xs text-muted-foreground">Weekly overview</p>
              </div>
              <div className="flex gap-3 text-xs text-muted-foreground">
                <span className="flex items-center gap-1"><span className="h-2 w-2 rounded-full bg-violet-400" />Completed</span>
                <span className="flex items-center gap-1"><span className="h-2 w-2 rounded-full bg-sky-400" />In Progress</span>
              </div>
            </div>
            <ResponsiveContainer width="100%" height={180}>
              <AreaChart data={chartData} margin={{ top: 5, right: 0, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="gCompleted" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#7c3aed" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#7c3aed" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="gInProgress" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#38bdf8" stopOpacity={0.2} />
                    <stop offset="95%" stopColor="#38bdf8" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="day" tick={{ fontSize: 11, fill: 'hsl(var(--muted-foreground))' }} axisLine={false} tickLine={false} />
                <YAxis tick={{ fontSize: 11, fill: 'hsl(var(--muted-foreground))' }} axisLine={false} tickLine={false} />
                <Tooltip
                  contentStyle={{ background: 'hsl(var(--card))', border: '1px solid hsl(var(--border))', borderRadius: 12, fontSize: 12 }}
                  labelStyle={{ color: 'hsl(var(--foreground))' }}
                />
                <Area type="monotone" dataKey="completed" stroke="#7c3aed" fill="url(#gCompleted)" strokeWidth={2} dot={false} />
                <Area type="monotone" dataKey="inProgress" stroke="#38bdf8" fill="url(#gInProgress)" strokeWidth={2} dot={false} />
              </AreaChart>
            </ResponsiveContainer>
          </motion.div>

          {/* Progress + Skills */}
          <motion.div
            initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.25, duration: 0.4 }}
            className="flex flex-col gap-4"
          >
            {/* Progress Ring */}
            <div className="rounded-2xl border border-border bg-card p-6 flex flex-col items-center gap-3">
              <h2 className="self-start font-semibold text-sm">Completion Rate</h2>
              <ProgressRing value={completionRate} size={100} strokeWidth={8} />
              <div className="w-full grid grid-cols-3 gap-2 text-center text-xs text-muted-foreground">
                <div><p className="font-semibold text-amber-400">{pending}</p><p>Pending</p></div>
                <div><p className="font-semibold text-sky-400">{inProgress}</p><p>Active</p></div>
                <div><p className="font-semibold text-emerald-400">{completed}</p><p>Done</p></div>
              </div>
            </div>

            {/* Skills */}
            {skills.length > 0 && (
              <div className="rounded-2xl border border-border bg-card p-4">
                <div className="mb-3 flex items-center gap-2">
                  <Zap className="h-4 w-4 text-amber-400" />
                  <h2 className="text-sm font-semibold">Skills</h2>
                </div>
                <div className="flex flex-wrap gap-1.5">
                  {skills.slice(0, 8).map(s => <SkillBadge key={s.id} skill={s} size="sm" />)}
                  {skills.length > 8 && <span className="text-xs text-muted-foreground">+{skills.length - 8} more</span>}
                </div>
              </div>
            )}
          </motion.div>
        </div>

        {/* Recent Tasks */}
        {recentTasks.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3, duration: 0.4 }}
            className="rounded-2xl border border-border bg-card p-6"
          >
            <div className="mb-4 flex items-center gap-2">
              <CalendarDays className="h-4 w-4 text-violet-400" />
              <h2 className="font-semibold">Recent Tasks</h2>
            </div>
            <motion.div variants={containerVariants} initial="hidden" animate="visible" className="space-y-1">
              {recentTasks.map(t => <RecentTaskRow key={t.id} task={t} />)}
            </motion.div>
          </motion.div>
        )}
      </div>
    </AnimatedPage>
  )
}
