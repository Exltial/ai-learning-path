import { BookOpen, Code, Database, Award, Clock, CheckCircle } from 'lucide-react'
import ProgressBar from '@/components/ProgressBar'
import CourseCard, { type Course } from '@/components/CourseCard'

const learningPaths = [
  {
    id: '1',
    title: '前端开发工程师',
    description: '从 HTML/CSS 基础到 React 高级应用，成为专业前端开发者',
    icon: Code,
    color: 'blue',
    totalCourses: 12,
    completedCourses: 5,
    totalHours: 120,
    courses: [
      {
        id: '101',
        title: 'HTML & CSS 基础',
        description: '掌握网页结构和样式设计基础',
        thumbnail: 'https://picsum.photos/seed/html/400/200',
        level: 'beginner' as const,
        duration: '10 小时',
        students: 2000,
        rating: 4.7,
        lessons: 30,
        tags: ['HTML', 'CSS', '前端基础'],
      },
      {
        id: '102',
        title: 'JavaScript 核心',
        description: '深入理解 JavaScript 语言和核心概念',
        thumbnail: 'https://picsum.photos/seed/js/400/200',
        level: 'beginner' as const,
        duration: '15 小时',
        students: 1800,
        rating: 4.8,
        lessons: 40,
        tags: ['JavaScript', '编程基础'],
      },
      {
        id: '103',
        title: 'React 前端开发',
        description: '学习现代 React 开发，包括 Hooks、Context 和最佳实践',
        thumbnail: 'https://picsum.photos/seed/react/400/200',
        level: 'intermediate' as const,
        duration: '12 小时',
        students: 892,
        rating: 4.7,
        lessons: 36,
        tags: ['React', '前端', 'JavaScript'],
      },
    ],
  },
  {
    id: '2',
    title: 'Python 数据科学家',
    description: '掌握 Python 数据分析、可视化和机器学习技能',
    icon: Database,
    color: 'green',
    totalCourses: 10,
    completedCourses: 2,
    totalHours: 100,
    courses: [
      {
        id: '201',
        title: 'Python 编程基础',
        description: '从零开始学习 Python，掌握编程基础语法和核心概念',
        thumbnail: 'https://picsum.photos/seed/python/400/200',
        level: 'beginner' as const,
        duration: '8 小时',
        students: 1234,
        rating: 4.8,
        lessons: 24,
        tags: ['Python', '编程基础', '入门'],
      },
      {
        id: '202',
        title: '数据分析与 Pandas',
        description: '使用 Pandas 进行高效的数据处理和分析',
        thumbnail: 'https://picsum.photos/seed/pandas/400/200',
        level: 'intermediate' as const,
        duration: '12 小时',
        students: 756,
        rating: 4.6,
        lessons: 32,
        tags: ['Python', 'Pandas', '数据分析'],
      },
    ],
  },
  {
    id: '3',
    title: 'AI 机器学习工程师',
    description: '从机器学习基础到深度学习，成为 AI 专家',
    icon: Award,
    color: 'purple',
    totalCourses: 15,
    completedCourses: 0,
    totalHours: 180,
    courses: [
      {
        id: '301',
        title: 'AI 机器学习入门',
        description: '了解机器学习基础，使用 Python 和 scikit-learn 构建模型',
        thumbnail: 'https://picsum.photos/seed/ai/400/200',
        level: 'intermediate' as const,
        duration: '15 小时',
        students: 567,
        rating: 4.9,
        lessons: 42,
        tags: ['AI', '机器学习', 'Python'],
      },
    ],
  },
]

export default function LearningPathPage() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-2">
          学习路径
        </h1>
        <p className="text-secondary-600 dark:text-secondary-400">
          系统化的学习路线，帮助你循序渐进地掌握技能
        </p>
      </div>

      {/* Learning Paths */}
      <div className="space-y-12">
        {learningPaths.map((path) => {
          const Icon = path.icon
          const progress = (path.completedCourses / path.totalCourses) * 100

          return (
            <section key={path.id}>
              {/* Path Header */}
              <div className="card mb-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    <div className={`p-4 bg-${path.color}-100 dark:bg-${path.color}-900 rounded-xl`}>
                      <Icon className={`h-8 w-8 text-${path.color}-600`} />
                    </div>
                    <div>
                      <h2 className="text-2xl font-bold text-secondary-800 dark:text-white">
                        {path.title}
                      </h2>
                      <p className="text-secondary-600 dark:text-secondary-400">
                        {path.description}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-2xl font-bold text-secondary-800 dark:text-white">
                      {path.completedCourses}/{path.totalCourses}
                    </div>
                    <div className="text-sm text-secondary-500 dark:text-secondary-400">
                      已完成课程
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                  <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                    <BookOpen className="h-5 w-5 mr-2" />
                    {path.totalCourses} 门课程
                  </div>
                  <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                    <Clock className="h-5 w-5 mr-2" />
                    {path.totalHours} 小时
                  </div>
                  <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                    <Award className="h-5 w-5 mr-2" />
                    完成获得证书
                  </div>
                </div>

                <ProgressBar progress={progress} size="lg" />
              </div>

              {/* Path Courses */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {path.courses.map((course, index) => (
                  <div key={course.id} className="relative">
                    <CourseCard course={course} />
                    {index < path.courses.length - 1 && (
                      <div className="hidden lg:block absolute top-1/2 -right-3 transform -translate-y-1/2 z-10">
                        <div className="w-6 h-0.5 bg-secondary-300 dark:bg-secondary-700" />
                      </div>
                    )}
                  </div>
                ))}
              </div>

              {path.courses.length < path.totalCourses && (
                <div className="mt-6 text-center">
                  <button className="btn-secondary">
                    查看全部 {path.totalCourses} 门课程
                  </button>
                </div>
              )}
            </section>
          )
        })}
      </div>
    </div>
  )
}
