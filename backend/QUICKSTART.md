# FilmFolk - Quick Start Guide

Get up and running in 5 minutes!

## Prerequisites

- Go 1.25.2 or higher installed
- PostgreSQL 14+ installed and running
- Terminal/command line access

## Setup Steps

### 1. Database Setup

```bash
# Create the database
createdb filmfolk

# Run the schema migration
psql filmfolk < migrations/001_initial_schema.sql
```

### 2. Configuration

Edit `configs/config.yaml` with your database credentials:

```yaml
db:
  host: localhost
  port: 5432              # Your PostgreSQL port
  user: your_username     # Your PostgreSQL user
  password: ""            # Your PostgreSQL password (if any)
  dbname: filmfolk
  sslmode: disable
```

### 3. Run the Server

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

You should see:
```
Loading configuration...
Configuration loaded from: YAML
Environment: development
Connecting to database...
Database connection established successfully
Running auto-migrations...
Auto-migrations completed successfully
Setting up routes...
🚀 FilmFolk server starting on http://localhost:8080
📝 API docs will be at http://localhost:8080/api/v1
❤️  Health check: http://localhost:8080/health
```

## Test It Works

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"service":"filmfolk-api","status":"ok"}
```

### 2. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

You'll get back an access token and user info!

### 3. Save Your Token

Copy the `access_token` from the response. You'll need it for authenticated requests.

```bash
# Set as environment variable (easier for testing)
export TOKEN="your_access_token_here"
```

### 4. Create a Movie

```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Shawshank Redemption",
    "release_year": 1994,
    "genres": ["Drama", "Crime"],
    "summary": "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
    "runtime_minutes": 142
  }'
```

### 5. Write a Review

```bash
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "rating": 10,
    "review_text": "One of the greatest films ever made! The story, acting, and cinematography are all perfect. A masterpiece that deserves every bit of praise it gets."
  }'
```

### 6. Browse Movies

```bash
# No authentication needed for browsing!
curl http://localhost:8080/api/v1/movies

# Search for a specific movie
curl "http://localhost:8080/api/v1/movies?search=shawshank"

# Get movie details
curl http://localhost:8080/api/v1/movies/1

# Get movie reviews
curl http://localhost:8080/api/v1/movies/1/reviews
```

## Common Issues

### "Database connection failed"

**Solution:** Check PostgreSQL is running and credentials in config.yaml are correct.

```bash
# Check if PostgreSQL is running
pg_ctl status

# Test connection manually
psql -U your_username -d filmfolk -c "SELECT 1;"
```

### "Movie already exists"

**Solution:** This is expected if you created the same movie twice. Try a different movie title.

### "Configuration loading error"

**Solution:** Make sure `configs/config.yaml` exists and is valid YAML.

### Token expired

**Solution:** Tokens expire after 15 minutes. Register/login again or use the refresh token endpoint.

## What's Next?

1. **Explore the API** - See [API_ENDPOINTS.md](API_ENDPOINTS.md) for all available endpoints
2. **Test with Postman** - See [API_TESTING.md](API_TESTING.md) for Postman setup
3. **Read the docs** - See [README.md](README.md) for full documentation
4. **Check project status** - See [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) for what's implemented

## Need Help?

- Check the logs in your terminal where the server is running
- All errors are returned in JSON format with descriptive messages
- Use the health check endpoint to verify the server is running

## Pro Tips

1. **Use Postman** for easier testing - import the collection from API_TESTING.md
2. **Set TOKEN as environment variable** to avoid copying it every time
3. **Check server logs** - they show all SQL queries in development mode
4. **Try the optional query parameters** on list endpoints (page, page_size, search, etc.)

---

**You're all set! Start building your movie review empire! 🎬🍿**
