import { Link } from 'react-router-dom'
import { ArrowRight, Code, BookOpen, Award, Zap } from 'lucide-react'
import CourseCard, { type Course } from '@/components/CourseCard'
import ProgressBar from '@/components/ProgressBar'

// Mock data
const featuredCourses: Course[] = [
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
]

const learningPaths = [
  {
    id: '1',
    title: '前端开发工程师',
    description: '从 HTML/CSS 到 React，成为专业前端开发者',
    progress: 45,
    totalCourses: 12,
    completedCourses: 5,
    icon: Code,
  },
  {
    id: '2',
    title: 'Python 数据科学家',
    description: '掌握 Python 数据分析、可视化和机器学习',
    progress: 20,
    totalCourses: 10,
    completedCourses: 2,
    icon: BookOpen,
  },
]

export default function HomePage() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Hero Section */}
      <section className="mb-12">
        <div className="bg-gradient-to-r from-primary-600 to-primary-800 rounded-2xl p-8 md:p-12 text-white">
          <div className="max-w-3xl">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              AI 陪你一起学编程
            </h1>
            <p className="text-xl text-primary-100 mb-8">
              智能化学习路径，个性化课程推荐，让编程学习更高效、更有趣
            </p>
            <div className="flex flex-wrap gap-4">
              <Link to="/courses" className="btn-primary bg-white text-primary-600 hover:bg-primary-50">
                开始学习
                <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
              <Link to="/learning-path" className="btn-primary bg-primary-700 hover:bg-primary-600">
                查看学习路径
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="mb-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="card flex items-start space-x-4">
            <div className="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg">
              <Zap className="h-6 w-6 text-primary-600" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-secondary-800 dark:text-white mb-2">
                智能推荐
              </h3>
              <p className="text-secondary-600 dark:text-secondary-400">
                AI 根据你的学习进度和兴趣，推荐最适合的课程
              </p>
            </div>
          </div>

          <div className="card flex items-start space-x-4">
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
              <Code className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-secondary-800 dark:text-white mb-2">
                实战练习
              </h3>
              <p className="text-secondary-600 dark:text-secondary-400">
                在线代码编辑器，边学边练，即时反馈
              </p>
            </div>
          </div>

          <div className="card flex items-start space-x-4">
            <div className="p-3 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
              <Award className="h-6 w-6 text-yellow-600" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-secondary-800 dark:text-white mb-2">
                学习认证
              </h3>
              <p className="text-secondary-600 dark:text-secondary-400">
                完成课程获得证书，展示你的学习成果
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Courses */}
      <section className="mb-12">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-secondary-800 dark:text-white">
            精选课程
          </h2>
          <Link to="/courses" className="text-primary-600 hover:text-primary-700 font-medium flex items-center">
            查看全部
            <ArrowRight className="ml-1 h-4 w-4" />
          </Link>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {featuredCourses.map((course) => (
            <CourseCard key={course.id} course={course} />
          ))}
        </div>
      </section>

      {/* Learning Paths */}
      <section>
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-secondary-800 dark:text-white">
            我的学习路径
          </h2>
          <Link to="/learning-path" className="text-primary-600 hover:text-primary-700 font-medium flex items-center">
            管理路径
            <ArrowRight className="ml-1 h-4 w-4" />
          </Link>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {learningPaths.map((path) => {
            const Icon = path.icon
            return (
              <Link key={path.id} to="/learning-path" className="card block">
                <div className="flex items-start space-x-4">
                  <div className="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg">
                    <Icon className="h-6 w-6 text-primary-600" />
                  </div>
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-secondary-800 dark:text-white mb-2">
                      {path.title}
                    </h3>
                    <p className="text-secondary-600 dark:text-secondary-400 text-sm mb-4">
                      {path.description}
                    </p>
                    <ProgressBar
                      progress={path.progress}
                      label={`${path.completedCourses}/${path.totalCourses} 课程完成`}
                      size="sm"
                    />
                  </div>
                </div>
              </Link>
            )
          })}
        </div>
      </section>
    </div>
  )
}
