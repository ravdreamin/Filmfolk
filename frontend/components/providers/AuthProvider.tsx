'use client';

import { useEffect } from 'react';
import { useAuthStore } from '@/store/authStore';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { fetchCurrentUser, accessToken } = useAuthStore();

  useEffect(() => {
    // Fetch current user if we have an access token
    if (accessToken) {
      fetchCurrentUser();
    } else {
      useAuthStore.setState({ isLoading: false });
    }
  }, []);

  return <>{children}</>;
}
