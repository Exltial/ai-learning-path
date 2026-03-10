import { useState, useEffect } from 'react'
import { Trophy, Medal, Award, Crown, TrendingUp, TrendingDown, Minus, Users, Calendar, Star } from 'lucide-react'

interface LeaderboardEntry {
  rank: number
  user_id: string
  username: string
  avatar_url?: string
  total_points: number
  level: number
  title: string
  achievements_count: number
}

interface LeaderboardProps {
  type?: 'weekly' | 'monthly' | 'all_time' | 'friends'
  limit?: number
  showTitle?: boolean
  compact?: boolean
}

const rankIcons = {
  1: Crown,
  2: Medal,
  3: Medal,
}

const rankColors = {
  1: 'from-yellow-400 to-yellow-600',
  2: 'from-gray-400 to-gray-600',
  3: 'from-amber-600 to-amber-800',
}

const typeLabels = {
  weekly: '周榜',
  monthly: '月榜',
  all_time: '总榜',
  friends: '好友榜',
}

export default function Leaderboard({ 
  type = 'all_time', 
  limit = 10, 
  showTitle = true,
  compact = false 
}: LeaderboardProps) {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [userRank, setUserRank] = useState<LeaderboardEntry | null>(null)
  const [selectedType, setSelectedType] = useState(type)

  useEffect(() => {
    fetchLeaderboard(selectedType)
  }, [selectedType])

  const fetchLeaderboard = async (leaderboardType: string) => {
    try {
      setLoading(true)
      const response = await fetch(`/api/leaderboard?type=${leaderboardType}&limit=${limit}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })
      const data = await response.json()
      if (data.success) {
        setEntries(data.data)
        
        // Find current user's rank if in top 100
        const userEntry = data.data.find((entry: LeaderboardEntry) => 
          entry.user_id === localStorage.getItem('user_id')
        )
        if (userEntry) {
          setUserRank(userEntry)
        }
      }
    } catch (error) {
      console.error('Failed to fetch leaderboard:', error)
    } finally {
      setLoading(false)
    }
  }

  const getRankDisplay = (rank: number) => {
    if (rank === 1) {
      return (
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-yellow-400 to-yellow-600 flex items-center justify-center shadow-lg">
          <Crown className="h-6 w-6 text-white" />
        </div>
      )
    }
    if (rank === 2) {
      return (
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-gray-400 to-gray-600 flex items-center justify-center shadow-lg">
          <span className="text-white font-bold text-lg">2</span>
        </div>
      )
    }
    if (rank === 3) {
      return (
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-amber-600 to-amber-800 flex items-center justify-center shadow-lg">
          <span className="text-white font-bold text-lg">3</span>
        </div>
      )
    }
    return (
      <div className="w-10 h-10 rounded-full bg-secondary-100 dark:bg-secondary-800 flex items-center justify-center">
        <span className="font-bold text-secondary-600 dark:text-secondary-400">{rank}</span>
      </div>
    )
  }

  const getTrendIcon = () => {
    // This would come from API in a real implementation
    return null
  }

  if (loading) {
    return (
      <div className="card">
        {showTitle && (
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold text-secondary-800 dark:text-white flex items-center">
              <Trophy className="h-6 w-6 mr-2 text-yellow-500" />
              排行榜
            </h2>
          </div>
        )}
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      {showTitle && (
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-secondary-800 dark:text-white flex items-center">
            <Trophy className="h-6 w-6 mr-2 text-yellow-500" />
            排行榜
          </h2>
          
          {/* Type Selector */}
          <div className="flex space-x-2">
            {(['weekly', 'monthly', 'all_time'] as const).map((t) => (
              <button
                key={t}
                onClick={() => setSelectedType(t)}
                className={`px-3 py-1 rounded-lg text-xs font-medium transition-colors ${
                  selectedType === t
                    ? 'bg-primary-600 text-white'
                    : 'bg-secondary-100 dark:bg-secondary-800 text-secondary-700 dark:text-secondary-300 hover:bg-secondary-200 dark:hover:bg-secondary-700'
                }`}
              >
                {typeLabels[t]}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Leaderboard List */}
      <div className="space-y-3">
        {entries.map((entry, index) => {
          const isTopThree = entry.rank <= 3
          const isCurrentUser = entry.user_id === localStorage.getItem('user_id')
          
          return (
            <div
              key={entry.user_id}
              className={`flex items-center p-4 rounded-lg transition-all ${
                isCurrentUser
                  ? 'bg-primary-50 dark:bg-primary-900/20 border-2 border-primary-500'
                  : isTopThree
                  ? 'bg-gradient-to-r from-yellow-50 to-transparent dark:from-yellow-900/10'
                  : 'bg-secondary-50 dark:bg-secondary-800/50'
              }`}
            >
              {/* Rank */}
              <div className="flex-shrink-0 mr-4">
                {getRankDisplay(entry.rank)}
              </div>

              {/* Avatar */}
              <div className="flex-shrink-0 mr-4">
                <img
                  src={entry.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(entry.username)}&background=random`}
                  alt={entry.username}
                  className="w-12 h-12 rounded-full object-cover border-2 border-white dark:border-secondary-700 shadow-sm"
                />
              </div>

              {/* User Info */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center space-x-2 mb-1">
                  <h3 className={`font-bold truncate ${
                    isCurrentUser ? 'text-primary-700 dark:text-primary-400' : 'text-secondary-800 dark:text-white'
                  }`}>
                    {entry.username}
                  </h3>
                  {entry.rank <= 3 && (
                    <span className={`px-2 py-0.5 rounded text-xs font-medium text-white bg-gradient-to-r ${rankColors[entry.rank as keyof typeof rankColors]}`}>
                      TOP {entry.rank}
                    </span>
                  )}
                  {isCurrentUser && (
                    <span className="px-2 py-0.5 rounded text-xs font-medium bg-primary-600 text-white">
                      我
                    </span>
                  )}
                </div>
                <div className="flex items-center space-x-3 text-xs text-secondary-500 dark:text-secondary-400">
                  <span className="flex items-center">
                    <Award className="h-3 w-3 mr-1" />
                    Lv.{entry.level}
                  </span>
                  <span>{entry.title}</span>
                  <span className="flex items-center">
                    <Star className="h-3 w-3 mr-1" />
                    {entry.achievements_count} 成就
                  </span>
                </div>
              </div>

              {/* Points */}
              <div className="flex-shrink-0 text-right">
                <div className="text-lg font-bold text-secondary-800 dark:text-white flex items-center justify-end">
                  <Star className="h-5 w-5 text-yellow-500 mr-1" />
                  {entry.total_points.toLocaleString()}
                </div>
                <div className="text-xs text-secondary-500 dark:text-secondary-400">
                  积分
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {entries.length === 0 && (
        <div className="text-center py-12">
          <Users className="h-16 w-16 text-secondary-300 dark:text-secondary-600 mx-auto mb-4" />
          <p className="text-secondary-500 dark:text-secondary-400">
            暂无排行榜数据
          </p>
        </div>
      )}

      {/* User's Rank Footer */}
      {userRank && userRank.rank > 10 && (
        <div className="mt-4 pt-4 border-t border-secondary-200 dark:border-secondary-700">
          <div className="flex items-center p-3 bg-primary-50 dark:bg-primary-900/20 rounded-lg border-2 border-primary-500">
            <div className="flex-shrink-0 mr-4">
              <div className="w-10 h-10 rounded-full bg-primary-600 flex items-center justify-center">
                <span className="text-white font-bold">{userRank.rank}</span>
              </div>
            </div>
            <div className="flex-1">
              <p className="font-bold text-primary-700 dark:text-primary-400">{userRank.username}</p>
              <p className="text-xs text-secondary-500 dark:text-secondary-400">
                Lv.{userRank.level} · {userRank.title}
              </p>
            </div>
            <div className="text-right">
              <p className="font-bold text-primary-700 dark:text-primary-400 flex items-center">
                <Star className="h-4 w-4 text-yellow-500 mr-1" />
                {userRank.total_points.toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Info Footer */}
      <div className="mt-4 pt-4 border-t border-secondary-200 dark:border-secondary-700">
        <div className="flex items-center justify-between text-xs text-secondary-500 dark:text-secondary-400">
          <div className="flex items-center space-x-4">
            <span className="flex items-center">
              <Calendar className="h-3 w-3 mr-1" />
              {selectedType === 'weekly' && '每周一重置'}
              {selectedType === 'monthly' && '每月 1 日重置'}
              {selectedType === 'all_time' && '历史累计'}
            </span>
          </div>
          <span>共 {entries.length} 人</span>
        </div>
      </div>
    </div>
  )
}
