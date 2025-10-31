# FilmFolk - Movie Review & Social Platform

A modern, full-stack movie review platform with animated UI, built with Go (backend) and Next.js (frontend).

## 🎬 Project Overview

FilmFolk is a comprehensive movie review and social platform featuring:
- User authentication with JWT
- Movie catalog with TMDB integration
- Review system with threaded comments
- Social features (friends, messaging, communities)
- Moderator and admin controls
- Gamification system
- Real-time notifications

## 📁 Project Structure

```
filmfolk/
├── backend/              # Go + PostgreSQL Backend
│   ├── cmd/server/       # Main application entry
│   ├── internal/
│   │   ├── models/       # Database models (13 files)
│   │   ├── services/     # Business logic
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/   # Auth & CORS
│   │   ├── routes/       # API routes
│   │   ├── db/           # Database connection
│   │   ├── config/       # Configuration
│   │   └── utils/        # JWT & password utilities
│   ├── migrations/       # Database schema
│   ├── configs/          # Config files
│   ├── go.mod
│   └── Documentation files
└── frontend/             # Next.js + TypeScript Frontend
    ├── app/              # Next.js App Router
    ├── components/       # React components
    ├── lib/              # Utilities, API, stores
    ├── types/            # TypeScript types
    ├── public/           # Static assets
    └── Configuration files
```

## ✅ Implemented Features (Backend)

### Complete & Working:
1. **Authentication System** ✅
   - User registration/login
   - JWT tokens (access + refresh)
   - Password hashing (bcrypt)
   - Role-based access control
   - Token refresh mechanism

2. **Movie Management** ✅
   - Browse/search movies
   - Submit new movies
   - TMDB API integration
   - Movie approval workflow
   - Rating aggregation

3. **Review System** ✅
   - Write/edit/delete reviews
   - Threaded comments (nested replies)
   - Thread locking
   - Like counters
   - Moderation hooks

4. **Database** ✅
   - Complete schema (21 tables)
   - All relationships defined
   - Indexes for performance
   - ENUM types for safety

### API Endpoints (20+):
- `POST /auth/register` - Create account
- `POST /auth/login` - Login
- `POST /auth/refresh` - Refresh token
- `GET /movies` - Browse movies
- `GET /movies/:id` - Movie details
- `POST /movies` - Submit movie
- `POST /reviews` - Write review
- `GET /reviews/:id` - Get review
- `POST /reviews/comments` - Add comment
- And many more...

## 🎨 Frontend Architecture (To Be Built)

### Tech Stack:
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Animations**: Framer Motion + GSAP
- **State**: Zustand
- **Forms**: React Hook Form + Zod
- **HTTP**: Axios
- **Icons**: Lucide React

### Key Features to Implement:

1. **Animated Homepage**
   - Hero with gradient background
   - Featured movies carousel
   - Latest reviews section
   - Smooth page transitions

2. **Movie Browse**
   - Animated grid layout
   - Filter sidebar
   - Infinite scroll
   - Search with debounce

3. **Movie Detail Page**
   - Parallax hero image
   - Stagger reveal animations
   - Reviews with threading
   - Cast horizontal scroll

4. **Authentication**
   - Login/register forms
   - Validation with animations
   - Error shake effects
   - Success transitions

5. **Review System**
   - Animated star rating
   - Threaded comments UI
   - Like button animations
   - Real-time updates

6. **Micro-Interactions**
   - Hover effects everywhere
   - Click feedback
   - Loading skeletons
   - Toast notifications
   - Modal animations

## 🚀 Quick Start

### Prerequisites
- Go 1.25.2+
- PostgreSQL 14+
- Node.js 18+
- npm or yarn

### Backend Setup

```bash
cd backend

# 1. Create database
createdb filmfolk
psql filmfolk < migrations/001_initial_schema.sql

# 2. Configure (edit configs/config.yaml)

# 3. Run backend
go run cmd/server/main.go
```

Backend runs on **http://localhost:8080**

### Frontend Setup

```bash
cd frontend

# 1. Install dependencies
npm install

# 2. Create .env.local
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local

# 3. Run development server
npm run dev
```

Frontend runs on **http://localhost:3000**

## 📚 Documentation Files

| File | Description |
|------|-------------|
| **QUICKSTART.md** | 5-minute setup guide |
| **FRONTEND_SETUP.md** | Complete frontend implementation guide |
| **API_ENDPOINTS.md** | All API endpoints with examples |
| **API_TESTING.md** | Postman/curl testing guide |
| **PROJECT_SUMMARY.md** | Implementation details & progress |

## 🎯 Current Status

