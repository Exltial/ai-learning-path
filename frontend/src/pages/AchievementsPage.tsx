import { useState, useEffect } from 'react'
import { Award, Trophy, Flame, BookOpen, CheckCircle, Lock, Star, Medal, Target, TrendingUp } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'

interface Achievement {
  id: string
  name: string
  description: string
  icon_url?: string
  points: number
  achievement_type: string
  tier: string
  criteria: {
    type: string
    threshold: number
  }
  is_unlocked: boolean
  earned_at?: string
  progress_percent: number
  current_count: number
}

interface AchievementSummary {
  user_id: string
  total_achievements: number
  unlocked_achievements: number
  total_points: number
  current_level: number
  current_title: string
  current_streak: number
  longest_streak: number
  rank: number
}

const tierColors = {
  bronze: 'from-amber-700 to-amber-900',
  silver: 'from-gray-400 to-gray-600',
  gold: 'from-yellow-400 to-yellow-600',
  platinum: 'from-cyan-400 to-cyan-600',
  diamond: 'from-purple-400 to-purple-600',
}

const tierLabels = {
  bronze: '青铜',
  silver: '白银',
  gold: '黄金',
  platinum: '铂金',
  diamond: '钻石',
}

const achievementIcons: Record<string, any> = {
  course: BookOpen,
  streak: Flame,
  exercise: Target,
  general: Award,
  social: Star,
  milestone: Trophy,
}

