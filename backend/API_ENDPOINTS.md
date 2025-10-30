# FilmFolk API Endpoints

Complete API documentation for all implemented endpoints.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

All authenticated endpoints require the Authorization header:
```
Authorization: Bearer <access_token>
```

---

## Auth Endpoints

### Register
`POST /auth/register`

Create a new user account.

**Request:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:** `201 Created`
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user"
  },
  "expires_in": 900
}
```

### Login
`POST /auth/login`

Authenticate with email and password.

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:** Same as Register

### Refresh Token
`POST /auth/refresh`

Get a new access token using refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

### Logout
`POST /auth/logout`

Revoke a refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

### Get Current User
`GET /auth/me`

Get currently authenticated user info.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user"
}
```

---

## Movie Endpoints

### List Movies
`GET /movies`

List and search movies with filtering.

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `page_size` (int): Items per page (default: 20, max: 100)
- `status` (string): Filter by status (pending_approval, approved, rejected)
- `genre` (string): Filter by genre
- `year` (int): Filter by release year
- `search` (string): Search by title
- `sort_by` (string): Sort order (rating, year, title, reviews)

**Example:**
```
GET /movies?search=inception&sort_by=rating&page=1
```

**Response:**
```json
{
  "movies": [
    {
      "id": 1,
      "title": "Inception",
      "release_year": 2010,
      "genres": ["Action", "Sci-Fi", "Thriller"],
      "summary": "A thief who steals corporate secrets...",
      "poster_url": "https://...",
      "average_rating": 8.5,
      "total_reviews": 1250,
      "status": "approved"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20
}
```

### Get Movie
`GET /movies/:id`

Get detailed information about a specific movie.

**Response:**
```json
{
  "id": 1,
  "title": "Inception",
  "release_year": 2010,
  "genres": ["Action", "Sci-Fi"],
  "summary": "...",
  "poster_url": "...",
  "backdrop_url": "...",
  "runtime_minutes": 148,
  "language": "en",
  "tmdb_id": 27205,
  "imdb_id": "tt1375666",
  "average_rating": 8.5,
  "total_reviews": 1250,
  "status": "approved",
  "created_at": "2025-01-15T10:00:00Z"
}
```

### Create Movie
`POST /movies` 🔒 **Authenticated**

Submit a new movie for approval.

**Request:**
```json
{
  "title": "The Matrix",
  "release_year": 1999,
  "genres": ["Action", "Sci-Fi"],
  "summary": "A computer hacker learns...",
  "poster_url": "https://...",
  "backdrop_url": "https://...",
  "runtime_minutes": 136,
  "language": "en",
  "tmdb_id": 603,
  "imdb_id": "tt0133093"
}
```

**Response:** `201 Created`

### Update Movie
`PUT /movies/:id` 🔒 **Moderator/Admin**

Update movie information.

**Request:**
```json
{
  "title": "Updated Title",
  "summary": "Updated summary..."
}
```

---

## Review Endpoints

### Get Movie Reviews
`GET /movies/:id/reviews`

Get all reviews for a movie.

**Query Parameters:**
- `page` (int): Page number
- `page_size` (int): Items per page (max: 50)

**Response:**
```json
{
  "reviews": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "johndoe",
        "avatar_url": "..."
      },
      "movie_id": 1,
      "rating": 9,
      "review_text": "Amazing movie! The plot twists...",
      "sentiment": "positive",
      "likes_count": 45,
      "comments_count": 12,
      "is_thread_locked": false,
      "created_at": "2025-01-15T10:00:00Z"
    }
  ],
  "total": 1250,
  "page": 1,
  "page_size": 20
}
```

### Get Single Review
`GET /reviews/:id`

Get a review with all its comments (threaded).

**Response:**
```json
{
  "id": 1,
  "user": {...},
  "movie": {...},
  "rating": 9,
  "review_text": "...",
  "comments": [
    {
      "id": 1,
      "user": {...},
      "comment_text": "Great review!",
      "likes_count": 5,
      "replies": [
        {
          "id": 2,
          "user": {...},
          "comment_text": "I agree!",
          "parent_comment_id": 1
        }
      ],
      "created_at": "..."
    }
  ],
  "created_at": "..."
}
```

### Create Review
`POST /reviews` 🔒 **Authenticated**

Write a review for a movie.

**Request:**
```json
{
  "movie_id": 1,
  "rating": 9,
  "review_text": "This movie was absolutely incredible! The cinematography..."
}
```

**Response:** `201 Created`

