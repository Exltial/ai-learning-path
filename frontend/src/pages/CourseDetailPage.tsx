import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  Clock,
  Users,
  Star,
  BookOpen,
  CheckCircle,
  Play,
  ChevronDown,
  ChevronUp,
  Lock,
  Video,
  FileText,
  HelpCircle,
} from 'lucide-react'
import CodeEditor from '@/components/CodeEditor'
import ProgressBar from '@/components/ProgressBar'
import { type Lesson } from '@/components/LessonCard'
import { api } from '@/services/api'
import { useAuth } from '@/contexts/AuthContext'

const mockCourse = {
  id: '1',
  title: 'Python 编程基础',
  description: '从零开始学习 Python，掌握编程基础语法和核心概念。本课程适合编程初学者，无需任何编程经验。',
  thumbnail: 'https://picsum.photos/seed/python/800/400',
  level: 'beginner' as 'beginner' | 'intermediate' | 'advanced',
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
  requirements: ['无需编程经验', '一台可以安装 Python 的电脑', '学习热情和耐心'],
}

const mockLessons: Lesson[] = [
  { id: '1', title: '课程介绍和环境搭建', duration: '15 分钟', completed: true, locked: false, type: 'video' },
  { id: '2', title: 'Python 基础语法', duration: '30 分钟', completed: true, locked: false, type: 'video' },
  { id: '3', title: '变量和数据类型', duration: '25 分钟', completed: false, locked: false, type: 'video' },
  { id: '4', title: '基础语法练习', duration: '20 分钟', completed: false, locked: false, type: 'exercise' },
  { id: '5', title: '控制流程：条件语句', duration: '30 分钟', completed: false, locked: false, type: 'video' },
  { id: '6', title: '控制流程：循环语句', duration: '35 分钟', completed: false, locked: false, type: 'video' },
  { id: '7', title: '函数的定义和使用', duration: '40 分钟', completed: false, locked: true, type: 'video' },
  { id: '8', title: '模块和包', duration: '30 分钟', completed: false, locked: true, type: 'video' },
]

const syllabusSections = [
  {
    title: '第一章：Python 入门',
    lessons: mockLessons.slice(0, 3),
  },
  {
    title: '第二章：基础语法',
    lessons: mockLessons.slice(3, 6),
  },
  {
    title: '第三章：进阶概念',
    lessons: mockLessons.slice(6, 8),
  },
]