export default function AchievementsPage() {
  const { user } = useAuth()
  const [achievements, setAchievements] = useState<Achievement[]>([])
  const [summary, setSummary] = useState<AchievementSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<'all' | 'unlocked' | 'locked'>('all')
  const [typeFilter, setTypeFilter] = useState<string>('all')

  useEffect(() => {
    fetchAchievements()
    fetchSummary()
  }, [])

  const fetchAchievements = async () => {
    try {
      const response = await fetch('/api/achievements', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })
      const data = await response.json()
      if (data.success) {
        setAchievements(data.data)
      }
    } catch (error) {
      console.error('Failed to fetch achievements:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchSummary = async () => {
    try {
      const response = await fetch('/api/achievements/summary', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })
      const data = await response.json()
      if (data.success) {
        setSummary(data.data)
      }
    } catch (error) {
      console.error('Failed to fetch summary:', error)
    }
  }

  const filteredAchievements = achievements.filter(ach => {
    if (filter === 'unlocked' && !ach.is_unlocked) return false
    if (filter === 'locked' && ach.is_unlocked) return false
    if (typeFilter !== 'all' && ach.achievement_type !== typeFilter) return false
    return true
  })

  const getProgressColor = (percent: number) => {
    if (percent >= 100) return 'bg-green-500'
    if (percent >= 75) return 'bg-blue-500'
    if (percent >= 50) return 'bg-yellow-500'
    if (percent >= 25) return 'bg-orange-500'
    return 'bg-red-500'
  }

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex items-center justify-center py-20">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-2">
          成就系统
        </h1>
        <p className="text-secondary-600 dark:text-secondary-400">
          解锁成就，提升等级，成为学习达人！
        </p>
      </div>

      {/* Summary Cards */}
      {summary && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="card bg-gradient-to-br from-primary-500 to-primary-700 text-white">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-primary-100 text-sm mb-1">当前等级</p>
                <p className="text-3xl font-bold">Lv.{summary.current_level}</p>
                <p className="text-primary-100 text-sm mt-1">{summary.current_title}</p>
              </div>
              <Medal className="h-16 w-16 text-primary-200 opacity-50" />
            </div>
          </div>

          <div className="card bg-gradient-to-br from-yellow-500 to-yellow-700 text-white">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-yellow-100 text-sm mb-1">总积分</p>
                <p className="text-3xl font-bold">{summary.total_points}</p>
                <p className="text-yellow-100 text-sm mt-1">继续学习赚取更多</p>
              </div>
              <Star className="h-16 w-16 text-yellow-200 opacity-50" />
            </div>
          </div>

          <div className="card bg-gradient-to-br from-orange-500 to-orange-700 text-white">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-orange-100 text-sm mb-1">学习 streak</p>
                <p className="text-3xl font-bold">{summary.current_streak} 天</p>
                <p className="text-orange-100 text-sm mt-1">最高：{summary.longest_streak} 天</p>
              </div>
              <Flame className="h-16 w-16 text-orange-200 opacity-50" />
            </div>
          </div>

          <div className="card bg-gradient-to-br from-purple-500 to-purple-700 text-white">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-purple-100 text-sm mb-1">成就进度</p>
                <p className="text-3xl font-bold">{summary.unlocked_achievements}/{summary.total_achievements}</p>
                <p className="text-purple-100 text-sm mt-1">
                  {summary.total_achievements > 0 
                    ? Math.round((summary.unlocked_achievements / summary.total_achievements) * 100) 
                    : 0}% 完成
                </p>
              </div>
              <Trophy className="h-16 w-16 text-purple-200 opacity-50" />
            </div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="card mb-6">
        <div className="flex flex-wrap gap-4">
          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium text-secondary-700 dark:text-secondary-300">状态:</span>
            <div className="flex space-x-2">
              <button
                onClick={() => setFilter('all')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  filter === 'all'
                    ? 'bg-primary-600 text-white'
                    : 'bg-secondary-100 dark:bg-secondary-800 text-secondary-700 dark:text-secondary-300 hover:bg-secondary-200 dark:hover:bg-secondary-700'
                }`}
              >
                全部
              </button>
              <button
                onClick={() => setFilter('unlocked')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center space-x-1 ${
                  filter === 'unlocked'
                    ? 'bg-green-600 text-white'
                    : 'bg-secondary-100 dark:bg-secondary-800 text-secondary-700 dark:text-secondary-300 hover:bg-secondary-200 dark:hover:bg-secondary-700'
                }`}
              >
                <CheckCircle className="h-4 w-4" />
                <span>已解锁</span>
              </button>
              <button
                onClick={() => setFilter('locked')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center space-x-1 ${
                  filter === 'locked'
                    ? 'bg-secondary-600 text-white'
                    : 'bg-secondary-100 dark:bg-secondary-800 text-secondary-700 dark:text-secondary-300 hover:bg-secondary-200 dark:hover:bg-secondary-700'
                }`}
              >
                <Lock className="h-4 w-4" />
                <span>未解锁</span>
              </button>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium text-secondary-700 dark:text-secondary-300">类型:</span>
            <select
              value={typeFilter}
              onChange={(e) => setTypeFilter(e.target.value)}
              className="px-4 py-2 rounded-lg border border-secondary-300 dark:border-secondary-700 bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="all">全部</option>
              <option value="course">课程</option>
              <option value="streak">连续学习</option>
              <option value="exercise">练习</option>
              <option value="milestone">里程碑</option>
              <option value="general">综合</option>
              <option value="social">社交</option>
            </select>
          </div>
        </div>
      </div>

      {/* Achievements Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredAchievements.map((achievement) => {
          const IconComponent = achievementIcons[achievement.achievement_type] || Award
          
          return (
            <div
              key={achievement.id}
              className={`card relative overflow-hidden transition-all duration-300 ${
                achievement.is_unlocked
                  ? 'hover:shadow-lg hover:-translate-y-1'
                  : 'opacity-75'
              }`}
            >
              {/* Tier Badge */}
              <div className={`absolute top-0 right-0 px-3 py-1 rounded-bl-lg text-xs font-medium text-white bg-gradient-to-r ${tierColors[achievement.tier as keyof typeof tierColors]}`}>
                {tierLabels[achievement.tier as keyof typeof tierLabels]}
              </div>

              {/* Icon */}
              <div className={`w-16 h-16 rounded-full flex items-center justify-center mb-4 bg-gradient-to-br ${achievement.is_unlocked ? tierColors[achievement.tier as keyof typeof tierColors] : 'from-gray-300 to-gray-500'}`}>
                <IconComponent className={`h-8 w-8 ${achievement.is_unlocked ? 'text-white' : 'text-gray-600'}`} />
              </div>

              {/* Content */}
              <h3 className={`text-lg font-bold mb-2 ${achievement.is_unlocked ? 'text-secondary-800 dark:text-white' : 'text-secondary-500 dark:text-secondary-400'}`}>
                {achievement.name}
              </h3>
              <p className={`text-sm mb-4 ${achievement.is_unlocked ? 'text-secondary-600 dark:text-secondary-400' : 'text-secondary-400 dark:text-secondary-500'}`}>
                {achievement.description}
              </p>

              {/* Progress */}
              {!achievement.is_unlocked && (
                <div className="mb-4">
                  <div className="flex justify-between text-xs text-secondary-500 dark:text-secondary-400 mb-1">
                    <span>进度：{achievement.current_count}/{achievement.criteria.threshold}</span>
                    <span>{Math.round(achievement.progress_percent)}%</span>
                  </div>
                  <div className="w-full bg-secondary-200 dark:bg-secondary-700 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full transition-all duration-500 ${getProgressColor(achievement.progress_percent)}`}
                      style={{ width: `${Math.min(achievement.progress_percent, 100)}%` }}
                    />
                  </div>
                </div>
              )}

              {/* Footer */}
              <div className="flex items-center justify-between pt-4 border-t border-secondary-200 dark:border-secondary-700">
                <div className="flex items-center space-x-1">
                  <Star className="h-4 w-4 text-yellow-500" />
                  <span className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                    +{achievement.points} 积分
                  </span>
                </div>
                {achievement.is_unlocked ? (
                  <div className="flex items-center space-x-1 text-green-600 dark:text-green-400">
                    <CheckCircle className="h-4 w-4" />
                    <span className="text-sm font-medium">已解锁</span>
                    {achievement.earned_at && (
                      <span className="text-xs text-secondary-500 dark:text-secondary-400 ml-2">
                        {new Date(achievement.earned_at).toLocaleDateString('zh-CN')}
                      </span>
                    )}
                  </div>
                ) : (
                  <div className="flex items-center space-x-1 text-secondary-400 dark:text-secondary-500">
                    <Lock className="h-4 w-4" />
                    <span className="text-sm font-medium">未解锁</span>
                  </div>
                )}
              </div>
            </div>
          )
        })}
      </div>

      {filteredAchievements.length === 0 && (
        <div className="text-center py-12">
          <Award className="h-16 w-16 text-secondary-300 dark:text-secondary-600 mx-auto mb-4" />
          <p className="text-secondary-500 dark:text-secondary-400">
            {filter === 'unlocked' 
              ? '还没有解锁任何成就，继续加油！' 
              : filter === 'locked'
              ? '所有成就都已解锁，太棒了！'
              : '暂无成就数据'}
          </p>
        </div>
      )}

      {/* Tips Section */}
      <div className="card mt-8 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20">
        <div className="flex items-start space-x-4">
          <TrendingUp className="h-6 w-6 text-blue-600 dark:text-blue-400 mt-1" />
          <div>
            <h3 className="font-bold text-secondary-800 dark:text-white mb-2">如何获得成就？</h3>
            <ul className="text-sm text-secondary-600 dark:text-secondary-400 space-y-1">
              <li>• 完成课程和课程来解锁课程类成就</li>
              <li>• 保持每日学习连续性来解锁 streak 成就</li>
              <li>• 完成大量练习题来获得练习达人成就</li>
              <li>• 积累积分可以解锁里程碑成就</li>
              <li>• 每天登录学习平台更新你的学习 streak</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}
