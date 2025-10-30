# FilmFolk - Complete Setup & Run Guide

Everything you need to get FilmFolk running in 10 minutes.

## 📋 What You Have

```
✅ Complete Go backend (fully functional)
✅ Database schema (21 tables)
✅ API endpoints (20+ working)
✅ Frontend structure (Next.js with config files)
✅ Full documentation
```

## 🚀 Step-by-Step Setup

### Step 1: Backend Setup (5 minutes)

```bash
# Navigate to backend
cd backend

# Create PostgreSQL database
createdb filmfolk

# Run database migrations
psql filmfolk < migrations/001_initial_schema.sql

# Verify database
psql filmfolk -c "\dt"
# You should see 21 tables

# Test backend build
go build ./cmd/server

# Run backend
go run cmd/server/main.go
```

**Expected output:**
```
Loading configuration...
Configuration loaded from: YAML
Environment: development
Connecting to database...
Database connection established successfully
Running auto-migrations...
Auto-migrations completed successfully
🚀 FilmFolk server starting on http://localhost:8080
```

**Test it:**
```bash
# In another terminal
curl http://localhost:8080/health
# Should return: {"service":"filmfolk-api","status":"ok"}
```

### Step 2: Frontend Setup (5 minutes)

```bash
# Open new terminal, navigate to frontend
cd frontend

# Install dependencies
npm install

# This will install:
# - Next.js 15
# - React 19
# - Framer Motion
# - GSAP
# - Tailwind CSS
# - Zustand
# - Axios
# - And more...

# Create environment file
cat > .env.local << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
EOF

# Run development server
npm run dev
```

**Expected output:**
```
  ▲ Next.js 15.0.3
  - Local:        http://localhost:3000

  ✓ Ready in 2.3s
```

### Step 3: Verify Everything Works

**Backend (Terminal 1):**
```bash
cd backend
go run cmd/server/main.go
# Keep running...
```

**Frontend (Terminal 2):**
```bash
cd frontend
npm run dev
# Keep running...
```

**Test (Terminal 3):**
```bash
# Test backend health
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# Visit frontend
open http://localhost:3000
```

## 🏗️ What's Implemented vs What's Next

### ✅ Backend (100% Done)

**Working Now:**
- User registration & login
- JWT authentication
- Movie browsing & submission
- Review system with comments
- Movie moderation
- Role-based access

**API Endpoints Ready:**
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/auth/me
POST   /api/v1/auth/refresh
GET    /api/v1/movies
GET    /api/v1/movies/:id
POST   /api/v1/movies
GET    /api/v1/movies/:id/reviews
POST   /api/v1/reviews
PUT    /api/v1/reviews/:id
DELETE /api/v1/reviews/:id
POST   /api/v1/reviews/comments
...and more!
```

### 📝 Frontend (To Build)

**What You Need to Create:**

1. **Authentication Pages** (`app/(auth)/`)
   - Login page
   - Register page
   - Forms with validation
   - Token management

2. **Movie Pages** (`app/movies/`)
   - Movie list/browse page
   - Movie detail page
   - Movie submission form

3. **Review Pages** (`app/reviews/`)
   - Review detail with comments
   - Review form
   - Comment threading UI

4. **Components** (`components/`)
   - MovieCard with animations
   - ReviewCard
   - CommentThread
   - Navbar
   - Footer
   - Buttons, Inputs, etc.

5. **API Integration** (`lib/api/`)
   - Axios setup (already configured in docs)
   - Auth API calls
   - Movie API calls
   - Review API calls

6. **State Management** (`lib/store/`)
   - Auth store (Zustand)
   - User preferences
   - UI state

## 📚 Where to Find Everything

### Documentation:
- **README.md** - Project overview
- **FRONTEND_SETUP.md** - Complete frontend guide with code examples
- **API_ENDPOINTS.md** - All API endpoints documented
- **API_TESTING.md** - How to test with Postman/curl
- **QUICKSTART.md** - Backend quick start
- **PROJECT_SUMMARY.md** - What's implemented, what's not

### Code Examples in FRONTEND_SETUP.md:
- ✅ API client configuration
- ✅ Auth store (Zustand)
- ✅ Animated MovieCard component
- ✅ GSAP scroll animations
- ✅ Page transitions
- ✅ Framer Motion variants

### Backend Code:
- `backend/cmd/server/main.go` - Server entry point
- `backend/internal/handlers/` - HTTP handlers
- `backend/internal/services/` - Business logic
- `backend/internal/models/` - Database models

### Frontend Config:
- `frontend/package.json` - Dependencies
- `frontend/tailwind.config.ts` - Tailwind + animations
- `frontend/tsconfig.json` - TypeScript config
- `frontend/next.config.js` - Next.js config

## 🎨 Frontend Development Guide

### Start Building:

1. **Create Basic Layout:**
```bash
# Create app structure
mkdir -p app/(auth)/login
mkdir -p app/(auth)/register
mkdir -p app/movies
mkdir -p components/ui
mkdir -p lib/api
```

2. **Copy Examples from FRONTEND_SETUP.md:**
   - API client setup
   - Auth store
   - MovieCard component
   - All animation patterns

3. **Install Dependencies:**
```bash
cd frontend
npm install
```

4. **Start Coding:**
```bash
npm run dev
# Visit http://localhost:3000
```

### Development Workflow:

1. **Create a page** in `app/`
2. **Build components** in `components/`
3. **Add animations** with Framer Motion/GSAP
4. **Connect to API** using axios
5. **Manage state** with Zustand
6. **Style with Tailwind**
7. **Test** with backend running

## 🎯 Quick Wins

### First Things to Build:

**1. Login Page (30 min)**
- Form with email/password
- Call `/api/v1/auth/login`
- Store tokens
- Redirect to homepage

**2. Movie List Page (1 hour)**
- Fetch from `/api/v1/movies`
- Display in grid
- Add Framer Motion animations
- Click to view details

**3. Movie Detail Page (1 hour)**
- Fetch from `/api/v1/movies/:id`
- Show movie info
- List reviews
- Add review button

**4. Write Review (45 min)**
- Form with rating + text
- POST to `/api/v1/reviews`
- Show success message
- Refresh reviews

## 🐛 Troubleshooting

### Backend Issues:

**"Database connection failed"**
```bash
# Check PostgreSQL is running
pg_ctl status

