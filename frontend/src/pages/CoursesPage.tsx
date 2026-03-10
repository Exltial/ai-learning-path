import { useState, useEffect } from 'react'
import { Search, Filter, Grid, List, X } from 'lucide-react'
import CourseCard, { type Course } from '@/components/CourseCard'
import { api } from '@/services/api'

const allTags = ['Python', 'React', 'JavaScript', 'TypeScript', 'Node.js', 'AI', '机器学习', '前端', '后端', '算法', '数据结构']

const mockCourses: Course[] = [
  {
    id: '1',
    title: 'Python 编程基础',
    description: '从零开始学习 Python，掌握编程基础语法和核心概念',
    thumbnail: 'https://picsum.photos/seed/python/400/200',
    level: 'beginner',
    duration: '8 小时',
    students: 1234,
    rating: 4.8,
    lessons: 24,
    tags: ['Python', '编程基础', '入门'],
  },
  {
    id: '2',
    title: 'React 前端开发',
    description: '学习现代 React 开发，包括 Hooks、Context 和最佳实践',
    thumbnail: 'https://picsum.photos/seed/react/400/200',
    level: 'intermediate',
    duration: '12 小时',
    students: 892,
    rating: 4.7,
    lessons: 36,
    tags: ['React', '前端', 'JavaScript'],
  },
  {
    id: '3',
    title: 'AI 机器学习入门',
    description: '了解机器学习基础，使用 Python 和 scikit-learn 构建模型',
    thumbnail: 'https://picsum.photos/seed/ai/400/200',
    level: 'intermediate',
    duration: '15 小时',
    students: 567,
    rating: 4.9,
    lessons: 42,
    tags: ['AI', '机器学习', 'Python'],
  },
  {
    id: '4',
    title: 'TypeScript 进阶',
    description: '深入理解 TypeScript 类型系统和高级特性',
    thumbnail: 'https://picsum.photos/seed/ts/400/200',
    level: 'advanced',
    duration: '10 小时',
    students: 445,
    rating: 4.6,
    lessons: 30,
    tags: ['TypeScript', 'JavaScript', '类型系统'],
  },
  {
    id: '5',
    title: 'Node.js 后端开发',
    description: '使用 Node.js 构建高性能后端服务',
    thumbnail: 'https://picsum.photos/seed/node/400/200',
    level: 'intermediate',
    duration: '14 小时',
    students: 678,
    rating: 4.7,
    lessons: 38,
    tags: ['Node.js', '后端', 'JavaScript'],
  },
  {
    id: '6',
    title: '数据结构与算法',
    description: '掌握核心数据结构和常用算法，提升编程能力',
    thumbnail: 'https://picsum.photos/seed/algo/400/200',
    level: 'advanced',
    duration: '20 小时',
    students: 923,
    rating: 4.9,
    lessons: 50,
    tags: ['算法', '数据结构', '面试'],
  },
  {
    id: '7',
    title: 'Vue.js 3 实战',
    description: '掌握 Vue.js 3 组合式 API 和现代前端开发',
    thumbnail: 'https://picsum.photos/seed/vue/400/200',
    level: 'intermediate',
    duration: '11 小时',
    students: 756,
    rating: 4.8,
    lessons: 32,
    tags: ['Vue', '前端', 'JavaScript'],
  },
  {
    id: '8',
    title: 'Docker 容器化技术',
    description: '学习 Docker 容器化部署和容器编排',
    thumbnail: 'https://picsum.photos/seed/docker/400/200',
    level: 'advanced',
    duration: '9 小时',
    students: 534,
    rating: 4.7,
    lessons: 28,
    tags: ['Docker', 'DevOps', '后端'],
  },
]