export default function CourseDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const [course, setCourse] = useState(mockCourse)
  const [lessons, setLessons] = useState(mockLessons)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'overview' | 'content' | 'code'>('overview')
  const [expandedSections, setExpandedSections] = useState<Record<number, boolean>>({ 0: true })
  const [codeValue, setCodeValue] = useState('')

  useEffect(() => {
    loadCourseData()
  }, [id])

  const loadCourseData = async () => {
    setLoading(true)
    try {
      const courseResponse = await api.getCourse(id!)
      const lessonsResponse = await api.getCourseLessons(id!)

      if (courseResponse.success && courseResponse.data) {
        setCourse({
          ...courseResponse.data,
          instructor: courseResponse.data.instructor || mockCourse.instructor,
          whatYouLearn: courseResponse.data.whatYouLearn || mockCourse.whatYouLearn,
          requirements: courseResponse.data.requirements || mockCourse.requirements,
        })
      }
      if (lessonsResponse.success && lessonsResponse.data) {
        setLessons(lessonsResponse.data)
      }
    } catch {
      // Use mock data on error
    }
    setLoading(false)
  }

  const completedLessons = lessons.filter((l) => l.completed).length
  const progress = (completedLessons / lessons.length) * 100

  const toggleSection = (index: number) => {
    setExpandedSections((prev) => ({ ...prev, [index]: !prev[index] }))
  }

  const handleStartLearning = () => {
    if (!isAuthenticated) {
      navigate('/login')
      return
    }

    // Find first uncompleted lesson
    const nextLesson = lessons.find((l) => !l.completed && !l.locked)
    if (nextLesson) {
      // In a real app, this would navigate to the lesson player
      setActiveTab('content')
      setExpandedSections({ 0: true })
    }
  }

  const handleRunCode = () => {
    // In a real app, this would send code to backend for execution
    console.log('Running code:', codeValue)
  }

  const handleResetCode = () => {
    setCodeValue('')
  }

  const getLevelLabel = (level: string) => {
    switch (level) {
      case 'beginner':
        return '入门'
      case 'intermediate':
        return '进阶'
      case 'advanced':
        return '高级'
      default:
        return level
    }
  }

  const getLevelColor = (level: string) => {
    switch (level) {
      case 'beginner':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
      case 'intermediate':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
      case 'advanced':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getLessonIcon = (type: Lesson['type']) => {
    switch (type) {
      case 'video':
        return <Video className="h-4 w-4" />
      case 'exercise':
        return <FileText className="h-4 w-4" />
      case 'quiz':
        return <HelpCircle className="h-4 w-4" />
    }
  }

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="animate-pulse space-y-8">
          <div className="bg-white dark:bg-secondary-800 rounded-xl h-96" />
          <div className="space-y-4">
            <div className="h-8 bg-secondary-200 dark:bg-secondary-700 rounded w-1/4" />
            <div className="h-32 bg-secondary-200 dark:bg-secondary-700 rounded" />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Breadcrumb */}
      <nav className="mb-6">
        <Link to="/courses" className="text-primary-600 hover:text-primary-700 text-sm flex items-center">
          ← 返回课程列表
        </Link>
      </nav>

      {/* Course Header */}
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-lg overflow-hidden mb-8">
        <div className="md:flex">
          <div className="md:w-1/3">
            <img
              src={course.thumbnail}
              alt={course.title}
              className="w-full h-64 md:h-full object-cover"
            />
          </div>
          <div className="p-6 md:w-2/3">
            <div className="flex items-center space-x-2 mb-4">
              <span
                className={`px-3 py-1 text-sm font-medium rounded-full ${getLevelColor(course.level)}`}
              >
                {getLevelLabel(course.level)}
              </span>
              <div className="flex items-center text-yellow-500">
                <Star className="h-5 w-5 fill-current" />
                <span className="ml-1 font-medium">{course.rating}</span>
              </div>
              <span className="text-secondary-500 dark:text-secondary-400">
                ({course.students.toLocaleString()} 名学生)
              </span>
            </div>

            <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-4">
              {course.title}
            </h1>

            <p className="text-secondary-600 dark:text-secondary-400 mb-6">{course.description}</p>

            <div className="flex flex-wrap gap-4 mb-6">
              <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                <Clock className="h-5 w-5 mr-2" />
                {course.duration}
              </div>
              <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                <BookOpen className="h-5 w-5 mr-2" />
                {course.lessons} 节课
              </div>
              <div className="flex items-center text-secondary-600 dark:text-secondary-400">
                <Users className="h-5 w-5 mr-2" />
                {course.students.toLocaleString()} 人在学
              </div>
            </div>

            {progress > 0 && (
              <div className="mb-6">
                <ProgressBar progress={progress} label="课程进度" />
              </div>
            )}

            <button onClick={handleStartLearning} className="btn-primary w-full md:w-auto">
              <Play className="h-5 w-5 mr-2" />
              {progress > 0 ? '继续学习' : '开始学习'}
            </button>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="mb-8 border-b border-secondary-200 dark:border-secondary-700">
        <div className="flex space-x-8 overflow-x-auto">
          <button
            onClick={() => setActiveTab('overview')}
            className={`pb-4 font-medium transition-colors whitespace-nowrap ${
              activeTab === 'overview'
                ? 'text-primary-600 border-b-2 border-primary-600'
                : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
            }`}
          >
            课程概述
          </button>
          <button
            onClick={() => setActiveTab('content')}
            className={`pb-4 font-medium transition-colors whitespace-nowrap ${
              activeTab === 'content'
                ? 'text-primary-600 border-b-2 border-primary-600'
                : 'text-secondary-500 hover:text-secondary-700 dark:text-secondary-400'
            }`}
          >
            课程大纲 ({lessons.length}节)
          </button>
          <button
            onClick={() => setActiveTab('code')}
            className={`pb-4 font-medium transition-colors whitespace-nowrap ${
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
              {course.whatYouLearn?.map((item, index) => (
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
              {course.requirements?.map((item, index) => (
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
            {course.instructor && (
              <div className="flex items-center space-x-4">
                <img
                  src={course.instructor.avatar}
                  alt={course.instructor.name}
                  className="h-16 w-16 rounded-full object-cover"
                />
                <div>
                  <h3 className="text-lg font-semibold text-secondary-800 dark:text-white">
                    {course.instructor.name}
                  </h3>
                  <p className="text-secondary-500 dark:text-secondary-400">
                    {course.instructor.title}
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'content' && (
        <div className="card">
          <h2 className="text-xl font-bold text-secondary-800 dark:text-white mb-6">
            课程大纲
          </h2>
          <div className="space-y-4">
            {syllabusSections.map((section, sectionIndex) => (
              <div key={sectionIndex} className="border border-secondary-200 dark:border-secondary-700 rounded-lg overflow-hidden">
                <button
                  onClick={() => toggleSection(sectionIndex)}
                  className="w-full flex items-center justify-between p-4 bg-secondary-50 dark:bg-secondary-800 hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors"
                >
                  <span className="font-medium text-secondary-800 dark:text-white">
                    {section.title}
                  </span>
                  <div className="flex items-center space-x-4">
                    <span className="text-sm text-secondary-500 dark:text-secondary-400">
                      {section.lessons.length} 节
                    </span>
                    {expandedSections[sectionIndex] ? (
                      <ChevronUp className="h-5 w-5 text-secondary-400" />
                    ) : (
                      <ChevronDown className="h-5 w-5 text-secondary-400" />
                    )}
                  </div>
                </button>
                {expandedSections[sectionIndex] && (
                  <div className="divide-y divide-secondary-200 dark:divide-secondary-700">
                    {section.lessons.map((lesson) => (
                      <div
                        key={lesson.id}
                        className="flex items-center justify-between p-4 hover:bg-secondary-50 dark:hover:bg-secondary-800 transition-colors"
                      >
                        <div className="flex items-center space-x-3 flex-1">
                          <div
                            className={`flex-shrink-0 ${
                              lesson.completed
                                ? 'text-green-500'
                                : lesson.locked
                                ? 'text-secondary-400'
                                : 'text-primary-500'
                            }`}
                          >
                            {lesson.completed ? (
                              <CheckCircle className="h-5 w-5" />
                            ) : lesson.locked ? (
                              <Lock className="h-5 w-5" />
                            ) : (
                              getLessonIcon(lesson.type)
                            )}
                          </div>
                          <div className="flex-1">
                            <h4 className="font-medium text-secondary-800 dark:text-white">
                              {lesson.title}
                            </h4>
                            <p className="text-sm text-secondary-500 dark:text-secondary-400">
                              {lesson.duration}
                            </p>
                          </div>
                        </div>
                        {!lesson.locked && (
                          <button className="btn-primary text-sm py-2 px-4">
                            {lesson.completed ? '复习' : '开始'}
                          </button>
                        )}
                      </div>
                    ))}
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
            在下方编辑器中练习今天学到的内容，支持 Python、JavaScript、TypeScript 等多种语言
          </p>
          <CodeEditor
            language="python"
            height="500px"
            value={codeValue}
            onChange={setCodeValue}
            placeholder="# 在这里编写你的 Python 代码\nprint('Hello, AI Learning!')"
          />
          <div className="mt-4 flex justify-end space-x-4">
            <button onClick={handleResetCode} className="btn-secondary">
              重置代码
            </button>
            <button onClick={handleRunCode} className="btn-primary">
              运行代码
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
