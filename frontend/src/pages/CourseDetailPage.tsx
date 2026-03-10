import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { Clock, Users, Star, BookOpen, CheckCircle, Play, ChevronDown, ChevronUp } from 'lucide-react'
import CodeEditor from '@/components/CodeEditor'
import ProgressBar from '@/components/ProgressBar'
import LessonCard, { type Lesson } from '@/components/LessonCard'

const courseData = {
  id: '1',
  title: 'Python 编程基础',
  description: '从零开始学习 Python，掌握编程基础语法和核心概念。本课程适合编程初学者，无需任何编程经验。',
  thumbnail: 'https://picsum.photos/seed/python/800/400',
  level: 'beginner' as const,
  duration: '8 小时',
  students: 1234,
  rating: 4.8,
  lessons: 24,
  tags: ['Python', '编程基础', '入门'],
  instructor: {
    name: '张老师',
    title: '资深 Python 开发工程师',
    avatar: 'https://picsum.photos/seed/instructor/100/100',
  },
  whatYouLearn: [
    'Python 基础语法和数据类型',
    '函数和模块的使用',
    '面向对象编程概念',
    '文件操作和异常处理',
    '常用标准库的使用',
  ],
  requirements: [
    '无需编程经验',
    '一台可以安装 Python 的电脑',
    '学习热情和耐心',
  ],
}

const lessons: Lesson[] = [
  { id: '1', title: '课程介绍和环境搭建', duration: '15 分钟', completed: true, locked: false, type: 'video' },
  { id: '2', title: 'Python 基础语法', duration: '30 分钟', completed: true, locked: false, type: 'video' },
  { id: '3', title: '变量和数据类型', duration: '25 分钟', completed: false, locked: false, type: 'video' },
  { id: '4', title: '基础语法练习', duration: '20 分钟', completed: false, locked: false, type: 'exercise' },
  { id: '5', title: '控制流程：条件语句', duration: '30 分钟', completed: false, locked: false, type: 'video' },
  { id: '6', title: '控制流程：循环语句', duration: '35 分钟', completed: false, locked: false, type: 'video' },
  { id: '7', title: '函数的定义和使用', duration: '40 分钟', completed: false, locked: true, type: 'video' },
  { id: '8', title: '模块和包', duration: '30 分钟', completed: false, locked: true, type: 'video' },
]

