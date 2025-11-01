import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

// API Client Configuration
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Enable cookies for CSRF tokens
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('accessToken') : null;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If 401 and not already retrying, try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = typeof window !== 'undefined' ? localStorage.getItem('refreshToken') : null;
        if (refreshToken) {
          const response = await axios.post(`${API_URL}/auth/refresh`, { refresh_token: refreshToken });
          const { access_token } = response.data;

          if (typeof window !== 'undefined') {
            localStorage.setItem('accessToken', access_token);
          }

          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return apiClient(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        if (typeof window !== 'undefined') {
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          window.location.href = '/login';
        }
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Types
export interface User {
  id: number;
  username: string;
  email: string;
  avatar_url?: string;
  bio?: string;
  followers_count: number;
  following_count: number;
  created_at: string;
}

export interface Movie {
  id: number;
  title: string;
  release_year: number;
  genres: string[];
  summary?: string;
  poster_url?: string;
  backdrop_url?: string;
  runtime_minutes?: number;
  language?: string;
  average_rating?: number;
  total_reviews: number;
  created_at: string;
}

export interface Review {
  id: number;
  user_id: number;
  movie_id: number;
  rating: number;
  review_text: string;
  is_thread_locked: boolean;
  likes_count: number;
  comments_count: number;
  created_at: string;
  updated_at: string;
  user?: User;
  movie?: Movie;
}

export interface ReviewComment {
  id: number;
  review_id: number;
  user_id: number;
  parent_comment_id?: number;
  comment_text: string;
  likes_count: number;
  created_at: string;
  updated_at: string;
  user?: User;
  replies?: ReviewComment[];
}

// Auth API
export const authAPI = {
  register: (data: { username: string; email: string; password: string }) =>
    apiClient.post('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    apiClient.post('/auth/login', data),

  logout: () => apiClient.post('/auth/logout'),

  refreshToken: (refreshToken: string) =>
    apiClient.post('/auth/refresh', { refresh_token: refreshToken }),

  getCurrentUser: () => apiClient.get<User>('/auth/me'),

  googleLogin: () => {
    window.location.href = `${API_URL}/auth/google`;
  },
};

// Movie API
export const movieAPI = {
  list: (params?: {
    genre?: string;
    year?: number;
    search?: string;
    sort_by?: string;
    page?: number;
    page_size?: number;
  }) => apiClient.get<{ movies: Movie[]; total: number; page: number; page_size: number }>('/movies', { params }),

  get: (id: number) => apiClient.get<Movie>(`/movies/${id}`),

  update: (id: number, data: Partial<Movie>) =>
    apiClient.put<Movie>(`/movies/${id}`, data),

  getReviews: (movieId: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<{ reviews: Review[]; total: number; page: number; page_size: number }>(
      `/movies/${movieId}/reviews`,
      { params }
    ),
};

// Review API
export const reviewAPI = {
  create: (data: { movie_id: number; rating: number; review_text: string }) =>
    apiClient.post<Review>('/reviews', data),

  get: (id: number) => apiClient.get<Review>(`/reviews/${id}`),

  update: (id: number, data: { rating?: number; review_text?: string }) =>
    apiClient.put<Review>(`/reviews/${id}`, data),

  delete: (id: number) => apiClient.delete(`/reviews/${id}`),

  lockThread: (id: number) => apiClient.post(`/reviews/${id}/lock`),

  unlockThread: (id: number) => apiClient.post(`/reviews/${id}/unlock`),

  createComment: (data: { review_id: number; comment_text: string; parent_comment_id?: number }) =>
    apiClient.post<ReviewComment>('/reviews/comments', data),

  deleteComment: (id: number) => apiClient.delete(`/reviews/comments/${id}`),
};

// Follower API
export const followerAPI = {
  follow: (userId: number) => apiClient.post(`/users/${userId}/follow`),

  unfollow: (userId: number) => apiClient.delete(`/users/${userId}/follow`),

  checkStatus: (userId: number) =>
    apiClient.get<{ is_following: boolean }>(`/users/${userId}/follow/status`),

  getFollowers: (userId: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<{ followers: User[]; total: number; page: number; page_size: number }>(
      `/users/${userId}/followers`,
      { params }
    ),

  getFollowing: (userId: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<{ following: User[]; total: number; page: number; page_size: number }>(
      `/users/${userId}/following`,
      { params }
    ),
};

export default apiClient;
