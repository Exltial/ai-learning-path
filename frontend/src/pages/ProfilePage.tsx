import { useState } from 'react'
import { User, Mail, Phone, MapPin, Award, BookOpen, Clock, Star, Settings, Bell, Shield } from 'lucide-react'
import ProgressBar from '@/components/ProgressBar'

const userData = {
  name: '张三',
  email: 'zhangsan@example.com',
  phone: '138****8888',
  location: '北京市朝阳区',
  avatar: 'https://picsum.photos/seed/avatar/200/200',
  bio: '热爱编程，正在学习前端开发和 Python',
  joinDate: '2024-01-15',
  level: '进阶学习者',
  points: 2580,
}

const learningStats = {
  totalLearningTime: '48 小时',
  completedCourses: 5,
  inProgressCourses: 3,
  totalPoints: 2580,
  averageRating: 4.7,
}

const recentCourses = [
  {
    id: '1',
    title: 'React 前端开发',
    progress: 65,
    lastAccessed: '2 小时前',
    thumbnail: 'https://picsum.photos/seed/react/400/200',
  },
  {
    id: '2',
    title: 'Python 编程基础',
    progress: 80,
    lastAccessed: '1 天前',
    thumbnail: 'https://picsum.photos/seed/python/400/200',
  },
  {
    id: '3',
    title: 'TypeScript 进阶',
    progress: 30,
    lastAccessed: '3 天前',
    thumbnail: 'https://picsum.photos/seed/ts/400/200',
  },
]

const achievements = [
  { id: '1', title: '初学者', description: '完成第一门课程', icon: '🎯', unlocked: true },
  { id: '2', title: '持之以恒', description: '连续学习 7 天', icon: '🔥', unlocked: true },
  { id: '3', title: '代码大师', description: '完成 100 次代码练习', icon: '💻', unlocked: true },
  { id: '4', title: '全能学霸', description: '完成 5 个学习路径', icon: '🏆', unlocked: false },
  { id: '5', title: '早起鸟', description: '连续 7 天早上学习', icon: '🌅', unlocked: false },
]

