'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { motion } from 'framer-motion'
import { Film, Mail, Lock, ArrowRight } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/Card'
import { authApi } from '@/lib/api'
import { useAuthStore } from '@/store/authStore'

export default function LoginPage() {
  const router = useRouter()
  const { setAuth } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})
    setIsLoading(true)

    try {
      const response = await authApi.login(formData)
      const { user, access_token, refresh_token } = response.data

      setAuth(user, access_token, refresh_token)
      toast.success(`Welcome back, ${user.username}!`)
      router.push('/dashboard')
    } catch (error: any) {
      const message = error.response?.data?.error || 'Login failed. Please check your credentials.'
      toast.error(message)

      if (error.response?.data?.errors) {
        setErrors(error.response.data.errors)
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen px-4 py-12">
      <motion.div
        className="w-full max-w-md"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <div className="text-center mb-8">
          <motion.div
            className="inline-flex items-center gap-2 mb-4"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.2, type: 'spring' }}
          >
            <Film className="h-10 w-10 text-primary" />
            <span className="text-2xl font-bold bg-gradient-to-r from-primary to-purple-400 bg-clip-text text-transparent">
              FilmFolk
            </span>
          </motion.div>
          <h1 className="text-3xl font-bold mb-2">Welcome Back</h1>
          <p className="text-muted-foreground">Sign in to continue your movie journey</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Login to your account</CardTitle>
            <CardDescription>Enter your credentials to access your account</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <Input
                type="email"
                label="Email"
                placeholder="john@example.com"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                error={errors.email}
                required
              />

              <Input
                type="password"
                label="Password"
                placeholder="Enter your password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                error={errors.password}
                required
              />

              <Button type="submit" className="w-full" isLoading={isLoading}>
                {!isLoading && <Mail className="h-4 w-4" />}
                Sign In
                {!isLoading && <ArrowRight className="h-4 w-4" />}
              </Button>
            </form>
          </CardContent>
          <CardFooter className="flex flex-col gap-4">
            <div className="text-sm text-center text-muted-foreground">
              Don't have an account?{' '}
              <Link href="/register" className="text-primary hover:underline font-medium">
                Sign up
              </Link>
            </div>
          </CardFooter>
        </Card>

        <motion.p
          className="text-center text-sm text-muted-foreground mt-8"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          By continuing, you agree to our Terms of Service and Privacy Policy
        </motion.p>
      </motion.div>
    </div>
  )
}
