# FilmFolk

A comprehensive movie review and social platform built with Go, PostgreSQL, and modern web technologies.

## Features

### Core Features
- **User Authentication**
  - Email/Password registration and login
  - OAuth2 support (Google, Facebook, Instagram, Twitter)
  - JWT-based authentication with refresh tokens
  - Guest mode (read-only access)
  - Three user roles: User, Moderator, Admin

- **Movie Management**
  - Comprehensive movie catalog with TMDB integration
  - Movie information with cast details
  - User-submitted movies with moderation workflow

- **Review System**
  - Write reviews with 1-10 ratings
  - Threaded comments on reviews
  - Review authors can lock threads
  - Like system for reviews and comments
  - AI-powered content moderation and sentiment analysis

- **User Lists**
  - Personal movie lists (Watched, Dropped, Plan to Watch)
  - Custom lists with privacy controls

- **Social Features**
  - Direct messaging between users
  - Friend system with taste-based recommendations
  - Community chat rooms (topic-based)
  - Global world chat

- **Gamification**
  - User titles that evolve based on engagement
  - Engagement scoring system

- **Moderation**
  - AI pre-moderation flags suspicious content
  - Moderators review flagged content
  - Warning system with escalation to admins
  - Ban/suspend functionality
  - Complete audit trail

## Tech Stack

- **Backend**: Go 1.25.2
- **Web Framework**: Gin
- **Database**: PostgreSQL 16+
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt (cost 12)
- **Logging**: Zerolog (structured JSON logging)
- **Configuration**: Viper + godotenv
- **Containerization**: Docker + Docker Compose
- **External APIs**:
  - TMDB (movie data)
  - OpenAI (content moderation & sentiment analysis)

### Production-Ready Features

- **Structured Logging**: JSON logging with zerolog for production observability
- **Rate Limiting**: Per-IP rate limiting (100 req/min global, 10 req/min auth)
- **Security Headers**: CSP, X-Frame-Options, HSTS, XSS Protection
- **CORS**: Configurable origin whitelisting for production
- **Health Checks**: Multiple endpoints for Kubernetes/Docker health probes
- **Graceful Shutdown**: Signal handling for zero-downtime deployments
- **Request Tracing**: Unique request IDs for distributed tracing
- **Docker Support**: Multi-stage builds, non-root user, health checks
- **Environment-based Config**: Full support for env vars and secrets management

## Project Structure

```
filmfolk/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── configs/
│   └── config.yaml              # Configuration file
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── db/
│   │   └── database.go          # Database connection & migrations
│   ├── handlers/
│   │   └── auth_handler.go      # HTTP request handlers
│   ├── middleware/
│   │   └── auth.go              # Authentication middleware
│   ├── models/                  # Database models
│   │   ├── user.go
│   │   ├── movie.go
│   │   ├── review.go
│   │   ├── community.go
│   │   └── ... (13 model files)
│   ├── routes/
│   │   └── routes.go            # API route definitions
│   ├── services/
│   │   └── auth_service.go      # Business logic
│   └── utils/
│       ├── jwt.go               # JWT utilities
│       └── password.go          # Password hashing
├── migrations/
│   └── 001_initial_schema.sql   # Database schema
├── .env.example                 # Environment variables template
├── go.mod                       # Go dependencies
└── README.md

```

## Getting Started

### Prerequisites

- Go 1.25.2 or higher
- PostgreSQL 14+
- Docker & Docker Compose (optional, recommended)
- Make (optional, for convenience commands)
- (Optional) TMDB API key
- (Optional) OpenAI API key

### Quick Start with Docker (Recommended)

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd filmfolk/backend
```

2. **Configure environment**
```bash
# Copy example env file
cp .env.example .env

# Edit .env with your settings
nano .env
```

3. **Start with Docker Compose**
```bash
# Build and start all services (API + PostgreSQL)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

The API will be available at http://localhost:8080

### Local Development (Without Docker)

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd filmfolk/backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up PostgreSQL**
```bash
# Create database
createdb filmfolk

# Run migrations
go run cmd/migrate/main.go
```