export default function ProfilePage() {
  const [activeTab, setActiveTab] = useState<'profile' | 'settings'>('profile')

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-2">
          个人中心
        </h1>
        <p className="text-secondary-600 dark:text-secondary-400">
          管理你的学习进度和个人信息
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Profile Card */}
        <div className="lg:col-span-1">
          <div className="card text-center">
            <img
              src={userData.avatar}
              alt={userData.name}
              className="h-32 w-32 rounded-full mx-auto mb-4 object-cover border-4 border-primary-100 dark:border-primary-900"
            />
            <h2 className="text-xl font-bold text-secondary-800 dark:text-white">
              {userData.name}
            </h2>
            <p className="text-primary-600 font-medium mb-2">{userData.level}</p>
            <p className="text-secondary-500 dark:text-secondary-400 text-sm mb-4">
              {userData.bio}
            </p>
            <div className="flex justify-center items-center space-x-2 text-secondary-500 dark:text-secondary-400 text-sm mb-6">
              <MapPin className="h-4 w-4" />
              {userData.location}
            </div>
            <div className="bg-secondary-50 dark:bg-secondary-800 rounded-lg p-4 mb-4">
              <div className="text-2xl font-bold text-primary-600">{userData.points}</div>
              <div className="text-sm text-secondary-500 dark:text-secondary-400">学习积分</div>
            </div>
            <button className="btn-secondary w-full">
              <Settings className="h-4 w-4 mr-2" />
              编辑资料
            </button>
          </div>

          {/* Stats */}
          <div className="card mt-6">
            <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-4">
              学习统计
            </h3>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                  <Clock className="h-5 w-5 mr-2" />
                  学习时长
                </div>
                <span className="font-medium text-secondary-800 dark:text-white">
                  {learningStats.totalLearningTime}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                  <BookOpen className="h-5 w-5 mr-2" />
                  完成课程
                </div>
                <span className="font-medium text-secondary-800 dark:text-white">
                  {learningStats.completedCourses}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                  <Award className="h-5 w-5 mr-2" />
                  进行中
                </div>
                <span className="font-medium text-secondary-800 dark:text-white">
                  {learningStats.inProgressCourses}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                  <Star className="h-5 w-5 mr-2" />
                  平均评分
                </div>
                <span className="font-medium text-secondary-800 dark:text-white">
                  {learningStats.averageRating}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="lg:col-span-2">
          {/* Tabs */}
          <div className="mb-6 border-b border-secondary-200 dark:border-secondary-700">
            <div className="flex space-x-8">
              <button
                onClick={() => setActiveTab('profile')}
                className={`pb-4 font-medium transition-colors ${
                  activeTab === 'profile'
                    ? 'text-primary-600 border-b-2 border-primary-600'
                    : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
                }`}
              >
                学习进度
              </button>
              <button
                onClick={() => setActiveTab('settings')}
                className={`pb-4 font-medium transition-colors ${
                  activeTab === 'settings'
                    ? 'text-primary-600 border-b-2 border-primary-600'
                    : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
                }`}
              >
                成就徽章
              </button>
              <button
                onClick={() => setActiveTab('account')}
                className={`pb-4 font-medium transition-colors ${
                  activeTab === 'account'
                    ? 'text-primary-600 border-b-2 border-primary-600'
                    : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
                }`}
              >
                账户设置
              </button>
            </div>
          </div>

          {/* Tab Content */}
          {activeTab === 'profile' && (
            <div className="space-y-6">
              <div className="card">
                <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-4">
                  正在学习
                </h3>
                <div className="space-y-4">
                  {recentCourses.map((course) => (
                    <div
                      key={course.id}
                      className="flex items-center space-x-4 p-4 bg-secondary-50 dark:bg-secondary-800 rounded-lg"
                    >
                      <img
                        src={course.thumbnail}
                        alt={course.title}
                        className="h-16 w-24 object-cover rounded"
                      />
                      <div className="flex-1">
                        <h4 className="font-medium text-secondary-800 dark:text-white mb-2">
                          {course.title}
                        </h4>
                        <ProgressBar progress={course.progress} size="sm" />
                        <p className="text-xs text-secondary-500 dark:text-secondary-400 mt-2">
                          最后学习：{course.lastAccessed}
                        </p>
                      </div>
                      <button className="btn-primary text-sm">继续</button>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {activeTab === 'settings' && (
            <div className="card">
              <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-6">
                成就徽章
              </h3>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {achievements.map((achievement) => (
                  <div
                    key={achievement.id}
                    className={`p-4 rounded-lg text-center border-2 ${
                      achievement.unlocked
                        ? 'bg-yellow-50 border-yellow-200 dark:bg-yellow-900/20 dark:border-yellow-800'
                        : 'bg-secondary-50 border-secondary-200 dark:bg-secondary-800 dark:border-secondary-700 opacity-50'
                    }`}
                  >
                    <div className="text-4xl mb-2">{achievement.icon}</div>
                    <h4 className="font-medium text-secondary-800 dark:text-white text-sm">
                      {achievement.title}
                    </h4>
                    <p className="text-xs text-secondary-500 dark:text-secondary-400 mt-1">
                      {achievement.description}
                    </p>
                    {achievement.unlocked && (
                      <div className="mt-2 text-xs text-green-600 font-medium">
                        已解锁 ✓
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'account' && (
            <div className="card">
              <h3 className="text-lg font-bold text-secondary-800 dark:text-white mb-6">
                账户信息
              </h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-2">
                    姓名
                  </label>
                  <input
                    type="text"
                    defaultValue={userData.name}
                    className="w-full px-4 py-2 border border-secondary-300 dark:border-secondary-700 rounded-lg bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-2">
                    邮箱
                  </label>
                  <div className="flex items-center">
                    <Mail className="h-5 w-5 text-secondary-400 mr-2" />
                    <input
                      type="email"
                      defaultValue={userData.email}
                      className="w-full px-4 py-2 border border-secondary-300 dark:border-secondary-700 rounded-lg bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-2">
                    手机
                  </label>
                  <div className="flex items-center">
                    <Phone className="h-5 w-5 text-secondary-400 mr-2" />
                    <input
                      type="tel"
                      defaultValue={userData.phone}
                      className="w-full px-4 py-2 border border-secondary-300 dark:border-secondary-700 rounded-lg bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white"
                    />
                  </div>
                </div>
                <div className="pt-4">
                  <button className="btn-primary">保存修改</button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
