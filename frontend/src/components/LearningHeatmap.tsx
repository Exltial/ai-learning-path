import { useState, useEffect } from 'react'
import { Calendar, TrendingUp, Flame, Clock, Award } from 'lucide-react'
import { api } from '@/services/api'

interface HeatmapData {
  date: string
  count: number
  level: number
}

interface DailyStats {
  date: string
  total_seconds: number
  lessons_completed: number
  courses_accessed: number
}

interface LearningHeatmapProps {
  userId?: string
  months?: number
  compact?: boolean
}

interface HeatmapWeek {
  days: HeatmapDay[]
}

interface HeatmapDay {
  date: string
  count: number
  level: number
  dayOfWeek: number
}

// Color levels for heatmap (GitHub-style)
const HEATMAP_COLORS = [
  'bg-secondary-200 dark:bg-secondary-700', // Level 0 - no activity
  'bg-green-200 dark:bg-green-900',         // Level 1 - light activity
  'bg-green-400 dark:bg-green-700',         // Level 2 - moderate activity
  'bg-green-600 dark:bg-green-500',         // Level 3 - heavy activity
  'bg-green-800 dark:bg-green-300',         // Level 4 - intense activity
]

const HEATMAP_LEVELS = [
  'No activity',
  'Light (1-15 min)',
  'Moderate (15-30 min)',
  'Heavy (30-60 min)',
  'Intense (60+ min)',
]

const DAYS_OF_WEEK = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

