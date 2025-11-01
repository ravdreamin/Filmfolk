# Authentication & UI Update

This document outlines the recent updates to FilmFolk's authentication system and user interface.

## Changes Summary

### 1. Minimalist Logo Design

**Location**: `frontend/components/ui/Logo.tsx`

The logo has been redesigned with a minimalist aesthetic inspired by Apple, Netflix, and Google:

- Clean typography-based design
- Split "FilmFolk" into "Film" (bold) and "Folk" (light) for visual hierarchy
- Subtle hover animation with underline effect
- Supports both light and dark themes
- Removed colorful gradients for a cleaner, professional look

**Usage**:
```tsx
import Logo from '@/components/ui/Logo'

// For dark backgrounds (default)
<Logo />

// For light backgrounds
<Logo theme="light" />
```

### 2. Redesigned Login & Register Pages

**Locations**:
- `frontend/app/(auth)/login/page.tsx`
- `frontend/app/(auth)/register/page.tsx`

The authentication pages have been completely redesigned with a minimalist aesthetic:

**Key Features**:
- Clean white background with ample whitespace
- Simple, borderless inputs with subtle focus states
- Black and white color scheme
- Smooth micro-animations
- Google OAuth integration button
- Minimal visual noise

**Design Principles**:
- Typography-first approach
- Generous spacing
- Subtle transitions
- Focus on usability
- Professional, modern aesthetic

### 3. Google OAuth Integration

#### Backend Implementation

**New Files**:
- `backend/internal/services/oauth_service.go` - OAuth service logic
- `backend/internal/handlers/oauth_handler.go` - OAuth HTTP handlers
- `backend/OAUTH_SETUP.md` - Setup guide for Google OAuth

**Updated Files**:
- `backend/internal/routes/routes.go` - Added OAuth routes
- `backend/.env.example` - Added OAuth environment variables

**New Endpoints**:
- `GET /api/v1/auth/google` - Initiates Google OAuth flow
- `GET /api/v1/auth/google/callback` - Handles OAuth callback

**Features**:
- Secure OAuth 2.0 flow with state verification
- Automatic user creation for new Google users
- Username generation from Google profile
- Support for existing users logging in with Google
- CSRF protection with state parameter

#### Frontend Implementation

**New Files**:
- `frontend/app/auth/callback/page.tsx` - OAuth callback handler

**Updated Files**:
- `frontend/app/(auth)/login/page.tsx` - Added Google login button
- `frontend/app/(auth)/register/page.tsx` - Added Google signup button
- `frontend/store/authStore.ts` - Added `setTokens` method for OAuth
- `frontend/components/ui/Logo.tsx` - Added theme support

**Features**:
- One-click Google authentication
- Secure token handling
- Automatic redirect after successful authentication
- Error handling and user feedback
- Loading states

### 4. Enhanced Auth Store

**Location**: `frontend/store/authStore.ts`

Added TypeScript types and new methods:

**New Types**:
```typescript
interface User {
  id: number;
  username: string;
  email: string;
  role: string;
}

interface AuthState {
  // ... existing fields
  setTokens: (accessToken: string, refreshToken: string) => Promise<void>;
}
```

**New Method**:
- `setTokens()` - Handles OAuth token storage and user info extraction

## Setup Instructions

### Google OAuth Setup

1. Follow the guide in `backend/OAUTH_SETUP.md`
2. Create a Google Cloud Project
3. Configure OAuth consent screen
4. Create OAuth 2.0 credentials
5. Update `.env` file with credentials

### Environment Variables

**Backend** (`backend/.env`):
```env
GOOGLE_CLIENT_ID=your_client_id_here
GOOGLE_CLIENT_SECRET=your_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

**Frontend** (`frontend/.env.local`):
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Running the Application

1. **Start Backend**:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **Start Frontend**:
   ```bash
   cd frontend
   npm run dev
   ```

3. Navigate to `http://localhost:3000/login`

## User Flow

### Google OAuth Login Flow

1. User clicks "Continue with Google" on login/register page
2. User is redirected to Google's authentication page
3. User grants permission
4. Google redirects back to `/api/v1/auth/google/callback`
5. Backend exchanges code for user info
6. Backend creates or finds user in database
7. Backend generates JWT tokens
8. Backend redirects to frontend `/auth/callback` with tokens
9. Frontend stores tokens and redirects to `/movies`

### Traditional Email/Password Flow

1. User enters email and password
2. Frontend sends credentials to backend
3. Backend validates and returns tokens
4. Frontend stores tokens and redirects to `/movies`

## Security Considerations

1. **OAuth State Parameter**: CSRF protection using random state
2. **Secure Cookies**: HttpOnly cookies for state storage
3. **Token Validation**: JWT verification on every request
4. **HTTPS Required**: Production must use HTTPS
5. **Scope Limitation**: Only request necessary Google scopes

## Design System

### Colors
- **Primary**: Black (`#000000`)
- **Background**: White (`#FFFFFF`)
- **Border**: Gray-300 (`#D1D5DB`)
- **Text**: Black/Gray-600

### Typography
- **Headings**: Light weight (300-400)
- **Body**: Normal weight (400)
- **Buttons**: Normal weight (400)

### Spacing
- Generous padding and margins
- Consistent 4px grid system
- Ample whitespace

### Animations
- Subtle scale transforms (1.01-1.02)
- Smooth transitions (0.2-0.4s)
- Ease-out timing function

## Future Enhancements

1. **Additional OAuth Providers**:
   - Facebook
   - GitHub
   - Apple

2. **Two-Factor Authentication**:
   - TOTP support
   - SMS verification

3. **Password Reset**:
   - Email-based reset flow
   - Secure token generation

4. **Session Management**:
   - Device tracking
   - Active session view
   - Remote logout

## Dependencies Added

### Backend
- `golang.org/x/oauth2` - OAuth 2.0 client
- `google.golang.org/api/oauth2/v2` - Google OAuth API

### Frontend
No new dependencies - using existing libraries

## Migration Notes

For existing users:
- Email/password authentication continues to work
- Users can link Google account to existing email account
- No breaking changes to existing authentication flow

## Testing

### Manual Testing Checklist

- [ ] Login with email/password
- [ ] Register with email/password
- [ ] Login with Google (new user)
- [ ] Login with Google (existing user)
- [ ] Logout functionality
- [ ] Token refresh
- [ ] Protected routes
- [ ] Error handling

### Test Accounts

For development, add test users in Google Cloud Console OAuth consent screen.

## Troubleshooting

### Common Issues

1. **"redirect_uri_mismatch"**
   - Check redirect URL matches exactly in Google Console and `.env`

2. **"invalid_client"**
   - Verify Client ID and Secret are correct
   - Check for whitespace in `.env` file

3. **OAuth callback not working**
   - Ensure backend is running on correct port
   - Check CORS settings in backend

4. **Tokens not saving**
   - Check browser localStorage
   - Verify authStore is properly initialized

## Support

For questions or issues:
1. Check `backend/OAUTH_SETUP.md` for Google OAuth setup
2. Review error messages in browser console
3. Check backend logs for authentication errors
