'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { motion } from 'framer-motion'
import { Film, User, Mail, Lock, ArrowRight, Check } from 'lucide-react'
import toast from 'react-hot-toast'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/Card'
import { authApi } from '@/lib/api'
import { useAuthStore } from '@/store/authStore'

export default function RegisterPage() {
  const router = useRouter()
  const { setAuth } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const passwordRequirements = [
    { label: 'At least 8 characters', met: formData.password.length >= 8 },
    { label: 'Contains a number', met: /\d/.test(formData.password) },
    { label: 'Contains uppercase', met: /[A-Z]/.test(formData.password) },
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})

    // Validate passwords match
    if (formData.password !== formData.confirmPassword) {
      setErrors({ confirmPassword: 'Passwords do not match' })
      return
    }

    setIsLoading(true)

    try {
      const response = await authApi.register({
        username: formData.username,
        email: formData.email,
        password: formData.password,
      })

      const { user, access_token, refresh_token } = response.data

      setAuth(user, access_token, refresh_token)
      toast.success(`Welcome to FilmFolk, ${user.username}!`)
      router.push('/dashboard')
    } catch (error: any) {
      const message = error.response?.data?.error || 'Registration failed. Please try again.'
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
          <h1 className="text-3xl font-bold mb-2">Join FilmFolk</h1>
          <p className="text-muted-foreground">Create your account and start your cinematic journey</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Create an account</CardTitle>
            <CardDescription>Enter your information to get started</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <Input
                type="text"
                label="Username"
                placeholder="johndoe"
                value={formData.username}
                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                error={errors.username}
                required
              />

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
                placeholder="Create a strong password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                error={errors.password}
                required
              />

              {formData.password && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  className="space-y-2"
                >
                  {passwordRequirements.map((req, index) => (
                    <motion.div
                      key={index}
                      className="flex items-center gap-2 text-sm"
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                    >
                      <div
                        className={`w-4 h-4 rounded-full flex items-center justify-center ${
                          req.met ? 'bg-green-500' : 'bg-muted'
                        }`}
                      >
                        {req.met && <Check className="h-3 w-3 text-white" />}
                      </div>
                      <span className={req.met ? 'text-green-500' : 'text-muted-foreground'}>
                        {req.label}
                      </span>
                    </motion.div>
                  ))}
                </motion.div>
              )}

              <Input
                type="password"
                label="Confirm Password"
                placeholder="Re-enter your password"
                value={formData.confirmPassword}
                onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                error={errors.confirmPassword}
                required
              />

              <Button type="submit" className="w-full" isLoading={isLoading}>
                {!isLoading && <User className="h-4 w-4" />}
                Create Account
                {!isLoading && <ArrowRight className="h-4 w-4" />}
              </Button>
            </form>
          </CardContent>
          <CardFooter className="flex flex-col gap-4">
            <div className="text-sm text-center text-muted-foreground">
              Already have an account?{' '}
              <Link href="/login" className="text-primary hover:underline font-medium">
                Sign in
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
          By creating an account, you agree to our Terms of Service and Privacy Policy
        </motion.p>
      </motion.div>
    </div>
  )
}
