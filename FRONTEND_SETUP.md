# FilmFolk Frontend Setup Guide

Complete guide to set up the modern, animated Next.js frontend.

## Project Structure

```
filmfolk/
├── backend/           # Go backend (already complete)
└── frontend/          # Next.js frontend (to be set up)
```

## Quick Setup

### Option 1: Automated Setup (Recommended)

```bash
cd frontend
npx create-next-app@latest . --typescript --tailwind --app --eslint
```

When prompted, choose:
- ✅ TypeScript
- ✅ ESLint
- ✅ Tailwind CSS
- ✅ App Router
- ✅ Import alias (@/*)
- ❌ src/ directory (No)

Then install additional dependencies:

```bash
npm install framer-motion gsap axios zustand react-hook-form zod @hookform/resolvers clsx tailwind-merge lucide-react date-fns react-hot-toast
```

### Option 2: Manual Setup

I've already created `frontend/package.json`. Now run:

```bash
cd frontend
npm install
```

## Frontend Architecture

### Folder Structure

```
frontend/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   ├── movies/
│   │   ├── page.tsx           # Movie list
│   │   └── [id]/
│   │       └── page.tsx       # Movie details
│   ├── reviews/
│   │   └── [id]/
│   │       └── page.tsx       # Review details
│   ├── profile/
│   │   └── page.tsx
│   ├── layout.tsx             # Root layout
│   ├── page.tsx               # Homepage
│   └── globals.css
├── components/
│   ├── ui/                    # Base UI components
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── card.tsx
│   │   └── modal.tsx
│   ├── movies/
│   │   ├── MovieCard.tsx
│   │   ├── MovieGrid.tsx
│   │   └── MovieFilters.tsx
│   ├── reviews/
│   │   ├── ReviewCard.tsx
│   │   ├── ReviewForm.tsx
│   │   └── CommentThread.tsx
│   ├── animations/
│   │   ├── FadeIn.tsx
│   │   ├── SlideIn.tsx
│   │   └── ScaleIn.tsx
│   └── layout/
│       ├── Navbar.tsx
│       ├── Footer.tsx
│       └── Sidebar.tsx
├── lib/
│   ├── api/
│   │   ├── axios.ts           # Axios instance
│   │   ├── auth.ts            # Auth API calls
│   │   ├── movies.ts          # Movie API calls
│   │   └── reviews.ts         # Review API calls
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useMovies.ts
│   │   └── useReviews.ts
│   ├── store/
│   │   ├── authStore.ts       # Zustand auth store
│   │   └── movieStore.ts
│   ├── utils/
│   │   ├── cn.ts              # Class name utility
│   │   └── format.ts
│   └── animations/
│       ├── variants.ts        # Framer Motion variants
│       └── gsap.ts            # GSAP animations
├── types/
│   ├── api.ts                 # API response types
│   ├── movie.ts
│   ├── review.ts
│   └── user.ts
├── public/
│   ├── images/
│   └── icons/
├── tailwind.config.ts
├── tsconfig.json
├── next.config.js
└── package.json
```

## Key Configuration Files

### 1. tailwind.config.ts

```typescript
import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#f0f9ff",
          100: "#e0f2fe",
          500: "#0ea5e9",
          600: "#0284c7",
          700: "#0369a1",
        },
      },
      animation: {
        "fade-in": "fadeIn 0.5s ease-in-out",
        "slide-up": "slideUp 0.5s ease-out",
        "scale-in": "scaleIn 0.3s ease-out",
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
        slideUp: {
          "0%": { transform: "translateY(20px)", opacity: "0" },
          "100%": { transform: "translateY(0)", opacity: "1" },
        },
        scaleIn: {
          "0%": { transform: "scale(0.95)", opacity: "0" },
          "100%": { transform: "scale(1)", opacity: "1" },
        },
      },
    },
  },
  plugins: [],
};
export default config;
```

### 2. next.config.js

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'image.tmdb.org',
      },
    ],
  },
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  },
};

module.exports = nextConfig;
```

### 3. .env.local

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Core Implementation Examples

### API Client (`lib/api/axios.ts`)

```typescript
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for auth token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const { data } = await axios.post(
          `${process.env.NEXT_PUBLIC_API_URL}/auth/refresh`,
          { refresh_token: refreshToken }
        );

        localStorage.setItem('access_token', data.access_token);
        originalRequest.headers.Authorization = `Bearer ${data.access_token}`;

        return apiClient(originalRequest);
      } catch (refreshError) {
        // Redirect to login
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
```

### Auth Store (`lib/store/authStore.ts`)

```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id: number;
  username: string;
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  login: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      login: (user, accessToken, refreshToken) =>
        set({ user, accessToken, refreshToken, isAuthenticated: true }),
      logout: () =>
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

### Animated Movie Card Component

```typescript
'use client';

import { motion } from 'framer-motion';
import Image from 'next/image';
import { Star } from 'lucide-react';

interface MovieCardProps {
  movie: {
    id: number;
    title: string;
    release_year: number;
    poster_url?: string;
    average_rating?: number;
    genres?: string[];
  };
}

export function MovieCard({ movie }: MovieCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ scale: 1.05, transition: { duration: 0.2 } }}
      className="relative group cursor-pointer"
    >
      <div className="relative h-[400px] rounded-lg overflow-hidden shadow-xl">
        {movie.poster_url ? (
          <Image
            src={movie.poster_url}
            alt={movie.title}
            fill
            className="object-cover transition-transform duration-300 group-hover:scale-110"
          />
        ) : (
          <div className="w-full h-full bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
            <span className="text-gray-500 text-6xl">🎬</span>
          </div>
        )}

        {/* Overlay on hover */}
        <motion.div
          initial={{ opacity: 0 }}
          whileHover={{ opacity: 1 }}
          className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/50 to-transparent flex flex-col justify-end p-4"
        >
          <h3 className="text-white font-bold text-lg mb-2">{movie.title}</h3>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-300">{movie.release_year}</span>
            {movie.average_rating && (
              <div className="flex items-center gap-1">
                <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                <span className="text-white font-semibold">
                  {movie.average_rating.toFixed(1)}
                </span>
              </div>
            )}
          </div>
          {movie.genres && movie.genres.length > 0 && (
            <div className="flex gap-2 mt-2 flex-wrap">
              {movie.genres.slice(0, 3).map((genre) => (
                <span
                  key={genre}
                  className="px-2 py-1 bg-white/20 backdrop-blur-sm rounded-full text-xs text-white"
                >
                  {genre}
                </span>
              ))}
            </div>
          )}
        </motion.div>
      </div>
    </motion.div>
  );
}
```

### GSAP Scroll Animations

```typescript
import { useEffect, useRef } from 'react';
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