### Backend: ✅ PRODUCTION READY
- All core features working
- **NEW**: Structured logging with Zerolog
- **NEW**: Rate limiting (100 req/min global, 10 req/min auth)
- **NEW**: Security headers (CSP, X-Frame-Options, HSTS)
- **NEW**: CORS with origin whitelisting
- **NEW**: Comprehensive health checks
- **NEW**: Request tracing with unique IDs
- **NEW**: Docker & Docker Compose support
- **NEW**: Makefile for development commands
- **NEW**: Environment-based configuration
- Graceful shutdown
- Full documentation
- Ready to deploy

### Frontend: ✅ BEAUTIFUL UI IMPLEMENTED
- **NEW**: Modern Next.js 15 with TypeScript
- **NEW**: Stunning homepage with animations
- **NEW**: Authentication pages (Login/Register)
- **NEW**: Beautiful UI components (Button, Input, Card)
- **NEW**: Framer Motion animations everywhere
- **NEW**: Responsive Navbar with scroll effects
- **NEW**: Dark mode with custom color scheme
- **NEW**: Glassmorphism effects
- **NEW**: State management with Zustand
- **NEW**: API client with auth interceptors
- **NEW**: Toast notifications
- Fully functional and ready to extend

## 🛠️ Development Workflow

### Testing Backend:

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'

# Save token and make authenticated requests
```

### Building Frontend:

1. Start with authentication pages
2. Build movie browsing UI
3. Implement review system
4. Add animations progressively
5. Connect to backend API
6. Test end-to-end flow

## 🎨 Design Philosophy

### Visual Style:
- **Dark theme primary** (with light mode support)
- **Modern glassmorphism** effects
- **Smooth animations** everywhere
- **Responsive design** mobile-first
- **Accessible** WCAG AA compliant

### Animation Guidelines:
- **Page transitions**: 300-500ms
- **Micro-interactions**: 150-300ms
- **Hover effects**: Subtle scale (1.02-1.05)
- **Loading states**: Skeleton screens
- **Feedback**: Immediate visual response

### Color Palette:
```
Primary: #0ea5e9 (Sky Blue)
Secondary: #8b5cf6 (Purple)
Accent: #f59e0b (Amber)
Background Dark: #0a0a0a
Card Dark: #1a1a1a
```

## 📦 Tech Stack Summary

### Backend:
- **Go 1.25.2** - Main language
- **Gin** - Web framework
- **GORM** - ORM
- **PostgreSQL** - Database
- **JWT** - Authentication
- **Bcrypt** - Password hashing

### Frontend:
- **Next.js 15** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Framer Motion** - Animations
- **GSAP** - Complex animations
- **Zustand** - State management
- **Axios** - HTTP client
- **React Hook Form** - Form handling
- **Zod** - Schema validation

## 🔐 Security Features

- JWT with short-lived access tokens
- Refresh token rotation
- Password hashing with bcrypt (cost: 12)
- CORS configured
- SQL injection protection (GORM)
- Input validation
- Role-based access control

## 🌟 Key Highlights

1. **Production-Ready Backend**
   - Clean architecture
   - Service layer pattern
   - Comprehensive error handling
   - Graceful shutdown
   - Auto-migrations (dev)

2. **Modern Frontend Setup**
   - App Router (Next.js 15)
   - TypeScript throughout
   - Animation library integration
   - State management ready
   - API client configured

3. **Complete Documentation**
   - Setup guides
   - API reference
   - Component examples
   - Animation patterns
   - Testing workflows

## 📖 Next Steps

### Immediate:
1. ✅ Backend is complete - test it!
2. 📝 Install frontend dependencies
3. 🎨 Start building UI components
4. 🔌 Connect to backend
5. ✨ Add animations

### Future Enhancements:
- User lists (watched, plan to watch)
- Direct messaging
- Friend system
- Communities & world chat
- OAuth login
- AI content moderation
- Push notifications
- Mobile app

## 💡 Learning Outcomes

By building this project, you'll master:
- Full-stack development (Go + React)
- RESTful API design
- JWT authentication
- Database design
- State management
- Advanced animations
- TypeScript
- Modern React patterns
- Production deployment

## 🤝 Contributing

This is a learning project. Feel free to:
- Add new features
- Improve animations
- Enhance UI/UX
- Optimize performance
- Write tests

## 📄 License

MIT License - Feel free to use for learning!

---

## 🎊 You're Ready!

**Backend**: Fully functional ✅
**Frontend**: Structured and ready to build 📋
**Documentation**: Complete and comprehensive 📚

**Start coding and bring FilmFolk to life!** 🚀🎬

Made with ❤️ for learning modern full-stack development.
