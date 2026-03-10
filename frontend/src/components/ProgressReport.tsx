import { useState, useEffect } from 'react'
import { 
  Download, 
  Calendar, 
  TrendingUp, 
  Award, 
  Clock, 
  BookOpen,
  ChevronLeft,
  ChevronRight,
  BarChart3,
  Target
} from 'lucide-react'
import { api } from '@/services/api'

interface CourseProgressSummary {
  course_id: string
  course_title: string
  progress_percent: number
  lessons_completed: number
  total_lessons: number
  time_spent_minutes: number
}

interface DailyStats {
  date: string
  total_seconds: number
  lessons_completed: number
  courses_accessed: number
}

interface WeeklyReport {
  week_start: string
  week_end: string
  total_hours: number
  lessons_completed: number
  courses_progress: CourseProgressSummary[]
  daily_stats: DailyStats[]
  avg_daily_minutes: number
}

interface MonthlyReport {
  month: string
  year: number
  total_hours: number
  lessons_completed: number
  courses_completed: number
  courses_progress: CourseProgressSummary[]
  daily_stats: DailyStats[]
  avg_daily_minutes: number
  best_day: string
  best_day_minutes: number
}

interface ProgressReportProps {
  userId?: string
  reportType?: 'weekly' | 'monthly'
}