# Check database exists
psql -l | grep filmfolk

# Verify config
cat backend/configs/config.yaml
```

**"Port 8080 already in use"**
```bash
# Find and kill process
lsof -ti:8080 | xargs kill -9
```

### Frontend Issues:

**"Module not found"**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**"EADDRINUSE: address already in use"**
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9
```

**"API call failing with CORS error"**
- Make sure backend is running on port 8080
- Check `.env.local` has correct API URL
- Backend has CORS middleware enabled (already done)

## ✨ Animation Tips

### Use Framer Motion for:
- Page transitions
- Component mount/unmount
- Hover effects
- Click feedback
- Stagger animations

### Use GSAP for:
- Scroll animations
- Complex timelines
- SVG animations
- Advanced tweens

### Use Tailwind for:
- Simple hover effects
- Loading skeletons
- Basic transitions

## 📊 Project Status

```
Backend:     ████████████████████ 100% ✅
Frontend:    ░░░░░░░░░░░░░░░░░░░░  0%  📋
Database:    ████████████████████ 100% ✅
Docs:        ████████████████████ 100% ✅
```

## 🎓 Learning Path

### Day 1-2: Setup & Understanding
- ✅ Set up backend
- ✅ Set up frontend
- ✅ Read documentation
- ✅ Test API with Postman

### Day 3-4: Basic UI
- 📝 Create login/register pages
- 📝 Build movie list page
- 📝 Add basic animations

### Day 5-6: Core Features
- 📝 Movie detail page
- 📝 Review system
- 📝 Comment threading

### Day 7-8: Polish
- 📝 Advanced animations
- 📝 Micro-interactions
- 📝 Responsive design

### Day 9-10: Extra Features
- 📝 User profile
- 📝 Search & filters
- 📝 Loading states

## 🎊 You're All Set!

**Everything is ready:**
- ✅ Backend running on http://localhost:8080
- ✅ Frontend dev server on http://localhost:3000
- ✅ Database with complete schema
- ✅ API documented and working
- ✅ Frontend structure prepared
- ✅ Component examples provided
- ✅ Animation patterns documented

**Now start building! 🚀**

Follow **FRONTEND_SETUP.md** for detailed implementation guide.

---

**Questions?** Check the documentation files!
**Stuck?** Review the code examples!
**Need inspiration?** Look at the animation patterns!

**Happy coding! 🎬✨**
