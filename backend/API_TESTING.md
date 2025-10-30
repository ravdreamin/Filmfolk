# FilmFolk API Testing Guide

This guide shows you how to test the API using curl or Postman.

## Prerequisites

1. Start the server:
```bash
go run cmd/server/main.go
```

2. Make sure PostgreSQL is running and database is created

## Testing with curl

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "filmfolk-api"
}
```

### 2. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

Expected response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "engagement_score": 0,
    "total_reviews": 0,
    "created_at": "2025-10-31T..."
  },
  "expires_in": 900
}
```

**Save the access_token for next requests!**

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Get Current User (Protected Route)

```bash
# Replace YOUR_ACCESS_TOKEN with the token from register/login
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

Expected response:
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "role": "user"
}
```

### 5. Refresh Access Token

```bash
# When access token expires (after 15 min), use refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

### 6. Logout

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Testing with Postman

### Setup

1. Open Postman
2. Create a new Collection called "FilmFolk API"
3. Add a variable `{{baseUrl}}` = `http://localhost:8080`
4. Add a variable `{{accessToken}}` (will be set automatically)

### Create Requests

#### 1. Health Check
- Method: `GET`
- URL: `{{baseUrl}}/health`

#### 2. Register
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/auth/register`
- Headers:
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```
- Tests tab (to save token automatically):
```javascript
if (pm.response.code === 201) {
    const response = pm.response.json();
    pm.collectionVariables.set("accessToken", response.access_token);
    pm.collectionVariables.set("refreshToken", response.refresh_token);
}
```

#### 3. Login
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/auth/login`
- Headers:
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```
- Tests tab:
```javascript
if (pm.response.code === 200) {
    const response = pm.response.json();
    pm.collectionVariables.set("accessToken", response.access_token);
    pm.collectionVariables.set("refreshToken", response.refresh_token);
}
```

#### 4. Get Current User
- Method: `GET`
- URL: `{{baseUrl}}/api/v1/auth/me`
- Headers:
  - `Authorization`: `Bearer {{accessToken}}`

#### 5. Refresh Token
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/auth/refresh`
- Headers:
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "refresh_token": "{{refreshToken}}"
}
```
- Tests tab:
```javascript
if (pm.response.code === 200) {
    const response = pm.response.json();
    pm.collectionVariables.set("accessToken", response.access_token);
}
```

#### 6. Logout
- Method: `POST`
- URL: `{{baseUrl}}/api/v1/auth/logout`
- Headers:
  - `Content-Type`: `application/json`
- Body (raw JSON):
```json
{
  "refresh_token": "{{refreshToken}}"
}
```

## Error Testing

### Test Invalid Credentials
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }'
```

Expected: `401 Unauthorized`

### Test Duplicate Registration
```bash
# Register same email twice
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

Expected: `400 Bad Request` with error "email already registered"

### Test Missing Authorization Header
```bash
curl http://localhost:8080/api/v1/auth/me
```

Expected: `401 Unauthorized` with error "Authorization header required"

### Test Invalid Token
```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer invalid_token_here"
```

Expected: `401 Unauthorized` with error "Invalid or expired token"

## Common HTTP Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created (registration)
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication failed
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Tips

1. **Tokens expire** - Access tokens last 15 minutes, use refresh token to get new one
2. **Case sensitive** - Usernames and emails are case-sensitive
3. **Password requirements** - Minimum 8 characters
4. **Username requirements** - 3-50 characters
5. **Save tokens** - Store access and refresh tokens to maintain session

## Next Steps

Once basic auth is working, you can extend testing to:
- Movie management endpoints (when implemented)
- Review CRUD operations (when implemented)
- Social features (when implemented)

## Troubleshooting

### "Database connection failed"
- Check if PostgreSQL is running
- Verify database exists: `psql -l | grep filmfolk`
- Check credentials in `config.yaml`

### "Configuration loading error"
- Ensure `configs/config.yaml` exists
- Or create `.env` file with required variables

### "Invalid token"
- Token may have expired (15 min)
- Use refresh token to get new access token
- Or login again

### "Email already registered"
- User already exists
- Use a different email
- Or test login instead

---

Happy testing! 🚀
