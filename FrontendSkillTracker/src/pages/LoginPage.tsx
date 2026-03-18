import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Eye, EyeOff, Sparkles, Loader2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from '@/hooks/use-toast'

const schema = z.object({
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
})
type FormData = z.infer<typeof schema>

export default function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [shake, setShake] = useState(false)

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const onSubmit = async (data: FormData) => {
    setIsLoading(true)
    try {
      await login(data)
      navigate('/dashboard')
    } catch {
      setShake(true)
      setTimeout(() => setShake(false), 600)
      toast({ title: 'Invalid credentials', description: 'Check your username and password.', variant: 'destructive' })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-[#0a0a0f]">
      {/* Animated background blobs */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="animate-blob absolute -left-40 -top-40 h-96 w-96 rounded-full bg-violet-600/20 blur-3xl" />
        <div className="animate-blob-delay absolute -right-40 top-1/3 h-80 w-80 rounded-full bg-purple-500/15 blur-3xl" />
        <div className="animate-blob-delay-2 absolute bottom-0 left-1/3 h-72 w-72 rounded-full bg-indigo-600/15 blur-3xl" />
        {/* Grid overlay */}
        <div
          className="absolute inset-0 opacity-[0.03]"
          style={{
            backgroundImage: 'linear-gradient(#fff 1px, transparent 1px), linear-gradient(90deg, #fff 1px, transparent 1px)',
            backgroundSize: '40px 40px',
          }}
        />
      </div>

      {/* Card */}
      <motion.div
        initial={{ opacity: 0, y: 32, scale: 0.95 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.5, ease: [0.25, 0.46, 0.45, 0.94] }}
        className="relative z-10 w-full max-w-sm px-4"
      >
        <motion.div
          animate={shake ? { x: [-6, 6, -6, 6, -3, 3, 0] } : { x: 0 }}
          transition={{ duration: 0.4 }}
          className="rounded-2xl border border-white/10 bg-white/[0.04] p-8 shadow-2xl backdrop-blur-2xl"
        >
          {/* Logo */}
          <div className="mb-8 flex flex-col items-center">
            <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-violet shadow-lg glow">
              <Sparkles className="h-7 w-7 text-white" />
            </div>
            <h1 className="text-2xl font-bold tracking-tight text-white">
              Welcome back
            </h1>
            <p className="mt-1 text-sm text-white/50">Sign in to SkillTracker</p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="username" className="text-white/70">Username</Label>
              <Input
                id="username"
                placeholder="admin"
                {...register('username')}
                className="border-white/10 bg-white/5 text-white placeholder:text-white/30 focus-visible:ring-violet-500"
                autoComplete="username"
              />
              {errors.username && (
                <p className="text-xs text-rose-400">{errors.username.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="password" className="text-white/70">Password</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  {...register('password')}
                  className="border-white/10 bg-white/5 pr-10 text-white placeholder:text-white/30 focus-visible:ring-violet-500"
                  autoComplete="current-password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-white/40 hover:text-white/70 transition-colors"
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              {errors.password && (
                <p className="text-xs text-rose-400">{errors.password.message}</p>
              )}
            </div>

            <Button
              type="submit"
              className="mt-2 w-full bg-gradient-violet text-white shadow-lg hover:opacity-90 active:scale-[0.98]"
              size="lg"
              disabled={isLoading}
            >
              {isLoading ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Signing in...
                </>
              ) : (
                'Sign In'
              )}
            </Button>
          </form>

          <p className="mt-6 text-center text-xs text-white/30">
            Default credentials: <span className="font-mono text-white/50">admin / admin123</span>
          </p>
        </motion.div>
      </motion.div>
    </div>
  )
}