export function useScrollAnimation() {
  const ref = useRef(null);

  useEffect(() => {
    const element = ref.current;

    gsap.fromTo(
      element,
      {
        opacity: 0,
        y: 50,
      },
      {
        opacity: 1,
        y: 0,
        duration: 1,
        scrollTrigger: {
          trigger: element,
          start: 'top 80%',
          end: 'top 20%',
          scrub: 1,
        },
      }
    );
  }, []);

  return ref;
}
```

## Animation Patterns

### 1. Page Transitions (Framer Motion)

```typescript
const pageVariants = {
  initial: { opacity: 0, x: -20 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 20 },
};

const pageTransition = {
  type: 'tween',
  ease: 'anticipate',
  duration: 0.5,
};
```

### 2. Stagger Animations

```typescript
const containerVariants = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0 },
};
```

### 3. Micro-Interactions

- **Hover effects** on buttons and cards
- **Click animations** with scale
- **Loading skeletons** with shimmer effect
- **Toast notifications** with slide-in
- **Modal animations** with backdrop blur
- **Form field focus** with border glow
- **Scroll progress bar** at top
- **Infinite scroll** loading indicator

## Color Palette

### Dark Theme (Primary)

```css
--background: #0a0a0a
--foreground: #ffffff
--card: #1a1a1a
--primary: #0ea5e9
--secondary: #8b5cf6
--accent: #f59e0b
--muted: #374151
--destructive: #ef4444
```

### Light Theme

```css
--background: #ffffff
--foreground: #0a0a0a
--card: #f3f4f6
--primary: #0284c7
--secondary: #7c3aed
--accent: #d97706
--muted: #9ca3af
--destructive: #dc2626
```

## Running the Project

### 1. Start Backend

```bash
cd backend
go run cmd/server/main.go
```

Backend runs on: http://localhost:8080

### 2. Start Frontend

```bash
cd frontend
npm run dev
```

Frontend runs on: http://localhost:3000

## Key Features to Implement

### Homepage
- Hero section with animated gradient background
- Featured movies carousel (GSAP)
- Latest reviews with stagger animation
- Call-to-action buttons with hover effects

### Movie Browse Page
- Grid layout with animated movie cards
- Filter sidebar with smooth transitions
- Search with debounce
- Infinite scroll with loading animation
- Sort dropdown with animations

### Movie Detail Page
- Hero section with backdrop image and parallax
- Movie info with stagger reveal
- Reviews section with threading
- "Write Review" floating action button
- Cast members horizontal scroll

### Review System
- Star rating with hover animation
- Textarea with character count
- Threaded comments with indentation
- Like button with heart animation
- Lock thread toggle

### Authentication
- Login/register forms with validation
- Password strength indicator
- Social login buttons (prepared for OAuth)
- Form errors with shake animation

### Profile Page
- User stats with count-up animation
- Reviews list with filters
- Edit profile modal
- Avatar upload with preview

## Next Steps

1. Run `cd frontend && npm install`
2. Create `.env.local` with API URL
3. Start implementing components from examples above
4. Test with backend running
5. Add more animations as you go

## Resources

- [Framer Motion Docs](https://www.framer.com/motion/)
- [GSAP Docs](https://greensock.com/docs/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Next.js Docs](https://nextjs.org/docs)
- [Zustand](https://docs.pmnd.rs/zustand/getting-started/introduction)

---

**The frontend setup is ready to go! Start building from the examples above.** 🚀