export default function CourseDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [activeTab, setActiveTab] = useState<'overview' | 'content' | 'code'>('overview')
  const [expandedSections, setExpandedSections] = useState<Record<number, boolean>>({ 0: true })

  const completedLessons = lessons.filter(l => l.completed).length
  const progress = (completedLessons / lessons.length) * 100

  const toggleSection = (index: number) => {
    setExpandedSections(prev => ({ ...prev, [index]: !prev[index] }))
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Breadcrumb */}
      <nav className="mb-6">
        <Link to="/courses" className="text-primary-600 hover:text-primary-700 text-sm">
          ← 返回课程列表
        </Link>
      </nav>

      {/* Course Header */}
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg overflow-hidden mb-8">
        <div className="md:flex">
          <div className="md:w-1/3">
            <img
              src={courseData.thumbnail}
              alt={courseData.title}
              className="w-full h-64 md:h-full object-cover"
            />
          </div>
          <div className="p-6 md:w-2/3">
            <div className="flex items-center space-x-2 mb-4">
              <span className="px-3 py-1 bg-green-100 text-green-800 text-sm font-medium rounded-full">
                入门
              </span>
              <div className="flex items-center text-yellow-500">
                <Star className="h-5 w-5 fill-current" />
                <span className="ml-1 font-medium">{courseData.rating}</span>
              </div>
              <span className="text-secondary-500 dark:text-secondary-400">
                ({courseData.students} 名学生)
              </span>
            </div>

            <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-4">
              {courseData.title}
            </h1>

            <p className="text-secondary-600 dark:text-secondary-400 mb-6">
              {courseData.description}
            </p>

            <div className="flex flex-wrap gap-4 mb-6">
              <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                <Clock className="h-5 w-5 mr-2" />
                {courseData.duration}
              </div>
              <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                <BookOpen className="h-5 w-5 mr-2" />
                {courseData.lessons} 节课
              </div>
              <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                <Users className="h-5 w-5 mr-2" />
                {courseData.students} 人在学
              </div>
            </div>

            <div className="mb-6">
              <ProgressBar progress={progress} label="课程进度" />
            </div>

            <button className="btn-primary w-full md:w-auto">
              <Play className="h-5 w-5 mr-2" />
              继续学习
            </button>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="mb-8 border-b border-secondary-200 dark:border-secondary-700">
        <div className="flex space-x-8">
          <button
            onClick={() => setActiveTab('overview')}
            className={`pb-4 font-medium transition-colors ${
              activeTab === 'overview'
                ? 'text-primary-600 border-b-2 border-primary-600'
                : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
            }`}
          >
            课程概述
          </button>
          <button
            onClick={() => setActiveTab('content')}
            className={`pb-4 font-medium transition-colors ${
              activeTab === 'content'
                ? 'text-primary-600 border-b-2 border-primary-600'
                : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
            }`}
          >
            课程大纲
          </button>
          <button
            onClick={() => setActiveTab('code')}
            className={`pb-4 font-medium transition-colors ${
              activeTab === 'code'
                ? 'text-primary-600 border-b-2 border-primary-600'
                : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
            }`}
          >
            代码练习
          </button>
        </div>
      </div>

      {/* Tab Content */}
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div className="card">
            <h2 className="text-xl font-bold text-secondary-800 dark:text-white mb-4">
              你将学到
            </h2>
            <ul className="space-y-3">
              {courseData.whatYouLearn.map((item, index) => (
                <li key={index} className="flex items-start">
                  <CheckCircle className="h-5 w-5 text-green-500 mr-3 mt-0.5 flex-shrink-0" />
                  <span className="text-secondary-600 dark:text-secondary-400">{item}</span>
                </li>
              ))}
            </ul>
          </div>

          <div className="card">
            <h2 className="text-xl font-bold text-secondary-800 dark:text-white mb-4">
              课程要求
            </h2>
            <ul className="space-y-3">
              {courseData.requirements.map((item, index) => (
                <li key={index} className="flex items-start">
                  <div className="h-2 w-2 bg-secondary-400 rounded-full mr-3 mt-2 flex-shrink-0" />
                  <span className="text-secondary-600 dark:text-secondary-400">{item}</span>
                </li>
              ))}
            </ul>
          </div>

          <div className="card md:col-span-2">
            <h2 className="text-xl font-bold text-secondary-800 dark:text-white mb-4">
              讲师介绍
            </h2>
            <div className="flex items-center space-x-4">
              <img
                src={courseData.instructor.avatar}
                alt={courseData.instructor.name}
                className="h-16 w-16 rounded-full object-cover"
              />
              <div>
                <h3 className="text-lg font-semibold text-secondary-800 dark:text-white">
                  {courseData.instructor.name}
                </h3>
                <p className="text-secondary-500 dark:text-secondary-400">
                  {courseData.instructor.title}
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'content' && (
        <div className="card">
          <h2 className="text-xl font-bold text-secondary-800 dark:text-white mb-6">
            课程大纲
          </h2>
          <div className="space-y-4">
            {lessons.map((lesson, index) => (
              <div key={lesson.id}>
                <button
                  onClick={() => toggleSection(index)}
                  className="w-full flex items-center justify-between p-4 bg-secondary-50 dark:bg-secondary-800 rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors"
                >
                  <span className="font-medium text-secondary-800 dark:text-white">
                    第 {index + 1} 章
                  </span>
                  {expandedSections[index] ? (
                    <ChevronUp className="h-5 w-5 text-secondary-400" />
                  ) : (
                    <ChevronDown className="h-5 w-5 text-secondary-400" />
                  )}
                </button>
                {expandedSections[index] && (
                  <div className="mt-2 space-y-2 pl-4">
                    <LessonCard lesson={lesson} />
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'code' && (
        <div className="card">
          <h2 className="text-xl font-bold text-secondary-800 dark:text-white mb-4">
            代码练习
          </h2>
          <p className="text-secondary-600 dark:text-secondary-400 mb-4">
            在下方编辑器中练习今天学到的内容
          </p>
          <CodeEditor
            language="python"
            height="500px"
            placeholder="# 在这里编写你的 Python 代码\nprint('Hello, AI Learning!')"
          />
          <div className="mt-4 flex justify-end space-x-4">
            <button className="btn-secondary">重置代码</button>
            <button className="btn-primary">运行代码</button>
          </div>
        </div>
      )}
    </div>
  )
}