export default function LearningHeatmap({ months = 6, compact = false }: LearningHeatmapProps) {
  const [heatmapData, setHeatmapData] = useState<HeatmapData[]>([])
  const [dailyStats, setDailyStats] = useState<DailyStats[]>([])
  const [loading, setLoading] = useState(true)
  const [totalStats, setTotalStats] = useState({
    totalMinutes: 0,
    averageDaily: 0,
    maxStreak: 0,
    currentStreak: 0,
    totalActiveDays: 0,
    bestDay: '',
    bestDayMinutes: 0,
  })

  useEffect(() => {
    loadHeatmapData()
  }, [months])

  const loadHeatmapData = async () => {
    setLoading(true)
    try {
      const [heatmapResponse, statsResponse] = await Promise.all([
        api.get(`/api/v1/progress/heatmap?months=${months}`),
        api.get(`/api/v1/progress/daily-stats?days=${months * 30}`),
      ])

      if (heatmapResponse.success && heatmapResponse.data?.heatmap) {
        setHeatmapData(heatmapResponse.data.heatmap)
      }

      if (statsResponse.success && statsResponse.data?.daily_stats) {
        setDailyStats(statsResponse.data.daily_stats)
        calculateStats(statsResponse.data.daily_stats)
      }
    } catch (error) {
      console.error('Failed to load heatmap data:', error)
      // Use mock data for development
      const mockData = generateMockHeatmapData(months)
      setHeatmapData(mockData)
      calculateStats(mockData.map(d => ({
        date: d.date,
        total_seconds: d.count * 60,
        lessons_completed: Math.floor(d.count / 10),
        courses_accessed: d.count > 0 ? 1 : 0,
      })))
    }
    setLoading(false)
  }

  const generateMockHeatmapData = (months: number): HeatmapData[] => {
    const data: HeatmapData[] = []
    const today = new Date()
    const startDate = new Date(today)
    startDate.setMonth(startDate.getMonth() - months)

    for (let d = new Date(startDate); d <= today; d.setDate(d.getDate() + 1)) {
      const dateStr = d.toISOString().split('T')[0]
      // Random activity with higher chance on weekdays
      const isWeekend = d.getDay() === 0 || d.getDay() === 6
      const baseChance = isWeekend ? 0.3 : 0.7
      const hasActivity = Math.random() < baseChance
      
      if (hasActivity) {
        const minutes = Math.floor(Math.random() * 120) + 5
        const level = minutes < 15 ? 1 : minutes < 30 ? 2 : minutes < 60 ? 3 : 4
        data.push({
          date: dateStr,
          count: minutes,
          level,
        })
      } else {
        data.push({
          date: dateStr,
          count: 0,
          level: 0,
        })
      }
    }

    return data
  }

  const calculateStats = (stats: DailyStats[]) => {
    let totalMinutes = 0
    let maxStreak = 0
    let currentStreak = 0
    let activeDays = 0
    let bestDay = ''
    let bestDayMinutes = 0

    const sortedStats = [...stats].sort((a, b) => 
      new Date(a.date).getTime() - new Date(b.date).getTime()
    )

    sortedStats.forEach((stat, index) => {
      const minutes = stat.total_seconds / 60
      totalMinutes += minutes

      if (minutes > 0) {
        activeDays++
        currentStreak++
        
        if (minutes > bestDayMinutes) {
          bestDay = stat.date
          bestDayMinutes = minutes
        }
      } else {
        if (currentStreak > maxStreak) {
          maxStreak = currentStreak
        }
        currentStreak = 0
      }
    })

    // Check if current streak is still active
    if (currentStreak > maxStreak) {
      maxStreak = currentStreak
    }

    const averageDaily = activeDays > 0 ? totalMinutes / activeDays : 0

    setTotalStats({
      totalMinutes: Math.round(totalMinutes),
      averageDaily: Math.round(averageDaily),
      maxStreak,
      currentStreak,
      totalActiveDays: activeDays,
      bestDay,
      bestDayMinutes: Math.round(bestDayMinutes),
    })
  }

  const organizeHeatmapData = (): HeatmapWeek[] => {
    const weeks: HeatmapWeek[] = []
    let currentWeek: HeatmapDay[] = []

    // Fill in missing dates
    const dateMap = new Map(heatmapData.map(d => [d.date, d]))
    const today = new Date()
    const startDate = new Date(today)
    startDate.setMonth(startDate.getMonth() - months)

    // Pad the first week with empty days if needed
    const startDayOfWeek = startDate.getDay()
    for (let i = 0; i < startDayOfWeek; i++) {
      const padDate = new Date(startDate)
      padDate.setDate(padDate.getDate() - (startDayOfWeek - i))
      currentWeek.push({
        date: padDate.toISOString().split('T')[0],
        count: 0,
        level: 0,
        dayOfWeek: i,
      })
    }

    for (let d = new Date(startDate); d <= today; d.setDate(d.getDate() + 1)) {
      const dateStr = d.toISOString().split('T')[0]
      const data = dateMap.get(dateStr) || { date: dateStr, count: 0, level: 0 }
      
      currentWeek.push({
        date: data.date,
        count: data.count,
        level: data.level,
        dayOfWeek: d.getDay(),
      })

      if (d.getDay() === 6) {
        weeks.push({ days: [...currentWeek] })
        currentWeek = []
      }
    }

    // Add remaining days
    if (currentWeek.length > 0) {
      while (currentWeek.length < 7) {
        const lastDate = new Date(currentWeek[currentWeek.length - 1].date)
        lastDate.setDate(lastDate.getDate() + 1)
        currentWeek.push({
          date: lastDate.toISOString().split('T')[0],
          count: 0,
          level: 0,
          dayOfWeek: currentWeek.length,
        })
      }
      weeks.push({ days: currentWeek })
    }

    return weeks
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
  }

  const formatMonthLabel = (weekDays: HeatmapDay[]) => {
    if (weekDays.length === 0) return ''
    const firstDay = new Date(weekDays[0].date)
    return MONTHS[firstDay.getMonth()]
  }

  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="h-48 bg-secondary-200 dark:bg-secondary-700 rounded-xl" />
      </div>
    )
  }

  const weeks = organizeHeatmapData()

  if (compact) {
    return (
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-bold text-secondary-800 dark:text-white flex items-center">
            <Flame className="h-5 w-5 mr-2 text-orange-500" />
            学习热力图
          </h3>
          <span className="text-sm text-secondary-500 dark:text-secondary-400">
            过去 {months} 个月
          </span>
        </div>

        {/* Stats Summary */}
        <div className="grid grid-cols-3 gap-4 mb-4">
          <div className="text-center">
            <p className="text-2xl font-bold text-primary-600 dark:text-primary-400">
              {totalStats.totalMinutes}
            </p>
            <p className="text-xs text-secondary-500 dark:text-secondary-400">总分钟数</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold text-green-600 dark:text-green-400">
              {totalStats.currentStreak}
            </p>
            <p className="text-xs text-secondary-500 dark:text-secondary-400">当前连续</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold text-purple-600 dark:text-purple-400">
              {totalStats.totalActiveDays}
            </p>
            <p className="text-xs text-secondary-500 dark:text-secondary-400">活跃天数</p>
          </div>
        </div>

        {/* Mini Heatmap */}
        <div className="overflow-x-auto">
          <div className="flex gap-1 min-w-max">
            {weeks.slice(-8).map((week, weekIdx) => (
              <div key={weekIdx} className="flex flex-col gap-1">
                {week.days.map((day, dayIdx) => (
                  <div
                    key={day.date}
                    className={`w-3 h-3 rounded-sm ${HEATMAP_COLORS[day.level]}`}
                    title={`${formatDate(day.date)}: ${day.count} 分钟`}
                  />
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-secondary-800 dark:text-white flex items-center">
            <Calendar className="h-6 w-6 mr-2 text-primary-600" />
            学习热力图
          </h2>
          <p className="text-secondary-500 dark:text-secondary-400 mt-1">
            过去 {months} 个月的学习活动
          </p>
        </div>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Clock className="h-5 w-5 text-blue-600" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {Math.round(totalStats.totalMinutes / 60)}h {totalStats.totalMinutes % 60}m
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">总学习时长</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <TrendingUp className="h-5 w-5 text-green-600" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.averageDaily}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">日均分钟数</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Flame className="h-5 w-5 text-orange-600" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.currentStreak}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">当前连续天数</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Award className="h-5 w-5 text-purple-600" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.maxStreak}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">最长连续天数</p>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow p-4">
          <div className="flex items-center justify-between mb-2">
            <Calendar className="h-5 w-5 text-pink-600" />
          </div>
          <p className="text-2xl font-bold text-secondary-800 dark:text-white">
            {totalStats.totalActiveDays}
          </p>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">活跃天数</p>
        </div>
      </div>

      {/* Heatmap */}
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg p-6">
        {/* Month labels */}
        <div className="flex mb-2 ml-10">
          {weeks.map((week, idx) => {
            const monthLabel = formatMonthLabel(week.days)
            const showLabel = idx === 0 || 
              (idx > 0 && formatMonthLabel(weeks[idx - 1].days) !== monthLabel)
            return (
              <div key={idx} className="flex-1 text-xs text-secondary-500 dark:text-secondary-400">
                {showLabel ? monthLabel : ''}
              </div>
            )
          })}
        </div>

        {/* Heatmap grid */}
        <div className="flex gap-1 overflow-x-auto pb-4">
          {/* Day labels */}
          <div className="flex flex-col gap-1 pt-4">
            {['Sun', 'Mon', 'Wed', 'Fri'].map((day, idx) => (
              <div 
                key={day} 
                className="h-3 flex items-center text-xs text-secondary-500 dark:text-secondary-400"
                style={{ height: '12px' }}
              >
                {idx === 0 ? 'Sun' : idx === 1 ? 'Mon' : idx === 2 ? 'Wed' : 'Fri'}
              </div>
            ))}
          </div>

          {/* Week columns */}
          {weeks.map((week, weekIdx) => (
            <div key={weekIdx} className="flex flex-col gap-1">
              {week.days.map((day, dayIdx) => (
                <div
                  key={day.date}
                  className={`w-3 h-3 rounded-sm ${HEATMAP_COLORS[day.level]} 
                    transition-all hover:scale-125 cursor-pointer
                    ${day.count > 0 ? 'hover:ring-2 hover:ring-primary-500' : ''}`}
                  title={`${formatDate(day.date)}: ${day.count} 分钟`}
                  style={{ height: '12px' }}
                />
              ))}
            </div>
          ))}
        </div>

        {/* Legend */}
        <div className="flex items-center justify-between mt-4 pt-4 border-t border-secondary-200 dark:border-secondary-700">
          <div className="flex items-center gap-2">
            <span className="text-xs text-secondary-500 dark:text-secondary-400">Less</span>
            <div className="flex gap-1">
              {HEATMAP_COLORS.map((color, idx) => (
                <div
                  key={idx}
                  className={`w-3 h-3 rounded-sm ${color}`}
                  title={HEATMAP_LEVELS[idx]}
                />
              ))}
            </div>
            <span className="text-xs text-secondary-500 dark:text-secondary-400">More</span>
          </div>
          
          {totalStats.bestDay && (
            <div className="text-xs text-secondary-500 dark:text-secondary-400">
              🏆 最佳日期：{formatDate(totalStats.bestDay)} ({totalStats.bestDayMinutes} 分钟)
            </div>
          )}
        </div>
      </div>

      {/* Daily Stats Chart */}
      {dailyStats.length > 0 && (
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg p-6">
          <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-4">
            每日学习时长
          </h3>
          <div className="h-48 flex items-end gap-1 overflow-x-auto">
            {dailyStats.slice(-30).map((stat) => {
              const minutes = stat.total_seconds / 60
              const heightPercent = Math.min((minutes / 120) * 100, 100)
              return (
                <div
                  key={stat.date}
                  className="flex-1 min-w-[4px] flex flex-col items-center group"
                >
                  <div 
                    className="w-full bg-primary-500 dark:bg-primary-400 rounded-t transition-all group-hover:bg-primary-600"
                    style={{ height: `${Math.max(heightPercent, 2)}%` }}
                    title={`${formatDate(stat.date)}: ${Math.round(minutes)} 分钟`}
                  />
                  <div className="text-[8px] text-secondary-500 dark:text-secondary-400 mt-1 rotate-45 origin-top-left">
                    {new Date(stat.date).getDate()}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}
