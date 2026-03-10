// User types
export interface User {
  id: string
  username: string
  email: string
  avatar?: string
  createdAt: string
}

export interface RegisterData {
  username: string
  email: string
  password: string
  confirmPassword: string
}

export interface LoginData {
  email: string
  password: string
}

export interface AuthResponse {
  user: User
  token: string
}

// Course types
export interface Course {
  id: string
  title: string
  description: string
  thumbnail: string
  level: 'beginner' | 'intermediate' | 'advanced'
  duration: string
  students: number
  rating: number
  lessons: number
  tags: string[]
  instructor?: Instructor
  whatYouLearn?: string[]
  requirements?: string[]
}

export interface Instructor {
  name: string
  title: string
  avatar: string
}

export interface Lesson {
  id: string
  title: string
  duration: string
  completed: boolean
  locked: boolean
  type: 'video' | 'exercise' | 'quiz'
}

export interface CourseProgress {
  courseId: string
  completedLessons: number
  totalLessons: number
  progress: number
  lastAccessedAt: string
}

// API Response types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface CourseFilterParams {
  search?: string
  level?: string
  tags?: string[]
  page?: number
  pageSize?: number
}