4. **Configure environment**
```bash
# Copy example config
cp .env.example .env

# Edit .env with your settings
# OR edit configs/config.yaml
```

5. **Run the server**
```bash
# Using Make
make run

# Or directly with go
go run cmd/server/main.go

# Or build and run
make build
./bin/server
```

The server will start on http://localhost:8080

### Using Makefile Commands

The project includes a Makefile for common tasks:

```bash
make help              # Show all available commands
make run               # Run the application locally
make build             # Build the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make docker-up         # Start services with Docker
make docker-down       # Stop Docker services
make docker-logs       # View API logs
make migrate-up        # Run database migrations
make fmt               # Format code
make lint              # Run linter
make clean             # Clean build artifacts
```

### Configuration

The application supports two configuration methods:

1. **YAML Configuration** (`configs/config.yaml`)
   - Preferred for development
   - Easier to read and edit

2. **Environment Variables** (`.env` file)
   - Preferred for production
   - Used in Docker/Kubernetes deployments

Priority: YAML > Environment Variables

## API Endpoints

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user"
  },
  "expires_in": 900
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Logout
```http
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

### Health Check Endpoints

```http
# Basic health check
GET /health

# Detailed health check with dependencies
GET /health/detailed

# Kubernetes readiness probe
GET /health/ready

# Kubernetes liveness probe
GET /health/live
```

Response (detailed):
```json
{
  "status": "healthy",
  "timestamp": "2025-10-31T10:00:00Z",
  "service": "filmfolk-api",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "checks": {
    "database": {
      "status": "healthy",
      "message": "Database connection is healthy",
      "latency": "2.5ms"
    },
    "memory": {
      "status": "healthy",
      "message": "Memory usage is normal (Alloc: 45MB, Sys: 72MB)"
    }
  }
}
```

## Authentication Flow

1. **Registration/Login** → Receive access token + refresh token
2. **API Requests** → Include access token in Authorization header
   ```
   Authorization: Bearer <access_token>
   ```
3. **Token Expiration** → Use refresh token to get new access token
4. **Logout** → Revoke refresh token

### Token Lifetimes
- **Access Token**: 15 minutes (short-lived for security)
- **Refresh Token**: 7 days (stored in database, can be revoked)

## Development

### Running Migrations

This project uses `golang-migrate` to manage database schema changes. The migration files are located in the `migrations` directory.

First, ensure all dependencies are up to date by running:
```bash
go mod tidy
```

To run the migrations, use the `migrate` command:

```bash
# Apply all available migrations
go run cmd/migrate/main.go up

# Roll back the last migration
go run cmd/migrate/main.go down

# Drop the entire database (use with caution!)
go run cmd/migrate/main.go drop

# Check the current migration version
go run cmd/migrate/main.go version
```


### Database Schema

The database uses PostgreSQL-specific features:
- ENUM types for type safety
- Array columns for genres
- Indexes for performance
- Foreign keys with CASCADE for data integrity

See [migrations/001_initial_schema.sql](migrations/001_initial_schema.sql) for complete schema.

## Security Features

- **Password Hashing**: bcrypt with cost factor 12
- **JWT Tokens**: HS256 algorithm, signed with secret
- **Token Refresh**: Separate refresh tokens stored in database
- **Role-Based Access Control**: User, Moderator, Admin roles
- **CORS**: Configurable cross-origin resource sharing
- **SQL Injection Protection**: GORM parameterized queries

## Next Steps

This is a foundational implementation. To complete the project, implement:

1. **Movie Management**
   - Movie CRUD handlers
   - TMDB API integration
   - Movie search and filtering

2. **Review System**
   - Review CRUD
   - Comment threading
   - Like/unlike functionality

3. **Social Features**
   - Direct messaging (WebSocket for real-time)
   - Friend system
   - Community chats
   - World chat

4. **Advanced Features**
   - AI content moderation integration
   - Friend recommendation algorithm
   - User title calculation
   - Notification system

5. **Frontend**
   - Build React/Vue/Next.js frontend
   - Connect to API
   - Real-time features with WebSockets

## Contributing

This is a learning project. Feel free to extend it with additional features!

## License

MIT License

---

**Built with ❤️ for learning Go and backend development**