**Constraints:**
- One review per user per movie
- Rating: 1-10
- Review text: minimum 10 characters
- Cannot review unapproved movies

### Update Review
`PUT /reviews/:id` 🔒 **Authenticated** (Own review only)

Update your own review.

**Request:**
```json
{
  "rating": 8,
  "review_text": "Updated my opinion after rewatching..."
}
```

### Delete Review
`DELETE /reviews/:id` 🔒 **Authenticated/Moderator**

Delete a review. Users can delete their own reviews, moderators can delete any review.

**Response:**
```json
{
  "message": "Review deleted successfully"
}
```

### Lock Review Thread
`POST /reviews/:id/lock` 🔒 **Authenticated** (Review author only)

Lock your review thread to prevent further comments.

**Response:**
```json
{
  "message": "Thread locked successfully"
}
```

### Unlock Review Thread
`POST /reviews/:id/unlock` 🔒 **Authenticated** (Review author only)

Unlock your review thread to allow comments again.

### Create Comment
`POST /reviews/comments` 🔒 **Authenticated**

Add a comment to a review or reply to another comment.

**Request:**
```json
{
  "review_id": 1,
  "comment_text": "Great review!"
}
```

**For nested reply:**
```json
{
  "review_id": 1,
  "parent_comment_id": 5,
  "comment_text": "I completely agree with this point!"
}
```

**Response:** `201 Created`

**Constraints:**
- Cannot comment on locked threads
- Parent comment must belong to the same review

### Delete Comment
`DELETE /reviews/comments/:id` 🔒 **Authenticated/Moderator**

Delete a comment. Users can delete their own comments, moderators can delete any.

---

## Moderator Endpoints

All moderator endpoints require `moderator` or `admin` role.

### Get Pending Movies
`GET /moderator/movies/pending` 🔒 **Moderator**

Get list of movies awaiting approval.

**Query Parameters:**
- `page` (int): Page number

**Response:**
```json
{
  "movies": [
    {
      "id": 5,
      "title": "New Movie",
      "status": "pending_approval",
      "submitted_by": {
        "id": 10,
        "username": "moviefan"
      },
      "created_at": "..."
    }
  ],
  "total": 25,
  "page": 1
}
```

### Approve Movie
`POST /moderator/movies/:id/approve` 🔒 **Moderator**

Approve a pending movie submission.

**Response:**
```json
{
  "id": 5,
  "title": "New Movie",
  "status": "approved",
  "approved_by": {
    "id": 2,
    "username": "moderator1"
  }
}
```

### Reject Movie
`POST /moderator/movies/:id/reject` 🔒 **Moderator**

Reject a pending movie submission.

**Response:**
```json
{
  "message": "Movie rejected"
}
```

---

## Admin Endpoints

All admin endpoints require `admin` role.

### Delete Movie
`DELETE /admin/movies/:id` 🔒 **Admin**

Permanently delete a movie and all associated data.

**Response:**
```json
{
  "message": "Movie deleted successfully"
}
```

**Warning:** This cascades to reviews, comments, etc.

---

## Error Responses

All endpoints return consistent error format:

### 400 Bad Request
```json
{
  "error": "Validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 403 Forbidden
```json
{
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Rate Limiting

Currently not implemented. Recommended for production:
- Auth endpoints: 5 requests/minute
- Read endpoints: 100 requests/minute
- Write endpoints: 20 requests/minute

---

## Testing Workflow

### 1. Create Account & Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'

# Save the access_token from response
TOKEN="eyJhbGci..."
```

### 2. Submit a Movie
```bash
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Shawshank Redemption",
    "release_year": 1994,
    "genres": ["Drama"],
    "summary": "Two imprisoned men bond over..."
  }'
```

### 3. Write a Review
```bash
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "rating": 10,
    "review_text": "One of the greatest films ever made!"
  }'
```

### 4. Browse Movies
```bash
# No auth needed for browsing
curl http://localhost:8080/api/v1/movies?sort_by=rating

# Get specific movie
curl http://localhost:8080/api/v1/movies/1

# Get movie reviews
curl http://localhost:8080/api/v1/movies/1/reviews
```

---

## Next Features (TODO)

The following features are designed but not yet implemented:

- User Lists (watched, plan to watch, etc.)
- Direct Messaging
- Friend System
- Communities
- World Chat (WebSocket)
- Notifications
- OAuth Login (Google, Facebook, etc.)
- AI Content Moderation
- Sentiment Analysis
- User Gamification
- Friend Recommendations

See code comments marked with `TODO` for integration points.
