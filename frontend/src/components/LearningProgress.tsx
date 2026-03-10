import { useState, useEffect } from 'react'
import { TrendingUp, Clock, Award, BookOpen, ChevronRight, Calendar } from 'lucide-react'
import ProgressBar from './ProgressBar'
import { api } from '@/services/api'
import type { CourseProgress } from '@/types'

interface LearningProgressProps {
  userId?: string
  compact?: boolean
}

interface CourseWithProgress {
  id: string
  title: string
  thumbnail: string
  progress: number
  totalLessons: number
  completedLessons: number
  lastAccessedAt: string
}

const mockProgressData: CourseWithProgress[] = [
  {
    id: '1',
    title: 'Python 编程基础',
    thumbnail: 'https://picsum.photos/seed/python/400/200',
    progress: 45,
    totalLessons: 24,
    completedLessons: 11,
    lastAccessedAt: '2024-03-10T10:30:00Z',
  },
  {
    id: '2',
    title: 'React 前端开发',
    thumbnail: 'https://picsum.photos/seed/react/400/200',
    progress: 78,
    totalLessons: 36,
    completedLessons: 28,
    lastAccessedAt: '2024-03-09T15:20:00Z',
  },
  {
    id: '3',
    title: '数据结构与算法',
    thumbnail: 'https://picsum.photos/seed/algo/400/200',
    progress: 23,
    totalLessons: 50,
    completedLessons: 12,
    lastAccessedAt: '2024-03-08T09:15:00Z',
  },
]

export default function LearningProgress({ compact = false }: LearningProgressProps) {
  const [progressData, setProgressData] = useState<CourseWithProgress[]>([])
  const [loading, setLoading] = useState(true)
  const [totalStats, setTotalStats] = useState({
    totalCourses: 0,
    completedCourses: 0,
    totalLessons: 0,
    completedLessons: 0,
    totalHours: 0,
  })

  useEffect(() => {
    loadProgressData()
  }, [])

  const loadProgressData = async () => {
    setLoading(true)
    try {
      const response = await api.getProgress()
      if (response.success && response.data) {
        // Transform API data to match our interface
        const transformedData = response.data.map((p: CourseProgress) => ({
          id: p.courseId,
          title: `Course ${p.courseId}`,
          thumbnail: `https://picsum.photos/seed/${p.courseId}/400/200`,
          progress: p.progress,
          totalLessons: p.totalLessons,
          completedLessons: p.completedLessons,
          lastAccessedAt: p.lastAccessedAt,
        }))
        setProgressData(transformedData)
      } else {
        setProgressData(mockProgressData)
      }
    } catch {
      setProgressData(mockProgressData)
    }
    setLoading(false)
  }

  useEffect(() => {
    if (progressData.length > 0) {
      const totalCourses = progressData.length
      const completedCourses = progressData.filter((p) => p.progress === 100).length
      const totalLessons = progressData.reduce((sum, p) => sum + p.totalLessons, 0)
      const completedLessons = progressData.reduce((sum, p) => sum + p.completedLessons, 0)
      const totalHours = Math.round((completedLessons / totalLessons) * 20) // Estimate

      setTotalStats({
        totalCourses,
        completedCourses,
        totalLessons,
        completedLessons,
        totalHours,
      })
    }
  }, [progressData])

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

    if (diffDays === 0) {
      return '今天'
    } else if (diffDays === 1) {
      return '昨天'
    } else if (diffDays < 7) {
      return `${diffDays}天前`
    } else {
      return date.toLocaleDateString('zh-CN')
    }
  }

  if (loading) {
    return (
      <div className="animate-pulse space-y-4">
        <div className="h-32 bg-secondary-200 dark:bg-secondary-700 rounded-xl" />
        <div className="h-48 bg-secondary-200 dark:bg-secondary-700 rounded-xl" />
      </div>
    )
  }

  if (compact) {
    return (
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg p-6">
        <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-4">学习进度</h3>
        <div className="space-y-4">
          {progressData.slice(0, 3).map((course) => (
            <div key={course.id} className="flex items-center space-x-3">
              <img
                src={course.thumbnail}
                alt={course.title}
                className="w-12 h-12 rounded-lg object-cover"
              />
              <div className="flex-1 min-w-0">
                <h4 className="text-sm font-medium text-secondary-800 dark:text-white truncate">
                  {course.title}
                </h4>
                <ProgressBar progress={course.progress} size="sm" showPercentage={false} />
              </div>
              <span className="text-sm font-medium text-primary-600 dark:text-primary-400">
                {Math.round(course.progress)}%
              </span>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Stats Overview */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <BookOpen className="h-5 w-5 text-primary-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.totalCourses}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">学习中的课程</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Award className="h-5 w-5 text-yellow-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.completedCourses}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">已完成课程</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Clock className="h-5 w-5 text-blue-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.totalHours}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">学习时长 (小时)</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Calendar className="h-5 w-5 text-purple-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.completedLessons}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">已完成课时</p>
        </div>
      </div>

      {/* Course Progress List */}
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg overflow-hidden">
        <div className="p-6 border-b border-secondary-200 dark:border-secondary-700">
          <h3 className="text-lg font-bold text-secondary-800 dark:text-white">
            学习进度
          </h3>
          <p className="text-sm text-secondary-500 dark:text-secondary-400 mt-1">
            继续你的学习之旅
          </p>
        </div>

        <div className="divide-y divide-secondary-200 dark:divide-secondary-700">
          {progressData.length > 0 ? (
            progressData.map((course) => (
              <div
                key={course.id}
                className="p-4 hover:bg-secondary-50 dark:hover:bg-secondary-800 transition-colors"
              >
                <div className="flex items-center space-x-4">
                  <img
                    src={course.thumbnail}
                    alt={course.title}
                    className="w-20 h-20 rounded-lg object-cover flex-shrink-0"
                  />
                  <div className="flex-1 min-w-0">
                    <h4 className="text-base font-semibold text-secondary-800 dark:text-white truncate">
                      {course.title}
                    </h4>
                    <div className="flex items-center space-x-4 mt-2 text-sm text-secondary-500 dark:text-secondary-400">
                      <span>
                        {course.completedLessons}/{course.totalLessons} 课时
                      </span>
                      <span>•</span>
                      <span>上次学习：{formatDate(course.lastAccessedAt)}</span>
                    </div>
                    <div className="mt-3">
                      <ProgressBar progress={course.progress} size="md" />
                    </div>
                  </div>
                  <div className="flex-shrink-0">
                    <button className="btn-primary text-sm py-2 px-4">
                      继续学习
                      <ChevronRight className="h-4 w-4 ml-1" />
                    </button>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="text-center py-12">
              <BookOpen className="h-12 w-12 text-secondary-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-secondary-800 dark:text-white mb-2">
                还没有开始学习
              </h3>
              <p className="text-secondary-500 dark:text-secondary-400 mb-4">
                去课程列表中选择一门课程开始学习吧
              </p>
              <a href="/courses" className="btn-primary inline-block">
                浏览课程
              </a>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
