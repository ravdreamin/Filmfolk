# Google OAuth Setup Guide

This guide will help you set up Google OAuth authentication for FilmFolk.

## Prerequisites

- A Google account
- Access to the Google Cloud Console

## Steps

### 1. Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click "Select a project" at the top
3. Click "NEW PROJECT"
4. Enter project name: `FilmFolk` (or your preferred name)
5. Click "CREATE"

### 2. Enable Google+ API

1. In your project, go to "APIs & Services" > "Library"
2. Search for "Google+ API"
3. Click on it and press "ENABLE"

### 3. Configure OAuth Consent Screen

1. Go to "APIs & Services" > "OAuth consent screen"
2. Select "External" user type
3. Click "CREATE"
4. Fill in the required fields:
   - App name: `FilmFolk`
   - User support email: Your email
   - Developer contact email: Your email
5. Click "SAVE AND CONTINUE"
6. Add scopes (click "ADD OR REMOVE SCOPES"):
   - `.../auth/userinfo.email`
   - `.../auth/userinfo.profile`
7. Click "SAVE AND CONTINUE"
8. Add test users if needed (for development)
9. Click "SAVE AND CONTINUE"

### 4. Create OAuth 2.0 Credentials

1. Go to "APIs & Services" > "Credentials"
2. Click "CREATE CREDENTIALS" > "OAuth client ID"
3. Select "Web application"
4. Configure:
   - Name: `FilmFolk Web Client`
   - Authorized JavaScript origins:
     - `http://localhost:3000` (for development)
     - `http://localhost:8080` (for development)
   - Authorized redirect URIs:
     - `http://localhost:8080/api/v1/auth/google/callback`
5. Click "CREATE"
6. Copy the **Client ID** and **Client Secret**

### 5. Update Environment Variables

1. Open your `.env` file in the `backend` directory
2. Add your credentials:

```env
GOOGLE_CLIENT_ID=your_client_id_here
GOOGLE_CLIENT_SECRET=your_client_secret_here
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

### 6. Production Setup

For production, you'll need to:

1. Update the OAuth consent screen to "Production" status
2. Add your production domain to authorized origins
3. Update the redirect URI with your production domain
4. Update the environment variables with production URLs

Example production `.env`:

```env
GOOGLE_CLIENT_ID=your_production_client_id
GOOGLE_CLIENT_SECRET=your_production_client_secret
GOOGLE_REDIRECT_URL=https://api.yourproductiondomain.com/api/v1/auth/google/callback
```

## Testing

1. Start your backend server
2. Navigate to `http://localhost:3000/login`
3. Click "Continue with Google"
4. You should be redirected to Google's login page
5. After successful authentication, you'll be redirected back to your app

## Troubleshooting

### Error: redirect_uri_mismatch

- Ensure the redirect URI in your `.env` file exactly matches the one configured in Google Cloud Console
- Check that you're using the correct protocol (http vs https)

### Error: invalid_client

- Verify your Client ID and Client Secret are correct
- Ensure there are no extra spaces in your `.env` file

### Error: access_denied

- Check that your email is added as a test user if the app is in development mode
- Verify the OAuth consent screen is properly configured

## Security Notes

- Never commit your `.env` file with real credentials
- Use environment variables for all sensitive data
- In production, use HTTPS for all OAuth redirects
- Regularly rotate your client secrets
- Limit scopes to only what you need
