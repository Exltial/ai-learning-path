import type {
  ApiResponse,
  User,
  LoginData,
  RegisterData,
  AuthResponse,
  Course,
  CourseFilterParams,
  PaginatedResponse,
  CourseProgress,
  Lesson,
} from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

class ApiError extends Error {
  constructor(public status: number, public message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

class ApiService {
  private token: string | null = null

  constructor() {
    this.token = localStorage.getItem('auth_token')
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    }
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }
    return headers
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new ApiError(response.status, errorData.message || '请求失败')
    }
    const data = await response.json()
    return data
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`
    const response = await fetch(url, {
      ...options,
      headers: {
        ...this.getHeaders(),
        ...options?.headers,
      },
    })
    return this.handleResponse<T>(response)
  }

  // Auth APIs
  async register(data: RegisterData): Promise<ApiResponse<AuthResponse>> {
    try {
      const response = await this.request<ApiResponse<AuthResponse>>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      if (response.success && response.data?.token) {
        this.token = response.data.token
        localStorage.setItem('auth_token', response.data.token)
      }
      return response
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async login(data: LoginData): Promise<ApiResponse<AuthResponse>> {
    try {
      const response = await this.request<ApiResponse<AuthResponse>>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      if (response.success && response.data?.token) {
        this.token = response.data.token
        localStorage.setItem('auth_token', response.data.token)
      }
      return response
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async logout(): Promise<void> {
    this.token = null
    localStorage.removeItem('auth_token')
    try {
      await this.request('/auth/logout', { method: 'POST' })
    } catch {
      // Ignore logout API errors
    }
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    try {
      return await this.request<ApiResponse<User>>('/auth/me')
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  // Course APIs
  async getCourses(params?: CourseFilterParams): Promise<ApiResponse<PaginatedResponse<Course>>> {
    try {
      const queryParams = new URLSearchParams()
      if (params?.search) queryParams.append('search', params.search)
      if (params?.level) queryParams.append('level', params.level)
      if (params?.tags) params.tags.forEach(tag => queryParams.append('tags', tag))
      if (params?.page) queryParams.append('page', params.page.toString())
      if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString())

      const query = queryParams.toString()
      const url = `/courses${query ? `?${query}` : ''}`
      return await this.request<ApiResponse<PaginatedResponse<Course>>>(url)
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async getCourse(id: string): Promise<ApiResponse<Course>> {
    try {
      return await this.request<ApiResponse<Course>>(`/courses/${id}`)
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async getCourseLessons(courseId: string): Promise<ApiResponse<Lesson[]>> {
    try {
      return await this.request<ApiResponse<Lesson[]>>(`/courses/${courseId}/lessons`)
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  // Progress APIs
  async getProgress(): Promise<ApiResponse<CourseProgress[]>> {
    try {
      return await this.request<ApiResponse<CourseProgress[]>>('/progress')
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async updateProgress(courseId: string, lessonId: string): Promise<ApiResponse<CourseProgress>> {
    try {
      return await this.request<ApiResponse<CourseProgress>>('/progress', {
        method: 'POST',
        body: JSON.stringify({ courseId, lessonId }),
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  // Discussion APIs
  async getDiscussions(courseId: string, params?: {
    page?: number
    limit?: number
    sort_by?: string
    search?: string
    tag_ids?: string[]
  }): Promise<any> {
    try {
      const queryParams = new URLSearchParams({ course_id: courseId })
      if (params?.page) queryParams.append('page', params.page.toString())
      if (params?.limit) queryParams.append('limit', params.limit.toString())
      if (params?.sort_by) queryParams.append('sort_by', params.sort_by)
      if (params?.search) queryParams.append('search', params.search)
      if (params?.tag_ids) params.tag_ids.forEach(tagId => queryParams.append('tag_ids', tagId))

      const query = queryParams.toString()
      const url = `/discussions${query ? `?${query}` : ''}`
      return await this.request(url)
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async getDiscussion(id: string, withReplies?: boolean): Promise<any> {
    try {
      const params = withReplies ? '?with_replies=true&max_depth=10' : ''
      return await this.request(`/discussions/${id}${params}`)
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async createDiscussion(data: {
    course_id: string
    lesson_id?: string
    title?: string
    content: string
    parent_id?: string
    is_anonymous?: boolean
    tag_ids?: string[]
  }): Promise<any> {
    try {
      return await this.request('/discussions', {
        method: 'POST',
        body: JSON.stringify(data),
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async updateDiscussion(id: string, data: {
    title?: string
    content?: string
    is_resolved?: boolean
    is_locked?: boolean
    is_pinned?: boolean
    tag_ids?: string[]
  }): Promise<any> {
    try {
      return await this.request(`/discussions/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async deleteDiscussion(id: string): Promise<any> {
    try {
      return await this.request(`/discussions/${id}`, {
        method: 'DELETE',
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async toggleLike(id: string, likeType: 'upvote' | 'downvote' = 'upvote'): Promise<any> {
    try {
      return await this.request(`/discussions/${id}/like?like_type=${likeType}`, {
        method: 'POST',
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async toggleFavorite(id: string): Promise<any> {
    try {
      return await this.request(`/discussions/${id}/favorite`, {
        method: 'POST',
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async resolveDiscussion(id: string): Promise<any> {
    try {
      return await this.request(`/discussions/${id}/resolve`, {
        method: 'POST',
      })
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async getDiscussionTags(): Promise<any> {
    try {
      return await this.request('/discussions/tags')
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }

  async getHotDiscussions(courseId: string, limit?: number): Promise<any> {
    try {
      const params = limit ? `?limit=${limit}` : ''
      return await this.request(`/discussions/hot?course_id=${courseId}${params}`)
    } catch (error) {
      if (error instanceof ApiError) {
        return { success: false, error: error.message }
      }
      return { success: false, error: '网络错误，请稍后重试' }
    }
  }
}

export const api = new ApiService()
export { ApiError }