export default function ProgressReport({ reportType = 'monthly' }: ProgressReportProps) {
  const [report, setReport] = useState<WeeklyReport | MonthlyReport | null>(null)
  const [loading, setLoading] = useState(true)
  const [offset, setOffset] = useState(0)
  const [exporting, setExporting] = useState(false)

  useEffect(() => {
    loadReport()
  }, [reportType, offset])

  const loadReport = async () => {
    setLoading(true)
    try {
      const endpoint = reportType === 'weekly' 
        ? `/api/v1/progress/reports/weekly?offset=${offset}`
        : `/api/v1/progress/reports/monthly?offset=${offset}`
      
      const response = await api.get(endpoint)
      
      if (response.success && response.data) {
        setReport(response.data)
      }
    } catch (error) {
      console.error('Failed to load report:', error)
      // Use mock data for development
      setReport(generateMockReport())
    }
    setLoading(false)
  }

  const generateMockReport = (): MonthlyReport => {
    const now = new Date()
    const months = ['一月', '二月', '三月', '四月', '五月', '六月', 
                    '七月', '八月', '九月', '十月', '十一月', '十二月']
    
    return {
      month: months[now.getMonth()],
      year: now.getFullYear(),
      total_hours: 24.5,
      lessons_completed: 18,
      courses_completed: 1,
      courses_progress: [
        {
          course_id: '1',
          course_title: 'Python 编程基础',
          progress_percent: 75,
          lessons_completed: 18,
          total_lessons: 24,
          time_spent_minutes: 890,
        },
        {
          course_id: '2',
          course_title: 'React 前端开发',
          progress_percent: 45,
          lessons_completed: 16,
          total_lessons: 36,
          time_spent_minutes: 580,
        },
      ],
      daily_stats: generateMockDailyStats(),
      avg_daily_minutes: 47,
      best_day: new Date().toISOString().split('T')[0],
      best_day_minutes: 125,
    }
  }

  const generateMockDailyStats = (): DailyStats[] => {
    const stats: DailyStats[] = []
    const today = new Date()
    const monthStart = new Date(today.getFullYear(), today.getMonth(), 1)
    
    for (let d = new Date(monthStart); d <= today; d.setDate(d.getDate() + 1)) {
      const isWeekend = d.getDay() === 0 || d.getDay() === 6
      const baseMinutes = isWeekend ? Math.random() * 30 : Math.random() * 90
      const hasActivity = Math.random() > 0.3
      
      stats.push({
        date: d.toISOString().split('T')[0],
        total_seconds: hasActivity ? Math.round(baseMinutes * 60) : 0,
        lessons_completed: hasActivity ? Math.floor(Math.random() * 3) : 0,
        courses_accessed: hasActivity ? Math.floor(Math.random() * 2) + 1 : 0,
      })
    }
    
    return stats
  }

  const handleExport = async (format: 'csv' | 'json' = 'csv') => {
    setExporting(true)
    try {
      const response = await api.post('/api/v1/progress/reports/export', {
        report_type: reportType,
        offset: offset,
        format: format,
      }, {
        responseType: 'blob',
      })
      
      // Create download link
      const blob = new Blob([response.data], { 
        type: format === 'csv' ? 'text/csv' : 'application/json' 
      })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `learning_report_${reportType}_${offset}.${format}`
      link.click()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('Failed to export report:', error)
      alert('导出失败，请重试')
    }
    setExporting(false)
  }

  const handlePrevious = () => {
    setOffset(prev => prev + 1)
  }

  const handleNext = () => {
    if (offset > 0) {
      setOffset(prev => prev - 1)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
  }

  const getReportTitle = () => {
    if (!report) return ''
    
    if (reportType === 'weekly') {
      const weeklyReport = report as WeeklyReport
      return `${formatDate(weeklyReport.week_start)} - ${formatDate(weeklyReport.week_end)}`
    } else {
      const monthlyReport = report as MonthlyReport
      return `${monthlyReport.year}年 ${monthlyReport.month}`
    }
  }

  if (loading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="h-32 bg-secondary-200 dark:bg-secondary-700 rounded-xl" />
        <div className="h-64 bg-secondary-200 dark:bg-secondary-700 rounded-xl" />
        <div className="h-48 bg-secondary-200 dark:bg-secondary-700 rounded-xl" />
      </div>
    )
  }

  if (!report) {
    return (
      <div className="text-center py-12">
        <Calendar className="h-12 w-12 text-secondary-300 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-secondary-800 dark:text-white mb-2">
          暂无学习报告
        </h3>
        <p className="text-secondary-500 dark:text-secondary-400">
          开始学习后，这里会显示你的学习报告
        </p>
      </div>
    )
  }

  const weeklyReport = reportType === 'weekly' ? report as WeeklyReport : null
  const monthlyReport = reportType === 'monthly' ? report as MonthlyReport : null

  return (
    <div className="space-y-6">
      {/* Header with Navigation */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-secondary-800 dark:text-white flex items-center">
            <BarChart3 className="h-6 w-6 mr-2 text-primary-600" />
            学习报告
          </h2>
          <p className="text-secondary-500 dark:text-secondary-400 mt-1">
            {getReportTitle()}
          </p>
        </div>

        <div className="flex items-center gap-2">
          {/* Navigation */}
          <div className="flex items-center bg-white dark:bg-secondary-800 rounded-lg shadow">
            <button
              onClick={handlePrevious}
              className="p-2 hover:bg-secondary-100 dark:hover:bg-secondary-700 rounded-l-lg transition-colors"
            >
              <ChevronLeft className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
            </button>
            <div className="px-4 py-2 text-sm font-medium text-secondary-800 dark:text-white">
              {offset === 0 ? '本期' : offset === 1 ? '上期' : `${offset} 期前`}
            </div>
            <button
              onClick={handleNext}
              disabled={offset <= 0}
              className="p-2 hover:bg-secondary-100 dark:hover:bg-secondary-700 rounded-r-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <ChevronRight className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
            </button>
          </div>

          {/* Export */}
          <button
            onClick={() => handleExport('csv')}
            disabled={exporting}
            className="btn-primary flex items-center gap-2"
          >
            <Download className="h-4 w-4" />
            {exporting ? '导出中...' : '导出 CSV'}
          </button>
        </div>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Clock className="h-5 w-5 text-blue-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {report.total_hours.toFixed(1)}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">总学习时长 (小时)</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <BookOpen className="h-5 w-5 text-green-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {report.lessons_completed}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">完成课时数</p>
        </div>

        {monthlyReport && (
          <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
            <div className="flex items-center justify-between mb-2">
              <Award className="h-5 w-5 text-yellow-600" />
              <TrendingUp className="h-4 w-4 text-green-500" />
            </div>
            <p className="text-2xl font-bold text-secondary-800 dark:text-white">
              {monthlyReport.courses_completed}
            </p>
            <p className="text-sm text-secondary-500 dark:text-secondary-400">完成课程数</p>
          </div>
        )}

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Target className="h-5 w-5 text-purple-600" />
            <TrendingUp className="h-4 w-4 text-green-500" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {Math.round(report.avg_daily_minutes)}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">日均学习 (分钟)</p>
        </div>

        {monthlyReport && monthlyReport.best_day && (
          <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
            <div className="flex items-center justify-between mb-2">
              <Award className="h-5 w-5 text-orange-600" />
            </div>
            <p className="text-2xl font-bold text-secondary-800 dark:text-white">
              {Math.round(monthlyReport.best_day_minutes)}
            </p>
            <p className="text-sm text-secondary-500 dark:text-secondary-400">最佳单日 (分钟)</p>
          </div>
        )}
      </div>

      {/* Course Progress */}
      {report.courses_progress.length > 0 && (
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg overflow-hidden">
          <div className="p-6 border-b border-secondary-200 dark:border-secondary-700">
            <h3 className="text-lg font-bold text-secondary-800 dark:text-white">
              课程进度
            </h3>
            <p className="text-sm text-secondary-500 dark:text-secondary-400 mt-1">
              本期学习的课程详情
            </p>
          </div>

          <div className="divide-y divide-secondary-200 dark:divide-secondary-700">
            {report.courses_progress.map((course) => (
              <div
                key={course.course_id}
                className="p-4 hover:bg-secondary-50 dark:hover:bg-secondary-800 transition-colors"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h4 className="text-base font-semibold text-secondary-800 dark:text-white">
                      {course.course_title}
                    </h4>
                    <div className="flex items-center gap-4 mt-2 text-sm text-secondary-500 dark:text-secondary-400">
                      <span>
                        {course.lessons_completed}/{course.total_lessons} 课时
                      </span>
                      <span>•</span>
                      <span>{course.time_spent_minutes} 分钟</span>
                    </div>
                    <div className="mt-3">
                      <div className="flex items-center justify-between text-sm mb-1">
                        <span className="text-secondary-600 dark:text-secondary-400">进度</span>
                        <span className="font-medium text-primary-600 dark:text-primary-400">
                          {course.progress_percent.toFixed(1)}%
                        </span>
                      </div>
                      <div className="w-full bg-secondary-200 dark:bg-secondary-700 rounded-full h-2">
                        <div
                          className="bg-primary-600 h-2 rounded-full transition-all"
                          style={{ width: `${course.progress_percent}%` }}
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Daily Activity Chart */}
      {report.daily_stats.length > 0 && (
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg p-6">
          <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-4">
            每日学习分布
          </h3>
          <div className="h-64 flex items-end gap-1 overflow-x-auto">
            {report.daily_stats.map((stat) => {
              const minutes = stat.total_seconds / 60
              const maxMinutes = Math.max(...report.daily_stats.map(s => s.total_seconds / 60))
              const heightPercent = maxMinutes > 0 ? (minutes / maxMinutes) * 100 : 0
              const date = new Date(stat.date)
              
              return (
                <div
                  key={stat.date}
                  className="flex-1 min-w-[6px] flex flex-col items-center group"
                >
                  <div className="relative w-full flex flex-col items-center">
                    {/* Tooltip */}
                    <div className="absolute bottom-full mb-2 opacity-0 group-hover:opacity-100 transition-opacity bg-secondary-800 dark:bg-secondary-900 text-white text-xs rounded py-1 px-2 whitespace-nowrap z-10 pointer-events-none">
                      {date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })}
                      <br />
                      {Math.round(minutes)} 分钟
                      {stat.lessons_completed > 0 && (
                        <>{' '}• {stat.lessons_completed} 课时</>
                      )}
                    </div>
                    
                    {/* Bar */}
                    <div 
                      className={`w-full rounded-t transition-all ${
                        minutes > 0 
                          ? 'bg-primary-500 dark:bg-primary-400 group-hover:bg-primary-600' 
                          : 'bg-secondary-200 dark:bg-secondary-700'
                      }`}
                      style={{ height: `${Math.max(heightPercent, 2)}%` }}
                    />
                  </div>
                  
                  {/* Date label */}
                  {date.getDate() % 5 === 0 && (
                    <div className="text-[8px] text-secondary-500 dark:text-secondary-400 mt-1 rotate-45 origin-top-left">
                      {date.getMonth() + 1}/{date.getDate()}
                    </div>
                  )}
                </div>
              )
            })}
          </div>
        </div>
      )}

      {/* Insights */}
      <div className="bg-gradient-to-r from-primary-500 to-primary-600 dark:from-primary-600 dark:to-primary-700 rounded-xl shadow-lg p-6 text-white">
        <h3 className="text-lg font-bold mb-4 flex items-center">
          <TrendingUp className="h-5 w-5 mr-2" />
          学习洞察
        </h3>
        <div className="grid md:grid-cols-2 gap-4">
          <div className="bg-white/10 rounded-lg p-4">
            <p className="text-sm opacity-90 mb-1">平均每天学习</p>
            <p className="text-2xl font-bold">{Math.round(report.avg_daily_minutes)} 分钟</p>
            <p className="text-xs opacity-75 mt-1">
              {report.avg_daily_minutes >= 60 ? '🔥 非常棒！' : 
               report.avg_daily_minutes >= 30 ? '👍 继续保持！' : 
               '💪 可以再多一点！'}
            </p>
          </div>
          
          {monthlyReport && monthlyReport.best_day && (
            <div className="bg-white/10 rounded-lg p-4">
              <p className="text-sm opacity-90 mb-1">最佳学习日</p>
              <p className="text-2xl font-bold">{Math.round(monthlyReport.best_day_minutes)} 分钟</p>
              <p className="text-xs opacity-75 mt-1">
                {formatDate(monthlyReport.best_day)}
              </p>
            </div>
          )}
          
          <div className="bg-white/10 rounded-lg p-4">
            <p className="text-sm opacity-90 mb-1">完成课时</p>
            <p className="text-2xl font-bold">{report.lessons_completed} 个</p>
            <p className="text-xs opacity-75 mt-1">
              {report.lessons_completed >= 10 ? '🎯 目标达成！' : '📚 继续加油！'}
            </p>
          </div>
          
          {monthlyReport && monthlyReport.courses_completed > 0 && (
            <div className="bg-white/10 rounded-lg p-4">
              <p className="text-sm opacity-90 mb-1">完成课程</p>
              <p className="text-2xl font-bold">{monthlyReport.courses_completed} 门</p>
              <p className="text-xs opacity-75 mt-1">
                🏆 恭喜完成！
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
