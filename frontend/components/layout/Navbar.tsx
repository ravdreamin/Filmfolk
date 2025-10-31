'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { motion, useScroll, useTransform } from 'framer-motion'
import { Film, User, LogOut, Menu, X, Search } from 'lucide-react'
import { useAuthStore } from '@/store/authStore'
import Button from '@/components/ui/Button'
import { cn } from '@/lib/utils'

export default function Navbar() {
  const [isOpen, setIsOpen] = useState(false)
  const { scrollY } = useScroll()
  const { isAuthenticated, user, clearAuth } = useAuthStore()

  const backgroundColor = useTransform(
    scrollY,
    [0, 100],
    ['rgba(0, 0, 0, 0)', 'rgba(0, 0, 0, 0.8)']
  )

  const backdropBlur = useTransform(
    scrollY,
    [0, 100],
    ['blur(0px)', 'blur(10px)']
  )

  const handleLogout = () => {
    clearAuth()
    window.location.href = '/login'
  }

  return (
    <motion.nav
      className="fixed top-0 left-0 right-0 z-50 border-b border-white/10"
      style={{
        backgroundColor,
        backdropFilter: backdropBlur,
        WebkitBackdropFilter: backdropBlur,
      }}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2 group">
            <motion.div
              whileHover={{ rotate: 360 }}
              transition={{ duration: 0.5 }}
            >
              <Film className="h-8 w-8 text-primary" />
            </motion.div>
            <span className="text-xl font-bold bg-gradient-to-r from-primary to-purple-400 bg-clip-text text-transparent">
              FilmFolk
            </span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-6">
            <NavLink href="/movies">Movies</NavLink>
            <NavLink href="/reviews">Reviews</NavLink>
            {isAuthenticated && <NavLink href="/dashboard">Dashboard</NavLink>}

            {/* Search */}
            <motion.div
              className="relative"
              initial={{ width: 40 }}
              whileHover={{ width: 200 }}
              transition={{ duration: 0.3 }}
            >
              <Search className="absolute left-2 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <input
                type="text"
                placeholder="Search..."
                className="w-full h-9 pl-8 pr-3 rounded-md border border-input bg-background/50 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </motion.div>

            {/* Auth Buttons */}
            {isAuthenticated ? (
              <div className="flex items-center gap-3">
                <Link href="/profile">
                  <motion.div
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    className="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-accent transition-colors"
                  >
                    <User className="h-4 w-4" />
                    <span className="text-sm font-medium">{user?.username}</span>
                  </motion.div>
                </Link>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={handleLogout}
                  title="Logout"
                >
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-3">
                <Link href="/login">
                  <Button variant="ghost">Login</Button>
                </Link>
                <Link href="/register">
                  <Button>Sign Up</Button>
                </Link>
              </div>
            )}
          </div>

          {/* Mobile menu button */}
          <Button
            variant="ghost"
            size="icon"
            className="md:hidden"
            onClick={() => setIsOpen(!isOpen)}
          >
            {isOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
        </div>

        {/* Mobile Navigation */}
        {isOpen && (
          <motion.div
            className="md:hidden py-4 space-y-2"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
          >
            <MobileNavLink href="/movies" onClick={() => setIsOpen(false)}>
              Movies
            </MobileNavLink>
            <MobileNavLink href="/reviews" onClick={() => setIsOpen(false)}>
              Reviews
            </MobileNavLink>
            {isAuthenticated && (
              <>
                <MobileNavLink href="/dashboard" onClick={() => setIsOpen(false)}>
                  Dashboard
                </MobileNavLink>
                <MobileNavLink href="/profile" onClick={() => setIsOpen(false)}>
                  Profile
                </MobileNavLink>
                <button
                  onClick={handleLogout}
                  className="block w-full text-left px-3 py-2 rounded-md text-sm hover:bg-accent transition-colors"
                >
                  Logout
                </button>
              </>
            )}
            {!isAuthenticated && (
              <>
                <MobileNavLink href="/login" onClick={() => setIsOpen(false)}>
                  Login
                </MobileNavLink>
                <MobileNavLink href="/register" onClick={() => setIsOpen(false)}>
                  Sign Up
                </MobileNavLink>
              </>
            )}
          </motion.div>
        )}
      </div>
    </motion.nav>
  )
}

function NavLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <Link href={href}>
      <motion.span
        className="text-sm font-medium text-foreground/80 hover:text-foreground transition-colors relative group"
        whileHover={{ scale: 1.05 }}
      >
        {children}
        <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary group-hover:w-full transition-all duration-300" />
      </motion.span>
    </Link>
  )
}

function MobileNavLink({ href, onClick, children }: { href: string; onClick: () => void; children: React.ReactNode }) {
  return (
    <Link href={href} onClick={onClick}>
      <motion.div
        className="block px-3 py-2 rounded-md text-sm hover:bg-accent transition-colors"
        whileHover={{ x: 4 }}
      >
        {children}
      </motion.div>
    </Link>
  )
}