export default function CoursesPage() {
  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedLevel, setSelectedLevel] = useState<string>('all')
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [sortBy, setSortBy] = useState<'popular' | 'rating' | 'newest'>('popular')

  useEffect(() => {
    loadCourses()
  }, [])

  const loadCourses = async () => {
    setLoading(true)
    try {
      // Try to fetch from API, fallback to mock data
      const response = await api.getCourses()
      if (response.success && response.data) {
        setCourses(response.data.data)
      } else {
        setCourses(mockCourses)
      }
    } catch {
      setCourses(mockCourses)
    }
    setLoading(false)
  }

  const filteredCourses = courses.filter((course) => {
    const matchesSearch =
      searchTerm === '' ||
      course.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      course.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      course.tags.some((tag) => tag.toLowerCase().includes(searchTerm.toLowerCase()))

    const matchesLevel = selectedLevel === 'all' || course.level === selectedLevel

    const matchesTags =
      selectedTags.length === 0 || selectedTags.every((tag) => course.tags.includes(tag))

    return matchesSearch && matchesLevel && matchesTags
  })

  const sortedCourses = [...filteredCourses].sort((a, b) => {
    if (sortBy === 'popular') {
      return b.students - a.students
    } else if (sortBy === 'rating') {
      return b.rating - a.rating
    } else {
      return a.id.localeCompare(b.id)
    }
  })

  const toggleTag = (tag: string) => {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]
    )
  }

  const clearFilters = () => {
    setSearchTerm('')
    setSelectedLevel('all')
    setSelectedTags([])
  }

  const hasActiveFilters = searchTerm || selectedLevel !== 'all' || selectedTags.length > 0

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-2">全部课程</h1>
        <p className="text-secondary-600 dark:text-secondary-400">
          探索我们的课程库，找到适合你的学习内容
        </p>
      </div>

      {/* Filters */}
      <div className="mb-8 space-y-4">
        {/* Search and Sort */}
        <div className="flex flex-col md:flex-row gap-4">
          {/* Search */}
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-secondary-400" />
            <input
              type="text"
              placeholder="搜索课程名称、描述或标签..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-secondary-300 dark:border-secondary-700 rounded-lg bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>

          {/* Sort */}
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as 'popular' | 'rating' | 'newest')}
            className="px-4 py-2 border border-secondary-300 dark:border-secondary-700 rounded-lg bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          >
            <option value="popular">最受欢迎</option>
            <option value="rating">评分最高</option>
            <option value="newest">最新发布</option>
          </select>

          {/* Level Filter */}
          <div className="flex items-center space-x-2">
            <Filter className="h-5 w-5 text-secondary-400" />
            <select
              value={selectedLevel}
              onChange={(e) => setSelectedLevel(e.target.value)}
              className="px-4 py-2 border border-secondary-300 dark:border-secondary-700 rounded-lg bg-white dark:bg-secondary-800 text-secondary-800 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
              <option value="all">全部难度</option>
              <option value="beginner">入门</option>
              <option value="intermediate">进阶</option>
              <option value="advanced">高级</option>
            </select>
          </div>

          {/* View Mode */}
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setViewMode('grid')}
              className={`p-2 rounded-lg ${
                viewMode === 'grid'
                  ? 'bg-primary-100 text-primary-600 dark:bg-primary-900 dark:text-primary-300'
                  : 'text-secondary-400 hover:bg-secondary-100 dark:hover:bg-secondary-800'
              }`}
            >
              <Grid className="h-5 w-5" />
            </button>
            <button
              onClick={() => setViewMode('list')}
              className={`p-2 rounded-lg ${
                viewMode === 'list'
                  ? 'bg-primary-100 text-primary-600 dark:bg-primary-900 dark:text-primary-300'
                  : 'text-secondary-400 hover:bg-secondary-100 dark:hover:bg-secondary-800'
              }`}
            >
              <List className="h-5 w-5" />
            </button>
          </div>
        </div>

        {/* Tags */}
        <div className="flex flex-wrap gap-2">
          {allTags.map((tag) => (
            <button
              key={tag}
              onClick={() => toggleTag(tag)}
              className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                selectedTags.includes(tag)
                  ? 'bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300'
                  : 'bg-secondary-100 text-secondary-600 dark:bg-secondary-800 dark:text-secondary-400 hover:bg-secondary-200 dark:hover:bg-secondary-700'
              }`}
            >
              {tag}
            </button>
          ))}
        </div>

        {/* Active Filters & Results Count */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 flex-wrap">
            {hasActiveFilters && (
              <>
                <span className="text-sm text-secondary-500 dark:text-secondary-400">
                  活跃筛选：
                </span>
                {searchTerm && (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300">
                    搜索：{searchTerm}
                    <button onClick={() => setSearchTerm('')} className="ml-2 hover:text-primary-900">
                      <X className="h-3 w-3" />
                    </button>
                  </span>
                )}
                {selectedLevel !== 'all' && (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300">
                    难度：{selectedLevel === 'beginner' ? '入门' : selectedLevel === 'intermediate' ? '进阶' : '高级'}
                    <button onClick={() => setSelectedLevel('all')} className="ml-2 hover:text-primary-900">
                      <X className="h-3 w-3" />
                    </button>
                  </span>
                )}
                {selectedTags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300"
                  >
                    {tag}
                    <button onClick={() => toggleTag(tag)} className="ml-2 hover:text-primary-900">
                      <X className="h-3 w-3" />
                    </button>
                  </span>
                ))}
                <button
                  onClick={clearFilters}
                  className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  清除全部
                </button>
              </>
            )}
          </div>
          <p className="text-sm text-secondary-500 dark:text-secondary-400">
            找到 {sortedCourses.length} 门课程
          </p>
        </div>
      </div>

      {/* Course Grid */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="h-40 bg-secondary-200 dark:bg-secondary-700 rounded-t-lg" />
              <div className="p-4 space-y-3">
                <div className="h-4 bg-secondary-200 dark:bg-secondary-700 rounded w-3/4" />
                <div className="h-3 bg-secondary-200 dark:bg-secondary-700 rounded w-full" />
                <div className="h-3 bg-secondary-200 dark:bg-secondary-700 rounded w-2/3" />
              </div>
            </div>
          ))}
        </div>
      ) : sortedCourses.length > 0 ? (
        <div
          className={
            viewMode === 'grid'
              ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
              : 'space-y-4'
          }
        >
          {sortedCourses.map((course) => (
            <CourseCard key={course.id} course={course} />
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <Search className="h-12 w-12 text-secondary-300 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-secondary-800 dark:text-white mb-2">
            没有找到匹配的课程
          </h3>
          <p className="text-secondary-500 dark:text-secondary-400 mb-4">
            尝试调整搜索条件或筛选器
          </p>
          {hasActiveFilters && (
            <button onClick={clearFilters} className="btn-primary">
              清除全部筛选
            </button>
          )}
        </div>
      )}
    </div>
  )
}
