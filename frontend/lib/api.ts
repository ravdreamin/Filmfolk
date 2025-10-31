import axios from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add request interceptor to attach auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Add response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('refreshToken')
        if (refreshToken) {
          const response = await axios.post(
            `${API_BASE_URL}/api/v1/auth/refresh`,
            { refresh_token: refreshToken }
          )

          const { access_token } = response.data
          localStorage.setItem('accessToken', access_token)

          originalRequest.headers.Authorization = `Bearer ${access_token}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        localStorage.removeItem('accessToken')
        localStorage.removeItem('refreshToken')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

// Auth API
export const authApi = {
  register: (data: { username: string; email: string; password: string }) =>
    api.post('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),

  logout: (refreshToken: string) =>
    api.post('/auth/logout', { refresh_token: refreshToken }),

  getCurrentUser: () =>
    api.get('/auth/me'),
}

// Movies API
export const moviesApi = {
  list: (params?: { search?: string; genre?: string; page?: number }) =>
    api.get('/movies', { params }),

  getById: (id: number) =>
    api.get(`/movies/${id}`),

  create: (data: any) =>
    api.post('/movies', data),

  getReviews: (id: number) =>
    api.get(`/movies/${id}/reviews`),
}

// Reviews API
export const reviewsApi = {
  create: (data: { movie_id: number; rating: number; review_text: string }) =>
    api.post('/reviews', data),

  getById: (id: number) =>
    api.get(`/reviews/${id}`),

  update: (id: number, data: { rating?: number; review_text?: string }) =>
    api.put(`/reviews/${id}`, data),

  delete: (id: number) =>
    api.delete(`/reviews/${id}`),

  addComment: (reviewId: number, comment: string, parentId?: number) =>
    api.post('/reviews/comments', {
      review_id: reviewId,
      comment_text: comment,
      parent_comment_id: parentId,
    }),
}

export default api
